package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewClient(addr string, db int, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		DB:          db,
		Password:    password,
		DialTimeout: time.Second * 30,
	})
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	if res := client.Ping(ctx); res.Err() != nil {
		return nil, res.Err()
	}
	return client, nil
}

func NewCluster(addrs []string) (*redis.ClusterClient, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:       addrs,
		DialTimeout: time.Second * 30,
	})
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	if res := client.Ping(ctx); res.Err() != nil {
		return nil, res.Err()
	}
	return client, nil
}
