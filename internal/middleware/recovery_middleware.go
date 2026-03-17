package middleware

import (
	"bytes"
	"net/http"
	"regexp"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RecoveryMiddleware(logger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				stack := debug.Stack()
				stack_at := ExtractFirstAppStackLine(stack)
				logger.Error().Str("path", ctx.Request.URL.Path).Str("method", ctx.Request.Method).Err(err.(error)).Str("stack_at", stack_at).Msg("Recovered from panic")
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "INTERNAL_SERVER_ERROR",
					"message": "An unexpected error occurred. Please try again later.",
				})
			
			}
		}()
		ctx.Next()
	}
}

var stackLineRegex = regexp.MustCompile(`^(.+?):(\d+):(\d+)$`)

func ExtractFirstAppStackLine(stack []byte) string {
	lines := bytes.Split(stack, []byte("\n"))
	for _, line := range lines {
		if bytes.Contains(line, []byte(".go")) && !bytes.Contains(line, []byte("/runtime/")) && !bytes.Contains(line, []byte("/internal/")) && !bytes.Contains(line, []byte("/debug/")) {
			cleanLine := bytes.TrimSpace(line)
			matches := stackLineRegex.FindSubmatch(cleanLine)
			if len(matches) > 0 {
				return string(matches[1])
			} else {
				return string(cleanLine)
			}
		}
	}
	return ""
}