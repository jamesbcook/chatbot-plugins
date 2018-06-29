package main

import (
	"fmt"
	"io"

	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
)

const (
	header = "ChatBot v%s\n*Accepted Commands:*\n"
)

var (
	version     = "0.17.0"
	msg         string
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
	return "/help"
}

//Help is what will show in the help menu
func (a activePlugin) Help() string {
	return "/help this message"
}

//Get export method that satisfies an interface in the main program.
//This Get method will generate a help message if input is not empty.
//If input is empty it will return the formated help message.
func (a activePlugin) Get(input string) (string, error) {
	if input != "" {
		fmtHeader := fmt.Sprintf(header, version)
		start := "```\n"
		end := "```"
		msg = fmtHeader + start + input + end
		return "", nil
	}
	return msg, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (a activePlugin) Send(subscription kbchat.SubscriptionMessage, msg string) error {
	debugPrintf("Starting kbchat\n")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Help Error] sending message %v", err)
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

func main() {}
