package gredis

import (
	"context"
	"fmt"
	"password-self-service/pkg/config"
	"password-self-service/pkg/logging"

	"github.com/redis/go-redis/v9"
)

// Client 定义全局Client
var (
	Client *redis.Client
	Ctx    = context.Background()
)

// Init Initialize the Redis instance
func Init() {
	address := fmt.Sprintf("%v:%v", config.Setting.Redis.Host, config.Setting.Redis.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: config.Setting.Redis.Password,
		DB:       config.Setting.Redis.DB,
	})

	pong, err := client.Ping(Ctx).Result()
	if err != nil {
		logging.Logger().Sugar().Error("Redis connect failed, err: %v", err)
	} else {
		logging.Logger().Sugar().Infof("Redis connected %s DB: %d, response %v.", address, config.Setting.Redis.DB, pong)
		Client = client
	}

}
