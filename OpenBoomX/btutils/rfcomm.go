package btutils

import (
	"obx/utils"
	"syscall"

	"golang.org/x/sys/unix"
)

func SendRfcommMsg(message []byte, address string, channel uint8) error {
	addr := utils.Str2ba(address)

	fd, err := unix.Socket(syscall.AF_BLUETOOTH, syscall.SOCK_STREAM, unix.BTPROTO_RFCOMM)
	if err != nil {
		return err
	}

	sockAddr := &unix.SockaddrRFCOMM{Addr: addr, Channel: channel}

	err = unix.Connect(fd, sockAddr)
	if err != nil {
		return err
	}
	defer unix.Close(fd)

	_, err = unix.Write(fd, message)
	if err != nil {
		return err
	}

	return nil
}
