package actor

import "time"

const (
	defaultInboxSize    = 1024
	defaultMaxRestarts  = 3
	defaultRestartDelay = 500 * time.Millisecond
)

type (
	ReceiveFunc    = func(*Context)
	MiddlewareFunc = func(ReceiveFunc) ReceiveFunc
	Options        struct {
		Producer     Producer
		Name         string
		Tags         []string
		MaxRestarts  int32
		RestartDelay time.Duration
		InboxSize    int
		Middleware   []MiddlewareFunc
	}
	Option func(*Options)
)

func DefaultOpts(producer Producer) Options {
	return Options{
		Producer:     producer,
		MaxRestarts:  defaultMaxRestarts,
		RestartDelay: defaultRestartDelay,
		InboxSize:    defaultInboxSize,
		Middleware:   []MiddlewareFunc{},
	}
}

func WithMaxRestarts(n int) Option {
	return func(opts *Options) {
		opts.MaxRestarts = int32(n)
	}
}

func WithRestartDelay(delay time.Duration) Option {
	return func(opts *Options) {
		opts.RestartDelay = delay
	}
}

func WithInboxSize(size int) Option {
	return func(opts *Options) {
		opts.InboxSize = size
	}
}

func WithMiddleware(middleware ...MiddlewareFunc) Option {
	return func(opts *Options) {
		opts.Middleware = append(opts.Middleware, middleware...)
	}
}

func WithTags(tags ...string) Option {
	return func(opts *Options) {
		opts.Tags = tags
	}
}
