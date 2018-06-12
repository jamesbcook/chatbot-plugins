package main

import (
	"io"
	"os"
	"testing"
	"time"
)

const (
	chatID = ""
)

func TestDebug(t *testing.T) {
	var out io.Writer
	out = os.Stdout
	Getter.Debug(true, &out)
	time.Sleep(2 * time.Second)
	res, err := Getter.Get("rvn")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results are less than or equal to 0")
	}
	t.Log(res)
}

func TestGet(t *testing.T) {
	time.Sleep(2 * time.Second)
	res, err := Getter.Get("bitcoin")
	if err != nil {
		t.Fatalf("Error in get request %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results are less than or equal to 0")
	}
	t.Log(res)
}

func TestSend(t *testing.T) {
	time.Sleep(2 * time.Second)
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
