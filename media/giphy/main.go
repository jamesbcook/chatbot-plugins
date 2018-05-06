package main

import (
	"fmt"
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
	//CMD that keybase will use to execute this plugin
	CMD = "/giphy"
	//Help is what will show in the help menu
	Help = "/giphy {string}"
)

type getting string

//Getter export symbol
var Getter getting

//Sender export symbol
var Sender getting

func init() {
	rand.Seed(time.Now().UnixNano())
}

//giphy grabs top 100 Gifs from Giphy chooses a random one and downloads it
func giphy(query string) ([]byte, error) {

	giphy := libgiphy.NewGiphy(os.Getenv("CHATBOT_GIPHY"))

	dataSearch, err := giphy.GetSearch(query, giphySearchLimit, -1, "", "", false)
	if err != nil {
		return nil, fmt.Errorf("[Giphy Error] Giphy search error: %v", err)
	}
	returnLen := len(dataSearch.Data)
	if returnLen <= 0 {
		return nil, fmt.Errorf("[Giphy Error] No gifs found :(")
	}
	gifURL := dataSearch.Data[rand.Intn(returnLen)].Images.Downsized.Url

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
func (g getting) Get(input string) (string, error) {
	f, err := media.Setup(input, giphy)
	if err != nil {
		return "", fmt.Errorf("[Gihpy Error] in Get request %v", err)
	}
	return f, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will upload the results to the message ID that sent the request,
//once the file is uploaded it will delete the file.
func (g getting) Send(msgID, msg string) error {
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Giphy Error] in send request %v", err)
	}

	if _, err = os.Stat(msg); os.IsNotExist(err) {
		if err := w.SendMessage(msgID, "No Picture Available"); err != nil {
			return w.Proc.Kill()
		}
		return w.Proc.Kill()
	}
	if err := w.Upload(msgID, msg, "Chatbot-Giphy"); err != nil {
		return w.Proc.Kill()
	}
	return w.Proc.Kill()
}

func init() {
	if giphy := os.Getenv("CHATBOT_GIPHY"); giphy == "" {
		log.Println("Missing CHATBOT_GIPHY environment variable")
	}
}

func main() {}
