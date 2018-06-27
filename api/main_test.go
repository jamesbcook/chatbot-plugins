package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/jamesbcook/chatbot-external-api/crypto"
	"github.com/jamesbcook/chatbot/kbchat"
)

var (
	chatID      = ""
	chatChannel = ""
)

func TestDebugExport(t *testing.T) {
	var output io.Writer
	output = os.Stdout
	AP.Debug(true, &output)
}

func TestDebugInternal(t *testing.T) {
	var output io.Writer
	output = os.Stdout
	AP.Debug(true, &output)
	debug("A debug statement")
}

func TestSendExport(t *testing.T) {
	sub := kbchat.SubscriptionMessage{}
	sub.Conversation.ID = chatID
	output, err := AP.Get("info")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
	if err := AP.Send(sub, output); err != nil {
		t.Fatalf("Error sending command to keybase %v", err)
	}
}

func TestSendInternal(t *testing.T) {
	output, err := AP.Get("info")
	if err != nil {
		t.Fatalf("Error getting info %v", err)
	}
	if len(output) <= 0 {
		t.Fatalf("Error in output no length %v", output)
	}
	if err := send(dm, chatChannel, output); err != nil {
		t.Fatalf("Error sending command to keybase %v", err)
	}
}

func TestInfo(t *testing.T) {
	out, err := AP.Get("info")
	if err != nil {
		t.Fatal(err)
	}
	if len(out) <= 0 {
		t.Fatalf("Len if output is %d", len(out))
	}
}

func TestInvalidCommand(t *testing.T) {
	_, err := AP.Get("Something")
	if err == nil {
		t.Fatal("This command should have failed")
	}
}

func TestAdd(t *testing.T) {
	var c crypto.ED25519
	if err := c.CreateKeys(); err != nil {
		t.Fatal(err)
	}
	out, err := AP.Get(fmt.Sprintf("add %s", hex.EncodeToString(c.PublicKey[:])))
	if err != nil {
		t.Fatal(err)
	}
	if len(out) <= 0 {
		t.Fatalf("Len if output is %d", len(out))
	}
	out2, err := AP.Get("info")
	if err != nil {
		t.Fatal(err)
	}
	if len(out2) <= 0 {
		t.Fatalf("Len if output is %d", len(out))
	}
	_, err = AP.Get(fmt.Sprintf("add %s", "123"))
	if err == nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	var c crypto.ED25519
	if err := c.CreateKeys(); err != nil {
		t.Fatal(err)
	}
	var c2 crypto.ED25519
	if err := c2.CreateKeys(); err != nil {
		t.Fatal(err)
	}
	_, err := AP.Get(fmt.Sprintf("add %s", hex.EncodeToString(c.PublicKey[:])))
	if err != nil {
		t.Fatal(err)
	}
	_, err = AP.Get(fmt.Sprintf("add %s", hex.EncodeToString(c2.PublicKey[:])))
	if err != nil {
		t.Fatal(err)
	}
	out, err := AP.Get("info")
	if err != nil {
		t.Fatal(err)
	}
	if len(out) <= 0 {
		t.Fatalf("Len if output is %d", len(out))
	}
	out2, err := AP.Get(fmt.Sprintf("delete %s", hex.EncodeToString(c.PublicKey[:])))
	if len(out2) <= 0 {
		t.Fatalf("Len if output is %d", len(out))
	}
	if out == out2 {
		t.Fatalf("Output before delete is the same once the delete was attempted")
	}
	_, err = AP.Get(fmt.Sprintf("delete %s", "123"))
	if err == nil {
		t.Fatal(err)
	}
}
