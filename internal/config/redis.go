package config

import (
	"context"
	"log"
	"time"
	"user-management-api/internal/utils"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string
	Password string
	User     string
	DB       int
}

func NewRedisClient() *redis.Client {
	cfg := &RedisConfig{
		Addr:     utils.GetEnv("REDIS_ADDR", "localhost:6379"),
		Password: utils.GetEnv("REDIS_PASSWORD", ""),
		User:     utils.GetEnv("REDIS_USER", ""),
		DB:     utils.GetIntEnv("REDIS_DB", 0),
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		Username:     cfg.User,
		DB:       cfg.DB,
		PoolSize: 10,
		MinIdleConns: 5,
		MaxIdleConns: 10,
		MaxActiveConns: 20,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		
	})

	ctx,cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	_,err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}

	log.Println("Connected to Redis")

	return client
}