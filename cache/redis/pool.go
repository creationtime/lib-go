package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

//初始化单机redis池
func NewPool(req *NewPoolReq) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:     req.RedisMaxIdle,                                  //空闲连接数
		MaxActive:   req.RedisMaxActive,                                //最大活跃连接数，=0没有限制
		IdleTimeout: time.Duration(req.RedisIdleTimeout) * time.Second, //过期时间
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", req.RedisAddr)
		},
	}
	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("ping"); err != nil {
		return nil, err
	}
	return pool, nil
}
