package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jamesbcook/chat-bot/kbchat"
)

const (
	userAgent = "KeyBase Chatbot"
	urlFMT    = "https://api.coinmarketcap.com/v1/ticker/%s/?convert=USD"
)

var (
	//CMD that keybase will use to execute this plugin
	CMD = "/crypto"
	//Help is what will show in the help menu
	Help = "/crypto {cryptocurrency}"
)

type getting string

//Getter export symbol
var Getter getting

//Sender export symbol
var Sender getting

//CryptoOutput is the results from the api
type CryptoOutput struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Rank     string `json:"rank"`
	PriceUSD string `json:"price_usd"`
	PriceBTC string `json:"price_btc"`
}

//Get export method that satisfies an interface in the main program.
//This Get method will query the coinmarketcap api.
func (g getting) Get(input string) (string, error) {
	url := fmt.Sprintf(urlFMT, input)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("Error creating request %v", err)
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Errror sending request %v", err)
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading from resp body %v", err)
	}
	crypto := []CryptoOutput{}
	if err := json.Unmarshal(buf.Bytes(), &crypto); err != nil {
		return "", fmt.Errorf("Error unmarshaling response %v", err)
	}
	output := fmt.Sprintf("Name: %-8s\tUSD: %-6s\tBTC: %-10s", crypto[0].Name, crypto[0].PriceUSD, crypto[0].PriceBTC)
	return output, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (g getting) Send(msgID, msg string) error {
	w, err := kbchat.Start("chat")
	if err != nil {
		return err
	}
	return w.SendMessage(msgID, msg)
}

func main() {}
