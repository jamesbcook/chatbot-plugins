package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
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

//Debug output
func (b backgroundPlugin) Debug(set bool, writer *io.Writer) {
	areDebugging = set
	debugWriter = writer
}

func debug(input string) {
	if areDebugging && *debugWriter != nil {
		output := fmt.Sprintf("[DEBUG] %s\n", input)
		(*debugWriter).Write([]byte(output))
	}
}

//Name that keybase will use for background plugins
func (b backgroundPlugin) Name() string {
	return "auth"
}

//Start the process of gathering uses based on the user names seperated by ","
//in the env var.
func (a authenticator) Start() {
	for {
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
		log.Println("Missing CHATBOT_USERS environment variable")
	}
}

func main() {}
