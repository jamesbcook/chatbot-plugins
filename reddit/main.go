package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	userAgent = "KeyBase Chatbot"
	urlFMT    = "https://www.reddit.com/r/%s/hot/.json"
)

var (
	areDebugging = false
	debugWriter  *io.Writer
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
	return "/reddit"
}

//Help is what will show in the help menu
func (a activePlugin) Help() string {
	return "/reddit {subreddit}"
}

//Get export method that satisfies an interface in the main program.
//This Get method will query reddit json api.
func (a activePlugin) Get(input string) (string, error) {
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
	if len(k.Data.Children) <= 10 {
		numOfLinks = len(k.Data.Children)
	} else {
		numOfLinks = 10
	}
	msg := "Top Posts\n"
	x := 0
	for _, child := range k.Data.Children {
		if x < numOfLinks && child.Data.Distinguished == nil {
			msg += fmt.Sprintf("Title: %-16s\n", child.Data.Title)
			x++
		} else if x == 10 {
			break
		}
	}
	debug(fmt.Sprintf("Message sending to user\n%s", msg))
	return msg, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (a activePlugin) Send(msgID, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Reddit Error] in send request %v", err)
	}
	debug(fmt.Sprintf("Sending this message to messageID: %s\n%s", msgID, msg))
	if err := w.SendMessage(msgID, msg); err != nil {
		if err := w.Proc.Kill(); err != nil {
			return err
		}
		return err
	}
	debug("Killing child process")
	return w.Proc.Kill()
}

func main() {}
