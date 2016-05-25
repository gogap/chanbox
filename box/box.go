package box

import (
	"fmt"
	"sync"
)

var (
	inboxesFn  = make(map[string]NewInboxFunc)
	outboxesFn = make(map[string]NewOutboxFunc)

	inboxLocker  sync.RWMutex
	outboxLocker sync.RWMutex
)

type (
	NewInboxFunc  func(opts Options) (box Inbox, err error)
	NewOutboxFunc func(opts Options) (box Outbox, err error)
)

type Inbox interface {
	String() string
	Recv() (msg interface{}, err error)
}

type Outbox interface {
	String() string
	Send(msg interface{}) error
}

func RegisterInbox(name string, fn NewInboxFunc) {
	inboxLocker.Lock()
	defer inboxLocker.Unlock()

	if fn == nil {
		panic("NewInboxFunc could not be nil")
	}
	if _, dup := inboxesFn[name]; dup {
		panic("register called twice for inbox " + name)
	}

	inboxesFn[name] = fn
}

func RegisterOutbox(name string, fn NewOutboxFunc) {
	outboxLocker.Lock()
	defer outboxLocker.Unlock()

	if fn == nil {
		panic("NewOutboxFunc could not be nil")
	}
	if _, dup := outboxesFn[name]; dup {
		panic("register called twice for outbox " + name)
	}

	outboxesFn[name] = fn
}

func NewInbox(name string, opts ...Option) (box Inbox, err error) {

	options := Options{}

	if opts != nil {
		for _, opt := range opts {
			opt(&options)
		}
	}

	if fn, exist := inboxesFn[name]; exist {
		box, err = fn(options)
		return
	}

	err = fmt.Errorf("inbox new func of %s not exist", name)

	return
}

func NewOutbox(name string, opts ...Option) (box Outbox, err error) {

	options := Options{}

	if opts != nil {
		for _, opt := range opts {
			opt(&options)
		}
	}

	if fn, exist := outboxesFn[name]; exist {
		box, err = fn(options)
		return
	}

	err = fmt.Errorf("outbox new func of %s not exist", name)

	return
}
