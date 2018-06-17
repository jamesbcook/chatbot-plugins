package main

import (
	"os"
	"testing"
)

const (
	chatID = ""
)

func TestSend(t *testing.T) {
	err := os.Setenv("CHATBOT_URL_SHORTEN", "")
	if err != nil {
		t.Fatalf("Error setting env var %v", err)
	}
	output, err := AP.Get("https://google.com")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
	if err := AP.Send(chatID, output); err != nil {
		t.Fatalf("Error sending command to keybase %v", err)
	}
}

func TestGet(t *testing.T) {
	err := os.Setenv("CHATBOT_URL_SHORTEN", "")
	if err != nil {
		t.Fatalf("Error setting env var %v", err)
	}
	output, err := AP.Get("https://google.com")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
}
