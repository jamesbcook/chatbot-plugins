package main

import (
	"os"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	go Auth.Start()
	time.Sleep(2 * time.Second)
}

func TestValidate(t *testing.T) {
	expected := ""
	notExpected := "asdjafs;dflkjafsd;lkjf"
	if err := os.Setenv("CHATBOT_USERS", ""); err != nil {
		t.Fatalf("Error getting env var %v", err)
	}
	go Auth.Start()
	time.Sleep(2 * time.Second)
	if !Auth.Validate(expected) {
		t.Fatalf("Error validating user %s", expected)
	}

	if Auth.Validate(notExpected) {
		t.Fatalf("Error validated an unexpected user %s", notExpected)
	}

}
