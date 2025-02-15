package utility

import (
	"goweb/common/webframe"
	"time"

	"github.com/Streamlet/gohttp"
	"github.com/redis/go-redis/v9"
)

type session struct {
	sessionId string
	provider  gohttp.CacheProvider
}

func NewFakeSession(client *redis.Client, prefix, sessionId string) gohttp.Session {
	return &session{sessionId, webframe.NewSessionProvider(client, prefix)}
}

func (s *session) Exists(key string) bool {
	return s.provider.HExists(s.sessionId, key)
}

func (s *session) Get(key string) interface{} {
	return s.provider.HGet(s.sessionId, key)
}

func (s *session) Set(key string, value interface{}, expiration time.Duration) {
	s.provider.HSet(s.sessionId, key, value, expiration)
}

func (s *session) Delete(key string) bool {
	if !s.provider.HExists(s.sessionId, key) {
		return true
	}
	if s.provider.HDelete(s.sessionId, key) {
		return true
	}
	return false
}
