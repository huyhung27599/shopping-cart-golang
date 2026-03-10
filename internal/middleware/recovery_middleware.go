package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RecoveryMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error().Str("path", ctx.Request.URL.Path).Str("method", ctx.Request.Method).Err(err.(error)).Msg("Recovered from panic")
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "INTERNAL_SERVER_ERROR",
					"message": "An unexpected error occurred. Please try again later.",
				})
			
			}
		}()
		ctx.Next()
	}
}