package context

import (
	"fmt"
	"os"
	"strconv"
)

func NewUserDbContext() *UserDbContext {
	udc := &UserDbContext{}
	return udc
}

type UserDbContext struct {
	Dc SqlDbContext

	s []func()
	e []func()
}

func (c *UserDbContext) Connect() {
	/*
		// TODO: shard
		i, _ := strconv.Atoi(c.userId.Value().String())
		shardNum, _ := strconv.Atoi(os.Getenv("USER_DB_SHARD_NUM"))
		idx := (i % shardNum) + 1
	*/
	if c.Dc.Db != nil {
		return
	}

	user := os.Getenv("USER_DB_USER")
	password := os.Getenv("USER_DB_PASS")
	host := os.Getenv(fmt.Sprintf("USER_DB_HOST_%d", 1))
	port := os.Getenv(fmt.Sprintf("USER_DB_PORT_%d", 1))
	dbname := os.Getenv("USER_DB_NAME")
	maxOpenConn, _ := strconv.Atoi(os.Getenv("USER_DB_MAX_OPEN_CONN"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("USER_DB_MAX_IDLE_CONN"))
	connLifeTime, _ := strconv.Atoi(os.Getenv("USER_DB_CONN_LIFE_TIME"))

	c.Dc.Connect(
		user,
		password,
		host,
		port,
		dbname,
		maxOpenConn,
		maxIdleConn,
		connLifeTime,
	)
}

func (c *UserDbContext) TransactionScope(cb func() error) error {
	c.Connect()
	defer c.Dc.Close()
	return c.Dc.TransactionScope(cb)
}
