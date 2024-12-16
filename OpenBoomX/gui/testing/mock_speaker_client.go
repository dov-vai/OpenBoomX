package testing

type MockSpeakerClient struct {
}

func (client *MockSpeakerClient) SetCustomEQ(bands string) error {
	return nil
}

func (client *MockSpeakerClient) SetOluvMode(mode string) error {
	return nil
}

func (client *MockSpeakerClient) HandleLightAction(action string, solid bool) error {
	return nil
}

func (client *MockSpeakerClient) SetShutdownTimeout(timeout string) error {
	return nil
}

func (client *MockSpeakerClient) PowerOffSpeaker() error {
	return nil
}

func (client *MockSpeakerClient) SetBluetoothPairing(mode string) error {
	return nil
}

func (client *MockSpeakerClient) SetBeepVolume(volume int) error {
	return nil
}

func (client *MockSpeakerClient) SendMessage(hexMsg string) error {
	return nil
}

func (client *MockSpeakerClient) CloseConnection() error {
	return nil
}

func (client *MockSpeakerClient) ReceiveMessage(bufferSize int) ([]byte, int, error) {
	return nil, 0, nil
}

func (client *MockSpeakerClient) ReadBatteryLevel() (int, error) {
	return 0, nil
}
