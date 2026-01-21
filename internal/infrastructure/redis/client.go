package redis

import (
	"OpsGo/internal/infrastructure/config"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func InitRedis() error {
	cfg := config.AppConfig.Redis
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %v", err)
	}

	return nil
}

func CloseRedis() {
	if Client != nil {
		Client.Close()
	}
}
