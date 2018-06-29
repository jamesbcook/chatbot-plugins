package main

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
)

var (
	debugPrintf func(format string, v ...interface{})
	minute      = regexp.MustCompile("minute|minutes")
	hour        = regexp.MustCompile("hour|hours")
	day         = regexp.MustCompile("day|days")
	reminders   = []*remindBucket{}
)

type activePlugin string

//AP for export
var AP activePlugin

type remindBucket struct {
	T       time.Time
	Sender  string
	Message []string
}

type duration func() time.Duration

func (a activePlugin) Debug(set bool, writer *io.Writer) {
	debugPrintf = print.Debugf(set, writer)
}

//CMD that keybase will use to execute this plugin
func (a activePlugin) CMD() string {
	return "/remindme"
}

//Help is what will show in the help menu
func (a activePlugin) Help() string {
	return "/remindme {time} {message}"
}

func minuteFunc() time.Duration {
	debugPrintf("minute reminder\n")
	return time.Minute
}

func hourFunc() time.Duration {
	debugPrintf("hour reminder\n")
	return time.Hour
}

func dayFunc() time.Duration {
	debugPrintf("day reminder\n")
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
func (a activePlugin) Get(input string) (string, error) {
	debugPrintf("Got input %s\n", input)
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
	debugPrintf("Checking if %s is valid\n", duration)
	if minute.MatchString(duration) ||
		hour.MatchString(duration) ||
		day.MatchString(duration) {
		return input, nil
	}
	return "", fmt.Errorf("%s is not a valid duration", duration)
}

//Send export method that satisfies an interface in the main program.
//This Send method will send the results to the message ID that sent the request.
func (a activePlugin) Send(subscription kbchat.SubscriptionMessage, msg string) error {
	debugPrintf("Got message %s from ID %s\n", msg, subscription.Conversation.ID)
	args := strings.Split(msg, " ")
	num := args[0]
	duration := args[1]
	message := args[2:]
	numInt, err := strconv.Atoi(num)
	if err != nil {
		return fmt.Errorf("%s is not an int", num)
	}
	reminder := &remindBucket{}
	debugPrintf("Finding correct function for %s\n", duration)
	if minute.MatchString(duration) {
		reminder = setRemindMe(numInt, message, minuteFunc)
	} else if hour.MatchString(duration) {
		reminder = setRemindMe(numInt, message, hourFunc)
	} else if day.MatchString(duration) {
		reminder = setRemindMe(numInt, message, dayFunc)
	}
	debugPrintf("Adding %v to reminders\n", reminder)
	reminder.Sender = subscription.Conversation.ID
	reminders = append(reminders, reminder)
	debugPrintf("Number of reminders now %d\n", len(reminders))
	t := fmt.Sprintf("Your reminder is set for %s", reminder.T.Format("2006 Jan 2 15:04:05 UTC"))
	debugPrintf("Sending %s to user\n", t)
	return send(subscription.Conversation.ID, t)
}

func send(msgID, msg string) error {
	debugPrintf("Starting kbchat\n")
	w, err := kbchat.Start("chat")
	if err != nil {
		return fmt.Errorf("[Remindme Error] in send request %v", err)
	}
	debugPrintf("Sending this message to messageID: %s\n%s\n", msgID, msg)
	if err := w.SendMessage(msgID, msg); err != nil {
		if err := w.Proc.Kill(); err != nil {
			return err
		}
		return err
	}
	debugPrintf("Killing child process\n")
	return w.Proc.Kill()
}

func init() {
	go getReminder()
}

func main() {}
