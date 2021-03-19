package db

import "time"

type Options struct {
	maxOpen     int
	maxIdle     int
	maxLifetime time.Duration
}

type Option func(o *Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		maxOpen:     20,
		maxIdle:     5,
		maxLifetime: time.Minute,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func WithMaxOpen(v int) Option {
	return func(o *Options) {
		o.maxOpen = v
	}
}

func WithMaxIdle(v int) Option {
	return func(o *Options) {
		o.maxIdle = v
	}
}

func WithMaxLifetime(v time.Duration) Option {
	return func(o *Options) {
		o.maxLifetime = v
	}
}
