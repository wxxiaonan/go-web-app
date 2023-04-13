package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rdb *redis.Client

func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("database.redis.host"),
			viper.GetInt("database.redis.port"),
		),

		Password: viper.GetString("database.redis.password"),
		DB:       viper.GetInt("database.redis.dbname"),
		PoolSize: viper.GetInt("database.redis.pool_size"),
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
