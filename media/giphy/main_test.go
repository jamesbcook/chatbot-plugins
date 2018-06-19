package main

import (
	"log"
	"os"
	"testing"

	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	chatID = ""
)

func init() {
	err := os.Setenv("CHATBOT_GIPHY", "")
	if err != nil {
		log.Fatalf("Error in getting env var %v", err)
	}
}

func TestGiphy(t *testing.T) {
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
