package routes

import (
	"net/http"
	"user-management-api/internal/middleware"
	v1routes "user-management-api/internal/routes/v1"
	"user-management-api/internal/utils"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, authService auth.TokenService, cacheService cache.RedisCacheService, routes ...Route) {
	
	httpLogger := utils.NewLoggerWithPath("http.log", "info")

	recoveryLogger := utils.NewLoggerWithPath("recovery.log", "error")

	rateLimiterLogger := utils.NewLoggerWithPath("rate_limiter.log", "info")

	r.Use(
		middleware.RateLimiterMiddleware(rateLimiterLogger),
		middleware.CorsMiddleware(),
		middleware.TraceMiddleware(),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
		
		
	)

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	v1api := r.Group("/api/v1")

	middleware.InitAuthMiddleware(authService, cacheService)

	protected := v1api.Group("")
	protected.Use(middleware.AuthMiddleware())

	for _, route := range routes {
		switch route.(type) {
		case *v1routes.AuthRoutes:
			route.Register(v1api)
		default:
			route.Register(protected)
		}
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
}



