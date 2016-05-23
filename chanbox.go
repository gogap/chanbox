package chanbox

import (
	"io"
	"sync"
	"time"

	"github.com/gogap/chanbox/box"
)

type Chanbox struct {
	options Options

	in  chan interface{}
	out chan interface{}
	err chan error

	isClosed bool

	closeLocker sync.Mutex
}

func (p *Chanbox) In() <-chan interface{} {
	if len(p.options.Inboxes) == 0 {
		panic("non inbox bound")
	}
	return p.in
}

func (p *Chanbox) Out() chan<- interface{} {
	if len(p.options.Outboxes) == 0 {
		panic("non outbox bound")
	}
	return p.out
}

func (p *Chanbox) Error() chan<- error {
	return p.err
}

func (p *Chanbox) Close() {
	p.closeLocker.Lock()
	defer p.closeLocker.Unlock()

	if p.isClosed {
		panic("chanbox already closed")
	}

	p.isClosed = true

	close(p.out)
	close(p.err)
}

func (p *Chanbox) send(msg interface{}) {
	for i := 0; i < len(p.options.Outboxes); i++ {
		if err := p.options.Outboxes[i].Send(msg); err != nil {
			p.onError(err)
		}
	}
}

func (p *Chanbox) onMsgReceived(msg interface{}) {
	if p.options.Timeout == 0 {
		p.in <- msg
	} else {
		select {
		case p.in <- msg:
			{
			}
		case <-time.After(p.options.Timeout):
			{
			}
		}
	}
}

func (p *Chanbox) onError(err error) {
	if err != nil {

		if err == io.ErrNoProgress {
			return
		}

		if p.options.Timeout == 0 {
			p.err <- err
		} else {
			select {
			case p.err <- err:
				{
				}
			default:
			}
		}
	}
}

func (p *Chanbox) run() {
	if len(p.options.Inboxes) > 0 {
		wg := sync.WaitGroup{}
		for _, inbox := range p.options.Inboxes {
			wg.Add(1)
			go func(box box.Inbox) {
				defer wg.Done()
				for {
					if p.isClosed {
						return
					}

					if msg, err := box.Recv(); err != nil {
						if err == io.EOF {
							return
						}
						p.onError(err)
					} else {
						p.onMsgReceived(msg)
					}
				}
			}(inbox)
		}
		go func(wg *sync.WaitGroup) {
			wg.Wait()
			close(p.in)
		}(&wg)
	}

	if len(p.options.Outboxes) > 0 {
		go func() {
			for {
				select {
				case msg, ok := <-p.out:
					{
						if !ok {
							return
						}
						p.send(msg)
					}
				}
			}
		}()
	}
}

func New(opts ...Option) (cbox *Chanbox, err error) {

	options := Options{
		Buffer:  1,
		Timeout: 0,
	}

	for _, opt := range opts {
		opt(&options)
	}

	cbox = &Chanbox{
		options:  options,
		in:       make(chan interface{}, options.Buffer),
		out:      make(chan interface{}, options.Buffer),
		err:      make(chan error, options.Buffer),
		isClosed: false,
	}

	cbox.run()

	return
}
