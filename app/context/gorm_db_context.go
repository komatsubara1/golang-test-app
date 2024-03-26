package context

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type GormDbContext struct {
	Db *gorm.DB

	s []func()
	e []func()
}

func (c *GormDbContext) Connect(
	user string,
	password string,
	host string,
	port string,
	dbname string,
	maxOpenConn int,
	maxIdleConn int,
	connLifeTime int,
) {
	dsn := c.generateDsn(user, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	c.Db = db

	sqlDb, err := c.Db.DB()
	if err != nil {
		return
	}
	sqlDb.SetMaxOpenConns(maxOpenConn)
	sqlDb.SetMaxIdleConns(maxIdleConn)
	sqlDb.SetConnMaxLifetime(time.Duration(connLifeTime))
}

func (c *GormDbContext) generateDsn(user string, password string, host string, port string, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", user, password, host, port, dbname)
}

func (c *GormDbContext) Close() {
	sqlDb, err := c.Db.DB()
	if err != nil {
		return
	}
	_ = sqlDb.Close()
}

func (c *GormDbContext) TransactionScope(cb func() error) error {
	tx := c.Db.Begin()

	c.addSuccess(func() {
		_ = tx.Commit()
	})

	c.addError(func() {
		_ = tx.Rollback()
	})

	err := cb()
	if err != nil {
		c.onError()
		return err
	}

	c.onSuccess()
	return nil
}

func (c *GormDbContext) onSuccess() {
	for _, s := range c.s {
		s()
	}
}

func (c *GormDbContext) addSuccess(f func()) {
	c.s = append(c.s, f)
}

func (c *GormDbContext) onError() {
	for _, e := range c.e {
		e()
	}
}

func (c *GormDbContext) addError(f func()) {
	c.e = append(c.e, f)
}
