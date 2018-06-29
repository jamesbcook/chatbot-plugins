package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
	"golang.org/x/crypto/sha3"
)

const (
	resolution = "1200,800"
	userAgent  = "Keybase Chatbot"
)

type activePlugin string

//AP for export
var AP activePlugin

var (
	debugPrintf func(format string, v ...interface{})
	chrome      = &chromeData{}
	paths       = []string{
		"/usr/bin/chromium",
		"/usr/bin/chromium-browser",
		"/usr/bin/google-chrome-stable",
		"/usr/bin/google-chrome",
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"C:/Program Files (x86)/Google/Chrome/Application/chrome.exe",
	}
)

type chromeData struct {
	resolution     string
	timeout        int
	path           string
	userAgent      string
	screenshotPath string
}

func (a activePlugin) Debug(set bool, writer *io.Writer) {
	debugPrintf = print.Debugf(set, writer)
}

//CMD that keybase will use to execute this plugin
func (a activePlugin) CMD() string {
	return "/screenshot"
}

//Help is what will show in the help menu
func (a activePlugin) Help() string {
	return "/screenshot {url}"
}

func shaFileName(fileName string) string {
	digest := make([]byte, 32)
	sha3.ShakeSum256(digest, []byte(fileName))
	return hex.EncodeToString(digest)
}

func (c *chromeData) locateChrome() {
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		c.path = path
		return
	}

	log.Fatal("Could not find chrome")
}

//Get export method that satisfies an interface in the main program.
//This Get method will take a screen shot of the url using headless chrome
//and return the file path.
func (a activePlugin) Get(query string) (string, error) {
	tmpfn := filepath.Join("/tmp", shaFileName(query))
	basicArguments := []string{
		"--headless", "--disable-gpu", "--hide-scrollbars",
		"--disable-crash-reporter",
		"--user-agent=" + chrome.userAgent,
		"--window-size=" + chrome.resolution, "--screenshot=" + tmpfn,
	}
	basicArguments = append(basicArguments, query)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(90*time.Second))
	defer cancel()
	debugPrintf("Executing %s with arguments %v\n", chrome.path, basicArguments)
	cmd := exec.CommandContext(ctx, chrome.path, basicArguments...)
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("[Screenshot Error] starting the chrome command %v", err)
	}
	if err := cmd.Wait(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("[Screenshot Error] Context time out")
		}
	}
	debugPrintf("Sending filename %s to user\n", tmpfn)
	return tmpfn, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will upload the results to the message ID that sent the request,
//once the file is uploaded it will delete the file.
func (a activePlugin) Send(subscription kbchat.SubscriptionMessage, msg string) error {
	debugPrintf("Starting kbchat\n")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Screenshot Error] in send request %v", err)
	}
	debugPrintf("Checking if file exists\n")
	if _, err = os.Stat(msg); os.IsNotExist(err) {
		debugPrintf("File didn't exist\n")
		if err := w.SendMessage(subscription.Conversation.ID, "No Picture Available"); err != nil {
			return w.Proc.Kill()
		}
		return w.Proc.Kill()
	}
	debugPrintf("Uploading %s to msgID: %s\n", msg, subscription.Conversation.ID)
	if err := w.Upload(subscription.Conversation.ID, msg, "Chatbot-Media"); err != nil {
		return w.Proc.Kill()
	}
	debugPrintf("Killing child process\n")
	return w.Proc.Kill()
}

func init() {
	chrome.locateChrome()
	chrome.resolution = resolution
	chrome.userAgent = userAgent
}

func main() {}
