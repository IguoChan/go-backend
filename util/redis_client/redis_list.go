package redis_client

import "github.com/gomodule/redigo/redis"

func (c *Client) BRPOP(key string, timeoutSeconds int64) ([]string, error) {
	return redis.Strings(c.ExecCommand("BRPOP", key, timeoutSeconds))
}
