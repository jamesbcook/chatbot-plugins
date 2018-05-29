package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	userAgent = "KeyBase Chatbot"
	urlFMT    = "https://www.reddit.com/r/%s/top/.json"
)

var (
	//CMD that keybase will use to execute this plugin
	CMD = "/reddit"
	//Help is what will show in the help menu
	Help         = "/reddit {subreddit}"
	areDebugging = false
)

type getting string

//Getter export symbol
var Getter getting

//Sender export symbol
var Sender getting

//Debugger export Symbol
var Debugger getting

//Kind of response from reddit
type Kind struct {
	Data `json:"data"`
}

//Data contains an array of children
type Data struct {
	Childrens []Children `json:"children"`
}

//Children contains data
type Children struct {
	Data InnerData `json:"data"`
}

//InnerData contains info about a post
type InnerData struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Permalink string `json:"permalink"`
}

func (g getting) Debug(set bool) {
	areDebugging = set
}

func debug(input string) {
	if areDebugging {
		fmt.Printf("[DEBUG] %s\n", input)
	}
}

//Get export method that satisfies an interface in the main program.
//This Get method will query reddit json api.
func (g getting) Get(input string) (string, error) {
	url := fmt.Sprintf(urlFMT, input)
	client := &http.Client{}
	debug(fmt.Sprintf("Creating GET request to %s", url))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("[Reddit Error] creating request %v", err)
	}
	req.Header.Set("User-Agent", userAgent)
	debug(fmt.Sprintf("Sending request %v", req))
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("[Reddit Error] sending request %v", err)
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	debug(fmt.Sprintf("Reading resp.Body"))
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", fmt.Errorf("[Reddit Error] reading from resp body %v", err)
	}
	debug(fmt.Sprintf("Unmarshalling json with length of %d", len(buf.Bytes())))
	k := &Kind{}
	if err := json.Unmarshal(buf.Bytes(), k); err != nil {
		return "", fmt.Errorf("[Reddit Error] unmarshalling response %v", err)
	}
	var numOfLinks int
	if len(k.Childrens) <= 10 {
		numOfLinks = len(k.Childrens)
	} else {
		numOfLinks = 10
	}
	msg := "Top Posts\n"
	for x := 0; x < numOfLinks; x++ {
		msg += fmt.Sprintf("Title: %-16s\n", k.Childrens[x].Data.Title)
	}
	debug(fmt.Sprintf("Message sending to user\n%s", msg))
	return msg, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (g getting) Send(msgID, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Reddit Error] in send request %v", err)
	}
	debug(fmt.Sprintf("Sending this message to messageID: %s\n%s", msgID, msg))
	if err := w.SendMessage(msgID, msg); err != nil {
		return w.Proc.Kill()
	}
	debug("Killing child process")
	return w.Proc.Kill()
}

func main() {}
