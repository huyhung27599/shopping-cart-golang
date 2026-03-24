package v1service

import (
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"
	"user-management-api/pkg/auth"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepository repository.UserRepository
	tokenService auth.TokenService
}

func NewAuthService(userRepository repository.UserRepository, tokenService auth.TokenService) *authService {
	return &authService{
		userRepository: userRepository,
		tokenService: tokenService,
	}
}

func (as *authService) Login(ctx *gin.Context, email string, password string)  (string, int, error) {
	context := ctx.Request.Context()
	email = utils.NormalizeString(email)
	user, err := as.userRepository.FindByEmail(context, email)
	if err != nil {
		return "", 0, utils.NewError("User not found", utils.ErrCodeUnauthorized)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password))
	if err != nil {
		return "", 0, utils.NewError(" Invalid password or email", utils.ErrCodeUnauthorized)
	}

 accessToken, err := as.tokenService.GenerateAccessToken(user)
 if err != nil {
	return "", 0, utils.WrapError(err, "Failed to generate access token", utils.ErrCodeInternal)
 }

	
	return accessToken, int(auth.AccessTokenExpiration.Seconds()) ,nil
}

func (as *authService) Logout(ctx *gin.Context) error {

}