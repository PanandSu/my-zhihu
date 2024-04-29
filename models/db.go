package models

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
	"log"
	"my-zhihu/config"
)

var (
	db        *sql.DB
	redisPool *redis.Pool
)

func initDB() {
	var err error
	db, err = sql.Open("mysql", config.Cfg.Database.DSN)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func initRedisPool() {
	redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", config.Cfg.Redis.Addr)
		},
		MaxIdle: config.Cfg.Redis.MaxIdle,
	}
}
func init() {
	initDB()
	initRedisPool()
}

const (
	zrem      = "ZREM"
	zadd      = "ZADD"
	scard     = "SCARD"
	zrevrange = "ZREVRANGE"
	zscore    = "ZSCORE"
	sadd      = "SADD"
	srem      = "SREM"
)
