package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	addr := "10.4.61.59:30379"
	cli, err := NewClient(addr, 1, "")
	assert.NoError(t, err)
	assert.NotNil(t, cli)
}

func BenchmarkNewClient(b *testing.B) {
	addr := "10.4.61.59:30379"
	cli, err := NewClient(addr, 0, "")
	assert.NoError(b, err)
	assert.NotNil(b, cli)
	ctx := context.TODO()
	err = cli.Set(ctx, "a", 1, time.Minute).Err()
	assert.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if v, err := cli.Get(ctx, "a").Result(); err != nil {
			assert.NoError(b, err)
			assert.EqualValues(b, 1, v)
		}
	}
}

func TestNewCluster(t *testing.T) {
	addrs := []string{"10.4.61.59:30379"}
	cli, err := NewCluster(addrs)
	assert.NoError(t, err)
	assert.NotNil(t, cli)
}
