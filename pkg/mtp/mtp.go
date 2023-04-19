package mtp

type RxChannelData struct {
	Rec     []byte
	MtpType string
}

var rxC chan RxChannelData

func SetRxChannel(rxChannel chan RxChannelData) {
	rxC = rxChannel
}
