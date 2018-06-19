package main

import (
	"testing"

	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	chatID = ""
)

func TestGet(t *testing.T) {
	res, err := AP.Get("Helper\nTesting")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if res != "" {
		t.Fatalf("Results should be empty")
	}
}

func TestSend(t *testing.T) {
	sub := kbchat.SubscriptionMessage{}
	sub.Conversation.ID = chatID
	_, err := AP.Get("Helper\nTesting")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	res, err := AP.Get("")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if res == "" {
		t.Fatalf("Results should not be empty the second time")
	}
	if err := AP.Send(sub, res); err != nil {
		t.Fatalf("Error sending to keybase %v", err)
	}
}
