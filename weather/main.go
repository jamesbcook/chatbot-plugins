package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	Help = "/weather {city}"
)

type getting string

//Getter export symbol
var Getter getting

//Sender export symbol
var Sender getting

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
	Pressure int     `json:"pressure"`
	Humidity int     `json:"humidity"`
	Min      float32 `json:"temp_min"`
	Max      float32 `json:"temp_max"`
}

//Wind speed for a city
type Wind struct {
	Speed float32 `json:"speed"`
}

//Get export method that satisfies an interface in the main program.
//This Get method will query the openweathermap api.
func (g getting) Get(input string) (string, error) {
	url := fmt.Sprintf(urlFMT, input, os.Getenv("CHATBOT_WEATHER"), unit)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("[Weather Error] creating request %v", err)
	}
	req.Header.Set("User-Agent", userAgent)
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
	if err := json.Unmarshal(buf.Bytes(), o); err != nil {
		return "", fmt.Errorf("[Weather Error] unmarshalling response %v", err)
	}
	if len(o.Weather) <= 0 {
		return "", fmt.Errorf("[Weather Error] no weather found")
	}
	output := fmt.Sprintf("Weather Description: %s\nTemperature: %.2f\tHumidity: %d\tMin: %.2f\tMax: %.2f\n",
		o.Weather[0].Description, o.Temperature.Temp, o.Temperature.Humidity, o.Temperature.Min, o.Temperature.Max)
	return output, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (g getting) Send(msgID, msg string) error {
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Weather Error] in send request %v", err)
	}
	return w.SendMessage(msgID, msg)
}

func main() {}
