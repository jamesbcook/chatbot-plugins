package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	userAgent = "KeyBase Chatbot"
	urlFMT    = "https://api.coinmarketcap.com/v1/ticker/%s/?convert=USD"
)

var (
	//CMD that keybase will use to execute this plugin
	CMD = "/crypto"
	//Help is what will show in the help menu
	Help         = "/crypto {cryptocurrency}"
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

//CryptoOutput is the results from the api
type CryptoOutput struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Rank     string `json:"rank"`
	PriceUSD string `json:"price_usd"`
	PriceBTC string `json:"price_btc"`
}

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

//Get export method that satisfies an interface in the main program.
//This Get method will query the coinmarketcap api.
func (g getting) Get(input string) (string, error) {
	url := fmt.Sprintf(urlFMT, input)
	client := &http.Client{}
	debug(fmt.Sprintf("Creating GET request to %s", url))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("[Crypto Error] creating request %v", err)
	}
	req.Header.Set("User-Agent", userAgent)
	debug(fmt.Sprintf("Sending request %v", req))
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("[Crypto Error] sending request %v", err)
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	debug(fmt.Sprintf("Reading resp.Body"))
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", fmt.Errorf("[Crypto Error] reading from resp body %v", err)
	}
	crypto := []CryptoOutput{}
	debug(fmt.Sprintf("Unmarshalling json with length of %d", len(buf.Bytes())))
	if err := json.Unmarshal(buf.Bytes(), &crypto); err != nil {
		return "", fmt.Errorf("[Crypto Error] unmarshalling response %v", err)
	}
	output := fmt.Sprintf("Name: %-8s\tUSD: %-6s\tBTC: %-10s", crypto[0].Name, crypto[0].PriceUSD, crypto[0].PriceBTC)
	debug(fmt.Sprintf("Message sending to user\n%s", output))
	return output, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (g getting) Send(msgID, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Crypto Error] sending message %v", err)
	}
	debug(fmt.Sprintf("Sending this message to messageID: %s\n%s", msgID, msg))
	if err := w.SendMessage(msgID, msg); err != nil {
		return w.Proc.Kill()
	}
	debug("Killing child process")
	return w.Proc.Kill()
}

func main() {}
