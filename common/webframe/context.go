package webframe

import (
	"net/http"

	"github.com/Streamlet/gohttp"
	"github.com/Streamlet/gosql"
	"github.com/redis/go-redis/v9"
)

func NewContextFactory(cache *redis.Client, db *gosql.Connection, debug bool) gohttp.ContextFactory[HttpContext] {
	return &contextFactory{gohttp.NewSessionManager(NewSessionProvider(cache, "SESSION_")), cache, db, debug}
}

type contextFactory struct {
	sm    gohttp.SessionManager
	cache *redis.Client
	db    *gosql.Connection
	debug bool
}

func (cf *contextFactory) NewContext(w http.ResponseWriter, r *http.Request) HttpContext {
	return &httpContext{
		HttpContext: gohttp.NewHttpContext(w, r, cf.sm),
		cache:       cf.cache,
		db:          cf.db,
		debug:       cf.debug,
	}
}

type HttpContext interface {
	gohttp.HttpContext
	Success(data interface{})
	Error(errorCode int, errorMessage string)
	Cache() *redis.Client
	DB() *gosql.Connection
}

type httpContext struct {
	gohttp.HttpContext
	cache *redis.Client
	db    *gosql.Connection
	debug bool
}

const (
	ErrorSuccess int = 0
)

type response struct {
	Success bool        `json:"success,omitempty"`
	Error   int         `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func (c *httpContext) Success(data interface{}) {
	c.Json(response{true, ErrorSuccess, data, ""})
}

func (c *httpContext) Error(errorCode int, errorMessage string) {
	if c.debug {
		c.Json(response{false, errorCode, nil, errorMessage})
	} else {
		c.Json(response{false, errorCode, nil, ""})
	}
}

func (c *httpContext) Cache() *redis.Client {
	return c.cache
}

func (c *httpContext) DB() *gosql.Connection {
	return c.db.Clone()
}
