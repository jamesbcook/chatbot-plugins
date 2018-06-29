package main

import (
	"os"
	"testing"

	"github.com/jamesbcook/chatbot/kbchat"
)

var (
	chatID = os.Getenv("CHATBOT_TEST_CHATID")
)

func TestURL(t *testing.T) {
	AP.Debug(false, nil)
	res, err := url("https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png")
	if err != nil {
		t.Fatalf("Error getting url %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results of url results is 0 or less")
	}
}

func TestGet(t *testing.T) {
	AP.Debug(false, nil)
	testURL := "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"
	res, err := AP.Get(testURL)
	if err != nil {
		t.Fatalf("Error performing get request %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results of direct results is 0 or less")
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
	testURL := "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"
	res, err := AP.Get(testURL)
	if err != nil {
		t.Fatalf("Error performing get request %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results of direct results is 0 or less")
	}
	if err := AP.Send(sub, res); err != nil {
		t.Fatalf("Error sending attachment %v", err)
	}
}
