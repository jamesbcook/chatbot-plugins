package main

import (
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
	user1 := "bob"
	if !Validate(user1) {
		t.Error("First time user failed")
	}
	if Validate(user1) {
		t.Error("User should have been blocked")
	}
	for x := 0; x < 15; x++ {
		Validate(user1)
	}
	if Validate(user1) {
		t.Error("User should have been blocked")
	}
	time.Sleep(time.Second * 10)
	if !Validate(user1) {
		t.Error("User should not have been blocked")
	}
}
