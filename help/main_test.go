package main

import "testing"

const (
	chatID = ""
)

func TestGet(t *testing.T) {
	res, err := Getter.Get("Helper\nTesting")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if res != "" {
		t.Fatalf("Results should be empty")
	}
}

func TestSend(t *testing.T) {
	_, err := Getter.Get("Helper\nTesting")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	res, err := Getter.Get("")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if res == "" {
		t.Fatalf("Results should not be empty the second time")
	}
	if err := Sender.Send(chatID, res); err != nil {
		t.Fatalf("Error sending to keybase %v", err)
	}
}
