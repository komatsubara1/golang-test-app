package context

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/exp/slog"
)

type DbContext struct {
	Db *sql.DB

	s []func()
	e []func()
}

func (c *DbContext) Connect(
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

	slog.Info("mysql connect.", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	c.Db = db

	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(time.Duration(connLifeTime))
}

func (c *DbContext) generateDsn(user string, password string, host string, port string, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", user, password, host, port, dbname)
}

func (c *DbContext) Close() {
	_ = c.Db.Close()
}

func (c *DbContext) TransactionScope(cb func() error) error {
	tx, err := c.Db.Begin()
	if err != nil {
		return err
	}

	c.addSuccess(func() {
		_ = tx.Commit()
	})

	c.addError(func() {
		_ = tx.Rollback()
	})

	err = cb()
	if err != nil {
		c.onError()
		return err
	}

	c.onSuccess()
	return nil
}

func (c *DbContext) onSuccess() {
	for _, s := range c.s {
		s()
	}
}

func (c *DbContext) addSuccess(f func()) {
	c.s = append(c.s, f)
}

func (c *DbContext) onError() {
	for _, e := range c.e {
		e()
	}
}

func (c *DbContext) addError(f func()) {
	c.e = append(c.e, f)
}
