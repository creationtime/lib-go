package redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

	err = cli.Set("a", 1, time.Minute).Err()
	assert.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if v, err := cli.Get("a").Result(); err != nil {
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
