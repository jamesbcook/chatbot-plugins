package main

import (
	"testing"

	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	chatID = ""
)

func TestSend(t *testing.T) {
	AP.Debug(false, nil)
	sub := kbchat.SubscriptionMessage{}
	sub.Conversation.ID = chatID
	output, err := AP.Get("8.8.8.8")
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
	AP.Debug(false, nil)
	output, err := AP.Get("8.8.8.8")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
}
