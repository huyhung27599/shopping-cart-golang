package middleware

import (
	"net/http"
	"strings"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"

	"github.com/gin-gonic/gin"
)

var (
	jwtService auth.TokenService
	cacheService cache.RedisCacheService
)

func InitAuthMiddleware(tokenService auth.TokenService, cache cache.RedisCacheService) {
	jwtService = tokenService
	cacheService = cache
}

func AuthMiddleware() gin.HandlerFunc {
	return func (ctx *gin.Context)  {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}
		token := strings.Split(authHeader, " ")[1]
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		_, claims, err := jwtService.ParseToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if jti, ok := claims["jti"].(string); ok  {
			
			key := "blacklist:" + jti
			exists, err := cacheService.Exists(key)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			if exists {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is revoked"})
				return
			}
		}
		encryptedPayload, err := jwtService.DecryptToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set("user_uuid", encryptedPayload.UserUUID)
		ctx.Set("email", encryptedPayload.Email)
		ctx.Set("role", encryptedPayload.Role)
		ctx.Next()
	}
}