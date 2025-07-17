package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DatabaseURL string
	RedisURL string
	JWTSecret string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port: getEnv("PORT"),
		DatabaseURL: getEnv("DATABASE_URL"),
		RedisURL: getEnv("REDIS_URL"),
		JWTSecret: getEnv("JWT_SECRET"),
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return ""
}