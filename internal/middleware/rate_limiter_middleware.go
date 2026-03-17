package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"
	"user-management-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*Client)
)

func getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}

	return ip
}

func getRateLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clients[ip]
	if !exists {
		requestSecStr := utils.GetEnv("RATE_LIMIT_REQUESTS_SEC", "5")
		brustStr := utils.GetEnv("RATE_LIMIT_BRUST", "10")


		requestSec, err  :=  strconv.Atoi(requestSecStr)
		if err != nil {
			panic(err)
		}
		brust, err  :=  strconv.Atoi(brustStr)
		if err != nil {
			panic(err)
		}
		limiter := rate.NewLimiter(rate.Limit(requestSec), brust)
		newClient := &Client{limiter, time.Now()}
		clients[ip] = newClient
		return limiter
	}

	client.lastSeen = time.Now()
	return client.limiter
}

func CleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

// ab -n 20 -c 1 -H "X-API-Key:ab2a7a8a-d601-4bf7-b0e2-dd00e5459392" http://localhost:8080/api/v1/categories/golang
func RateLimiterMiddleware(logger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := getClientIP(ctx)

		limiter := getRateLimiter(ip)

		if !limiter.Allow() {
			if shouldLogRateLimit(ip) {
			logger.Warn().Str("client_ip", ctx.ClientIP()).Str("user_agent", ctx.Request.UserAgent()).Str("referer", ctx.Request.Referer()).Msg("Too many requests") }

			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many request",
				"message": "Bạn đã gửi quá nhiêu request. Hãy thử lại sau",
			})

			return
		}

		ctx.Next()
	}
}


var rateLimitLogCache = sync.Map{}

const rateLimitLogTTL = 10 * time.Second

func shouldLogRateLimit(ip string) bool {
 now := time.Now()
   if val, ok := rateLimitLogCache.Load(ip); ok {
	 if t, ok := val.(time.Time); ok && now.Sub(t) < rateLimitLogTTL {
		return false
	 }
	 rateLimitLogCache.Store(ip, now)
	 return true
   }
   return true
}