package mock

import (
	"sync"

	"github.com/gogap/chanbox/box"
)

var (
	_ box.Inbox  = (*MockBox)(nil)
	_ box.Outbox = (*MockBox)(nil)
)

type MockBox struct {
	lastMsg interface{}
	locker  sync.Mutex
}

var (
	mockBox = MockBox{}
)

func init() {
	box.RegisterInbox("MockBox", NewMockInbox)
	box.RegisterOutbox("MockBox", NewMockOutbox)
}

func NewMockInbox(opts box.Options) (boxs box.Inbox, err error) {
	return &mockBox, nil
}

func NewMockOutbox(opts box.Options) (boxs box.Outbox, err error) {
	return &mockBox, nil
}

func (p *MockBox) String() string {
	return "MockBox"
}

func (p *MockBox) Send(msg interface{}) (err error) {
	p.locker.Lock()
	defer p.locker.Unlock()

	p.lastMsg = msg

	return
}

func (p *MockBox) Recv() (msg interface{}, err error) {
	p.locker.Lock()
	defer p.locker.Unlock()

	msg = p.lastMsg
	p.lastMsg = nil

	return
}
