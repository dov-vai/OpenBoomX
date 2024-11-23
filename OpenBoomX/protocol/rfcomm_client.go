package protocol

type RfcommClient interface {
	SendMessage(hexMsg string) error
	CloseSocket() error
}
