package v1service

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"
	"user-management-api/pkg/mail"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

type authService struct {
	userRepository repository.UserRepository
	tokenService auth.TokenService
	cache cache.RedisCacheService
	mailService mail.EmailProviderService
}

type LoginAttempt struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*LoginAttempt)
	LoginAttemptTTL = 3 * time.Minute
	MaxLoginAttempt = 5
)

func NewAuthService(userRepository repository.UserRepository, tokenService auth.TokenService, cache cache.RedisCacheService, mailService mail.EmailProviderService) *authService {
	return &authService{
		userRepository: userRepository,
		tokenService: tokenService,
		cache: cache,
		mailService: mailService,
	}
}


func (as *authService) getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}

	return ip
}

func (as *authService) getLoginAttempt(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clients[ip]
	if !exists {
		


		requestSec := float64(MaxLoginAttempt) / LoginAttemptTTL.Seconds()
	
		
		limiter := rate.NewLimiter(rate.Limit(requestSec), MaxLoginAttempt)
		newClient := &LoginAttempt{limiter, time.Now()}
		clients[ip] = newClient
		return limiter
	}

	client.lastSeen = time.Now()
	return client.limiter
}

func (as *authService) checkLoginAttempt(ip string) error {
	
	limiter := as.getLoginAttempt(ip)
	if(!limiter.Allow()) {
		return utils.NewError("Too many login attempts", utils.ErrCodeTooManyRequests)
	}
	return nil
}

func (as *authService) CleanupClients(ip string) {
	
		
		mu.Lock()
		defer	mu.Unlock()
	
		delete(clients, ip)
		
	
	
}

func (as *authService) Login(ctx *gin.Context, email string, password string)  (string, string, int, error) {
	context := ctx.Request.Context()

	ip := as.getClientIP(ctx)
	err := as.checkLoginAttempt(ip)
	if err != nil {
		return "", "", 0, err
	}
	email = utils.NormalizeString(email)
	user, err := as.userRepository.FindByEmail(context, email)
	if err != nil {
		as.getLoginAttempt(ip)
		return "", "", 0, utils.NewError("User not found", utils.ErrCodeUnauthorized)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password))
	if err != nil {
		as.getLoginAttempt(ip)
		return "", "", 0, utils.NewError(" Invalid password or email", utils.ErrCodeUnauthorized)
	}

 accessToken, err := as.tokenService.GenerateAccessToken(user)
 
 if err != nil {
	return "", "", 0, utils.WrapError(err, "Failed to generate access token", utils.ErrCodeInternal)
 }
 refreshToken, err := as.tokenService.GenerateRefreshToken(user)
 if err != nil {
	return "", "", 0, utils.WrapError(err, "Failed to generate refresh token", utils.ErrCodeInternal)
 }
 if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
	return "", "", 0, utils.WrapError(err, "Failed to store refresh token", utils.ErrCodeInternal)
 }
	as.CleanupClients(ip)
	return accessToken, refreshToken.Token, int(auth.AccessTokenExpiration.Seconds()) ,nil
}

func (as *authService) Logout(ctx *gin.Context , refreshToken string) error {
 authHeader := ctx.GetHeader("Authorization")
 if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
	return utils.NewError("Missing Authorization header", utils.ErrCodeUnauthorized)
 }
 token := strings.Split(authHeader, " ")[1]
 _, claims,err  := as.tokenService.ParseToken(token)
if err != nil {
	return utils.NewError("Failed to parse token", utils.ErrCodeUnauthorized)
 }

 if jti, ok := claims["jti"].(string); ok  {
	expUnix, _ := claims["exp"].(float64)
	expTime := time.Unix(int64(expUnix), 0)
key := "blacklist:" + jti
ttl := time.Until(expTime)
as.cache.Set(key, "revoked", ttl)
 }

 _, err = as.tokenService.ValidateRefreshToken(refreshToken)
 if err != nil {
	return utils.WrapError(err, "Failed to validate refresh token", utils.ErrCodeUnauthorized)
 }
  err = as.tokenService.RevokeRefreshToken(refreshToken)
  if err != nil {
	return utils.WrapError(err, "Failed to revoke refresh token", utils.ErrCodeInternal)
  }
 return nil
}

