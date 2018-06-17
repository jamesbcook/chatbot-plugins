package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/chatbot/kbchat/team"
)

type backgroundPlugin string
type authenticator string

//BP for export
var BP backgroundPlugin

//Auth for export
var Auth authenticator

var (
	users        = []string{}
	areDebugging = false
	debugWriter  *io.Writer
)

//Name that keybase will use for background plugins
func (b backgroundPlugin) Name() string {
	return "auth"
}

//Debug output
func (b backgroundPlugin) Debug(set bool, writer *io.Writer) {
	areDebugging = set
	debugWriter = writer
}

func debug(input string) {
	if areDebugging && debugWriter != nil {
		output := fmt.Sprintf("[DEBUG] %s\n", input)
		(*debugWriter).Write([]byte(output))
	}
}

//Start the process of gathering uses based on the team name in the env var.
func (a authenticator) Start() {
	for {
		w, err := kbchat.Start("team")
		if err != nil {
			debug(fmt.Sprintf("[Team Error] getting team api %v", err.Error()))
			time.Sleep(5 * time.Minute)
			continue
		}
		teamName := os.Getenv("CHATBOT_TEAM")
		output, err := team.Get(w, teamName, team.Members)
		if err != nil {
			debug(fmt.Sprintf("[Team Error] getting team members %v", err.Error()))
			continue
		}
		users = make([]string, len(output))
		for x, user := range output {
			users[x] = user.Username
		}
		if err := w.Proc.Kill(); err != nil {
			debug(fmt.Sprintf("[Team Error] killing process %v", err))
		}
		time.Sleep(5 * time.Minute)
	}
}

//Validate that a user is part of a team and allowed to send commands to the bot
func (a authenticator) Validate(user string) bool {
	for _, u := range users {
		if user == u {
			return true
		}
	}
	return false
}

func init() {
	teamName := os.Getenv("CHATBOT_TEAM")
	if teamName == "" {
		log.Println("Missing CHATBOT_TEAM environment variable")
	}
}

func main() {}
