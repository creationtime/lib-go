package jaeger

type Options struct {
	serverName string
	address    string
}

type Option func(o *Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		serverName: "go.micro.srv",
		address:    "http://localhost:8080",
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func ServerName(k string) Option {
	return func(o *Options) {
		o.serverName = k
	}
}

func Address(sep string) Option {
	return func(o *Options) {
		o.address = sep
	}
}
