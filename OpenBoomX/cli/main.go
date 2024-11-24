package main

import (
	"flag"
	"fmt"
	"obx/protocol"
	"obx/utils"
	"obx/utils/bluetooth"
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

	address, err := bluetooth.GetUBoomXAddress()
	utils.Must("find uboomx address", err)

	rfcomm, err := protocol.NewRfcommClient(address)
	utils.Must("create rfcomm client", err)
	defer rfcomm.CloseSocket()

	client := protocol.NewSpeakerClient(rfcomm)

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
