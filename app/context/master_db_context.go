package context

import (
	"fmt"
	"os"
	"strconv"
)

type MasterDbContext struct {
	Dc GormDbContext
}

func NewMasterDbContext() *MasterDbContext {
	mdc := &MasterDbContext{}
	return mdc
}

func (c *MasterDbContext) Connect() {
	if c.Dc.Db != nil {
		return
	}

	user := os.Getenv("MASTER_DB_USER")
	password := os.Getenv("MASTER_DB_PASS")
	host := os.Getenv("MASTER_DB_HOST")
	port := os.Getenv("MASTER_DB_PORT")
	dbname := os.Getenv("MASTER_DB_NAME")

	maxOpenConn, _ := strconv.Atoi(os.Getenv("USER_DB_MAX_OPEN_CONN"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("USER_DB_MAX_IDLE_CONN"))
	connLifeTime, _ := strconv.Atoi(os.Getenv("USER_DB_CONN_LIFE_TIME"))

	c.Dc.Connect(user, password, host, port, dbname, maxOpenConn, maxIdleConn, connLifeTime)
}

func (c *MasterDbContext) generateDsn(user string, password string, host string, port string, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", user, password, host, port, dbname)
}

func (c *MasterDbContext) TransactionScope(cb func() error) error {
	c.Connect()
	defer c.Dc.Close()
	return c.Dc.TransactionScope(cb)
}
