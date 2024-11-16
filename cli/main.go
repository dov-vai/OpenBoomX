package main

import (
	"BoomX/btutils"
	"BoomX/protocol"
	"BoomX/utils"
	"flag"
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
)

func main() {
	const UBoomXName = "EarFun UBOOM X"

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

	address, err := btutils.FindDeviceAddress(adapter, UBoomXName, 5*time.Second)
	utils.Must("find device", err)

	switch {
	case *lightAction != "":
		err = protocol.HandleLightAction(*lightAction, *solidLight, address)
	case *eq != "":
		err = protocol.SetCustomEQ(*eq, address)
	case *oluvMode != "":
		err = protocol.SetOluvMode(*oluvMode, address)
	case *shutdown != "":
		err = protocol.SetShutdownTimeout(*shutdown, address)
	case *poweroff:
		err = protocol.PowerOffSpeaker(address)
	case *pairing != "":
		err = protocol.SetBluetoothPairing(*pairing, address)
	case *volume != -1:
		err = protocol.SetBeepVolume(*volume, address)
	case *custom != "":
		err = protocol.SendMessage(*custom, address)
	default:
		fmt.Println("No valid action specified")
		flag.Usage()
		return
	}
	utils.Must("send message", err)

	fmt.Println("Command executed successfully")
}
