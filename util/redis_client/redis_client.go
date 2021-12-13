package redis_client

import (
	"github.com/gomodule/redigo/redis"
	"go-backend/config"
	"go-backend/util/log"
	"time"
)

type Client struct {
	redisPool *redis.Pool
}

type ClientParam struct {
	Password       string
	Address        string
	DBId           int
	MaxIdle        int
	MaxActive      int
	ConnectTimeout time.Duration
	IdleTimeout    time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

/*
	NewRedis: create a redis connect pool
*/
func NewRedis(param *ClientParam) *Client {
	return &Client{
		redisPool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", param.Address,
					redis.DialConnectTimeout(param.ConnectTimeout),
					redis.DialReadTimeout(param.ReadTimeout),
					redis.DialWriteTimeout(param.WriteTimeout))
				if err != nil {
					return nil, err
				}

				if config.GetRedisCfg().Log {
					c = redis.NewLoggingConn(c, log.GetLogLogger("[REDIS] "), "REDIS")
				}

				if _, err = c.Do("AUTH", param.Password); err != nil {
					return nil, err
				}

				if _, err = c.Do("SELECT", param.DBId); err != nil {
					return nil, err
				}

				return c, nil
			},
			MaxIdle:     param.MaxIdle,     // 最大的空闲连接数，表示即使没有redis连接时毅然可以保持的空闲连接数，随时待命
			MaxActive:   param.MaxActive,   // 最大的激活连接数
			IdleTimeout: param.IdleTimeout, // 空闲连接最大的存活时间
		},
	}
}

func (c *Client) ExecCommand(command string, args ...interface{}) (interface{}, error) {
	rc := c.redisPool.Get()
	defer rc.Close()
	return rc.Do(command, args...)
}

func (c *Client) Close() {
	if c.redisPool != nil {
		c.redisPool.Close()
	}
}
