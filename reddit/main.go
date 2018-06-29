package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
)

const (
	userAgent = "KeyBase Chatbot"
	baseURL   = "https://reddit.com/r/"
	jsonURL   = baseURL + "%s/hot/.json"
)

var (
	debugPrintf func(format string, v ...interface{})
)

type activePlugin string

//AP for export
var AP activePlugin

//Kind of response from reddit
type Kind struct {
	Data struct {
		Children []struct {
			Data struct {
				Author        string  `json:"author"`
				Title         string  `json:"title"`
				Permalink     string  `json:"permalink"`
				Distinguished *string `json:"distinguished"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func (a activePlugin) Debug(set bool, writer *io.Writer) {
	debugPrintf = print.Debugf(set, writer)
}

//CMD that keybase will use to execute this plugin
func (a activePlugin) CMD() string {
	return "/reddit"
}

//Help is what will show in the help menu
func (a activePlugin) Help() string {
	return "/reddit {subreddit}"
}

//Get export method that satisfies an interface in the main program.
//This Get method will query reddit json api.
func (a activePlugin) Get(input string) (string, error) {
	url := fmt.Sprintf(jsonURL, input)
	client := &http.Client{}
	debugPrintf("Creating GET request to %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("[Reddit Error] creating request %v", err)
	}
	req.Header.Set("User-Agent", userAgent)
	debugPrintf("Sending request %v\n", req)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("[Reddit Error] sending request %v", err)
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	debugPrintf("Reading resp.Body\n")
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", fmt.Errorf("[Reddit Error] reading from resp body %v", err)
	}
	debugPrintf("Unmarshalling json with length of %d\n", len(buf.Bytes()))
	k := &Kind{}
	if err := json.Unmarshal(buf.Bytes(), k); err != nil {
		return "", fmt.Errorf("[Reddit Error] unmarshalling response %v", err)
	}
	var numOfLinks int
	if len(k.Data.Children) <= 10 {
		numOfLinks = len(k.Data.Children)
	} else {
		numOfLinks = 10
	}
	if len(k.Data.Children) == 0 {
		return "", fmt.Errorf("Subreddit %s not found", input)
	}
	msg := fmt.Sprintf("Top Posts for %s%s\n", baseURL, input)
	x := 0
	for _, child := range k.Data.Children {
		if x < numOfLinks && child.Data.Distinguished == nil {
			msg += fmt.Sprintf("Title: %-16s\n", child.Data.Title)
			x++
		} else if x == 10 {
			break
		}
	}
	debugPrintf("Message sending to user\n%s\n", msg)
	return msg, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (a activePlugin) Send(subscription kbchat.SubscriptionMessage, msg string) error {
	debugPrintf("Starting kbchat\n")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Reddit Error] in send request %v", err)
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

func main() {}
