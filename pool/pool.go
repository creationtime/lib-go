package pool

import (
	"time"

	"github.com/panjf2000/ants/v2"
)

// NewWorkerPool instantiates a non-blocking *WorkerPool with the capacity of DefaultAntsPoolSize.
func New(size int, expire time.Duration) *ants.Pool {
	ants.Release()
	options := ants.Options{ExpiryDuration: expire, PreAlloc: true, Nonblocking: true}
	defaultAntsPool, _ := ants.NewPool(size, ants.WithOptions(options))
	return defaultAntsPool
}
