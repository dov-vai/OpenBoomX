package btutils

import (
	"context"
	"time"

	"tinygo.org/x/bluetooth"
)

func FindDeviceAddress(adapter *bluetooth.Adapter, deviceName string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	deviceAddrCh := make(chan string, 1)

	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if ctx.Err() != nil {
			adapter.StopScan()
			return
		}

		name := device.LocalName()
		if name == deviceName {
			deviceAddrCh <- device.Address.String()
			adapter.StopScan()
		}

	})

	if err != nil {
		return "", nil
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case address := <-deviceAddrCh:
		return address, nil
	}
}
