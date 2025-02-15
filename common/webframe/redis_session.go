package webframe

import (
	"context"
	"time"

	"github.com/Streamlet/gohttp"
	"github.com/redis/go-redis/v9"
)

func NewSessionProvider(client *redis.Client, sessionKeyPrefix string) gohttp.CacheProvider {
	return &redisCache{client, sessionKeyPrefix}
}

type redisCache struct {
	client           *redis.Client
	sessionKeyPrefix string
}

func (rc *redisCache) Exists(key string) bool {
	if r, err := rc.client.Exists(context.Background(), rc.sessionKeyPrefix+key).Result(); err == nil && r > 0 {
		return true
	} else {
		return false
	}
}

func (rc *redisCache) HExists(key, field string) bool {
	if r, err := rc.client.HExists(context.Background(), rc.sessionKeyPrefix+key, field).Result(); err == nil && r {
		return true
	} else {
		return false
	}
}

func (rc *redisCache) HGet(key, field string) interface{} {
	s, err := rc.client.HGet(context.Background(), rc.sessionKeyPrefix+key, field).Result()
	if err != nil {
		return nil
	}
	return s
}

func (rc *redisCache) HSet(key, field string, value interface{}, expiration time.Duration) {
	if r, err := rc.client.HSet(context.Background(), rc.sessionKeyPrefix+key, field, value).Result(); err != nil || r <= 0 {
		return
	}

	if expiration > 0 {
		_, _ = rc.client.HExpire(context.Background(), rc.sessionKeyPrefix+key, expiration, field).Result()
	}
}

func (rc *redisCache) HDelete(key, field string) bool {
	if r, err := rc.client.HExists(context.Background(), rc.sessionKeyPrefix+key, field).Result(); err == nil && !r {
		return true
	}
	if r, err := rc.client.HDel(context.Background(), rc.sessionKeyPrefix+key, field).Result(); err == nil && r > 0 {
		return true
	}
	return false
}
