package chanbox

import (
	"github.com/gogap/chanbox/box/mock"

	"testing"
)

func TestSendAndReceive(t *testing.T) {

	mockbox := new(mock.MockBox)

	cbox, err := New(
		Inboxes(mockbox),
		Outboxes(mockbox),
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
