package v1routes

import (
	v1handler "user-management-api/internal/handler/v1"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	handler *v1handler.AuthHandler
}

func NewAuthRoutes(handler *v1handler.AuthHandler) *AuthRoutes {
	return &AuthRoutes{
		handler: handler,
	}
}

func (ur *AuthRoutes) Register(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", ur.handler.Login)
		auth.POST("/logout", ur.handler.Logout)
		auth.POST("/refresh", ur.handler.RefreshToken)
	}
}
