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
)

const (
	userAgent = "KeyBase Chatbot"
	unit      = "imperial"
	urlFMT    = "https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=%s"
)

var (
	//CMD that keybase will use to execute this plugin
	CMD = "/weather"
	//Help is what will show in the help menu
	Help         = "/weather {city}"
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
//This Get method will query the openweathermap api.
func (g getting) Get(input string) (string, error) {
	url := fmt.Sprintf(urlFMT, input, os.Getenv("CHATBOT_WEATHER"), unit)
	client := &http.Client{}
	debug(fmt.Sprintf("Creating GET request to %s", url))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("[Weather Error] creating request %v", err)
	}
	req.Header.Set("User-Agent", userAgent)
	debug(fmt.Sprintf("Sending request %v", req))
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
	debug(fmt.Sprintf("Unmarshalling json with length of %d", len(buf.Bytes())))
	if err := json.Unmarshal(buf.Bytes(), o); err != nil {
		return "", fmt.Errorf("[Weather Error] unmarshalling response %v", err)
	}
	if len(o.Weather) <= 0 {
		return "", fmt.Errorf("[Weather Error] no weather found")
	}
	output := fmt.Sprintf("Weather Description: %s\nTemperature: %.2f\tHumidity: %d\tMin: %.2f\tMax: %.2f\n",
		o.Weather[0].Description, o.Temperature.Temp, o.Temperature.Humidity, o.Temperature.Min, o.Temperature.Max)
	debug(fmt.Sprintf("Message sending to user\n%s", output))
	return output, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (g getting) Send(msgID, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Weather Error] in send request %v", err)
	}
	debug(fmt.Sprintf("Sending this message to messageID: %s\n%s", msgID, msg))
	if err := w.SendMessage(msgID, msg); err != nil {
		if err := w.Proc.Kill(); err != nil {
			return err
		}
		return err
	}
	debug("Killing child process")
	return w.Proc.Kill()
}

func init() {
	if api := os.Getenv("CHATBOT_WEATHER"); api == "" {
		log.Println("Missing CHATBOT_WEATHER environment variable")
	}
}

func main() {}
