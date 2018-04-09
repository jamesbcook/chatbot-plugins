package main

import (
	"testing"
	"time"
)

const (
	chatID = ""
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
	res, err := Getter.Get("ajitvpai@gmail.com")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results are less than or equal to 0")
	}
}

func TestSend(t *testing.T) {
	time.Sleep(2 * time.Second) //was getting limited
	res, err := Getter.Get("ajitvpai@gmail.com")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results are less than or equal to 0")
	}
	if err := Sender.Send(chatID, res); err != nil {
		t.Fatalf("Error sending message to keybase %v", err)
	}
}
