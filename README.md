# ChanBox

chanbox is a chan usage like and could with different transports

### Usage

We need import box driver first

```go
	import _ "github.com/gogap/chanbox/box/mock"
```

```go
	intbox, _ := box.NewInbox("MockBox")
	outbox, _ := box.NewOutbox("MockBox")

	cbox, err := New(
		Inboxes(intbox),
		Outboxes(outbox),
	)

	if err != nil {
		t.Error(err)
		return
	}

	cbox.Out() <- "hello"

	msg := <-cbox.In()

	if msg != "hello" {
		t.Error("send and receive message not match")
		return
	}
```

we could implement `box` with golang `chan`, `message queue` and `http`, and we could get messages from multi `inboxes`, send messages to multi `outboxes`