package redis_client

import "github.com/gomodule/redigo/redis"

var deleteAndRPUSHScript = redis.NewScript(2, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		redis.call("DEL", KEYS[1])
		if redis.call("LLEN", KEYS[2]) == 0 then
			redis.call("RPUSH", KEYS[2], 1)
			redis.call("EXPIRE", KEYS[2], 5)
			return 1
		else
			return 0
		end
	else
		return 0
	end
`)

func (c *Client) TryGetLock(key, value string, expireTimeSeconds int64) (bool, error) {
	ok, err := c.SET(key, value, SetWithEx(expireTimeSeconds), SetWithNx())
	if ok == "OK" && err == nil {
		return true, nil
	}
	return false, err
}

func (c *Client) WaitForGetLock(waitKey string, waitTimeSeconds int64) (bool, error) {
	_, err := c.BRPOP(waitKey, waitTimeSeconds)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *Client) ReleaseLockAndRPUSH(key, waitKey, value string) error {
	rc := c.redisPool.Get()
	defer rc.Close()

	_, err := deleteAndRPUSHScript.Do(rc, key, waitKey, value)
	return err
}
