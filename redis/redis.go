package redis

import (
	"context"
	"log/slog"
	"main/configs"

	"github.com/go-redis/redis/v8"
)

func ConnectRedis(sLog *slog.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     configs.Envs.RedisAddr + configs.Envs.RedisPort,
		Password: configs.Envs.RedisPassword,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		sLog.Error("redis connect fail", "err", err)
	}

	return client
}
