package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTExpirationInSeconds int64
	JWTSecret              string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBName                 string
	DBPort                 string
	RedisPassword          string
	RedisAddr              string
	RedisPort              string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "secret"),
		Port:                   getEnv("PORT", ":3000"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "root"),
		RedisPassword:          getEnv("REDIS_PASSWORD", ""),
		RedisPort:              getEnv("REDIS_PORT", ":6379"),
		RedisAddr:              getEnv("REDIS_ADDR", "redis"),
		DBName:                 getEnv("DB_NAME", "rest_api"),
		DBPort:                 getEnv("DB_PORT", "5432"),
	}
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
