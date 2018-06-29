package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/jamesbcook/chatbot-external-api/api"
	"github.com/jamesbcook/chatbot-external-api/network"
	"github.com/jamesbcook/chatbot/kbchat"
)

var (
	port   = 40001
	chatID = os.Getenv("CHATBOT_TEST_CHATID")
)

func scan(args []byte) ([]byte, error) {
	cmd := exec.Command("nmap", strings.Split(string(args), " ")...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func setuplistener() error {
	server := fmt.Sprintf("localhost:%d", port)
	listener, err := network.Listen("tcp", server)
	if err != nil {
		return fmt.Errorf("Couldn't set up listener %v", err)
	}
	defer listener.Close()
	s, err := listener.Accept()
	if err != nil {
		return err
	}
	msg, err := s.ReceiveEncryptedMsg()
	if err != nil {
		return err
	}
	res, err := scan(msg.IO)
	if err != nil {
		return err
	}
	m := &api.Message{}
	m.ID = api.MessageID_Response
	m.IO = []byte(res)
	if err := s.SendEncryptedMsg(m); err != nil {
		return err
	}
	s.Close()
	return nil
}

func TestInfo(t *testing.T) {
	AP.Debug(false, nil)
	output, err := AP.Get("info")
	if err != nil {
		t.Fatal(err)
	}
	if len(output) <= 0 {
		t.Fatalf("Size of output %d", len(output))
	}
	t.Logf("Output:\n%s", output)
}

func TestGet(t *testing.T) {
	AP.Debug(false, nil)
	go func() {
		err := setuplistener()
		if err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(1 * time.Second)
	server := fmt.Sprintf("localhost:%d", port)
	args := fmt.Sprintf("-p %d localhost", port)
	input := fmt.Sprintf("%s %s", server, args)
	t.Log(input)
	results, err := AP.Get(input)
	if err != nil {
		t.Errorf("Error in Get %v", err)
	}
	t.Log(results)
	if len(results) == 0 {
		t.Errorf("Length of results is 0, this shouldn't be")
	}
}

func TestSend(t *testing.T) {
	AP.Debug(false, nil)
	sub := kbchat.SubscriptionMessage{}
	sub.Conversation.ID = chatID
	go func() {
		err := setuplistener()
		if err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(1 * time.Second)
	server := fmt.Sprintf("localhost:%d", port)
	args := fmt.Sprintf("-p %d localhost", port)
	input := fmt.Sprintf("%s %s", server, args)
	results, err := AP.Get(input)
	if err != nil {
		t.Errorf("Error in Get %v", err)
	}
	t.Log(results)
	if len(results) == 0 {
		t.Errorf("Length of results is 0, this shouldn't be")
	}
	if err := AP.Send(sub, results); err != nil {
		t.Fatalf("Error sending command to keybase %v", err)
	}
}

func TestRandomSecretKey(t *testing.T) {
	AP.Debug(false, nil)
	if err := randomSecretKey(); err != nil {
		t.Fatal(err)
	}
}
