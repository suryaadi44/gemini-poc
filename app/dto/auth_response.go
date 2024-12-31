package dto

type AuthResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	HakAkses     []string `json:"hakAkses"`
}
