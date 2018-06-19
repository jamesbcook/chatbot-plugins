package main

import (
	"fmt"
	"io"

	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	header = "ChatBot v%s\n*Accepted Commands:*\n"
)

var (
	version      = "0.16.0"
	msg          string
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
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Help Error] sending message %v", err)
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

func main() {}
