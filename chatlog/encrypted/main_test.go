package main

import (
	"bufio"
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

func TestExportedWrite(t *testing.T) {
	written, err := Logger.Write([]byte("hello world"))
	if err != nil {
		t.Fatalf("Error writing to file %v", err)
	}
	if written <= 0 {
		t.Fatalf("Bytes written was 0 or less")
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

func TestDecrypt(t *testing.T) {
	written, err := Logger.Write([]byte("hello world"))
	if err != nil {
		t.Fatalf("Error writing to file %v", err)
	}
	if written <= 0 {
		t.Fatalf("Bytes written was 0 or less")
	}
	f, err := os.OpenFile(chatlog.LogFile, os.O_RDONLY, 0644)
	if err != nil {
		t.Fatalf("Error opening file as readonly %v", err)
	}
	buffer := make([]byte, 32)
	_, err = f.Read(buffer)
	if err != nil {
		t.Fatalf("Error reading 32 bytes from file %v", err)
	}
	var offset int64
	for x := range buffer {
		if buffer[x] == ']' {
			offset = int64(x)
			break
		}
	}
	_, err = f.Seek(offset+2, 0)
	if err != nil {
		t.Fatalf("Error seeking offset %d, %v", offset+2, err)
	}
	reader := bufio.NewReader(f)
	line, _, err := reader.ReadLine()
	if err != nil {
		t.Fatalf("Error reading line %v", err)
	}
	plaintext, err := Logger.Decrypt(line)
	if err != nil {
		t.Fatalf("Error decrypting line %v", err)
	}
	if string(plaintext) != "hello world" {
		t.Fatalf("Expected 'hello world', Got %s", plaintext)
	}
	if err := os.Remove(chatlog.LogFile); err != nil {
		t.Fatalf("Error removing log file %v", err)
	}
}
