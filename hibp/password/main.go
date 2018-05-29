package main

import (
	"fmt"

	"github.com/jamesbcook/chatbot-plugins/hibp"
	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	passwordURL      = "https://api.pwnedpasswords.com/"
	pwnedPasswordURL = "pwnedpassword/"
	rangePasswordURL = "range/"
)

var (
	//CMD that keybase will use to execute this plugin
	CMD = "/hibp-password"
	//Help is what will show in the help menu
	Help         = "/hibp-password {passsword}"
	areDebugging = false
)

type getting string

//Getter export symbol
var Getter getting

//Sender export symbol
var Sender getting

//Debugger export Symbol
var Debugger getting

func (g getting) Debug(set bool) {
	areDebugging = set
}

func debug(input string) {
	if areDebugging {
		fmt.Printf("[DEBUG] %s\n", input)
	}
}

//Get export method that satisfies an interface in the main program.
//This Get method will query the hibp password api.
func (g getting) Get(input string) (string, error) {
	debug(fmt.Sprintf("Sending %s to HIBP Password API", input))
	res, err := hibp.Get(input, pwnedPassword)
	if err != nil {
		return "", fmt.Errorf("[HIBP-Password Error] There was an error with your request")
	}
	msg := fmt.Sprintf("Password has been seen %s times", string(res))
	debug(fmt.Sprintf("Returning the following message to user\n%s", msg))
	return msg, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (g getting) Send(msgID, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[HIBP-Password Error] sending message %v", err)
	}
	debug(fmt.Sprintf("Sending this message to messageID: %s\n%s", msgID, msg))
	if err := w.SendMessage(msgID, msg); err != nil {
		return w.Proc.Kill()
	}
	debug("Killing child process")
	return w.Proc.Kill()
}

//pwnedPassword returns the number of times this password has been seen in all breaches
func pwnedPassword(password string) string {
	fullURL := passwordURL + pwnedPasswordURL + password
	return fullURL
}

//pwnedPasswordRange returns a list of hashes that match the first 5 chars that were sent
func pwnedPasswordRange(password string) string {
	fullURL := passwordURL + rangePasswordURL + password
	return fullURL
}

func main() {}
