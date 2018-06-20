package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/jamesbcook/chatbot-plugins/media"
	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/sanzaru/go-giphy"
)

const (
	giphySearchLimit = 100
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
	if areDebugging && *debugWriter != nil {
		output := fmt.Sprintf("[DEBUG] %s\n", input)
		(*debugWriter).Write([]byte(output))
	}
}

//CMD that keybase will use to execute this plugin
func (a activePlugin) CMD() string {
	return "/giphy"
}

//Help is what will show in the help menu
func (a activePlugin) Help() string {
	return "/giphy {string}"
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

//giphy grabs top 100 Gifs from Giphy chooses a random one and downloads it
func giphy(query string) ([]byte, error) {

	giphy := libgiphy.NewGiphy(os.Getenv("CHATBOT_GIPHY"))
	debug(fmt.Sprintf("Looking for random GIF for %s", query))
	dataSearch, err := giphy.GetSearch(query, giphySearchLimit, -1, "", "", false)
	if err != nil {
		return nil, fmt.Errorf("[Giphy Error] Giphy search error: %v", err)
	}
	debug(fmt.Sprintf("Found %d Gifs", len(dataSearch.Data)))
	returnLen := len(dataSearch.Data)
	if returnLen <= 0 {
		return nil, fmt.Errorf("[Giphy Error] No gifs found :(")
	}
	gifURL := dataSearch.Data[rand.Intn(returnLen)].Images.Downsized.Url

	debug(fmt.Sprintf("Sending GET request to %s", gifURL))
	// Get the data
	resp, err := http.Get(gifURL)
	if err != nil {
		return nil, fmt.Errorf("[Giphy Error] Unable to retrieve gif %v", err)
	}
	defer resp.Body.Close()
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[Giphy Error] Buffer read error %v", err)
	}
	return buffer, nil
}

//Get export method that satisfies an interface in the main program.
//This Get method will query the giphy api.
func (a activePlugin) Get(input string) (string, error) {
	f, err := media.Setup(input, giphy)
	if err != nil {
		return "", fmt.Errorf("[Giphy Error] in Get request %v", err)
	}
	debug(fmt.Sprintf("Sending filename %s to user", f))
	return f, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will upload the results to the message ID that sent the request,
//once the file is uploaded it will delete the file.
func (a activePlugin) Send(subscription kbchat.SubscriptionMessage, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Giphy Error] in send request %v", err)
	}

	debug("Checking if file exists")
	if _, err = os.Stat(msg); os.IsNotExist(err) {
		debug("File didn't exist")
		if err := w.SendMessage(subscription.Conversation.ID, msg); err != nil {
			if err := w.Proc.Kill(); err != nil {
				return err
			}
			return err
		}
		return w.Proc.Kill()
	}
	debug(fmt.Sprintf("Uploading %s to msgID: %s", msg, subscription.Conversation.ID))
	if err := w.Upload(subscription.Conversation.ID, msg, "Chatbot-Giphy"); err != nil {
		if err := w.Proc.Kill(); err != nil {
			return err
		}
		return err
	}
	debug("Killing child process")
	return w.Proc.Kill()
}

func init() {
	if giphy := os.Getenv("CHATBOT_GIPHY"); giphy == "" {
		log.Println("Missing CHATBOT_GIPHY environment variable")
	}
}

func main() {}
