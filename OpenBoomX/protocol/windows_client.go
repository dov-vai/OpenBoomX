//go:build windows

package protocol

import (
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/sys/windows"
)

type WindowsClient struct {
	fd      int
	address string
}

func NewRfcommClient(address string) (*WindowsClient, err) {
	client := &WindowsClient{}
	client.address = address
	fd, err := NewRfcommClient(address)
	if err != nil {
		return nil, err
	}
	client.fd = fd
	return client, nil
}

func (client *WindowsClient) SendMessage(hexMsg string) error {
	message, err := hex.DecodeString(hexMsg)
	if err != nil {
		return fmt.Errorf("failed to decode hex message: %w", err)
	}
	// FIXME: fails writing, "The parameter is incorrect.", which parameter?!
	_, err = windows.Write(fd, message)
	if err != nil {
		return err
	}
	return nil
}

func (client *WindowsClient) CloseSocket() error {
	return windows.Closesocket(fd)
}

func NewRfcommSocket(address string, channel uint8) (int, error) {
	addr, err := addrToUint64(address)
	if err != nil {
		return err
	}

	fd, err := windows.Socket(windows.AF_BTH, windows.SOCK_STREAM, windows.BTHPROTO_RFCOMM)
	if err != nil {
		return err
	}

	sppGuid, err := windows.GUIDFromString("{00001101-0000-1000-8000-00805f9b34fb}")
	if err != nil {
		return -1, err
	}

	sockAddr := &windows.SockaddrBth{BtAddr: addr, ServiceClassId: sppGuid, Port: uint32(channel)}

	err = windows.Connect(fd, sockAddr)
	if err != nil {
		return -1, err
	}

	return fd, nil
}

func addrToUint64(addr string) (uint64, error) {
	addr = strings.ReplaceAll(addr, ":", "")

	bytes, err := hex.DecodeString(addr)
	if err != nil {
		return 0, err
	}

	var btAddr uint64
	for i := 0; i < len(bytes); i++ {
		btAddr = (btAddr << 8) | uint64(bytes[i])
	}

	return btAddr, nil
}
