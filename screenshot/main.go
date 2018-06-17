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
	areDebugging = false
	debugWriter  *io.Writer
	chrome       = &chromeData{}
	paths        = []string{
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
	debug(fmt.Sprintf("Executing %s with arguments %v", chrome.path, basicArguments))
	cmd := exec.CommandContext(ctx, chrome.path, basicArguments...)
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("[Screenshot Error] starting the chrome command %v", err)
	}
	if err := cmd.Wait(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("[Screenshot Error] Context time out")
		}
	}
	debug(fmt.Sprintf("Sending filename %s to user", tmpfn))
	return tmpfn, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will upload the results to the message ID that sent the request,
//once the file is uploaded it will delete the file.
func (a activePlugin) Send(msgID, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Screenshot Error] in send request %v", err)
	}
	debug("Checking if file exists")
	if _, err = os.Stat(msg); os.IsNotExist(err) {
		debug("File didn't exist")
		if err := w.SendMessage(msgID, "No Picture Available"); err != nil {
			return w.Proc.Kill()
		}
		return w.Proc.Kill()
	}
	debug(fmt.Sprintf("Uploading %s to msgID: %s", msg, msgID))
	if err := w.Upload(msgID, msg, "Chatbot-Media"); err != nil {
		return w.Proc.Kill()
	}
	debug("Killing child process")
	return w.Proc.Kill()
}

func init() {
	chrome.locateChrome()
	chrome.resolution = resolution
	chrome.userAgent = userAgent
}

func main() {}
