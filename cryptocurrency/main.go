package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
)

const (
	userAgent  = "KeyBase Chatbot"
	listingURL = "https://api.coinmarketcap.com/v2/listings/"
	priceURL   = "https://api.coinmarketcap.com/v2/ticker/%d/?convert=BTC"
)

type activePlugin string

//AP for export
var AP activePlugin

var (
	debugPrintf func(format string, v ...interface{})
	nameMap     = make(map[string]int)
	symbolMap   = make(map[string]int)
)

//Listing results
type listing struct {
	Data []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Symbol      string `json:"symbol"`
		WebsiteSlug string `json:"website_slug"`
	} `json:"data"`
	MetaData struct {
		TimeStamp           int64       `json:"timestamp"`
		NumCryptoCurrencies int         `json:"num_cryptocurrencies"`
		Error               interface{} `json:"error"`
	} `json:"metadata"`
}

//TickerSpecific for a currency
type tickerSpecific struct {
	Data struct {
		ID                int     `json:"id"`
		Name              string  `json:"name"`
		Symbol            string  `json:"symbol"`
		WebsiteSlug       string  `json:"website_slug"`
		Rank              int     `json:"rank"`
		CirculatingSupply float64 `json:"circulating_supply"`
		TotalSupply       float64 `json:"total_supply"`
		MaxSupply         float64 `json:"max_supply"`
		Quotes            struct {
			USD struct {
				Price            float64 `json:"price"`
				Volume24H        float64 `json:"volume_24h"`
				MarketCap        float64 `json:"market_cap"`
				PercentChange1H  float64 `json:"percent_change_1h"`
				PercentChange24H float64 `json:"percent_change_24h"`
				PercentChange7D  float64 `json:"percent_change_7d"`
			} `json:"USD"`
			BTC struct {
				Price            float64 `json:"price"`
				Volume24H        float64 `json:"volume_24h"`
				MarketCap        float64 `json:"market_cap"`
				PercentChange1H  float64 `json:"percent_change_1h"`
				PercentChange24H float64 `json:"percent_change_24h"`
				PercentChange7D  float64 `json:"percent_change_7d"`
			} `json:"BTC"`
		} `json:"quotes"`
		LastUpdated int64 `json:"last_updated"`
	} `json:"data"`
	Metadata struct {
		Timestamp int64       `json:"timestamp"`
		Error     interface{} `json:"error"`
	} `json:"metadata"`
}

func (a activePlugin) Debug(set bool, writer *io.Writer) {
	debugPrintf = print.Debugf(set, writer)
}

//CMD that keybase will use to execute this plugin
func (a activePlugin) CMD() string {
	return "/crypto"
}

//Help is what will show in the help menu
func (a activePlugin) Help() string {
	return "/crypto {cryptocurrency}"
}

func updateListing() {
	for {
		l := &listing{}
		req, err := http.NewRequest("GET", listingURL, nil)
		if err != nil {
			errMsg := fmt.Errorf("[Crypto Error] making http request %v", err)
			debugPrintf("%v\n", errMsg.Error())
			continue
		}
		req.Header.Set("User-Agent", userAgent)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errMsg := fmt.Errorf("[Crypto Error] making do request %v", err)
			debugPrintf("%v\n", errMsg.Error())
			continue
		}
		defer resp.Body.Close()
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errMsg := fmt.Errorf("[Crypto Error] reading resp body %v", err)
			debugPrintf("%v\n", errMsg.Error())
			continue
		}
		if err := json.Unmarshal(buf, l); err != nil {
			errMsg := fmt.Errorf("[Crypto Error] unmarshal request %v", err)
			debugPrintf("%v\n", errMsg.Error())
			continue
		}
		for _, coin := range l.Data {
			name := strings.ToLower(coin.Name)
			symbol := strings.ToLower(coin.Symbol)
			nameMap[name] = coin.ID
			symbolMap[symbol] = coin.ID
		}
		//API updates every five mintues
		time.Sleep(5 * time.Minute)
	}
}

//Get export method that satisfies an interface in the main program.
//This Get method will query the coinmarketcap api.
func (a activePlugin) Get(input string) (string, error) {
	debugPrintf("Got Input: %s\n", input)
	var id int
	lowerInput := strings.ToLower(input)
	if _, ok := nameMap[lowerInput]; ok {
		id = nameMap[lowerInput]
	} else if _, ok := symbolMap[lowerInput]; ok {
		id = symbolMap[lowerInput]
	} else {
		err := fmt.Errorf("%s not found", input)
		debugPrintf("%v\n", err.Error())
		return "", err
	}
	url := fmt.Sprintf(priceURL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errMsg := fmt.Errorf("[Crypto Error] new request %v", err)
		debugPrintf("%v\n", errMsg.Error())
		return "", errMsg
	}
	req.Header.Set("User-Agent", userAgent)
	client := &http.Client{}
	debugPrintf("Sending req %v\n", req)
	resp, err := client.Do(req)
	if err != nil {
		errMsg := fmt.Errorf("[Crypto Error] do request %v", err)
		debugPrintf("%v\n", errMsg.Error())
		return "", errMsg
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errMsg := fmt.Errorf("[Crypto Error] reading resp body %v", err)
		debugPrintf("%v\n", errMsg.Error())
		return "", errMsg
	}
	debugPrintf("Response size %d\n", len(body))
	t := &tickerSpecific{}
	if err := json.Unmarshal(body, t); err != nil {
		errMsg := fmt.Errorf("[Crypto Error] unmarshal %v", err)
		debugPrintf("%v\n", errMsg.Error())
		return "", errMsg
	}
	output := fmt.Sprintf("Name: %-8s\tUSD: %-6.2f\tBTC: %-10.9f", t.Data.Name, t.Data.Quotes.USD.Price, t.Data.Quotes.BTC.Price)
	debugPrintf("Returning %s\n", output)
	return output, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (a activePlugin) Send(subscription kbchat.SubscriptionMessage, msg string) error {
	debugPrintf("Starting kbchat\n")
	w, err := kbchat.Start("chat")
	if err != nil {
		errMsg := fmt.Errorf("[Crypto Error] sending message %v", err)
		debugPrintf("%v\n", errMsg.Error())
		return errMsg
	}
	debugPrintf("Sending this message to messageID: %s\n%s\n", subscription.Conversation.ID, msg)
	if err := w.SendMessage(subscription.Conversation.ID, msg); err != nil {
		if err := w.Proc.Kill(); err != nil {
			return err
		}
		return err
	}
	debugPrintf("Killing child process\n")
	return w.Proc.Kill()
}

func init() {
	debugPrintf = func(format string, v ...interface{}) {
		return
	}
	go updateListing()
}

func main() {}
