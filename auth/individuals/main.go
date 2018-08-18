package main

import (
	"io"
	"os"
	"strings"
	"time"

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

//Debug output
func (b backgroundPlugin) Debug(set bool, writer *io.Writer) {
	debugPrintf = print.Debugf(set, writer)
}

//Name that keybase will use for background plugins
func (b backgroundPlugin) Name() string {
	return "auth"
}

//Start the process of gathering uses based on the user names seperated by ","
//in the env var.
func (a authenticator) Start() {
	for {
		users = []string{}
		userEnv := os.Getenv("CHATBOT_USERS")
		for _, user := range strings.Split(userEnv, ",") {
			users = append(users, user)
		}
		time.Sleep(5 * time.Minute)
	}
}

//Validate that a user is  allowed to send commands to the bot
func (a authenticator) Validate(user string) bool {
	for _, u := range users {
		if user == u {
			return true
		}
	}
	return false
}

func init() {
	userEnv := os.Getenv("CHATBOT_USERS")
	if userEnv == "" {
		print.Warningln("Missing CHATBOT_USERS environment variable")
	}
}

func main() {}
