package main

import (
	"fmt"
	"io"

	"github.com/jamesbcook/chatbot/kbchat"
)

var (
	//CMD that keybase will use to execute this plugin
	CMD = "/help"
	//Help is what will show in the help menu
	Help         = "/help this message"
	version      = "0.12.0"
	msg          string
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

const (
	header = "ChatBot v%s\n*Accepted Commands:*\n"
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
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Help Error] sending message %v", err)
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

func main() {}
