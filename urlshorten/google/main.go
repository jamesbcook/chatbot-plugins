package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
)

const (
	url       = "https://www.googleapis.com/urlshortener/v1/url?key=%s"
	userAgent = "Keybase Chatbot"
)

var (
	debugPrintf func(format string, v ...interface{})
)

type activePlugin string

//AP for export
var AP activePlugin

type shortenAPI struct {
	Kind    string      `json:"kind"`
	ID      string      `json:"id"`
	LongURL string      `json:"longUrl"`
	Error   customError `json:"error"`
}

type customError struct {
	Errors []struct {
		Domain  string `json:"domain"`
		Reason  string `json:"reason"`
		Message string `json:"message"`
	} `json:"errors"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (a activePlugin) Debug(set bool, writer *io.Writer) {
	debugPrintf = print.Debugf(set, writer)
}

//CMD that keybase will use to execute this plugin
func (a activePlugin) CMD() string {
	return "/url-shorten"
}

//Help is what will show in the help menu
func (a activePlugin) Help() string {
	return "/url-shorten {url}"
}

//Get export method that satisfies an interface in the main program.
//This Get method will query the Google URL shortener api.
func (a activePlugin) Get(input string) (string, error) {
	debugPrintf("Setting up POST request\n")
	finalURL := fmt.Sprintf(url, os.Getenv("CHATBOT_URL_SHORTEN"))
	urlToShorten := fmt.Sprintf(`{"longUrl": "%s"}`, input)
	req, err := http.NewRequest("POST", finalURL, bytes.NewBufferString(urlToShorten))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	client := &http.Client{}
	debugPrintf("Sending post request for %s\n", urlToShorten)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	debugPrintf("Reading %d of response data\n", resp.ContentLength)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	shorten := &shortenAPI{}
	if err := json.Unmarshal(body, shorten); err != nil {
		return "", err
	}
	debugPrintf("ShortenAPI %v\n", shorten)
	if shorten.Error.Code != 0 {
		return "", fmt.Errorf("Non 0 error code")
	}
	return fmt.Sprintf("URL: %s", shorten.ID), nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (a activePlugin) Send(subscription kbchat.SubscriptionMessage, msg string) error {
	debugPrintf("Starting kbchat\n")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[URL Short Error] in send request %v", err)
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

func init() {
	if api := os.Getenv("CHATBOT_URL_SHORTEN"); api == "" {
		print.Warningln("Missing CHATBOT_URL_SHORTEN environment variable")
	}
}

func main() {}
