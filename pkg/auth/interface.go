package auth

import (
	"user-management-api/internal/db/sqlc"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateAccessToken(user sqlc.User) (string, error)
	GenerateRefreshToken(user sqlc.User) (RefreshToken, error)
	ParseToken(token string) (*jwt.Token, jwt.MapClaims, error)
	DecryptToken(token string) (*EncryptedPayload, error)
	StoreRefreshToken(refreshToken RefreshToken) error
	ValidateRefreshToken(refreshToken string) (RefreshToken, error)
	RevokeRefreshToken(refreshToken string) error
}