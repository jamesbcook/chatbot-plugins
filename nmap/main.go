package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
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
	//CMD that keybase will use to execute this plugin
	CMD = "/nmap"
	//Help is what will show in the help menu
	Help         = `/nmap {info|apiIP:apiPort nmap args}`
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
	if areDebugging && debugWriter != nil {
		output := fmt.Sprintf("[DEBUG] %s\n", input)
		(*debugWriter).Write([]byte(output))
	}
}

//Get export method that satisfies an interface in the main program.
//This Get method will send a request to the chatbot-extern-api server with the
//given nmap arguments and return the results of the scan.
func (g getting) Get(input string) (string, error) {
	debug(fmt.Sprintf("Got %s for input", input))
	args := strings.Split(input, " ")
	if args[0] == "info" {
		var output string
		output = fmt.Sprintf("Use the following key for authentication\n")
		output += fmt.Sprintf("Public Key:%s", network.GetIdentityKey())
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
func (g getting) Send(msgID, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Nmap Error] in send request %v", err)
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

func saveFile(fileName string, input []byte) error {
	return ioutil.WriteFile(fileName, input, 0400)
}

func randomSecretKey() error {
	debug("Generating random secret key pair")
	if err := network.GenerateSecretKeyPair(); err != nil {
		return err
	}
	skFile, err := filesystem.GetPrivateKeyFile(app)
	if err != nil {
		return err
	}
	pkFile, err := filesystem.GetPublicKeyFile(app)
	if err != nil {
		return err
	}
	if err := saveFile(skFile, []byte(network.GetSecretKey())); err != nil {
		return err
	}
	if err := saveFile(pkFile, []byte(network.GetIdentityKey())); err != nil {
		return err
	}
	return nil
}

func loadFile(input string) ([]byte, error) {
	f, err := os.Open(input)
	if err != nil {
		return nil, err
	}
	output, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return decodeHex(output)
}

func decodeHex(input []byte) ([]byte, error) {
	output := make([]byte, hex.DecodedLen(len(input)))
	_, err := hex.Decode(output, input)
	if err != nil {
		return nil, err
	}
	return output, err
}

func init() {
	skFile, err := filesystem.GetPrivateKeyFile(app)
	if err != nil {
		log.Println(err)
		return
	}
	pkFile, err := filesystem.GetPublicKeyFile(app)
	if err != nil {
		log.Println(err)
		return
	}
	priv, err := loadFile(skFile)
	if err != nil {
		log.Println(err)
		if err := randomSecretKey(); err != nil {
			log.Println(err)
		}
		return
	}
	pub, err := loadFile(pkFile)
	if err != nil {
		log.Println(err)
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
