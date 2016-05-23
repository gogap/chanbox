package chanbox

import (
	"time"

	"github.com/gogap/chanbox/box"
	"golang.org/x/net/context"
)

type Option func(*Options)

type Options struct {
	Inboxes  []box.Inbox
	Outboxes []box.Outbox
	Buffer   int
	Timeout  time.Duration

	Context context.Context
}

func Buffer(buf int) Option {
	return func(opts *Options) {
		opts.Buffer = buf
	}
}

func Inboxes(inboxes ...box.Inbox) Option {
	return func(opts *Options) {
		opts.Inboxes = inboxes
	}
}

func Outboxes(outboxes ...box.Outbox) Option {
	return func(opts *Options) {
		opts.Outboxes = outboxes
	}
}

func Timeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.Timeout = timeout
	}
}
