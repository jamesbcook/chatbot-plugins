package main

import (
	"testing"

	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	chatID = ""
)

func TestPwnedPassword(t *testing.T) {
	pass := "password"
	expected := "https://api.pwnedpasswords.com/pwnedpassword/" + pass
	res := pwnedPassword(pass)
	if res != expected {
		t.Fatalf("Expected %s Got %s", expected, res)
	}
}

func TestPwnedPasswordRange(t *testing.T) {
	pass := "2aae6"
	expected := "https://api.pwnedpasswords.com/range/" + pass
	res := pwnedPasswordRange(pass)
	if res != expected {
		t.Fatalf("Expected %s Got %s", expected, res)
	}
}

func TestGet(t *testing.T) {
	res, err := AP.Get("hunter2")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results are less than or equal to 0")
	}
}

func TestSend(t *testing.T) {
	sub := kbchat.SubscriptionMessage{}
	sub.Conversation.ID = chatID
	res, err := AP.Get("hunter2")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results are less than or equal to 0")
	}
	if err := AP.Send(sub, res); err != nil {
		t.Fatalf("Error sending message to keybase %v", err)
	}
}
