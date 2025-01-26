package controllers

import (
	"fmt"
	"image/color"
	"log"
	"obx/protocol"
	"obx/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MessageListener interface {
	OnMessage(msg string)
}

type SpeakerController struct {
	client         protocol.ISpeakerClient
	debounceMutex  sync.Mutex
	debounceTimer  *time.Timer
	lastColor      color.NRGBA
	lastColorSolid bool
	firstColorSet  bool
	timeoutMap     []string
	listeners      []MessageListener
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
		sc.notifyListeners(fmt.Sprintf("Failed setting %s mode", mode))
		return
	}

	sc.notifyListeners(fmt.Sprintf("Successfully set %s mode", mode))
}

func (sc *SpeakerController) OnLightOffClicked() {
	err := sc.client.HandleLightAction(protocol.LightOff, false)
	if err != nil {
		log.Printf("OnLightOffClicked failed: %v", err)
		sc.notifyListeners("Failed turning lights off")
		return
	}

	sc.notifyListeners(fmt.Sprintf("Successfully turned lights off"))
}

func (sc *SpeakerController) OnLightDefaultClicked() {
	err := sc.client.HandleLightAction(protocol.LightDefault, false)
	if err != nil {
		log.Printf("OnLightDefaultClicked failed: %v", err)
		sc.notifyListeners("Failed setting default lights")
		return
	}
	sc.notifyListeners("Successfully set default lights")
}

func (sc *SpeakerController) OnColorChanged(color color.NRGBA, solidColor bool) {
	err := sc.client.HandleLightAction(utils.NrgbaToHex(color), solidColor)
	if err != nil {
		log.Printf("HandleLightAction failed: %v", err)
		sc.notifyListeners("Failed setting lights color")
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
		sc.notifyListeners(fmt.Sprintf("Failed setting beep volume to %d", 25*step))
		return
	}
	sc.notifyListeners(fmt.Sprintf("Successfully set beep volume to %d", 25*step))
}

func (sc *SpeakerController) OnOffButtonClicked() {
	err := sc.client.PowerOffSpeaker()
	if err != nil {
		log.Printf("PowerOffSpeaker failed: %v", err)
		sc.notifyListeners("Failed powering off speaker")
		return
	}
	sc.notifyListeners("Successfully powered off speaker")
}

func (sc *SpeakerController) OnVideoModeEnabled() {
	err := sc.client.SetVideoMode(protocol.VideoModeOn)
	if err != nil {
		log.Printf("SetVideoMode failed: %v", err)
		sc.notifyListeners("Failed turning video mode on")
		return
	}
	sc.notifyListeners("Successfully turned video mode on")
}

func (sc *SpeakerController) OnVideoModeDisabled() {
	err := sc.client.SetVideoMode(protocol.VideoModeOff)
	if err != nil {
		log.Printf("SetVideoMode failed: %v", err)
		sc.notifyListeners("Failed turning video mode off")
		return
	}
	sc.notifyListeners("Successfully turned video mode off")
}

func (sc *SpeakerController) OnShutdownStepChanged(step int) {
	err := sc.client.SetShutdownTimeout(sc.timeoutMap[step])
	if err != nil {
		log.Printf("SetShutdownTimeout failed: %v", err)
		sc.notifyListeners(fmt.Sprintf("Failed setting shutdown timeout to %s", sc.timeoutMap[step]))
		return
	}
	sc.notifyListeners(fmt.Sprintf("Successfully set shutdown timeout to %s", sc.timeoutMap[step]))
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
		sc.notifyListeners("Failed setting custom EQ")
		return
	}
}

func (sc *SpeakerController) UpdateBattery(onUpdate func(value int, err error)) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		batteryLevel, err := sc.client.ReadBatteryLevel()

		if err == nil {
			onUpdate(batteryLevel, nil)
			continue
		}

		fmt.Println("Error reading battery level:", err)

		// handling for unix and windows if device disconnected
		if protocol.IsSocketDisconnected(err) {
			onUpdate(0, fmt.Errorf("Is speaker not connected?: %w", err))
			err = sc.client.CloseConnection()
			if err != nil {
				log.Printf("Error closing speaker connection: %v", err)
			}
			break
		}
	}
}

func (sc *SpeakerController) GetFirmwareName() string {
	firmware, err := sc.client.ReadFirmwarePackageName()
	if err != nil {
		log.Println(err)
	}
	return firmware
}

func (sc *SpeakerController) RegisterListener(listener MessageListener) {
	sc.listeners = append(sc.listeners, listener)
}

func (sc *SpeakerController) RemoveListener(listener MessageListener) {
	for i, l := range sc.listeners {
		if l == listener {
			sc.listeners = append(sc.listeners[:i], sc.listeners[i+1:]...)
			break
		}
	}
}

func (sc *SpeakerController) notifyListeners(msg string) {
	for _, listener := range sc.listeners {
		listener.OnMessage(msg)
	}
}
