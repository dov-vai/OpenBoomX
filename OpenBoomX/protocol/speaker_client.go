package protocol

import (
	"fmt"
	"obx/utils"
	"strconv"
	"strings"
)

type ISpeakerClient interface {
	SetCustomEQ(bands string) error
	SetOluvMode(mode string) error
	HandleLightAction(action string, solid bool) error
	SetShutdownTimeout(timeout string) error
	PowerOffSpeaker() error
	SetBluetoothPairing(mode string) error
	SetBeepVolume(volume int) error
	SendMessage(hexMsg string) error
	CloseConnection() error
}

type SpeakerClient struct {
	rfcomm RfcommClient
}

func NewSpeakerClient(rfcomm RfcommClient) *SpeakerClient {
	client := &SpeakerClient{}
	client.rfcomm = rfcomm
	return client
}

// SetCustomEQ accepts a bands argument that is a comma-separated string of 10 integer values representing each band.
func (client *SpeakerClient) SetCustomEQ(bands string) error {
	bandValues := strings.Split(bands, ",")
	if len(bandValues) != 10 {
		return fmt.Errorf("invalid number of EQ bands, must be exactly 10 bands")
	}

	eqData := ""
	for i, band := range bandValues {
		bandValue, err := strconv.Atoi(band)
		if err != nil {
			return fmt.Errorf("invalid EQ band value: %s", band)
		}

		if bandValue < MinBandValue || bandValue > MaxBandValue {
			return fmt.Errorf("EQ band value must be between 0 (-10 dB) and 120 (+10 dB)")
		}

		eqData = fmt.Sprintf("%s%02x", eqData[:i*2], bandValue)
	}
	hexMsg := fmt.Sprintf("efb0450b01%s00fe", eqData)
	return client.SendMessage(hexMsg)
}

func (client *SpeakerClient) SetOluvMode(mode string) error {
	hexMsg, ok := EQModes[mode]
	if !ok {
		return fmt.Errorf("invalid Oluv's EQ mode: %s", mode)
	}
	return client.SendMessage(hexMsg)
}

func (client *SpeakerClient) HandleLightAction(action string, solid bool) error {
	hexMsg, ok := LightActions[action]
	if !ok {
		if len(action) == 6 && utils.IsValidHex(action) {
			mode := "02"
			if solid {
				mode = "01"
			}
			hexMsg = fmt.Sprintf("efb09504%s%s00fe", mode, action)
		} else {
			return fmt.Errorf("invalid light action or RGB value: %s", action)
		}
	}
	return client.SendMessage(hexMsg)
}

func (client *SpeakerClient) SetShutdownTimeout(timeout string) error {
	hexMsg, ok := ShutdownTimeouts[timeout]
	if !ok {
		return fmt.Errorf("invalid shutdown timeout: %s", timeout)
	}
	return client.SendMessage(hexMsg)
}

func (client *SpeakerClient) PowerOffSpeaker() error {
	return client.SendMessage(SpeakerPowerOff)
}

func (client *SpeakerClient) SetBluetoothPairing(mode string) error {
	hexMsg, ok := BluetoothPairing[mode]
	if !ok {
		return fmt.Errorf("invalid Bluetooth pairing mode: %s", mode)
	}
	return client.SendMessage(hexMsg)
}

func (client *SpeakerClient) SetBeepVolume(volume int) error {
	hexMsg, ok := BeepVolumes[volume]
	if !ok {
		return fmt.Errorf("invalid volume level: %d", volume)
	}
	return client.SendMessage(hexMsg)
}

func (client *SpeakerClient) SendMessage(hexMsg string) error {
	return client.rfcomm.SendMessage(hexMsg)
}

func (client *SpeakerClient) CloseConnection() error {
	return client.rfcomm.CloseSocket()
}
