package main

import (
	"encoding/json"
	"fmt"

	"github.com/jamesbcook/chatbot-plugins/hibp"
	"github.com/jamesbcook/chatbot/kbchat"
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
	//CMD that keybase will use to execute this plugin
	CMD = "/hibp-email"
	//Help is what will show in the help menu
	Help         = "/hibp-email {email}"
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
//This Get method will query the hibp account api.
func (g getting) Get(input string) (string, error) {
	debug(fmt.Sprintf("Sending %s to HIBP Breach API", input))
	breachRes, err := hibp.Get(input, allBreachesForAccount)
	if err != nil {
		return "", fmt.Errorf("[HIBP-Account Error] There was an error with your beaches request")
	}
	debug(fmt.Sprintf("Sending %s to HIBP Pastes API", input))
	pasteRes, err := hibp.Get(input, allPastesForAccount)
	if err != nil {
		return "", fmt.Errorf("[HIBP-Account Error] There was an error with your pastes request")
	}
	breaches := []breachedAccount{}
	debug(fmt.Sprintf("Unmarshaling json for breaches"))
	if err := json.Unmarshal(breachRes, &breaches); err != nil {
		return "", fmt.Errorf("[HIBP-Account Error] There was an error unmarshaling your request")
	}
	pastes := []pasteAccount{}
	debug(fmt.Sprintf("Unmarshaling json for pastes"))
	if err := json.Unmarshal(pasteRes, &pastes); err != nil {
		return "", fmt.Errorf("[HIBP-Account Error] There was an error unmarshaling your request")
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
	debug(fmt.Sprintf("Returning the following message to user\n%s", msg))
	return msg
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (g getting) Send(msgID, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[HIBP-Account Error] sending message %v", err)
	}
	debug(fmt.Sprintf("Sending this message to messageID: %s\n%s", msgID, msg))
	if err := w.SendMessage(msgID, msg); err != nil {
		return w.Proc.Kill()
	}
	debug("Killing child process")
	return w.Proc.Kill()
}

//AllBreachesForAccount returns an array of breaches the account has been seen in
func allBreachesForAccount(account string) string {
	fullURL := accountURL + breachedAccountURL + account
	return fullURL
}

//AllPastesForAccount retuns an array of pastes the account has been seen in
func allPastesForAccount(account string) string {
	fullURL := accountURL + pasteAccountURL + account
	return fullURL
}

func main() {}
