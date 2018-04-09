package main

import "testing"

const (
	chatID = ""
)

func TestGet(t *testing.T) {
	res, err := Getter.Get("bitcoin")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results are less than or equal to 0")
	}
}

func TestSend(t *testing.T) {
	res, err := Getter.Get("bitcoin")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results are less than or equal to 0")
	}
	if err := Sender.Send(chatID, res); err != nil {
		t.Fatalf("Error sending to keybase %v", err)
	}
}