func (as *authService) RefreshToken(ctx *gin.Context, refreshToken string) (string, string, int, error) {
	context := ctx.Request.Context()


	storedToken, err := as.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Failed to validate refresh token", utils.ErrCodeUnauthorized)
	}
	

	userUUID, _ := uuid.Parse(storedToken.UserUUID)
 user,err :=	as.userRepository.GetByUUID(context, userUUID)
 if err != nil {
	return "", "", 0, utils.WrapError(err, "Failed to get user", utils.ErrCodeInternal)
 }

 accessToken, err := as.tokenService.GenerateAccessToken(user)
 if err != nil {
	return "", "", 0, utils.WrapError(err, "Failed to generate access token", utils.ErrCodeInternal)
 }

 newRefreshToken, err := as.tokenService.GenerateRefreshToken(user)
 if err != nil {
	return "", "", 0, utils.WrapError(err, "Failed to generate refresh token", utils.ErrCodeInternal)
 }

  if err := as.tokenService.RevokeRefreshToken(refreshToken); err != nil {
	return "", "", 0, utils.WrapError(err, "Failed to revoke refresh token", utils.ErrCodeInternal)
 }

 if err := as.tokenService.StoreRefreshToken(newRefreshToken); err != nil {
	return "", "", 0, utils.WrapError(err, "Failed to store refresh token", utils.ErrCodeInternal)
 }

 return accessToken, newRefreshToken.Token, int(auth.AccessTokenExpiration.Seconds()), nil
}


func (as *authService) RequestForgotPassword(ctx *gin.Context, email string) error {
	context := ctx.Request.Context()

	rateLimitKey := fmt.Sprintf("rate_limit:forgot_password:%s", email)
	if exists, err := as.cache.Exists(rateLimitKey); err == nil && exists {
		return utils.NewError("Too many requests", utils.ErrCodeTooManyRequests)
	}
	
	user, err := as.userRepository.FindByEmail(context, email)
	if err != nil {
		return utils.WrapError(err, "Failed to find user", utils.ErrCodeInternal)
	}

	token, err := utils.GenerateRandomString(32)
	if err != nil {
		return utils.WrapError(err, "Failed to generate token", utils.ErrCodeInternal)
	}
	err = as.cache.Set("reset"+ token, user.UserUuid.String(), 60*time.Minute)
	if err != nil {
		return utils.WrapError(err, "Failed to set reset token", utils.ErrCodeInternal)
	}

	 if err := as.cache.Set(rateLimitKey, "1", 10*time.Minute); err != nil {
		return utils.WrapError(err, "Failed to set rate limit", utils.ErrCodeInternal)
	 }

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", utils.GetEnv("FRONTEND_URL", "http://localhost:3000"), token)



	return nil
}

func (as *authService) ResetPassword(ctx *gin.Context, token, newPassword string) error {
	context := ctx.Request.Context()
	var userUUIDStr string
	err := as.cache.Get("reset"+ token, &userUUIDStr)
	if err == redis.Nil  || userUUIDStr == "" {
		return utils.NewError("Invalid or expired token", utils.ErrCodeUnauthorized)
	}
	if err != nil {
		return utils.WrapError(err, "Failed to get reset token", utils.ErrCodeInternal)
	}
	userUUID, err := uuid.Parse(userUUIDStr)
	if err != nil {
		return utils.WrapError(err, "Failed to parse user UUID", utils.ErrCodeInternal)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return utils.WrapError(err, "Failed to hash password", utils.ErrCodeInternal)
	}
	hashedPasswordStr := string(hashedPassword)
	_, err = as.userRepository.UpdatePassword(context, userUUID, hashedPasswordStr)
	if err != nil {
		return utils.WrapError(err, "Failed to update password", utils.ErrCodeInternal)
	}

	if err := as.cache.Clear("reset" + token); err != nil {
		return utils.WrapError(err, "Failed to clear reset token", utils.ErrCodeInternal)
	}
	return nil
}