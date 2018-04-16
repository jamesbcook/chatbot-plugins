package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/jamesbcook/chatbot/kbchat"
	"golang.org/x/crypto/sha3"
)

type getting string

//Getter export symbol
var Getter getting

//Sender export symbol
var Sender getting

type chromeData struct {
	resolution     string
	timeout        int
	path           string
	userAgent      string
	screenshotPath string
}

var (
	//CMD that keybase will use to execute this plugin
	CMD = "/screenshot"
	//Help is what will show in the help menu
	Help = "/screenshot {url}"
)

const (
	resolution = "1200,800"
	userAgent  = "Keybase Chatbot/0.5.0"
)

var (
	chrome = &chromeData{}
	paths  = []string{
		"/usr/bin/chromium",
		"/usr/bin/google-chrome-stable",
		"/usr/bin/google-chrome",
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"C:/Program Files (x86)/Google/Chrome/Application/chrome.exe",
	}
)

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
func (g getting) Get(query string) (string, error) {
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

	cmd := exec.CommandContext(ctx, chrome.path, basicArguments...)
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("[Screenshot Error] starting the chrome command %v", err)
	}
	if err := cmd.Wait(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("[Screenshot Error] Context time out")
		}
	} else {
		return "", fmt.Errorf("[Screenshot Error] waiting on command %v", err)
	}
	return tmpfn, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will upload the results to the message ID that sent the request,
//once the file is uploaded it will delete the file.
func (g getting) Send(msgID, msg string) error {
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Screenshot Error] in send request %v", err)
	}
	return w.Upload(msgID, msg, "Chatbot-Media")
}

func init() {
	chrome.locateChrome()
	chrome.resolution = resolution
	chrome.userAgent = userAgent
}

func main() {}
