package config

import (
	"os"

	"gharabakloo/search/internal/entity"
)

func New() *entity.Config {
	return &entity.Config{
		HTTP: entity.HTTPConfig{
			IP:   os.Getenv("HTTP_IP"),
			Port: os.Getenv("HTTP_PORT"),
		},
		DB: entity.DBConfig{
			MySQL: entity.MySQLConfig{
				Driver: os.Getenv("DATABASE_DRIVER"),
				DBName: os.Getenv("DATABASE_NAME"),
				Host:   os.Getenv("DATABASE_HOST"),
				Port:   os.Getenv("DATABASE_PORT"),
				User:   os.Getenv("DATABASE_USER"),
				Pass:   os.Getenv("DATABASE_PASS"),
			},
			Redis: entity.RedisConfig{
				Host: os.Getenv("REDIS_HOST"),
				Pass: os.Getenv("REDIS_PASS"),
				DB:   os.Getenv("REDIS_DB"),
			},
		},
	}
}
