package main

import (
	"os"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	go Start(os.Stdout)
	time.Sleep(2 * time.Second)
}

func TestValidate(t *testing.T) {
	expected := ""
	notExpected := "asdjafs;dflkjafsd;lkjf"
	if err := os.Setenv("CHATBOT_TEAM", ""); err != nil {
		t.Fatalf("Error getting env var %v", err)
	}
	go Start(os.Stdout)
	time.Sleep(2 * time.Second)
	if !Validate(expected) {
		t.Fatalf("Error validating user %s", expected)
	}

	if Validate(notExpected) {
		t.Fatalf("Error validated an unexpected user %s", notExpected)
	}

}
