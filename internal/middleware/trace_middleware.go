package middleware

import (
	"context"
	"user-management-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)



func TraceMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		contextValue := context.WithValue(ctx.Request.Context(), logger.TraceIDContextKey, traceID)
		ctx.Request = ctx.Request.WithContext(contextValue)

		ctx.Writer.Header().Set("X-Trace-ID", traceID)

		ctx.Set(string(logger.TraceIDContextKey), traceID)
		ctx.Next()
	}
}