package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/chatbot/kbchat/team"
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

//Start the process of gathering uses based on the team name in the env var.
func Start(writer io.Writer) {
	for {
		w, err := kbchat.Start("team")
		if err != nil {
			errorWriter(writer, fmt.Errorf("[Team Error] getting team api %v", err.Error()))
			continue
		}
		teamName := os.Getenv("CHATBOT_TEAM")
		output, err := team.Get(w, teamName, team.Members)
		if err != nil {
			errorWriter(writer, fmt.Errorf("[Team Error] getting team members %v", err.Error()))
			continue
		}
		users = make([]string, len(output))
		for x, user := range output {
			users[x] = user.Username
		}
		time.Sleep(5 * time.Minute)
	}
}

//Validate that a user is part of a team and allowed to send commands to the bot
func Validate(user string) bool {
	for _, u := range users {
		if user == u {
			return true
		}
	}
	return false
}

func main() {}
