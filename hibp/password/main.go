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
	Help = "/hibp-password {passsword}"
)

type getting string

//Getter export symbol
var Getter getting

//Sender export symbol
var Sender getting

//Get export method that satisfies an interface in the main program.
//This Get method will query the hibp password api.
func (g getting) Get(input string) (string, error) {
	res, err := hibp.Get(input, pwnedPassword)
	if err != nil {
		return "", fmt.Errorf("[HIBP-Password Error] There was an error with your request")
	}
	msg := fmt.Sprintf("Password has been seen %s times", string(res))
	return msg, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (g getting) Send(msgID, msg string) error {
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[HIBP-Password Error] sending message %v", err)
	}
	return w.SendMessage(msgID, msg)
}

//pwnedPassword returns the number of times this password has been seen in all breaches
func pwnedPassword(password string) string {
	fullURL := passwordURL + pwnedPasswordURL + password
	return fullURL
}

//pwnedPasswordRange returns a list of hashes that match the first 5 chars that waere sent
func pwnedPasswordRange(password string) string {
	fullURL := passwordURL + rangePasswordURL + password
	return fullURL
}

func main() {}
