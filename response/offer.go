package response

type Offer struct {
	Offer    string `json:"offer"`
	DataMode bool   `json:"dataMode"`
	From     string `json:"from"`
}
