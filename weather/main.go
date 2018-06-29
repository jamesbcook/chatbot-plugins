package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
)

const (
	userAgent = "KeyBase Chatbot"
	unit      = "imperial"
	urlFMT    = "https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=%s"
)

var (
	debugPrintf func(format string, v ...interface{})
)

type activePlugin string

//AP for export
var AP activePlugin

//OpenWeather result struct
type OpenWeather struct {
	Weather     []Weather `json:"weather"`
	Temperature `json:"main"`
	Wind        `json:"wind"`
}

//Weather results from a request
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
}

//Temperature results of a city
type Temperature struct {
	Temp     float32 `json:"temp"`
	Pressure float32 `json:"pressure"`
	Humidity int     `json:"humidity"`
	Min      float32 `json:"temp_min"`
	Max      float32 `json:"temp_max"`
}

//Wind speed for a city
type Wind struct {
	Speed float32 `json:"speed"`
}

func (a activePlugin) Debug(set bool, writer *io.Writer) {
	debugPrintf = print.Debugf(set, writer)
}

//CMD that keybase will use to execute this plugin
func (a activePlugin) CMD() string {
	return "/weather"
}

//Help is what will show in the help menu
func (a activePlugin) Help() string {
	return "/weather {city}"
}

//Get export method that satisfies an interface in the main program.
//This Get method will query the openweathermap api.
func (a activePlugin) Get(input string) (string, error) {
	url := fmt.Sprintf(urlFMT, input, os.Getenv("CHATBOT_WEATHER"), unit)
	client := &http.Client{}
	debugPrintf("Creating GET request to %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("[Weather Error] creating request %v", err)
	}
	req.Header.Set("User-Agent", userAgent)
	debugPrintf("Sending request %v\n", req)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("[Weather Error] sending request %v", err)
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", fmt.Errorf("[Weather Error] reading from resp body %v", err)
	}
	o := &OpenWeather{}
	debugPrintf("Unmarshalling json with length of %d\n", len(buf.Bytes()))
	if err := json.Unmarshal(buf.Bytes(), o); err != nil {
		return "", fmt.Errorf("[Weather Error] unmarshalling response %v", err)
	}
	if len(o.Weather) <= 0 {
		return "", fmt.Errorf("[Weather Error] no weather found")
	}
	output := fmt.Sprintf("Weather Description: %s\nTemperature: %.2f\tHumidity: %d\tMin: %.2f\tMax: %.2f\n",
		o.Weather[0].Description, o.Temperature.Temp, o.Temperature.Humidity, o.Temperature.Min, o.Temperature.Max)
	debugPrintf("Message sending to user\n%s\n", output)
	return output, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (a activePlugin) Send(subscription kbchat.SubscriptionMessage, msg string) error {
	debugPrintf("Starting kbchat\n")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Weather Error] in send request %v", err)
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
	if api := os.Getenv("CHATBOT_WEATHER"); api == "" {
		log.Println("Missing CHATBOT_WEATHER environment variable")
	}
}

func main() {}
