package main

import (
	"fmt"
	"io"

	"github.com/jamesbcook/chatbot-plugins/hibp"
	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	passwordURL      = "https://api.pwnedpasswords.com/"
	pwnedPasswordURL = "pwnedpassword/"
	rangePasswordURL = "range/"
)

var (
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
	return "/hibp-password"
}

//Help is what will show in the help menu
func (a activePlugin) Help() string {
	return "/hibp-password {passsword}"
}

//Get export method that satisfies an interface in the main program.
//This Get method will query the hibp password api.
func (a activePlugin) Get(input string) (string, error) {
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
func (a activePlugin) Send(subscription kbchat.SubscriptionMessage, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[HIBP-Password Error] sending message %v", err)
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
