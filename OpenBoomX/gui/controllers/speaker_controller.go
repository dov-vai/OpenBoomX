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
	timeoutMap     []string
}

const debounceDelay = 200 * time.Millisecond

func NewSpeakerController(client protocol.ISpeakerClient) *SpeakerController {
	return &SpeakerController{
		client:     client,
		timeoutMap: utils.SortedKeysByValue(protocol.ShutdownTimeouts),
	}
}

func (sc *SpeakerController) OnModeClicked(mode string) {
	err := sc.client.SetOluvMode(mode)
	if err != nil {
		log.Printf("SetOluvMode failed: %v", err)
	}
}

func (sc *SpeakerController) OnLightOffClicked() {
	err := sc.client.HandleLightAction(protocol.LightOff, false)
	if err != nil {
		log.Printf("OnLightOffClicked failed: %v", err)
	}
}

func (sc *SpeakerController) OnLightDefaultClicked() {
	err := sc.client.HandleLightAction(protocol.LightDefault, false)
	if err != nil {
		log.Printf("OnLightDefaultClicked failed: %v", err)
	}
}

func (sc *SpeakerController) OnColorChanged(color color.NRGBA, solidColor bool) {
	// ignore the first update coming from the picker
	// it would set the default light picker color
	// on the speaker on app launch, we don't want that
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
	err := sc.client.SetBluetoothPairing(protocol.PairingOn)
	if err != nil {
		log.Printf("SetBluetoothPairing failed: %v", err)
	}
}

func (sc *SpeakerController) OnPairingOff() {
	err := sc.client.SetBluetoothPairing(protocol.PairingOff)
	if err != nil {
		log.Printf("SetBluetoothPairing failed: %v", err)
	}
}

func (sc *SpeakerController) OnShutdownStepChanged(step int) {
	err := sc.client.SetShutdownTimeout(sc.timeoutMap[step])
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
