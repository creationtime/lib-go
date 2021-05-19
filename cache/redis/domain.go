package redis

import (
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	OnlineExpireTime = 24 * time.Hour
	SaveExpire       = time.Minute * 5
	SetnxExpire      = 10 //秒级
)

var (
	cli        *redis.ClusterClient
	poolBuffer *sync.Pool
)

type NewPoolReq struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int64
}
