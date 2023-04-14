package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"go-web-app/settings"
	"go.uber.org/zap"
)

var rdb *redis.Client

func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),

		Password: cfg.Password,
		DB:       cfg.DbName,
		PoolSize: cfg.PoolSize,
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		zap.L().Error("connect to rdb failed", zap.Error(err))
		return
	}
	return
}

func Close() {
	_ = rdb.Close()
}
