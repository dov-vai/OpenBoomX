//go:build windows

package protocol

import (
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/sys/windows"
)

type WindowsClient struct {
	handle  windows.Handle
	address string
}

func NewRfcommClient(address string) (*WindowsClient, error) {
	client := &WindowsClient{}
	client.address = address
	handle, err := NewRfcommSocket(address, RfcommChannel)
	if err != nil {
		return nil, err
	}
	client.handle = handle
	return client, nil
}

func (c *WindowsClient) SendMessage(hexMsg string) error {
	message, err := hex.DecodeString(hexMsg)
	if err != nil {
		return fmt.Errorf("failed to decode hex message: %w", err)
	}

	var bytesSent uint32
	var overlapped windows.Overlapped

	err = windows.WSASend(c.handle, &windows.WSABuf{Len: uint32(len(message)), Buf: &message[0]}, 1, &bytesSent, 0, &overlapped, nil)
	if err != nil {
		return fmt.Errorf("WSASend failed: %w", err)
	}

	return nil
}

func (c *WindowsClient) ReceiveMessage(bufferSize int) ([]byte, int, error) {
	buf := make([]byte, bufferSize)
	var bytesReceived uint32
	var flags uint32
	var overlapped windows.Overlapped
	wsaBuf := windows.WSABuf{Len: uint32(bufferSize), Buf: &buf[0]}

	err := windows.WSARecv(c.handle, &wsaBuf, 1, &bytesReceived, &flags, &overlapped, nil)
	if err != nil {
		return nil, int(bytesReceived), fmt.Errorf("WSARecv failed: %w", err)
	}

	return buf, int(bytesReceived), nil
}

func (client *WindowsClient) CloseSocket() error {
	return windows.Closesocket(client.handle)
}

func NewRfcommSocket(address string, channel uint8) (windows.Handle, error) {
	addr, err := addrToUint64(address)
	if err != nil {
		return windows.InvalidHandle, err
	}

	handle, err := windows.Socket(windows.AF_BTH, windows.SOCK_STREAM, windows.BTHPROTO_RFCOMM)
	if err != nil {
		return windows.InvalidHandle, err
	}

	sppGuid, err := windows.GUIDFromString("{00001101-0000-1000-8000-00805f9b34fb}")
	if err != nil {
		return windows.InvalidHandle, err
	}

	sockAddr := &windows.SockaddrBth{BtAddr: addr, ServiceClassId: sppGuid, Port: uint32(channel)}

	err = windows.Connect(handle, sockAddr)
	if err != nil {
		return windows.InvalidHandle, err
	}

	if handle == windows.InvalidHandle {
		return windows.InvalidHandle, fmt.Errorf("invalid handle after connect: %v", err)
	}

	return handle, nil
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
