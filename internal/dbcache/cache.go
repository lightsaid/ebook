package dbcache

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lightsaid/ebook/internal/config"
	"github.com/redis/go-redis/v9"
)

func Open(conf config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})

	err := rdb.Ping(context.TODO()).Err()
	if err != nil {
		return rdb, err
	}

	return rdb, nil
}

func Close() {
	err := rdb.Close()
	if err != nil {
		slog.Warn("close redis error: " + err.Error())
	}
}
