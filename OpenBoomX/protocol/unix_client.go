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
	fd      int
	address string
}

func NewRfcommClient(address string) (*UnixClient, error) {
	client := &UnixClient{}
	client.address = address
	fd, err := NewRfcommSocket(address, RfcommChannel)
	if err != nil {
		return nil, err
	}
	client.fd = fd
	return client, nil
}

func (client *UnixClient) SendMessage(hexMsg string) error {
	message, err := hex.DecodeString(hexMsg)
	if err != nil {
		return fmt.Errorf("failed to decode hex message: %w", err)
	}

	_, err = unix.Write(client.fd, message)
	if err != nil {
		return err
	}
	return nil
}

func (client *UnixClient) ReceiveMessage(bufferSize int) ([]byte, int, error) {
	buf := make([]byte, bufferSize)
	n, err := unix.Read(client.fd, buf)
	if err != nil {
		return nil, n, err
	}
	return buf, n, nil
}

func (client *UnixClient) CloseSocket() error {
	return unix.Close(client.fd)
}

func NewRfcommSocket(address string, channel uint8) (int, error) {
	addr := str2ba(address)

	fd, err := unix.Socket(syscall.AF_BLUETOOTH, syscall.SOCK_STREAM, unix.BTPROTO_RFCOMM)
	if err != nil {
		return -1, err
	}

	sockAddr := &unix.SockaddrRFCOMM{Addr: addr, Channel: channel}

	err = unix.Connect(fd, sockAddr)
	if err != nil {
		return -1, err
	}

	return fd, nil
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
