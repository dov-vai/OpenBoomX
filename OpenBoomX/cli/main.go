package main

import (
	"flag"
	"fmt"
	"obx/btutils"
	"obx/protocol"
	"obx/utils"
	"runtime"
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

func main() {
	lightAction := flag.String("light", "", "Set light action: 'default', 'off', or RGB hex value")
	solidLight := flag.Bool("solid", false, "Set if the light should be solid. Otherwise it will dance. Must be used with -light.")
	eq := flag.String("eq", "", "Set custom eq bands: 10 comma separated values from 0 (-10 dB) to 120 (+10dB). E.g. 0,0,0,0,0,0,0,0,0,0")
	oluvMode := flag.String("oluv", "", "Set EQ mode: 'studio', 'indoor', 'indoor+', 'outdoor', 'outdoor+', 'boom', 'ground'")
	shutdown := flag.String("shutdown", "", "Set shutdown timeout: '5m', '10m', '30m', '60m', '90m', '120m', 'no'")
	poweroff := flag.Bool("poweroff", false, "Power off the speaker")
	pairing := flag.String("pairing", "", "Enable or disable Bluetooth pairing: 'on' or 'off'")
	volume := flag.Int("volume", -1, "Set beep volume: 0, 25, 50, 75, 100")
	custom := flag.String("custom", "", "Send custom hex message (advanced)")

	flag.Parse()

	adapter := bluetooth.DefaultAdapter
	utils.Must("enable BLE stack", adapter.Enable())

	var address string
	if runtime.GOOS == "windows" {
		var err error
		address, err = btutils.FindDeviceAddress(adapter, protocol.UBoomXName2, 5*time.Second)
		utils.Must("find device", err)
		// FIXME: a hack for getting the correct MAC address of the device, because scanning on windows doesn't seem to work correctly
		address = strings.Replace(address, protocol.UBoomXOUI2, protocol.UBoomXOUI, 1)
	} else {
		var err error
		address, err = btutils.FindDeviceAddress(adapter, protocol.UBoomXName, 5*time.Second)
		utils.Must("find device", err)
	}

	var rfcomm protocol.RfcommClient = protocol.NewRfcommClient(address)
	client := protocol.NewSpeakerClient(rfcomm)

	var err error
	switch {
	case *lightAction != "":
		err = client.HandleLightAction(*lightAction, *solidLight)
	case *eq != "":
		err = client.SetCustomEQ(*eq)
	case *oluvMode != "":
		err = client.SetOluvMode(*oluvMode)
	case *shutdown != "":
		err = client.SetShutdownTimeout(*shutdown)
	case *poweroff:
		err = client.PowerOffSpeaker()
	case *pairing != "":
		err = client.SetBluetoothPairing(*pairing)
	case *volume != -1:
		err = client.SetBeepVolume(*volume)
	case *custom != "":
		err = client.SendMessage(*custom)
	default:
		fmt.Println("No valid action specified")
		flag.Usage()
		return
	}
	utils.Must("send message", err)

	fmt.Println("Command executed successfully")
}
