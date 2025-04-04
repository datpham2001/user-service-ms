package dto

type UserLoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserGoogleLoginResponse struct {
	AccessToken string `json:"accessToken"`
	Email       string `json:"email"`
}
