package main

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	chatID = ""
)

func TestDebug(t *testing.T) {
	var out io.Writer
	out = os.Stdout
	AP.Debug(true, &out)
	time.Sleep(2 * time.Second)
	res, err := AP.Get("rvn")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results are less than or equal to 0")
	}
	t.Log(res)
}

func TestGet(t *testing.T) {
	AP.Debug(false, nil)
	time.Sleep(2 * time.Second)
	res, err := AP.Get("bitcoin")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results are less than or equal to 0")
	}
	t.Log(res)
}

func TestSend(t *testing.T) {
	AP.Debug(false, nil)
	sub := kbchat.SubscriptionMessage{}
	sub.Conversation.ID = chatID
	time.Sleep(2 * time.Second)
	res, err := AP.Get("bitcoin")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results are less than or equal to 0")
	}
	if err := AP.Send(sub, res); err != nil {
		t.Fatalf("Error sending to keybase %v", err)
	}
}
