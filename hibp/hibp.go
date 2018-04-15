package hibp

import (
	"bytes"
	"fmt"
	"net/http"
)

const (
	userAgent = "KeyBase Chatbot"
)

//Request is the request for HIBP
type Request func(string) string

//Get results from HIBP based on the input and function passed in
func Get(input string, kind Request) ([]byte, error) {
	fullURL := kind(input)
	client := &http.Client{}
	req, err := setupGetRequest(fullURL)
	if err != nil {
		return nil, fmt.Errorf("[HIBP Error] setuping up request %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[HIBP Error] sending request %v", err)
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[HIBP Error] reading from resp body %v", err)
	}
	return buf.Bytes(), nil
}

func setupGetRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("[HIBP Error] creating request %v", err)
	}
	req.Header.Set("User-Agent", userAgent)
	return req, nil
}
