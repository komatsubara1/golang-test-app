package lib

import "github.com/gomodule/redigo/redis"

func HGet[T any](conn redis.Conn, key string, field any) *T {
	defer func(Cache redis.Conn) {
		_ = Cache.Close()
	}(conn)
	do, err := conn.Do("HGET", key, field)
	if err != nil {
		return nil
	}
	if do == nil {
		return nil
	}

	return do.(*T)
}

func HSet[T any](conn redis.Conn, key string, field any, value *T) error {
	defer func(Cache redis.Conn) {
		_ = Cache.Close()
	}(conn)
	_, err := conn.Do("HSET", key, field, value)
	return err
}
