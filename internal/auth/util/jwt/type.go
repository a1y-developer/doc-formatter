package jwt

type TokenClaim struct {
	TokenPath string `json:"token_path"`
}

func NewTokenClaim(tokenPath string) *TokenClaim {
	return &TokenClaim{
		TokenPath: tokenPath,
	}
}
