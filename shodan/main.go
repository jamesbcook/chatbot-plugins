package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
	"gopkg.in/ns3777k/go-shodan.v3/shodan"
)

var (
	debugPrintf func(format string, v ...interface{})
)

type activePlugin string

//AP for export
var AP activePlugin

func (a activePlugin) Debug(set bool, writer *io.Writer) {
	debugPrintf = print.Debugf(set, writer)
}

//CMD that keybase will use to execute this plugin
func (a activePlugin) CMD() string {
	return "/shodan"
}

//Help is what will show in the help menu
func (a activePlugin) Help() string {
	return "/shodan {ip}"
}

//Get export method that satisfies an interface in the main program.
//This Get method will take a query virustotal with the given input
//and return the results of that file.
func (a activePlugin) Get(input string) (string, error) {
	api := os.Getenv("CHATBOT_SHODAN")
	sc := shodan.NewClient(nil, api)
	debugPrintf("Query Shodan API for %s\n", input)
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
	debugPrintf("Sending the following to user:\n%s\n", output)
	return output, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will respond with the results to the message ID that sent the request.
func (a activePlugin) Send(subscription kbchat.SubscriptionMessage, msg string) error {
	debugPrintf("Starting kbchat\n")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Shodan Error] in send request %v", err)
	}
	debugPrintf("Sending %s to %s\n", msg, subscription.Conversation.ID)
	if err := w.SendMessage(subscription.Conversation.ID, msg); err != nil {
		if err := w.Proc.Kill(); err != nil {
			return err
		}
		return err
	}
	debugPrintf("Killing child process\n")
	return w.Proc.Kill()
}

func init() {
	if api := os.Getenv("CHATBOT_SHODAN"); api == "" {
		print.Warningln("Missing CHATBOT_SHODAN environment variable")
	}
}

func main() {}
