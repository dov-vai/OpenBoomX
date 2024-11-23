package controllers

import (
	"image/color"
	"log"
	"obx/protocol"
	"obx/utils"
	"sync"
	"time"
)

type SpeakerController struct {
	client         *protocol.SpeakerClient
	debounceMutex  sync.Mutex
	debounceTimer  *time.Timer
	lastColor      color.NRGBA
	lastColorSolid bool
	firstColorSet  bool
}

const debounceDelay = 200 * time.Millisecond

func NewSpeakerController(client *protocol.SpeakerClient) *SpeakerController {
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
