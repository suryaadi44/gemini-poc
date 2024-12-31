package dto

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthRequest(username, password string) *AuthRequest {
	return &AuthRequest{
		Username: username,
		Password: password,
	}
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func NewRefreshTokenRequest(refreshToken string) *RefreshTokenRequest {
	return &RefreshTokenRequest{
		RefreshToken: refreshToken,
	}
}
