package main

import (
	"fmt"
	"io"

	"github.com/jamesbcook/chatbot-plugins/hibp"
	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
)

const (
	passwordURL      = "https://api.pwnedpasswords.com/"
	pwnedPasswordURL = "pwnedpassword/"
	rangePasswordURL = "range/"
)

var (
	debugPrintf func(format string, v ...interface{})
)

type activePlugin string

//AP for export
var AP activePlugin

func (a activePlugin) Debug(set bool, writer *io.Writer) {
	debugPrintf = print.Debugf(set, writer)
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
	debugPrintf("Sending %s to HIBP Password API\n", input)
	res, err := hibp.Get(input, pwnedPassword)
	if err != nil {
		return "", fmt.Errorf("[HIBP-Password Error] There was an error with your request")
	}
	msg := fmt.Sprintf("Password has been seen %s times", string(res))
	debugPrintf("Returning the following message to user\n%s\n", msg)
	return msg, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (a activePlugin) Send(subscription kbchat.SubscriptionMessage, msg string) error {
	debugPrintf("Starting kbchat\n")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[HIBP-Password Error] sending message %v", err)
	}
	debugPrintf("Sending this message to messageID: %s\n%s\n", subscription.Conversation.ID, msg)
	if err := w.SendMessage(subscription.Conversation.ID, msg); err != nil {
		if err := w.Proc.Kill(); err != nil {
			return err
		}
		return err
	}
	debugPrintf("Killing child process\n")
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
