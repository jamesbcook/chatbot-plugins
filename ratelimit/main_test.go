package main

import (
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
	user1 := "bob"
	if !Auth.Validate(user1) {
		t.Error("First time user failed")
	}
	if Auth.Validate(user1) {
		t.Error("User should have been blocked")
	}
	for x := 0; x < 15; x++ {
		Auth.Validate(user1)
	}
	if Auth.Validate(user1) {
		t.Error("User should have been blocked")
	}
	time.Sleep(time.Second * 10)
	if !Auth.Validate(user1) {
		t.Error("User should not have been blocked")
	}
}
