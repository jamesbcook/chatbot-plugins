package main

import (
	"testing"
)

const (
	chatID = ""
)

func TestSend(t *testing.T) {
	output, err := AP.Get("bd2c2cf0631d881ed382817afcce2b093f4e412ffb170a719e2762f250abfea4")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
	if err := AP.Send(chatID, output); err != nil {
		t.Fatalf("Error sending command to keybase %v", err)
	}
}

func TestGet(t *testing.T) {
	output, err := AP.Get("bd2c2cf0631d881ed382817afcce2b093f4e412ffb170a719e2762f250abfea4")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
}
