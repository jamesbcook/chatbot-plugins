package main

import (
	"os"
	"testing"

	"github.com/jamesbcook/chatbot-plugins/chatlog"
)

func TestStart(t *testing.T) {
	l, err := start()
	if err != nil {
		t.Fatalf("Error staring log %v", err)
	}
	if l == nil {
		t.Fatalf("Logger returned nil %v", err)
	}
}

func TestUnExportedWrite(t *testing.T) {
	l, err := start()
	if err != nil {
		t.Fatalf("Error staring log %v", err)
	}
	if l == nil {
		t.Fatalf("Logger returned nil %v", err)
	}
	written, err := l.write([]byte("hello world"))
	if err != nil {
		t.Fatalf("Error writing to file %v", err)
	}
	if written <= 0 {
		t.Fatalf("Bytes written was 0 or less")
	}
}

func TestExportedWrite(t *testing.T) {
	written, err := Logger.Write([]byte("hello world"))
	if err != nil {
		t.Fatalf("Error writing to file %v", err)
	}
	if written <= 0 {
		t.Fatalf("Bytes written was 0 or less")
	}
	if err := os.Remove(chatlog.LogFile); err != nil {
		t.Fatalf("Error removing log file %v", err)
	}
}
