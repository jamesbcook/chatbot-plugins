package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	url       = "https://www.googleapis.com/urlshortener/v1/url?key=%s"
	userAgent = "Keybase Chatbot"
)

var (
	//CMD that keybase will use to execute this plugin
	CMD = "/url-shorten"
	//Help is what will show in the help menu
	Help         = "/url-shorten {url}"
	areDebugging = false
	debugWriter  *io.Writer
)

type getting string

//Getter export symbol
var Getter getting

//Sender export symbol
var Sender getting

//Debugger export Symbol
var Debugger getting

func (g getting) Debug(set bool, writer *io.Writer) {
	areDebugging = set
	debugWriter = writer
}

func debug(input string) {
	if areDebugging {
		output := fmt.Sprintf("[DEBUG] %s\n", input)
		(*debugWriter).Write([]byte(output))
	}
}

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

//Get export method that satisfies an interface in the main program.
//This Get method will query the Google URL shortener api.
func (g getting) Get(input string) (string, error) {
	debug("Setting up POST request")
	finalURL := fmt.Sprintf(url, os.Getenv("CHATBOT_URL_SHORTEN"))
	urlToShorten := fmt.Sprintf(`{"longUrl": "%s"}`, input)
	req, err := http.NewRequest("POST", finalURL, bytes.NewBufferString(urlToShorten))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	client := &http.Client{}
	debug(fmt.Sprintf("Sending post request for %s", urlToShorten))
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	debug(fmt.Sprintf("Reading %d of response data", resp.ContentLength))
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	shorten := &shortenAPI{}
	if err := json.Unmarshal(body, shorten); err != nil {
		return "", err
	}
	debug(fmt.Sprintf("ShortenAPI %v", shorten))
	if shorten.Error.Code != 0 {
		return "", fmt.Errorf("Non 0 error code")
	}
	return fmt.Sprintf("URL: %s", shorten.ID), nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (g getting) Send(msgID, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[URL Short Error] in send request %v", err)
	}
	debug(fmt.Sprintf("Sending this message to messageID: %s\n%s", msgID, msg))
	if err := w.SendMessage(msgID, msg); err != nil {
		return w.Proc.Kill()
	}
	debug("Killing child process")
	return w.Proc.Kill()
}

func init() {
	if api := os.Getenv("CHATBOT_URL_SHORTEN"); api == "" {
		log.Println("Missing CHATBOT_URL_SHORTEN environment variable")
	}
}

func main() {}
