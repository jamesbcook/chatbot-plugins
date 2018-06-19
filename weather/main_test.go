package main

import (
	"os"
	"testing"

	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	chatID = ""
)

func TestSend(t *testing.T) {
	sub := kbchat.SubscriptionMessage{}
	sub.Conversation.ID = chatID
	err := os.Setenv("CHATBOT_WEATHER", "")
	if err != nil {
		t.Fatalf("Error setting env var %v", err)
	}
	output, err := AP.Get("phoenix")
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
	err := os.Setenv("CHATBOT_WEATHER", "")
	if err != nil {
		t.Fatalf("Error setting env var %v", err)
	}
	output, err := AP.Get("phoenix")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
}
