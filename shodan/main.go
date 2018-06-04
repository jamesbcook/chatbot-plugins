package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jamesbcook/chatbot/kbchat"
	"gopkg.in/ns3777k/go-shodan.v3/shodan"
)

type getting string

//Getter export symbol
var Getter getting

//Sender export symbol
var Sender getting

//Debugger export Symbol
var Debugger getting

var (
	//CMD that keybase will use to execute this plugin
	CMD = "/shodan"
	//Help is what will show in the help menu
	Help         = "/shodan {ip}"
	areDebugging = false
	debugWriter  *io.Writer
)

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
//This Get method will take a query virustotal with the given input
//and return the results of that file.
func (g getting) Get(input string) (string, error) {
	api := os.Getenv("CHATBOT_SHODAN")
	sc := shodan.NewClient(nil, api)
	debug(fmt.Sprintf("Query Shodan API for %s", input))
	res, err := sc.GetServicesForHost(context.Background(), input, nil)
	if err != nil {
		return "", fmt.Errorf("[Shodan Error] in get request %v", err)
	}
	var (
		output  string
		hosts   string
		service string
	)
	if res.OS != "" {
		output += fmt.Sprintf("OS: %s\n", res.OS)
	}
	if len(res.Hostnames) > 0 {
		output += "HostNames: "
		for _, hostname := range res.Hostnames {
			hosts += fmt.Sprintf("%s ", hostname)
		}
		output += hosts
		output += "\n"
	}
	if res.Organization != "" {
		output += fmt.Sprintf("Organization: %s\n", res.Organization)
	}
	if res.ASN != "" {
		output += fmt.Sprintf("ASN: %s\n", res.ASN)
	}
	if len(res.Data) > 0 {
		service += "Ports: "
		for _, data := range res.Data {
			service += fmt.Sprintf("%d/%s, ", data.Port, data.Transport)
		}
		service = strings.TrimSuffix(service, ", ")
		service += "\n"
		output += service
	}
	debug(fmt.Sprintf("Sending the following to user:\n%s", output))
	return output, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will respond with the results to the message ID that sent the request.
func (g getting) Send(msgID, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Shodan Error] in send request %v", err)
	}
	debug(fmt.Sprintf("Sending %s to %s", msg, msgID))
	if err := w.SendMessage(msgID, msg); err != nil {
		if err := w.Proc.Kill(); err != nil {
			return err
		}
		return err
	}
	debug("Killing child process")
	return w.Proc.Kill()
}

func init() {
	if api := os.Getenv("CHATBOT_SHODAN"); api == "" {
		log.Println("Missing CHATBOT_SHODAN environment variable")
	}
}

func main() {}
