package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jamesbcook/chatbot-plugins/virustotal"
	"github.com/jamesbcook/chatbot/kbchat"
)

type getting string

//Getter export symbol
var Getter getting

//Sender export symbol
var Sender getting

//Debugger export Symbol
var Debugger getting

var (
	//CMD that keybase will use to execute this plugin
	CMD = "/virustotal"
	//Help is what will show in the help menu
	Help         = "/virustotal {sha256 of file}"
	areDebugging = false
)

func (g getting) Debug(set bool) {
	areDebugging = set
}

func debug(input string) {
	if areDebugging {
		fmt.Printf("[DEBUG] %s\n", input)
	}
}

func getURL() string {
	return fmt.Sprintf("%s/%s", virustotal.BaseURL, "report?apikey=%s&resource=%s")
}

//Get export method that satisfies an interface in the main program.
//This Get method will take a query virustotal with the given input
//and return the results of that file.
func (g getting) Get(input string) (string, error) {
	vt := &virustotal.Response{}
	api := os.Getenv("CHATBOT_VIRUSTOTAL")
	query := fmt.Sprintf(getURL(), api, input)
	debug(fmt.Sprintf("Sending GET request to %s", query))
	resp, err := http.Get(query)
	if err != nil {
		return "", fmt.Errorf("[VirusTotal Error] in get request")
	}
	defer resp.Body.Close()
	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("[VirusTotal Error] reading body")
	}
	debug(fmt.Sprintf("Unmarshaling json with length of %d", len(out)))
	if err := json.Unmarshal(out, vt); err != nil {
		return "", fmt.Errorf("[VirusTotal Error] unmarshal json")
	}
	output := "VirusTotal Detection Results\n"
	output += fmt.Sprintf("Total Detected %d\n", vt.Positives)
	if vt.Positives > 0 {
		output += "```\n"
		for scanner, scan := range vt.Scans {
			if scan.Detected {
				output += fmt.Sprintf("%-10s %s\n", scanner, scan.Result)
			}
		}
		output += "```"
	}
	debug(fmt.Sprintf("Sending the fowlling to user\n%s", output))
	return output, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will respond with the results to the message ID that sent the request.
func (g getting) Send(msgID, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[VirusTotal Error] in send request %v", err)
	}
	debug(fmt.Sprintf("Sending this message to messageID: %s\n%s", msgID, msg))
	if err := w.SendMessage(msgID, msg); err != nil {
		return w.Proc.Kill()
	}
	debug("Killing child process")
	return w.Proc.Kill()
}

func init() {
	if api := os.Getenv("CHATBOT_VIRUSTOTAL"); api == "" {
		log.Println("Missing CHATBOT_VIRUSTOTAL environment variable")
	}
}

func main() {}
