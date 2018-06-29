package main

import (
	"os"
	"testing"

	"github.com/jamesbcook/chatbot/kbchat"
)

var (
	chatID = os.Getenv("CHATBOT_TEST_CHATID")
)

func TestGet(t *testing.T) {
	AP.Debug(false, nil)
	output, err := AP.Get("15832795")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
}

func TestGetCreate(t *testing.T) {
	AP.Debug(false, nil)
	output, err := AP.Get(`"This is my title" "something, something2, something3" "false" "normal" "true"`)
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
}

func TestSend(t *testing.T) {
	AP.Debug(false, nil)
	sub := kbchat.SubscriptionMessage{}
	sub.Conversation.ID = chatID
	output, err := AP.Get("15832795")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
	if err := AP.Send(sub, output); err != nil {
		t.Fatalf("Error sending command to keybase %v", err)
	}
}
