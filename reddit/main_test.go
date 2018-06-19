package main

import (
	"testing"

	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	chatID = ""
)

func TestSend(t *testing.T) {
	sub := kbchat.SubscriptionMessage{}
	sub.Conversation.ID = chatID
	output, err := AP.Get("netsec")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
	if err := AP.Send(sub, output); err != nil {
		t.Fatalf("Error sending command to keybase %v", err)
	}
}

func TestGet(t *testing.T) {
	output, err := AP.Get("netsec")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
	t.Log(output)
	_, err = AP.Get("asdfasdfadsf;aldska;lj")
	if err == nil {
		t.Fatalf("Error shouldn't be nil")
	}
}
