package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/utils"
	"user-management-api/pkg/cache"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct{
	cache cache.RedisCacheService
}



type EncryptedPayload struct {
	UserUUID string `json:"user_uuid"`
	Email    string `json:"email"`
	Role     int `json:"role"`
}

func NewJWTService(cache cache.RedisCacheService) TokenService {
	return &JWTService{
		cache: cache,
	}
}
var (jwtSecret = []byte(utils.GetEnv("JWT_SECRET", "secret"))
	encryptionKey = []byte(utils.GetEnv("JWT_ENCRYPTION_KEY", "encryption_key")))
const (
	AccessTokenExpiration = 1 * time.Hour
	RefreshTokenExpiration = 24 * time.Hour
)

func (js *JWTService) GenerateAccessToken(user sqlc.User) (string, error) {

	payload := &EncryptedPayload{
		UserUUID: user.UserUuid.String(),
		Email: user.UserEmail,
		Role: int(user.UserLevel),
	}


	rawData, err := json.Marshal(payload)
	if err != nil {
		return "", utils.WrapError(err, "Failed to marshal payload", utils.ErrCodeInternal)
	}

	encryptedData, err := utils.EncryptAES(rawData, encryptionKey)
	if err != nil {
		return "", utils.WrapError(err, "Failed to encrypt payload", utils.ErrCodeInternal)
	}

//  claims := &Claims{
// 	UserUUID: user.UserUuid.String(),
// 	Email: user.UserEmail,
// 	Role:  int(user.UserLevel),
// 	RegisteredClaims: jwt.RegisteredClaims{
// 		IssuedAt: jwt.NewNumericDate(time.Now()),
// 		ID: uuid.New().String(),
// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiration)),
// 		Issuer: "user-management-api",
// 	},
		
//  }

claims := jwt.MapClaims{
	"data": encryptedData,
	"jti": uuid.New().String(),
	"iat": jwt.NewNumericDate(time.Now()),
	"exp": jwt.NewNumericDate(time.Now().Add(AccessTokenExpiration)),
	"iss": "user-management-api",
}
 token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
 return token.SignedString(jwtSecret)
}



func (js *JWTService) ParseToken(token string) (*jwt.Token, jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, nil, utils.WrapError(err, "Failed to parse token", utils.ErrCodeInternal)
	}
	

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, utils.NewError("Invalid token claims", utils.ErrCodeInternal)
	}

	return parsedToken, claims, nil
}

func (js *JWTService) DecryptToken(token string) (*EncryptedPayload, error) {
	_, claims, err := js.ParseToken(token)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to decrypt token", utils.ErrCodeInternal)
	}

	encryptedData,ok := claims["data"].(string)
	if !ok {
		return nil, utils.NewError("Invalid token data", utils.ErrCodeInternal)
	}

	userData, err := utils.DecryptAES(encryptedData, encryptionKey)

	if err != nil {
		return nil, utils.WrapError(err, "Failed to decrypt token", utils.ErrCodeInternal)
	}

	var encryptedPayload EncryptedPayload
	if err := json.Unmarshal(userData, &encryptedPayload); err != nil {
		return nil, utils.WrapError(err, "Failed to unmarshal user data", utils.ErrCodeInternal)
	}

	return &encryptedPayload, nil
}

type RefreshToken struct {
	Token string `json:"token"`
	UserUUID string `json:"user_uuid"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked bool `json:"revoked"`
}

func (js *JWTService) GenerateRefreshToken(user sqlc.User) (RefreshToken, error) {
	tokenBytes := make([]byte, 32)
	 _, err := rand.Read(tokenBytes)
	 if err != nil {
		return RefreshToken{}, utils.WrapError(err, "Failed to generate refresh token", utils.ErrCodeInternal)
	 }
	 token := base64.StdEncoding.EncodeToString(tokenBytes)
	 return RefreshToken{
		Token: token,
		UserUUID: user.UserUuid.String(),
		ExpiresAt: time.Now().Add(RefreshTokenExpiration),
		Revoked: false,
	 }, nil
   }

func (js *JWTService) StoreRefreshToken(refreshToken RefreshToken) error {
	key := fmt.Sprintf("refresh_token:%s", refreshToken.Token)
	return js.cache.Set(key, refreshToken, RefreshTokenExpiration)
}

func (js *JWTService) ValidateRefreshToken(refreshToken string) (RefreshToken, error) {
	key := fmt.Sprintf("refresh_token:%s", refreshToken)
 var storedToken RefreshToken
 err := js.cache.Get(key, &storedToken)
 if err != nil {
	return RefreshToken{}, utils.WrapError(err, "Failed to validate refresh token", utils.ErrCodeInternal)
 }

 if storedToken.Revoked || storedToken.ExpiresAt.Before(time.Now()) {
	return RefreshToken{}, utils.NewError("Refresh token expired or revoked", utils.ErrCodeUnauthorized)
 }

 return storedToken, nil
}

func (js *JWTService) RevokeRefreshToken(refreshToken string) error {
	key := fmt.Sprintf("refresh_token:%s", refreshToken)
	var token RefreshToken
	err := js.cache.Get(key, &token)
	if err != nil {
		return utils.WrapError(err, "Failed to revoke refresh token", utils.ErrCodeInternal)
	}
	token.Revoked = true
	return js.cache.Set(key, token, time.Until(token.ExpiresAt))
}