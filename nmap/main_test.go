package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/jamesbcook/chatbot-external-api/api"
	"github.com/jamesbcook/chatbot-external-api/network"
)

var (
	port   = 40001
	chatID = ""
)

func scan(args []byte) ([]byte, error) {
	cmd := exec.Command("nmap", strings.Split(string(args), " ")...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func setuplistener(t *testing.T) {
	server := fmt.Sprintf("localhost:%d", port)
	listener, err := network.Listen("tcp", server)
	if err != nil {
		t.Errorf("Couldn't set up listener %v", err)
	}
	s, err := listener.Accept()
	if err != nil {
		t.Fatal(err)
	}
	msg, err := s.ReceiveEncryptedMsg()
	if err != nil {
		t.Fatal(err)
	}
	res, err := scan(msg.IO)
	if err != nil {
		t.Fatal(err)
	}
	m := &api.Message{}
	m.ID = api.MessageID_Response
	m.IO = []byte(res)
	if err := s.SendEncryptedMsg(m); err != nil {
		t.Fatal(err)
	}
	s.Close()
}

func TestInfo(t *testing.T) {
	output, err := Getter.Get("info")
	if err != nil {
		t.Fatal(err)
	}
	if len(output) <= 0 {
		t.Fatalf("Size of output %d", len(output))
	}
	t.Logf("Output:\n%s", output)
}

func TestGet(t *testing.T) {
	go setuplistener(t)
	server := fmt.Sprintf("localhost:%d", port)
	args := fmt.Sprintf("-p %d localhost", port)
	input := fmt.Sprintf("%s %s", server, args)
	results, err := Getter.Get(input)
	if err != nil {
		t.Errorf("Error in Get %v", err)
	}
	t.Log(results)
	if len(results) == 0 {
		t.Errorf("Length of results is 0, this shouldn't be")
	}
}

func TestSend(t *testing.T) {
	go setuplistener(t)
	server := fmt.Sprintf("localhost:%d", port)
	args := fmt.Sprintf("-p %d localhost", port)
	input := fmt.Sprintf("%s %s", server, args)
	results, err := Getter.Get(input)
	if err != nil {
		t.Errorf("Error in Get %v", err)
	}
	t.Log(results)
	if len(results) == 0 {
		t.Errorf("Length of results is 0, this shouldn't be")
	}
	if err := Sender.Send(chatID, results); err != nil {
		t.Fatalf("Error sending command to keybase %v", err)
	}
}

func TestRandomSecretKey(t *testing.T) {
	if err := randomSecretKey(); err != nil {
		t.Fatal(err)
	}
}
