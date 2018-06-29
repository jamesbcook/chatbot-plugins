package main

import (
	"os"
	"testing"

	"github.com/jamesbcook/chatbot/kbchat"
)

var (
	chatID = os.Getenv("CHATBOT_TEST_CHATID")
)

func TestGiphy(t *testing.T) {
	AP.Debug(false, nil)
	res, err := AP.Get("Hackers")
	if err != nil {
		t.Fatalf("Error getting gif from Giphy %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results of giphy results is 0 or less")
	}
	if _, err := os.Stat(res); os.IsNotExist(err) {
		t.Fatalf("Path does not exist %v", err)
	}
	if err := os.Remove(res); err != nil {
		t.Fatalf("Error removing file %v", err)
	}
}

func TestSend(t *testing.T) {
	AP.Debug(false, nil)
	sub := kbchat.SubscriptionMessage{}
	sub.Conversation.ID = chatID
	res, err := AP.Get("Hackers")
	if err != nil {
		t.Fatalf("Error getting gif from Giphy %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results of giphy results is 0 or less")
	}
	if err := AP.Send(sub, res); err != nil {
		t.Fatalf("Error sending attachment %v", err)
	}
}
