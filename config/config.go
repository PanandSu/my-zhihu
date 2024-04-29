package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server   ServerConfig
	Redis    RedisConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port          string `json:"port"`
	SessionKey    string `json:"session_key"`
	SessionSecret string `json:"session_secret"`
	Salt          string `json:"salt"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Addr     string
	MaxIdle  int `json:"max_idle"`
}

type DataBaseConfig struct {
	DriverName string `json:"driver_name"`
	DSN        string `json:"dsn"`
}

type DatabaseConfig struct {
	DriverName   string `json:"driver_name"`
	User         string `json:"user"`
	Password     string `json:"password"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	DatabaseName string `json:"database_name"`
	DSN          string
}

var (
	Cfg      *Config
	server   *ServerConfig
	redis    *RedisConfig
	database *DatabaseConfig
)

func initJson() {
	data, err := os.ReadFile("config/config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &Cfg)
	if err != nil {
		panic(err)
	}
}

func initServer() {
	server = &Cfg.Server
}

func initRedis() {
	redis = &Cfg.Redis
	redis.Addr = redis.Host + ":" + redis.Port
}

func initDatabase() {
	database = &Cfg.Database
	switch database.DriverName {
	case "mysql":
		database.DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			database.User,
			database.Password,
			database.Host,
			database.Port,
			database.DatabaseName,
		)
	}
}

func init() {
	initJson()
	initServer()
	initRedis()
	initDatabase()
}
