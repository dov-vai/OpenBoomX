package controllers

import (
	"fmt"
	"image/color"
	"log"
	"obx/gui/components"
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
	snackbar       *components.Snackbar
}

const debounceDelay = 200 * time.Millisecond

func NewSpeakerController(client protocol.ISpeakerClient, snackbar *components.Snackbar) *SpeakerController {
	return &SpeakerController{
		client:     client,
		timeoutMap: utils.SortedKeysByValue(protocol.ShutdownTimeouts),
		snackbar:   snackbar,
	}
}

func (c *SpeakerController) showMessage(msg string) {
	c.snackbar.ShowMessage(msg)
}

func (sc *SpeakerController) OnModeClicked(mode string) {
	err := sc.client.SetOluvMode(mode)
	if err != nil {
		log.Printf("SetOluvMode failed: %v", err)
		sc.showMessage(fmt.Sprintf("Failed setting %s mode", mode))
		return
	}

	sc.showMessage(fmt.Sprintf("Successfully set %s mode", mode))
}

func (sc *SpeakerController) OnLightOffClicked() {
	err := sc.client.HandleLightAction(protocol.LightOff, false)
	if err != nil {
		log.Printf("OnLightOffClicked failed: %v", err)
		sc.showMessage("Failed turning lights off")
		return
	}

	sc.showMessage(fmt.Sprintf("Successfully turned lights off"))
}

func (sc *SpeakerController) OnLightDefaultClicked() {
	err := sc.client.HandleLightAction(protocol.LightDefault, false)
	if err != nil {
		log.Printf("OnLightDefaultClicked failed: %v", err)
		sc.showMessage("Failed setting default lights")
		return
	}
	sc.showMessage("Successfully set default lights")
}

func (sc *SpeakerController) OnColorChanged(color color.NRGBA, solidColor bool) {
	err := sc.client.HandleLightAction(utils.NrgbaToHex(color), solidColor)
	if err != nil {
		log.Printf("HandleLightAction failed: %v", err)
		sc.showMessage("Failed setting lights color")
	}
}

func (sc *SpeakerController) OnColorChangedDebounced(color color.NRGBA, solidColor bool) {
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
		sc.OnColorChanged(color, solidColor)
	})
}

func (sc *SpeakerController) OnBeepStepChanged(step int) {
	err := sc.client.SetBeepVolume(25 * step)
	if err != nil {
		log.Printf("SetBeepVolume failed: %v", err)
		sc.showMessage(fmt.Sprintf("Failed setting beep volume to %d", 25*step))
		return
	}
	sc.showMessage(fmt.Sprintf("Successfully set beep volume to %d", 25*step))
}

func (sc *SpeakerController) OnOffButtonClicked() {
	err := sc.client.PowerOffSpeaker()
	if err != nil {
		log.Printf("PowerOffSpeaker failed: %v", err)
		sc.showMessage("Failed powering off speaker")
		return
	}
	sc.showMessage("Successfully powered off speaker")
}

func (sc *SpeakerController) OnVideoModeEnabled() {
	err := sc.client.SetVideoMode(protocol.VideoModeOn)
	if err != nil {
		log.Printf("SetVideoMode failed: %v", err)
		sc.showMessage("Failed turning video mode on")
		return
	}
	sc.showMessage("Successfully turned video mode on")
}

func (sc *SpeakerController) OnVideoModeDisabled() {
	err := sc.client.SetVideoMode(protocol.VideoModeOff)
	if err != nil {
		log.Printf("SetVideoMode failed: %v", err)
		sc.showMessage("Failed turning video mode off")
		return
	}
	sc.showMessage("Successfully turned video mode off")
}

func (sc *SpeakerController) OnShutdownStepChanged(step int) {
	err := sc.client.SetShutdownTimeout(sc.timeoutMap[step])
	if err != nil {
		log.Printf("SetShutdownTimeout failed: %v", err)
		sc.showMessage(fmt.Sprintf("Failed setting shutdown timeout to %s", sc.timeoutMap[step]))
		return
	}
	sc.showMessage(fmt.Sprintf("Successfully set shutdown timeout to %s", sc.timeoutMap[step]))
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
		sc.showMessage("Failed setting custom EQ")
		return
	}
}
