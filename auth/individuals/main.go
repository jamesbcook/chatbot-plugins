package main

import (
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var (
	//Name that keybase will use for background plugins
	Name = "auth"
)

var (
	users = []string{}
)

func errorWriter(writer io.Writer, err error) {
	output := []byte(err.Error())
	output = append(output, '\n')
	writer.Write(output)
}

//Start the process of gathering uses based on the user names seperated by ","
//in the env var.
func Start(writer io.Writer) {
	for {
		userEnv := os.Getenv("CHATBOT_USERS")
		for _, user := range strings.Split(userEnv, ",") {
			users = append(users, user)
		}
		time.Sleep(5 * time.Minute)
	}
}

//Validate that a user is  allowed to send commands to the bot
func Validate(user string) bool {
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
