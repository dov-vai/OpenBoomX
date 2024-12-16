package protocol

type RfcommClient interface {
	SendMessage(hexMsg string) error
	ReceiveMessage(bufferSize int) ([]byte, int, error)
	CloseSocket() error
}
