package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jamesbcook/chatbot-external-api/api"
	"github.com/jamesbcook/chatbot-external-api/network"
	"github.com/jamesbcook/chatbot/kbchat"
)

var (
	//CMD that keybase will use to execute this plugin
	CMD = "/nmap"
	//Help is what will show in the help menu
	Help         = `/nmap {"apiIP:apiPort" "nmap args" }`
	areDebugging = false
	debugWriter  *io.Writer
)

type getting string

//Getter export symbol
var Getter getting

//Sender export symbol
var Sender getting

//Debugger export Symbol
var Debugger getting

func (g getting) Debug(set bool, writer *io.Writer) {
	areDebugging = set
	debugWriter = writer
}

func debug(input string) {
	if areDebugging {
		output := fmt.Sprintf("[DEBUG] %s\n", input)
		(*debugWriter).Write([]byte(output))
	}
}

//Get export method that satisfies an interface in the main program.
//This Get method will send a request to the chatbot-extern-api server with the
//given nmap arguments and return the results of the scan.
func (g getting) Get(input string) (string, error) {
	debug(fmt.Sprintf("Got %s for input", input))
	output := strings.FieldsFunc(input, func(c rune) bool {
		if c != '"' {
			return false
		}
		return true
	})
	if len(output) != 3 {
		return "", fmt.Errorf("Not enough arguments")
	}
	server := output[0]
	args := output[2]
	debug(fmt.Sprintf("Connecting to %s", server))
	s, err := network.Dial("tcp", server)
	if err != nil {
		return "", err
	}
	defer s.Close()
	msg := &api.Message{}
	msg.ID = api.MessageID_Nmap
	msg.IO = []byte(args)
	debug("Sending encrypted message")
	if err := s.SendEncryptedMsg(msg); err != nil {
		return "", err
	}
	debug("Receive encrypted message")
	recv, err := s.ReceiveEncryptedMsg()
	if err != nil {
		return "", err
	}
	msg.ID = api.MessageID_Done
	msg.IO = make([]byte, 1)
	debug("Sending Done")
	if err := s.SendEncryptedMsg(msg); err != nil {
		return "", err
	}
	debug(fmt.Sprintf("Returning %s", recv.IO))
	return string(recv.IO), nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (g getting) Send(msgID, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Reddit Error] in send request %v", err)
	}
	debug(fmt.Sprintf("Sending this message to messageID: %s\n%s", msgID, msg))
	if err := w.SendMessage(msgID, msg); err != nil {
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
	return nil
}

func init() {
	//Give the debug var time to get set
	//time.Sleep(40 * time.Second)
	privFile, err := os.Open("key.priv")
	if err != nil {
		log.Println("Couldn't open key.priv file")
		if err := randomSecretKey(); err != nil {
			log.Println(err)
		}
		return
	}
	pubFile, err := os.Open("key.pub")
	if err != nil {
		log.Println("Couldn't open key.pub file")
		if err := randomSecretKey(); err != nil {
			log.Println(err)
		}
		return
	}
	priv, err := ioutil.ReadAll(privFile)
	if err != nil {
		log.Println("Couldn't read key.priv file")
		if err := randomSecretKey(); err != nil {
			log.Println(err)
		}
		return
	}
	pub, err := ioutil.ReadAll(pubFile)
	if err != nil {
		log.Println("Couldn't read key.pub file")
		if err := randomSecretKey(); err != nil {
			log.Println(err)
		}
		return
	}
	decodePriv, err := hex.DecodeString(string(priv))
	if err != nil {
		log.Println("Couldn't decode priv hex")
		if err := randomSecretKey(); err != nil {
			log.Println(err)
		}
		return
	}
	decodePub, err := hex.DecodeString(string(pub))
	if err != nil {
		log.Println("Couldn't decode pub hex")
		if err := randomSecretKey(); err != nil {
			log.Println(err)
		}
		return
	}
	debug("Setting key pair")
	if err := network.SetSecretKeyPair(decodePriv, decodePub); err != nil {
		log.Println("Couldn't set secret key pair")
		if err := randomSecretKey(); err != nil {
			log.Println(err)
		}
		return
	}
}

func main() {}
