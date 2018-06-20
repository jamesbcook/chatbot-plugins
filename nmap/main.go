package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"

	"github.com/jamesbcook/chatbot-external-api/api"
	"github.com/jamesbcook/chatbot-external-api/filesystem"
	"github.com/jamesbcook/chatbot-external-api/network"
	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	app = "nmap"
)

var (
	areDebugging = false
	debugWriter  *io.Writer
)

type activePlugin string

//AP for export
var AP activePlugin

func (a activePlugin) Debug(set bool, writer *io.Writer) {
	areDebugging = set
	debugWriter = writer
}

func debug(input string) {
	if areDebugging && debugWriter != nil {
		output := fmt.Sprintf("[DEBUG] %s\n", input)
		(*debugWriter).Write([]byte(output))
	}
}

//CMD that keybase will use to execute this plugin
func (a activePlugin) CMD() string {
	return "/nmap"
}

//Help is what will show in the help menu
func (a activePlugin) Help() string {
	return "/nmap {info|apiIP:apiPort nmap args}"
}

//Get export method that satisfies an interface in the main program.
//This Get method will send a request to the chatbot-extern-api server with the
//given nmap arguments and return the results of the scan.
func (a activePlugin) Get(input string) (string, error) {
	debug(fmt.Sprintf("Got %s for input", input))
	args := strings.Split(input, " ")
	if args[0] == "info" {
		var output string
		output = fmt.Sprintf("Use the following key for authentication\n")
		output += fmt.Sprintf("Public Key: %s", network.GetIdentityKey())
		return output, nil
	}
	server := args[0]
	nmapArgs := strings.Join(args[1:], " ")
	debug(fmt.Sprintf("Connecting to %s", server))
	s, err := network.Dial("tcp", server)
	if err != nil {
		return "", err
	}
	msg := &api.Message{}
	msg.ID = api.MessageID_Nmap
	msg.IO = []byte(nmapArgs)
	debug("Sending encrypted message")
	if err := s.SendEncryptedMsg(msg); err != nil {
		s.Close()
		return "", err
	}
	debug("Receive encrypted message")
	recv, err := s.ReceiveEncryptedMsg()
	if err != nil {
		s.Close()
		return "", err
	}
	msg.ID = api.MessageID_Done
	length := rand.Intn(48)
	buf := make([]byte, length)
	rand.Read(buf)
	msg.IO = buf
	debug("Sending Done")
	if err := s.SendEncryptedMsg(msg); err != nil {
		s.Close()
		return "", err
	}
	debug(fmt.Sprintf("Returning %s", recv.IO))
	s.Close()
	return string(recv.IO), nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (a activePlugin) Send(subscription kbchat.SubscriptionMessage, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Nmap Error] in send request %v", err)
	}
	debug(fmt.Sprintf("Sending this message to messageID: %s\n%s", subscription.Conversation.ID, msg))
	if err := w.SendMessage(subscription.Conversation.ID, msg); err != nil {
		if err := w.Proc.Kill(); err != nil {
			return err
		}
		return err
	}
	debug("Killing child process")
	return w.Proc.Kill()
}

func randomSecretKey() error {
	debug("Generating random secret key pair")
	if err := network.GenerateSecretKeyPair(); err != nil {
		return err
	}
	fs, err := filesystem.New(app)
	if err != nil {
		return err
	}
	if err := fs.SaveKeyToFile([]byte(network.GetIdentityKey()), fs.GetPublicKeyFile()); err != nil {
		return err
	}
	if err := fs.SaveKeyToFile([]byte(network.GetSecretKey()), fs.GetPrivateKeyFile()); err != nil {
		return err
	}
	return nil
}

func init() {
	fs, err := filesystem.New(app)
	if err != nil {
		log.Println(err)
		return
	}
	priv, err := fs.LoadPrivateKeyFile()
	if err != nil {
		if err := randomSecretKey(); err != nil {
			log.Println(err)
		}
		return
	}
	pub, err := fs.LoadPublicKeyFile()
	if err != nil {
		if err := randomSecretKey(); err != nil {
			log.Println(err)
		}
		return
	}
	debug("Setting key pair")
	if err := network.SetSecretKeyPair(priv, pub); err != nil {
		log.Println("Couldn't set secret key pair")
		if err := randomSecretKey(); err != nil {
			log.Println(err)
		}
	}
}

func main() {}
