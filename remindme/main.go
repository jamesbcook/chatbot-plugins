package main

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jamesbcook/chatbot/kbchat"
)

var (
	//CMD that keybase will use to execute this plugin
	CMD = "/remindme"
	//Help is what will show in the help menu
	Help         = `/remindme {time} {message}`
	areDebugging = false
	debugWriter  *io.Writer
	minute       = regexp.MustCompile("minute|minutes")
	hour         = regexp.MustCompile("hour|hours")
	day          = regexp.MustCompile("day|days")
	reminders    = []*remindBucket{}
)

type remindBucket struct {
	T       time.Time
	Sender  string
	Message []string
}

type getting string

type duration func() time.Duration

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

func minuteFunc() time.Duration {
	return time.Minute
}

func hourFunc() time.Duration {
	return time.Hour
}

func dayFunc() time.Duration {
	return time.Hour * 24
}

func setRemindMe(t int, message []string, f duration) *remindBucket {
	date := time.Now().Add(time.Duration(t) * f())
	return &remindBucket{T: date, Message: message}
}

func getReminder() {
	for {
		for x, r := range reminders {
			if time.Now().After(r.T) {
				s := strings.Join(r.Message, " ")
				send(r.Sender, s)
				copy(reminders[x:], reminders[x+1:])
				reminders = reminders[:len(reminders)-1]
			}
		}
		time.Sleep(1 * time.Second)
	}
}

//Get export method that satisfies an interface in the main program.
//This Get method will query reddit json api.
func (g getting) Get(input string) (string, error) {
	args := strings.Split(input, " ")
	if len(args) <= 2 {
		return "", fmt.Errorf("Not enough arguments")
	}
	num := args[0]
	duration := args[1]
	_, err := strconv.Atoi(num)
	if err != nil {
		return "", fmt.Errorf("%s is not an int", num)
	}
	if minute.MatchString(duration) ||
		hour.MatchString(duration) ||
		day.MatchString(duration) {
		return input, nil
	}
	return "", fmt.Errorf("%s is not a valid duration", duration)
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (g getting) Send(msgID, msg string) error {
	args := strings.Split(msg, " ")
	num := args[0]
	duration := args[1]
	message := args[2:]
	numInt, err := strconv.Atoi(num)
	if err != nil {
		return fmt.Errorf("%s is not an int", num)
	}
	reminder := &remindBucket{}
	if minute.MatchString(duration) {
		reminder = setRemindMe(numInt, message, minuteFunc)
	} else if hour.MatchString(duration) {
		reminder = setRemindMe(numInt, message, hourFunc)
	} else if day.MatchString(duration) {
		reminder = setRemindMe(numInt, message, dayFunc)
	}
	reminder.Sender = msgID
	reminders = append(reminders, reminder)
	t := fmt.Sprintf("Your reminder is set for %s", reminder.T.Format("2006 Jan 2 15:04:05 UTC"))
	return send(msgID, t)
}

func send(msgID, msg string) error {
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

func init() {
	go getReminder()
}

func main() {}
