package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCacheService struct{
	ctx context.Context
	rdb *redis.Client
}

func NewRedisCacheService(rdb *redis.Client) RedisCacheService {
	return &redisCacheService{ctx: context.Background(), rdb: rdb}

}

func (r *redisCacheService) Get(key string, dest any)  error {
	data, err := r.rdb.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return err
	}
	if err != nil {
		return err
	}
	
	return json.Unmarshal([]byte(data), dest)
}

func (r *redisCacheService) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	 return r.rdb.Set(r.ctx, key, data, ttl).Err()
	
}

func (r *redisCacheService) Clear(pattern string) error {
 cursor := uint64(0)
 for {
	keys, nextCursor, err := r.rdb.Scan(r.ctx, cursor, pattern, 100).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		r.rdb.Del(r.ctx, keys...)
	}

	cursor = nextCursor
	if nextCursor == 0 {
		break
	}
 }
 return nil
}

func (r *redisCacheService) Exists(key string) (bool, error) {
	 exists, err := r.rdb.Exists(r.ctx, key).Result()
	 if err != nil {
		return false, err
	 }
	 return exists > 0, nil
}