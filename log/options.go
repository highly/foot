package log

import "go.uber.org/zap"

var defaultOptions = &options{
	maxAge:     7,
	maxBackups: 10,
	maxSize:    100,
	level:      DebugLevel,
	path:       "",
	fields:     make([]zap.Field, 0),
}

type options struct {
	maxAge     int
	maxBackups int
	maxSize    int
	level      Level
	path       string
	fields     []zap.Field
}

func evaluateOptions(opts []Option) *options {
	opt := &options{}
	*opt = *defaultOptions
	for _, o := range opts {
		o(opt)
	}
	return opt
}

type Option func(*options)

func WithMaxAge(val int) Option {
	return func(args *options) {
		args.maxAge = val
	}
}

func WithMaxBackups(val int) Option {
	return func(args *options) {
		args.maxBackups = val
	}
}

func WithMaxSize(val int) Option {
	return func(args *options) {
		args.maxSize = val
	}
}

func WithLevel(val Level) Option {
	return func(args *options) {
		args.level = val
	}
}

func WithPath(val string) Option {
	return func(args *options) {
		args.path = val
	}
}

func WithField(key string, val interface{}) Option {
	return func(args *options) {
		args.fields = append(args.fields, zap.Any(key, val))
	}
}
