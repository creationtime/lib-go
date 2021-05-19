package redis

import "testing"

func Test(t *testing.T) {
	if pool, err := NewPool(&NewPoolReq{
		RedisAddr:        "127.0.0.1:6379",
		RedisMaxIdle:     8,
		RedisMaxActive:   0,
		RedisIdleTimeout: 300,
	}); err != nil {
		t.Error(err)
	} else {
		t.Log(pool)
	}
}
