package context

import (
	"os"

	"github.com/gomodule/redigo/redis"
)

type MasterCacheContext struct {
	Cache redis.Conn
}

func NewMasterCacheContext() *MasterCacheContext {
	return &MasterCacheContext{}
}

func (c *MasterCacheContext) Connect() redis.Conn {
	if c.Cache != nil {
		return c.Cache
	}

	addr := os.Getenv("MASTER_REDIS_ADDR")
	conn, err := redis.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	c.Cache = conn

	return conn
}
