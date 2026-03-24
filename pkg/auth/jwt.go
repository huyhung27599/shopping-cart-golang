package auth

import (
	"encoding/json"
	"time"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct{}



type EncryptedPayload struct {
	UserUUID string `json:"user_uuid"`
	Email    string `json:"email"`
	Role     int `json:"role"`
}

func NewJWTService() *JWTService {
	return &JWTService{}
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

func (js *JWTService) GenerateRefreshToken() (string, error) {
 return "", nil
}