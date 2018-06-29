package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/jamesbcook/chatbot-external-api/api"
	"github.com/jamesbcook/chatbot-external-api/network"
	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
)

var (
	debugPrintf func(format string, v ...interface{})
	keys        []string
	serv        = &server{}
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
	debugPrintf = print.Debugf(set, writer)
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
	debugPrintf("Got the input %s\n", input)
	args := strings.Split(input, " ")
	var output string
	switch args[0] {
	case "info":
		debugPrintf("Gathering server info\n")
		output = fmt.Sprintf("Server Info\nPublic %s:%s\n", serv.PublicIP, serv.Port)
		output += fmt.Sprintf("Private %s:%s\n", serv.PrivateIP, serv.Port)
		if len(keys) > 0 {
			output += fmt.Sprintf("Imported Keys\n")
			for x := range keys {
				output += fmt.Sprintf("Key%d: %s\n", x, keys[x])
			}
		}
	case "add":
		debugPrintf("Adding public key\n")
		if len(args[1]) != 64 {
			return "", fmt.Errorf("Invalid public key")
		}
		output = fmt.Sprintf("Adding %s", args[1])
		keys = append(keys, args[1])
	case "delete":
		debugPrintf("Deleting pubic key\n")
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
		debugPrintf("A wrong command was sent\n")
		return "", fmt.Errorf("Wrong command %s", args[0])
	}
	debugPrintf("Returning output %s\n", output)
	return output, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (a activePlugin) Send(subscription kbchat.SubscriptionMessage, msg string) error {
	debugPrintf("Starting kbchat\n")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[API Error] in send request %v", err)
	}
	debugPrintf("Sending this message to messageID: %s\n%s\n", subscription.Conversation.ID, msg)
	if err := w.SendMessage(subscription.Conversation.ID, msg); err != nil {
		if err := w.Proc.Kill(); err != nil {
			return err
		}
		return err
	}
	debugPrintf("Killing child process\n")
	return w.Proc.Kill()
}

func team(c kbchat.API, args ...string) error {
	name := args[0]
	channel := args[1]
	msg := args[2]
	debugPrintf("Sending this message to team %s in channel: %s\n%s\n", name, channel, msg)
	return c.SendMessageByTeamName(name, msg, &channel)
}

func dm(c kbchat.API, args ...string) error {
	name := args[0]
	msg := args[1]
	debugPrintf("Sending this message to channel: %s\n%s\n", name, msg)
	return c.SendMessageByTlfName(name, msg)
}

func send(f func(c kbchat.API, args ...string) error, args ...string) error {
	debugPrintf("Starting kbchat\n")
	c, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[API Error] in send request %v", err)
	}
	if err := f(*c, args...); err != nil {
		if err := c.Proc.Kill(); err != nil {
			return err
		}
		return err
	}
	debugPrintf("Killing child process\n")
	return c.Proc.Kill()
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
		debugPrintf("Getting Encrypted Message\n")
		msg, err := session.ReceiveEncryptedMsg()
		if err != nil {
			print.Warningln(err)
			session.Close()
			return
		}
		debugPrintf("Got message\n%v\n", msg)
		if !validKey(session.Keys.TheirIdentityKey[:]) {
			session.Close()
			return
		}
		switch msg.ID {
		case api.MessageID_Beacon, api.MessageID_Nmap:
		case api.MessageID_Done:
			debugPrintf("Got Done Message\n")
			session.Close()
			return
		default:
			debugPrintf("Not a matching ID, closing session\n")
			session.Close()
			return
		}
		debugPrintf("Got %s Message\n", msg.ID.String())
		if msg.ChatType == api.ChatType_Team {
			send(team, msg.Chat.Team, msg.Chat.Channel, string(msg.IO))
		} else {
			send(dm, msg.Chat.Channel, string(msg.IO))
		}
		length := rand.Intn(48)
		buf := make([]byte, length)
		rand.Read(buf)
		m := &api.Message{}
		m.ID = api.MessageID_Response
		m.IO = buf
		debugPrintf("Sending Encrypted Response\n")
		if err := session.SendEncryptedMsg(m); err != nil {
			print.Warningln(err)
			session.Close()
			return
		}
	}
}

func startListener() {
	host := fmt.Sprintf(":%s", serv.Port)
	debugPrintf("starting server\n")
	l, err := network.Listen("tcp", host)
	if err != nil {
		print.Warningln(err)
		return
	}
	defer l.Close()
	for {
		debugPrintf("Waiting for connection\n")
		s, err := l.Accept()
		if err != nil {
			print.Warningln(err)
			continue
		}
		debugPrintf("Got connection\n")
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
	debugPrintf = func(format string, v ...interface{}) {
	}
	port := os.Getenv("CHATBOT_LISTEN_PORT")
	if port == "" {
		print.Warningln("CHATBOT_LISTEN_PORT not set, using 55449")
		port = "55449"
	}
	serv.Port = port
	debugPrintf("Getting Public IP\n")
	pub, err := getPublicIP()
	if err != nil {
		print.Warningln(err)
		return
	}
	serv.PublicIP = strings.TrimSuffix(pub, "\n")
	debugPrintf("Getting Private IP\n")
	priv, err := getPrivateIP()
	if err != nil {
		print.Warningln(err)
		return
	}
	serv.PrivateIP = priv
	debugPrintf("Server info %v\n", serv)
	go startListener()
}

func main() {}
