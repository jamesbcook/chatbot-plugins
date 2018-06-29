package main

import (
	"os"
	"testing"
	"time"

	"github.com/jamesbcook/chatbot/kbchat"
)

var (
	chatID = os.Getenv("CHATBOT_TEST_CHATID")
)

func TestAllBreachesForAccount(t *testing.T) {
	accountName := "example@gmail.com"
	expected := "https://haveibeenpwned.com/api/v2/breachedaccount/example@gmail.com"
	res := allBreachesForAccount(accountName)
	if res != expected {
		t.Fatalf("Expected %s Got %s", expected, res)
	}
}

func TestAllPastesForAccount(t *testing.T) {
	accountName := "example@gmail.com"
	expected := "https://haveibeenpwned.com/api/v2/pasteaccount/example@gmail.com"
	res := allPastesForAccount(accountName)
	if res != expected {
		t.Fatalf("Expected %s Got %s", expected, res)
	}
}

func TestGet(t *testing.T) {
	AP.Debug(false, nil)
	res, err := AP.Get("ajitvpai@gmail.com")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results are less than or equal to 0")
	}
}

func TestSend(t *testing.T) {
	AP.Debug(false, nil)
	sub := kbchat.SubscriptionMessage{}
	sub.Conversation.ID = chatID
	time.Sleep(2 * time.Second) //was getting limited
	res, err := AP.Get("ajitvpai@gmail.com")
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
