package main

import (
	"fmt"

	"github.com/jamesbcook/chatbot/kbchat"
)

var (
	//CMD that keybase will use to execute this plugin
	CMD = "/help"
	//Help is what will show in the help menu
	Help    = "/help this message"
	version = "0.5.0"
	msg     string
)

type getting string

//Getter export symbol
var Getter getting

//Sender export symbol
var Sender getting

const (
	header = "ChatBot v%s\n*Accepted Commands:*\n"
)

//Get export method that satisfies an interface in the main program.
//This Get method will generate a help message if input is not empty.
//If input is empty it will return the formated help message.
func (g getting) Get(input string) (string, error) {
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
func (g getting) Send(msgID, msg string) error {
	w, err := kbchat.Start("chat")
	if err != nil {
		return err
	}
	return w.SendMessage(msgID, msg)
}

func main() {}
