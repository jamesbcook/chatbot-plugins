package main

import (
	"io"
	"os"
	"time"

	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/chatbot/kbchat/team"
	"github.com/jamesbcook/print"
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
	debugPrintf  func(format string, v ...interface{})
)

//Name that keybase will use for background plugins
func (b backgroundPlugin) Name() string {
	return "auth"
}

//Debug output
func (b backgroundPlugin) Debug(set bool, writer *io.Writer) {
	debugPrintf = print.Debugf(set, writer)
}

//Start the process of gathering uses based on the team name in the env var.
func (a authenticator) Start() {
	for {
		w, err := kbchat.Start("team")
		if err != nil {
			debugPrintf("[Team Error] getting team api %v\n", err.Error())
			time.Sleep(5 * time.Minute)
			continue
		}
		teamName := os.Getenv("CHATBOT_TEAM")
		output, err := team.Get(w, teamName, team.Members)
		if err != nil {
			debugPrintf("[Team Error] getting team members %v\n", err.Error())
			continue
		}
		users = make([]string, len(output))
		for x, user := range output {
			users[x] = user.Username
		}
		if err := w.Proc.Kill(); err != nil {
			debugPrintf("[Team Error] killing process %v\n", err)
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
		print.Warningln("Missing CHATBOT_TEAM environment variable")
	}
}

func main() {}
