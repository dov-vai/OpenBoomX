package bluetooth

import (
	"fmt"
	"obx/protocol"
	"runtime"
	"strings"
	"time"
	"tinygo.org/x/bluetooth"
)

func GetUBoomXAddress() (string, error) {
	adapter := bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		address, err := FindDeviceAddress(adapter, protocol.UBoomXName2, 5*time.Second)
		if err != nil {
			return "", err
		}
		// FIXME: a hack for getting the correct MAC address of the device, because scanning on windows doesn't seem to work correctly
		address = strings.Replace(address, protocol.UBoomXOUI2, protocol.UBoomXOUI, 1)
		return address, nil
	}

	address, err := FindDeviceAddress(adapter, protocol.UBoomXName, 5*time.Second)
	if err != nil {
		return "", err
	}
	return address, nil
}

func ConnectUBoomX() (protocol.ISpeakerClient, error) {
	address, err := GetUBoomXAddress()
	if err != nil {
		err = fmt.Errorf("is speaker not connected?: %w", err)
		return nil, err
	}

	rfcomm, err := protocol.NewRfcommClient(address)
	if err != nil {
		err = fmt.Errorf("is device already connected to speaker?: %w", err)
		return nil, err
	}

	client := protocol.NewSpeakerClient(rfcomm)

	return client, nil
}
