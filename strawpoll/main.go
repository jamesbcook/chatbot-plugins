package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/jamesbcook/chatbot/kbchat"
)

const (
	userAgent = "KeyBase Chatbot"
	baseURL   = "https://www.strawpoll.me/"
	getURL    = "https://www.strawpoll.me/api/v2/polls/"
	postURL   = "https://www.strawpoll.me/api/v2/polls"
)

var (
	//CMD that keybase will use to execute this plugin
	CMD = "/strawpoll"
	//Help is what will show in the help menu
	Help         = "/strawpoll {id | title [options] (multi) (dup) (captcha)}"
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

func (g getting) Debug(set bool, writer *io.Writer) {
	areDebugging = set
	debugWriter = writer
}

func debug(input string) {
	if areDebugging && debugWriter != nil {
		output := fmt.Sprintf("[DEBUG] %s\n", input)
		(*debugWriter).Write([]byte(output))
	}
}

type poll struct {
	ID int `json:"id"`
	newPoll
	Votes []int `json:"votes"`
}

type newPoll struct {
	Title    string   `json:"title"`
	Options  []string `json:"options"`
	Multi    bool     `json:"multi"`
	DupCheck string   `json:"dupcheck"`
	Captcha  bool     `json:"captcha"`
}

func getData(id string) (*poll, error) {
	url := fmt.Sprintf(getURL+"%s", id)
	debug(fmt.Sprintf("Sending GET to %s", url))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	debug(fmt.Sprintf("Response length %d", len(body)))
	p := &poll{}
	if err := json.Unmarshal(body, p); err != nil {
		return nil, err
	}
	return p, nil
}

func postData(np *newPoll) (*poll, error) {
	enc, err := json.Marshal(np)
	if err != nil {
		return nil, err
	}
	debug(fmt.Sprintf("Sending data %v to %s", string(enc), postURL))
	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(enc))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	debug(fmt.Sprintf("Got response size of %d", len(body)))
	p := &poll{}
	if err := json.Unmarshal(body, p); err != nil {
		return nil, err
	}
	return p, nil
}

func createPoll(arguments []string) (string, error) {
	np := &newPoll{}
	var fixArgs []string
	for x := range arguments {
		if arguments[x] == " " {
			continue
		}
		fixArgs = append(fixArgs, arguments[x])
	}
	switch len(fixArgs) {
	case 5:
		v, err := strconv.ParseBool(fixArgs[4])
		if err != nil {
			return "", err
		}
		np.Captcha = v
		fallthrough
	case 4:
		np.DupCheck = fixArgs[3]
		fallthrough
	case 3:
		v, err := strconv.ParseBool(fixArgs[2])
		if err != nil {
			return "", err
		}
		np.Multi = v
		fallthrough
	case 2:
		np.Title = fixArgs[0]
		np.Options = strings.Split(fixArgs[1], ",")
	}
	p, err := postData(np)
	if err != nil {
		return "", err
	}
	debug(fmt.Sprintf("poll results %v", p))
	output := fmt.Sprintf("Title: %s\nURL: ", p.Title)
	output += fmt.Sprintf(baseURL+"%d", p.ID)
	return output, nil
}

//Get export method that satisfies an interface in the main program.
//This Get method will query reddit json api.
func (g getting) Get(input string) (string, error) {
	var output string
	arguments := strings.FieldsFunc(input, func(c rune) bool {
		if c != '"' {
			return false
		}
		return true
	})
	if len(arguments) == 1 {
		poll, err := getData(arguments[0])
		if err != nil {
			return "", err
		}
		output = fmt.Sprintf("Title: %s\nURL: ", poll.Title)
		output += fmt.Sprintf(baseURL+"%s\n", arguments[0])
		output += "Option: Vote Count\n"
		for x := range poll.Options {
			output += fmt.Sprintf("%s: %d\n", poll.Options[x], poll.Votes[x])
		}
	} else {
		var err error
		output, err = createPoll(arguments)
		if err != nil {
			return "", err
		}
	}
	debug(fmt.Sprintf("Output sending to user %s", output))
	return output, nil
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (g getting) Send(msgID, msg string) error {
	debug("Starting kbchat")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[URL Short Error] in send request %v", err)
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

func main() {}
