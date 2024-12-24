package response

type TurnCredential struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	ExpiresAt int64  `json:"expiresAt"`
}
