package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jamesbcook/chatbot-plugins/media"
	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
)

const (
	ctSize = 512
)

var (
	validTypes  = []string{"image/jpeg", "image/jpg", "image/gif", "image/png"}
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
	return "/media"
}

//Help is what will show in the help menu
func (a activePlugin) Help() string {
	return "/media {url}"
}

//url where we will download the file from
func url(query string) ([]byte, error) {
	debugPrintf("Sending GET request to %s\n", query)
	resp, err := http.Get(query)
	if err != nil {
		return nil, fmt.Errorf("[Media Error] HTTP Get error %v", err)
	}
	defer resp.Body.Close()

	debugPrintf("Reading Resp Body\n")
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[Media Error] Buffer read error %v", err)
	}
	header := make([]byte, ctSize)
	copy(header, buffer)
	debugPrintf("Checking if file is a valid content type\n")
	if validContentType(header) {
		return buffer, nil
	}
	return nil, fmt.Errorf("[Media Error] Invalid ContentType")
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
func (a activePlugin) Get(input string) (string, error) {
	f, err := media.Setup(input, url)
	if err != nil {
		return "", fmt.Errorf("[Media Error] in Get request %v", err)
	}
	debugPrintf("Sending filename %s to user\n", f)
	return f, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will upload the results to the message ID that sent the request,
//once the file is uploaded it will delete the file.
func (a activePlugin) Send(subscription kbchat.SubscriptionMessage, msg string) error {
	debugPrintf("Starting kbchat\n")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Media Error] in send request %v", err)
	}
	debugPrintf("Checking if file exists\n")
	if _, err = os.Stat(msg); os.IsNotExist(err) {
		debugPrintf("File didn't exist\n")
		if err := w.SendMessage(subscription.Conversation.ID, msg); err != nil {
			return w.Proc.Kill()
		}
		if err := w.Proc.Kill(); err != nil {
			return err
		}
		return err
	}
	debugPrintf("Uploading %s to msgID: %s\n", msg, subscription.Conversation.ID)
	if err := w.Upload(subscription.Conversation.ID, msg, "Chatbot-Media"); err != nil {
		if err := w.Proc.Kill(); err != nil {
			return err
		}
		return err
	}
	debugPrintf("Killing child process\n")
	return w.Proc.Kill()
}

func main() {}
