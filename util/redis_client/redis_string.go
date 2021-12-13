package redis_client

import "github.com/gomodule/redigo/redis"

type setArgs []interface{}

type SetOption func(args setArgs) setArgs

// Set the specified expire time, in seconds
// >= v2.6.12
func SetWithEx(seconds int64) SetOption {
	return func(args setArgs) setArgs {
		args = append(args, "EX", seconds)
		return args
	}
}

// Set the specified expire time, in milliseconds
// >= v2.6.12
func SetWithPx(milliseconds int) SetOption {
	return func(args setArgs) setArgs {
		args = append(args, "PX", milliseconds)
		return args
	}
}

// Only set the key if it dose not already exist
// >= v2.6.12
func SetWithNx() SetOption {
	return func(args setArgs) setArgs {
		args = append(args, "NX")
		return args
	}
}

// Only set the key if it already exist
// >= v2.6.12
func SetWithXx() SetOption {
	return func(args setArgs) setArgs {
		args = append(args, "XX")
		return args
	}
}

func (c *Client) SET(key string, value interface{}, options ...SetOption) (string, error) {
	args := setArgs{key, value}
	for _, f := range options {
		args = f(args)
	}
	return redis.String(c.ExecCommand("SET", args...))
}

func (c *Client) GET(key string) (string, error) {
	return redis.String(c.ExecCommand("GET", key))
}

func (c *Client) EXPIRE(key string, seconds int64) (int, error) {
	return redis.Int(c.ExecCommand("EXPIRE", key, seconds))
}
