package main

import (
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	go Auth.Start()
	time.Sleep(2 * time.Second)
}

func TestValidate(t *testing.T) {
	expected := "chatbot2"
	notExpected := "asdjafs;dflkjafsd;lkjf"
	go Auth.Start()
	time.Sleep(2 * time.Second)
	if !Auth.Validate(expected) {
		t.Fatalf("Error validating user %s", expected)
	}

	if Auth.Validate(notExpected) {
		t.Fatalf("Error validated an unexpected user %s", notExpected)
	}

}
