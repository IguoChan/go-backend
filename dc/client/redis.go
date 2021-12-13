package client

import (
	"go-backend/config"
	"go-backend/util/redis_client"
	"time"
)

var DefaultRedisClient *redis_client.Client
var DefaultRedisLockClient *redis_client.Client

func CloseRedisDB() {
	if DefaultRedisClient != nil {
		DefaultRedisClient.Close()
	}
}

func init() {
	redisCfg := config.GetRedisCfg()
	defaultRedisClientParam := &redis_client.ClientParam{
		Password:       redisCfg.Password,
		Address:        redisCfg.Hostname + ":" + redisCfg.Port,
		DBId:           0, // 默认选择DB0
		MaxIdle:        redisCfg.MaxConnection,
		MaxActive:      redisCfg.MaxConnection,
		ConnectTimeout: redisCfg.Timeout,
		ReadTimeout:    redisCfg.Timeout,
		WriteTimeout:   redisCfg.Timeout,
		IdleTimeout:    180 * time.Second,
	}

	defaultRedisLockClientParam := &redis_client.ClientParam{
		Password:       redisCfg.Password,
		Address:        redisCfg.Hostname + ":" + redisCfg.Port,
		DBId:           1, // 默认选择DB1
		MaxIdle:        redisCfg.MaxConnection,
		MaxActive:      redisCfg.MaxConnection,
		ConnectTimeout: redisCfg.Timeout,
		ReadTimeout:    50 * time.Second,
		WriteTimeout:   redisCfg.Timeout,
		IdleTimeout:    180 * time.Second,
	}

	DefaultRedisClient = redis_client.NewRedis(defaultRedisClientParam)
	DefaultRedisLockClient = redis_client.NewRedis(defaultRedisLockClientParam)
}
