package request

type IceCandidate struct {
	IceCandidate string `json:"iceCandidate"`
	To           string `json:"to"`
}
