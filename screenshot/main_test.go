package main

import (
	"testing"
)

const (
	chatID = ""
)

func TestSend(t *testing.T) {
	output, err := Getter.Get("https://google.com")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
	if err := Sender.Send(chatID, output); err != nil {
		t.Fatalf("Error sending command to keybase %v", err)
	}
}

func TestGet(t *testing.T) {
	output, err := Getter.Get("https://google.com")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
}