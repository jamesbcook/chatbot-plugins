package main

import (
	"testing"
)

const (
	chatID = ""
)

func TestMinute(t *testing.T) {
	output, err := Getter.Get(`1 minute "something I want to know about"`)
	if err != nil {
		t.Fatal(err)
	}
	if err := Getter.Send("1234", output); err != nil {
		t.Fatal(err)
	}
	output, err = Getter.Get(`3 minutes "something I want to know about3"`)
	if err != nil {
		t.Fatal(err)
	}
	if err := Getter.Send("1234", output); err != nil {
		t.Fatal(err)
	}
}

func TestHour(t *testing.T) {
	output, err := Getter.Get(`1 hour "something I want to know about2"`)
	if err != nil {
		t.Fatal(err)
	}
	if err := Getter.Send("1234", output); err != nil {
		t.Fatal(err)
	}
}

func TestDay(t *testing.T) {
	output, err := Getter.Get(`4 days "something I want to know about4"`)
	if err != nil {
		t.Fatal(err)
	}
	if err := Getter.Send("1234", output); err != nil {
		t.Fatal(err)
	}
}
