package protocol

import (
	"encoding/hex"
	"fmt"
	"obx/btutils"
)

type UnixClient struct {
	address string
}

func NewUnixClient(address string) *UnixClient {
	return &UnixClient{address: address}
}

func (client *UnixClient) SendMessage(hexMsg string) error {
	message, err := hex.DecodeString(hexMsg)
	if err != nil {
		return fmt.Errorf("failed to decode hex message: %w", err)
	}
	return btutils.SendRfcommMsg(message, client.address, RfcommChannel)
}
