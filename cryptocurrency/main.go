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
)

const (
	userAgent  = "KeyBase Chatbot"
	listingURL = "https://api.coinmarketcap.com/v2/listings/"
	priceURL   = "https://api.coinmarketcap.com/v2/ticker/%d/?convert=BTC"
)

var (
	//CMD that keybase will use to execute this plugin
	CMD = "/crypto"
	//Help is what will show in the help menu
	Help         = "/crypto {cryptocurrency}"
	areDebugging = false
	debugWriter  *io.Writer
	nameMap      = make(map[string]int)
	symbolMap    = make(map[string]int)
)

type getting string

//Getter export symbol
var Getter getting

//Sender export symbol
var Sender getting

//Debugger export Symbol
var Debugger getting

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

func updateListing() {
	//If debug is passed this breaks it, so we wait some time until everything is loaded
	time.Sleep(40 * time.Second)
	for {
		l := &listing{}
		req, err := http.NewRequest("GET", listingURL, nil)
		if err != nil {
			errMsg := fmt.Errorf("[Crypto Error] making http request %v", err)
			debug(errMsg.Error())
			continue
		}
		req.Header.Set("User-Agent", userAgent)
		client := &http.Client{}
		debug(fmt.Sprintf("Sending request %v", req))
		resp, err := client.Do(req)
		if err != nil {
			errMsg := fmt.Errorf("[Crypto Error] making do request %v", err)
			debug(errMsg.Error())
			continue
		}
		defer resp.Body.Close()
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errMsg := fmt.Errorf("[Crypto Error] reading resp body %v", err)
			debug(errMsg.Error())
			continue
		}
		debug(fmt.Sprintf("response length %d", len(buf)))
		if err := json.Unmarshal(buf, l); err != nil {
			errMsg := fmt.Errorf("[Crypto Error] unmarshal request %v", err)
			debug(errMsg.Error())
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
func (g getting) Get(input string) (string, error) {
	debug(fmt.Sprintf("Got Input: %s", input))
	var id int
	lowerInput := strings.ToLower(input)
	if _, ok := nameMap[lowerInput]; ok {
		id = nameMap[lowerInput]
	} else if _, ok := symbolMap[lowerInput]; ok {
		id = symbolMap[lowerInput]
	} else {
		err := fmt.Errorf("%s not found", input)
		debug(err.Error())
		return "", err
	}
	url := fmt.Sprintf(priceURL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errMsg := fmt.Errorf("[Crypto Error] new request %v", err)
		debug(errMsg.Error())
		return "", errMsg
	}
	req.Header.Set("User-Agent", userAgent)
	client := &http.Client{}
	debug(fmt.Sprintf("Sending req %v", req))
	resp, err := client.Do(req)
	if err != nil {
		errMsg := fmt.Errorf("[Crypto Error] do request %v", err)
		debug(errMsg.Error())
		return "", errMsg
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errMsg := fmt.Errorf("[Crypto Error] reading resp body %v", err)
		debug(errMsg.Error())
		return "", errMsg
	}
	debug(fmt.Sprintf("Response size %d", len(body)))
	t := &tickerSpecific{}
	if err := json.Unmarshal(body, t); err != nil {
		errMsg := fmt.Errorf("[Crypto Error] unmarshal %v", err)
		debug(errMsg.Error())
		return "", errMsg
	}
	output := fmt.Sprintf("Name: %-8s\tUSD: %-6.2f\tBTC: %-10f", t.Data.Name, t.Data.Quotes.USD.Price, t.Data.Quotes.BTC.Price)
	debug(fmt.Sprintf("Returning %s", output))
	return output, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (g getting) Send(msgID, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		errMsg := fmt.Errorf("[Crypto Error] sending message %v", err)
		debug(errMsg.Error())
		return errMsg
	}
	debug(fmt.Sprintf("Sending this message to messageID: %s\n%s", msgID, msg))
	if err := w.SendMessage(msgID, msg); err != nil {
		return w.Proc.Kill()
	}
	debug("Killing child process")
	return w.Proc.Kill()
}

func init() {
	go updateListing()
}

func main() {}
