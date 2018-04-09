package main

import (
	"log"
	"os"
	"testing"

	"github.com/jamesbcook/chat-bot-plugins/chatlog"
)

func init() {
	if err := os.Setenv("CHATBOT_LOG_KEY", "43f287ac2487750aaf4b3cafa3f4c979"); err != nil {
		log.Fatal(err)
	}
}

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
}

func TestDecrypt(t *testing.T) {
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
