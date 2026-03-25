package v1dto

type LoginInput struct {
	Email    string `json:"email" binding:"required,email,email_advanced"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Expiration int `json:"expiration"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}