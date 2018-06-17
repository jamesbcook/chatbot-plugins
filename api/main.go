package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/jamesbcook/chatbot-external-api/api"
	"github.com/jamesbcook/chatbot-external-api/network"
	"github.com/jamesbcook/chatbot/kbchat"
)

var (
	areDebugging = false
	debugWriter  *io.Writer
	keys         []string
	serv         = &server{}
)

type activePlugin string

//AP for export
var AP activePlugin

type server struct {
	Port      string
	PrivateIP string
	PublicIP  string
}

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
	return "/api"
}

//Help is what will show in the help menu
func (a activePlugin) Help() string {
	return "/api {info|add {public key}|delete {public key}}"
}

//Get export method that satisfies an interface in the main program.
//This Get method will query reddit json api.
func (a activePlugin) Get(input string) (string, error) {
	debug(fmt.Sprintf("Got the input %s", input))
	args := strings.Split(input, " ")
	var output string
	switch args[0] {
	case "info":
		debug("Gathering server info")
		output = fmt.Sprintf("Server Info\nPublic %s:%s\n", serv.PublicIP, serv.Port)
		output += fmt.Sprintf("Private %s:%s\n", serv.PrivateIP, serv.Port)
		if len(keys) > 0 {
			output += fmt.Sprintf("Imported Keys\n")
			for x := range keys {
				output += fmt.Sprintf("Key%d: %s\n", x, keys[x])
			}
		}
	case "add":
		debug("Adding public key")
		if len(args[1]) != 64 {
			return "", fmt.Errorf("Invalid public key")
		}
		output = fmt.Sprintf("Adding %s", args[1])
		keys = append(keys, args[1])
	case "delete":
		debug("Deleting pubic key")
		if len(args[1]) != 64 {
			return "", fmt.Errorf("Invalid public key")
		}
		output = fmt.Sprintf("Deleting %s", args[1])
		for x, key := range keys {
			if key == args[1] {
				copy(keys[x:], keys[x+1:])
				keys[len(keys)-1] = ""
				keys = keys[:len(keys)-1]
			}
		}
	default:
		debug("A wrong command was sent")
		return "", fmt.Errorf("Wrong command %s", args[0])
	}
	debug(fmt.Sprintf("Returning output %s", output))
	return output, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (a activePlugin) Send(msgID, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[API Error] in send request %v", err)
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

func send(channelName, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[API Error] in send request %v", err)
	}
	debug(fmt.Sprintf("Sending this message to channel: %s\n%s", channelName, msg))
	if err := w.SendMessageByTlfName(channelName, msg); err != nil {
		if err := w.Proc.Kill(); err != nil {
			return err
		}
		return err
	}
	debug("Killing child process")
	return w.Proc.Kill()
}

func validKey(theirPK []byte) bool {
	for _, key := range keys {
		decoded, err := hex.DecodeString(key)
		if err != nil {
			return false
		}
		if bytes.Compare(theirPK, decoded) == 0 {
			return true
		}
	}
	return false
}

func handle(session *network.Session) {
	for {
		debug("Getting Encrypted Message")
		msg, err := session.ReceiveEncryptedMsg()
		if err != nil {
			log.Println(err)
			session.Close()
			return
		}
		if !validKey(session.Keys.TheirIdentityKey[:]) {
			session.Close()
			return
		}
		switch msg.ID {
		case api.MessageID_Beacon:
			debug("Got Beacon Message")
			send(msg.Channel, string(msg.IO))
			length := rand.Intn(48)
			buf := make([]byte, length)
			rand.Read(buf)
			m := &api.Message{}
			m.ID = api.MessageID_Response
			m.IO = buf
			debug("Sending Encrypted Response")
			if err := session.SendEncryptedMsg(m); err != nil {
				log.Println(err)
				session.Close()
				return
			}
		case api.MessageID_Hash:
			debug("Got Hash Message")
			fmt.Println("Not Implemented")
		case api.MessageID_Nmap:
			debug("Got Nmap Message")
			fmt.Println("Not Implemented")
		case api.MessageID_Done:
			debug("Got Done Message")
			session.Close()
			return
		}
	}
}

func startListener() {
	host := fmt.Sprintf(":%s", serv.Port)
	debug("starting server")
	l, err := network.Listen("tcp", host)
	if err != nil {
		log.Println(err)
		return
	}
	defer l.Close()
	for {
		debug("Waiting for connection")
		s, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		debug("Got connection")
		go handle(s)
	}
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://ifconfig.co/ip")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bod, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bod), nil
}

func getPrivateIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func init() {
	port := os.Getenv("CHATBOT_LISTEN_PORT")
	if port == "" {
		log.Println("CHATBOT_LISTEN_PORT not set, using 55449")
		port = "55449"
	}
	serv.Port = port
	debug("Getting Public IP")
	pub, err := getPublicIP()
	if err != nil {
		log.Println(err)
		return
	}
	serv.PublicIP = strings.TrimSuffix(pub, "\n")
	debug("Getting Private IP")
	priv, err := getPrivateIP()
	if err != nil {
		log.Println(err)
		return
	}
	serv.PrivateIP = priv
	debug(fmt.Sprintf("Server info %v", serv))
	go startListener()
}

func main() {}
