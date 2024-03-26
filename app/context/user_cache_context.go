package context

import (
	"os"

	"github.com/gomodule/redigo/redis"
)

type UserCacheContext struct {
	Conn redis.Conn
}

func NewUserCacheContext() *UserCacheContext {
	return &UserCacheContext{}
}

func (c *UserCacheContext) Connect() redis.Conn {
	if c.Conn != nil {
		return c.Conn
	}

	addr := os.Getenv("USER_REDIS_ADDR")
	conn, err := redis.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	c.Conn = conn

	return conn
}
