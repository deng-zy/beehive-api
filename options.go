package beehive

import "time"

type Option func(opts *Options)

type Options struct {
	PanicHandler     func(interface{})
	ExpiryDuration   time.Duration
	Logger           Logger
	DumpFile         string
	NonBlocking      bool
	MaxBlockingTasks int
}

func WithOptions(options Options) Option {
	return func(opts *Options) {
		*opts = options
	}
}

func WithExpriyDuration(expiryDuration time.Duration) Option {
	return func(opts *Options) {
		opts.ExpiryDuration = expiryDuration
	}
}

func WithLogger(logger Logger) Option {
	return func(opts *Options) {
		opts.Logger = logger
	}
}

func WithPanicHandler(panicHandler func(interface{})) Option {
	return func(opts *Options) {
		opts.PanicHandler = panicHandler
	}
}

func WithNonBlocking(nonBlocking bool) Option {
	return func(opts *Options) {
		opts.NonBlocking = nonBlocking
	}
}

func WithMaxBlockingTasks(maxBlockingTasks int) Option {
	return func(opts *Options) {
		opts.MaxBlockingTasks = maxBlockingTasks
	}
}

func WithDumpFile(file string) Option {
	return func(opts *Options) {
		opts.DumpFile = file
	}
}
