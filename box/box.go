package box

type Inbox interface {
	String() string
	Recv() (msg interface{}, err error)
}

type Outbox interface {
	String() string
	Send(msg interface{}) error
}
