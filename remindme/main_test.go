package main

import (
	"testing"

	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	chatID = ""
)

func TestMinute(t *testing.T) {
	AP.Debug(false, nil)
	sub := kbchat.SubscriptionMessage{}
	sub.Conversation.ID = chatID
	output, err := AP.Get(`1 minute "something I want to know about"`)
	if err != nil {
		t.Fatal(err)
	}
	if err := AP.Send(sub, output); err != nil {
		t.Fatal(err)
	}
	output, err = AP.Get(`3 minutes "something I want to know about3"`)
	if err != nil {
		t.Fatal(err)
	}
	if err := AP.Send(sub, output); err != nil {
		t.Fatal(err)
	}
}

func TestHour(t *testing.T) {
	AP.Debug(false, nil)
	sub := kbchat.SubscriptionMessage{}
	sub.Conversation.ID = chatID
	output, err := AP.Get(`1 hour "something I want to know about2"`)
	if err != nil {
		t.Fatal(err)
	}
	if err := AP.Send(sub, output); err != nil {
		t.Fatal(err)
	}
}

func TestDay(t *testing.T) {
	AP.Debug(false, nil)
	sub := kbchat.SubscriptionMessage{}
	sub.Conversation.ID = chatID
	output, err := AP.Get(`4 days "something I want to know about4"`)
	if err != nil {
		t.Fatal(err)
	}
	if err := AP.Send(sub, output); err != nil {
		t.Fatal(err)
	}
}
