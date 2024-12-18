package controllers

import (
	"image/color"
	"log"
	"obx/protocol"
	"obx/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SpeakerController struct {
	client         protocol.ISpeakerClient
	debounceMutex  sync.Mutex
	debounceTimer  *time.Timer
	lastColor      color.NRGBA
	lastColorSolid bool
	firstColorSet  bool
}

const debounceDelay = 200 * time.Millisecond

func NewSpeakerController(client protocol.ISpeakerClient) *SpeakerController {
	return &SpeakerController{client: client}
}

func (sc *SpeakerController) OnModeClicked(mode string) {
	err := sc.client.SetOluvMode(mode)
	if err != nil {
		log.Printf("SetOluvMode failed: %v", err)
	}
}

func (sc *SpeakerController) OnActionClicked(action string) {
	err := sc.client.HandleLightAction(action, false)
	if err != nil {
		log.Printf("HandleLightAction failed: %v", err)
	}
}

func (sc *SpeakerController) OnColorChanged(color color.NRGBA, solidColor bool) {
	// ignore the first update coming from the picker
	if !sc.firstColorSet {
		sc.firstColorSet = true
		return
	}

	sc.debounceMutex.Lock()
	defer sc.debounceMutex.Unlock()

	sc.lastColor = color
	sc.lastColorSolid = solidColor

	if sc.debounceTimer != nil {
		sc.debounceTimer.Stop()
	}

	sc.debounceTimer = time.AfterFunc(debounceDelay, func() {
		err := sc.client.HandleLightAction(utils.NrgbaToHex(color), solidColor)
		if err != nil {
			log.Printf("HandleLightAction failed: %v", err)
		}
	})
}

func (sc *SpeakerController) OnBeepStepChanged(step int) {
	err := sc.client.SetBeepVolume(25 * step)
	if err != nil {
		log.Printf("SetBeepVolume failed: %v", err)
	}
}

func (sc *SpeakerController) OnOffButtonClicked() {
	err := sc.client.PowerOffSpeaker()
	if err != nil {
		log.Printf("PowerOffSpeaker failed: %v", err)
	}
}

func (sc *SpeakerController) OnPairingOn() {
	err := sc.client.SetBluetoothPairing("on")
	if err != nil {
		log.Printf("SetBluetoothPairing failed: %v", err)
	}
}

func (sc *SpeakerController) OnPairingOff() {
	err := sc.client.SetBluetoothPairing("off")
	if err != nil {
		log.Printf("SetBluetoothPairing failed: %v", err)
	}
}

func (sc *SpeakerController) OnShutdownStepChanged(step int) {
	timeoutMap := []string{"no", "5m", "10m", "30m", "60m", "90m", "120m"}
	err := sc.client.SetShutdownTimeout(timeoutMap[step])
	if err != nil {
		log.Printf("SetShutdownTimeout failed: %v", err)
	}
}

func (sc *SpeakerController) OnEqValuesChanged(values []float32) {
	var sb strings.Builder
	for i, value := range values {
		// convert normalized value to range from 0 to 120
		converted := int((1 - value) * 120)
		sb.WriteString(strconv.Itoa(converted))
		if i != len(values)-1 {
			sb.WriteString(",")
		}
	}
	err := sc.client.SetCustomEQ(sb.String())
	if err != nil {
		log.Printf("SetCustomEQ failed: %v", err)
	}
}
