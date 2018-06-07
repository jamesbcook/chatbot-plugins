package main

import (
	"testing"
)

const (
	chatID = ""
)

func TestGet(t *testing.T) {
	output, err := Getter.Get("15832795")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
}

func TestGetCreate(t *testing.T) {
	output, err := Getter.Get(`"This is my title" "something, something2, something3" "false" "normal" "true"`)
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
}

func TestSend(t *testing.T) {
	output, err := Getter.Get("15832795")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
	if err := Sender.Send(chatID, output); err != nil {
		t.Fatalf("Error sending command to keybase %v", err)
	}
}
