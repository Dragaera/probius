package persistence

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func InitializeRedis(host string, port string) (*redis.Pool, error) {
	redisPool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
		},
	}

	return redisPool, nil
}
