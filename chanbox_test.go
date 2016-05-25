package chanbox

import (
	"testing"

	"github.com/gogap/chanbox/box"

	_ "github.com/gogap/chanbox/box/mock"
)

func TestSendAndReceive(t *testing.T) {

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
}
