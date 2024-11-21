//go:build unix

package protocol

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/sys/unix"
	"strconv"
	"strings"
	"syscall"
)

type UnixClient struct {
	address string
}

func NewRfcommClient(address string) *UnixClient {
	return &UnixClient{address: address}
}

func (client *UnixClient) SendMessage(hexMsg string) error {
	message, err := hex.DecodeString(hexMsg)
	if err != nil {
		return fmt.Errorf("failed to decode hex message: %w", err)
	}
	return SendRfcommMsg(message, client.address, RfcommChannel)
}

func SendRfcommMsg(message []byte, address string, channel uint8) error {
	addr := str2ba(address)

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

// str2ba converts MAC address string representation to little-endian byte array
func str2ba(addr string) [6]byte {
	a := strings.Split(addr, ":")
	var b [6]byte
	for i, tmp := range a {
		u, _ := strconv.ParseUint(tmp, 16, 8)
		b[len(b)-1-i] = byte(u)
	}
	return b
}
