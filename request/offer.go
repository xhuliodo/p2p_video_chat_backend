package request

type Offer struct {
	Offer    string `json:"offer"`
	DataMode bool   `json:"dataMode"`
	To       string `json:"to"`
}
