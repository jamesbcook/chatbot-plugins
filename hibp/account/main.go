package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/jamesbcook/chatbot-plugins/hibp"
	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
)

const (
	accountURL         = "https://haveibeenpwned.com/api/v2/"
	breachedAccountURL = "breachedaccount/"
	pasteAccountURL    = "pasteaccount/"
)

//BreachedAccount results from hibp api
type breachedAccount struct {
	Title        string
	Name         string
	Domain       string
	BreachDate   string
	AddedDate    string
	ModifiedDate string
	PwnCount     int
	Description  string
	DataClasses  []string
	IsVerified   bool
	IsFabricated bool
	IsSensitive  bool
	IsActive     bool
	IsRetired    bool
	IsSpamList   bool
	LogoType     string
}

//PasteAccount results from hibp api
type pasteAccount struct {
	Source     string
	ID         string
	Title      string
	Date       string
	EmailCount int
}

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
	return "/hibp-email"
}

//Help is what will show in the help menu
func (a activePlugin) Help() string {
	return "/hibp-email {email}"
}

//Get export method that satisfies an interface in the main program.
//This Get method will query the hibp account api.
func (a activePlugin) Get(input string) (string, error) {
	debugPrintf("Sending %s to HIBP Breach API\n", input)
	breachRes, err := hibp.Get(input, allBreachesForAccount)
	if err != nil {
		return "", fmt.Errorf("[HIBP-Account Error] There was an error with your beaches request")
	}
	debugPrintf("Sending %s to HIBP Pastes API\n", input)
	pasteRes, err := hibp.Get(input, allPastesForAccount)
	if err != nil {
		return "", fmt.Errorf("[HIBP-Account Error] There was an error with your pastes request")
	}
	breaches := []breachedAccount{}
	debugPrintf("Unmarshalling json for breaches\n")
	if err := json.Unmarshal(breachRes, &breaches); err != nil {
		return "", fmt.Errorf("[HIBP-Account Error] There was an error unmarshalling your request")
	}
	pastes := []pasteAccount{}
	debugPrintf("Unmarshalling json for pastes\n")
	if err := json.Unmarshal(pasteRes, &pastes); err != nil {
		return "", fmt.Errorf("[HIBP-Account Error] There was an error unmarshalling your request")
	}
	return formatOutput(breaches, pastes), nil
}

func formatOutput(breaches []breachedAccount, pastes []pasteAccount) string {
	msg := "Account has been seen in the following breaches\n```"
	for _, breach := range breaches {
		msg += fmt.Sprintf("Name %s\nWhat Leaked: ", breach.Name)
		for _, dataClass := range breach.DataClasses {
			msg += fmt.Sprintf("%s ", dataClass)
		}
		msg += "\n"
	}
	msg += "```\nAccount has been seen in the following pastes\n```"
	for _, paste := range pastes {
		msg += fmt.Sprintf("Name %s ID: %s Source %s\n", paste.Title, paste.ID, paste.Source)
	}
	msg += "```"
	debugPrintf("Returning the following message to user\n%s\n", msg)
	return msg
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (a activePlugin) Send(subscription kbchat.SubscriptionMessage, msg string) error {
	debugPrintf("Starting kbchat\n")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[HIBP-Account Error] sending message %v", err)
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

//AllBreachesForAccount returns an array of breaches the account has been seen in
func allBreachesForAccount(account string) string {
	fullURL := accountURL + breachedAccountURL + account
	return fullURL
}

//AllPastesForAccount returns an array of pastes the account has been seen in
func allPastesForAccount(account string) string {
	fullURL := accountURL + pasteAccountURL + account
	return fullURL
}

func main() {}
