package protocol

import (
	"encoding/hex"
	"fmt"
	"obx/btutils"
	"obx/utils"
	"strconv"
	"strings"
)

// Oluv's EQ Modes
var EQModes = map[string]string{
	"studio":   "efb046010102fe",
	"indoor":   "efb046010203fe",
	"indoor+":  "efb046010304fe",
	"outdoor":  "efb046010405fe",
	"outdoor+": "efb046010506fe",
	"boom":     "efb046010607fe",
	"ground":   "efb046010708fe",
}

// Light Actions
var LightActions = map[string]string{
	"default": "efb095040000000000fe",
	"off":     "efb095040100000000fe",
}

// Shutdown Timeout Modes
var ShutdownTimeouts = map[string]string{
	"5m":   "efb075010102fe",
	"10m":  "efb075010203fe",
	"30m":  "efb075010304fe",
	"60m":  "efb075010405fe",
	"90m":  "efb075010506fe",
	"120m": "efb075010607fe",
	"no":   "efb07501ff00fe",
}

const SpeakerPowerOff = "efb025010102fe"

var BluetoothPairing = map[string]string{
	"on":  "efb035010102fe",
	"off": "efb035010001fe",
}

// Beep Volume Levels
var BeepVolumes = map[int]string{
	0:   "efb065010102fe",
	25:  "efb065010203fe",
	50:  "efb065010304fe",
	75:  "efb065010405fe",
	100: "efb065010506fe",
}

// EQ Band Values
const (
	MaxBandValue = 120 // +10 dB is the max
	MinBandValue = 0   // -10 dB is the min
)

const RfcommChannel = 2

// The bands argument is a comma-separated string of 10 integer values representing each band.
func SetCustomEQ(bands string, address string) error {
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
	return SendMessage(hexMsg, address)
}

func SetOluvMode(mode string, address string) error {
	hexMsg, ok := EQModes[mode]
	if !ok {
		return fmt.Errorf("invalid Oluv's EQ mode: %s", mode)
	}
	return SendMessage(hexMsg, address)
}

func HandleLightAction(action string, solid bool, address string) error {
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
	return SendMessage(hexMsg, address)
}

func SetShutdownTimeout(timeout string, address string) error {
	hexMsg, ok := ShutdownTimeouts[timeout]
	if !ok {
		return fmt.Errorf("invalid shutdown timeout: %s", timeout)
	}
	return SendMessage(hexMsg, address)
}

func PowerOffSpeaker(address string) error {
	return SendMessage(SpeakerPowerOff, address)
}

func SetBluetoothPairing(mode string, address string) error {
	hexMsg, ok := BluetoothPairing[mode]
	if !ok {
		return fmt.Errorf("invalid Bluetooth pairing mode: %s", mode)
	}
	return SendMessage(hexMsg, address)
}

func SetBeepVolume(volume int, address string) error {
	hexMsg, ok := BeepVolumes[volume]
	if !ok {
		return fmt.Errorf("invalid volume level: %d", volume)
	}
	return SendMessage(hexMsg, address)
}

func SendMessage(hexMsg string, address string) error {
	message, err := hex.DecodeString(hexMsg)
	if err != nil {
		return fmt.Errorf("failed to decode hex message: %w", err)
	}
	return btutils.SendRfcommMsg(message, address, RfcommChannel)
}
