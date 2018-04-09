package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jamesbcook/chatbot-plugins/media"
	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	ctSize = 512
)

var (
	validTypes = []string{"image/jpeg", "image/jpg", "image/gif", "image/png"}
	//CMD that keybase will use to execute this plugin
	CMD = "/media"
	//Help is what will show in the help menu
	Help = "/media {url}"
)

type getting string

//Getter export symbol
var Getter getting

//Sender export symbol
var Sender getting

//url where we will download the file from
func url(query string) ([]byte, error) {
	resp, err := http.Get(query)
	if err != nil {
		return nil, fmt.Errorf("[Error] HTTP Get error %v", err)
	}
	defer resp.Body.Close()

	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[Error] Buffer read error %v", err)
	}
	header := make([]byte, ctSize)
	copy(header, buffer)
	if validContentType(header) {
		return buffer, nil
	}
	return nil, fmt.Errorf("[Error] Invalid ContentType")
}

func validContentType(buffer []byte) bool {
	for _, t := range validTypes {
		if t == http.DetectContentType(buffer) {
			return true
		}
	}
	return false
}

//Get export method that satisfies an interface in the main program.
//This Get method will query the raw url and return the results if it's a valid
//content type.
func (g getting) Get(input string) (string, error) {
	f, err := media.Setup(input, url)
	if err != nil {
		return "", err
	}
	return f, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will upload the results to the message ID that sent the request,
//once the file is uploaded it will delete the file.
func (g getting) Send(msgID, msg string) error {
	w, err := kbchat.Start("chat")
	if err != nil {
		return err
	}
	return w.Upload(msgID, msg, "Chatbot-Media")
}

func main() {}
