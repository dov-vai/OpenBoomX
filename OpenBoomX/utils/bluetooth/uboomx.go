package bluetooth

import (
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
