package response

type SignUpResponse struct {
	UserID string `json:"user_id"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	ExpiryUnix  int64  `json:"expiry_unix"`
}
