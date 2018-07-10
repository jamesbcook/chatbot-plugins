package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jamesbcook/chatbot-external-api/filesystem"

	"github.com/jamesbcook/chatbot-external-api/crypto"
	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
)

var (
	debugPrintf func(format string, v ...interface{})
	minute      = regexp.MustCompile("minute|minutes")
	hour        = regexp.MustCompile("hour|hours")
	day         = regexp.MustCompile("day|days")
	reminders   = []*remindBucket{}
	ourState    = &state{}
)

type activePlugin string

//AP for export
var AP activePlugin

type state struct {
	symmetric *crypto.Symmetric
	mutex     sync.RWMutex
	file      string
}

type remindBucket struct {
	Time    []byte
	Sender  []byte
	Message []byte
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

func encrypt(input []byte) []byte {
	nonce, err := crypto.GenerateNonce()
	if err != nil {
		print.Badln(err)
	}
	debugPrintf("Nonce %x\n", (*nonce)[:])
	copy(ourState.symmetric.Nonce[:], (*nonce)[:])
	encryptedDate, err := ourState.symmetric.Encrypt(input)
	if err != nil {
		print.Badln(err)
	}
	debugPrintf("Encrypted Data %x\n", encryptedDate)
	output := make([]byte, len(encryptedDate)+12)
	copy(output, ourState.symmetric.Nonce[:])
	copy(output[12:], encryptedDate)
	return output
}

func decrypt(input []byte) ([]byte, error) {
	debugPrintf("Data to be decrypted %x\n", input)
	data := make([]byte, len(input)-12)
	copy(ourState.symmetric.Nonce[:], input[:12])
	copy(data, input[12:])
	debugPrintf("Nonce %x\n", ourState.symmetric.Nonce)
	debugPrintf("Encrypted Data %x\n", data)
	res, err := ourState.symmetric.Decrypt(data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func getReminder() {
	for {
		for x, r := range reminders {
			t := time.Now()
			ourState.mutex.Lock()
			date, err := decrypt(r.Time)
			if err != nil {
				print.Badln(err)
			}
			if err := t.UnmarshalBinary(date); err != nil {
				print.Badln(err)
			}
			ourState.mutex.Unlock()
			if time.Now().After(t) {
				ourState.mutex.Lock()
				message, err := decrypt(r.Message)
				if err != nil {
					print.Badln(err)
				}
				sender, err := decrypt(r.Sender)
				if err != nil {
					print.Badln(err)
				}
				send(string(sender), string(message))
				reminders[x] = nil
				copy(reminders[x:], reminders[x+1:])
				reminders = reminders[:len(reminders)-1]
				buf, err := json.Marshal(&reminders)
				if err != nil {
					print.Badln(err)
				}
				if err := ioutil.WriteFile(ourState.file, buf, 0600); err != nil {
					print.Badln(err)
				}
				ourState.mutex.Unlock()
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
		return send(subscription.Conversation.ID, msg)
	}
	reminder := &remindBucket{}
	var timeMult time.Duration
	debugPrintf("Finding correct function for %s\n", duration)
	if minute.MatchString(duration) {
		timeMult = minuteFunc()
	} else if hour.MatchString(duration) {
		timeMult = hourFunc()
	} else if day.MatchString(duration) {
		timeMult = dayFunc()
	}
	date := time.Now().Add(time.Duration(numInt) * timeMult)
	ourState.mutex.Lock()
	reminder.Sender = encrypt([]byte(subscription.Conversation.ID))
	reminder.Message = encrypt([]byte(strings.Join(message, " ")))
	binTime, err := date.MarshalBinary()
	if err != nil {
		print.Badln(err)
	}
	reminder.Time = encrypt(binTime)
	ourState.mutex.Unlock()
	reminders = append(reminders, reminder)
	buf, err := json.Marshal(&reminders)
	if err != nil {
		print.Badln(err)
	}
	ourState.mutex.Lock()
	if err := ioutil.WriteFile(ourState.file, buf, 0600); err != nil {
		print.Badln(err)
	}
	ourState.mutex.Unlock()
	debugPrintf("%s\n", string(buf))
	debugPrintf("Number of reminders now %d\n", len(reminders))
	t := fmt.Sprintf("Your reminder is set for %s", date.Format("2006 Jan 2 15:04:05 UTC"))
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

func cryptoSetup() *crypto.Symmetric {
	symmetric := &crypto.Symmetric{}
	var password string
	var salt [32]byte
	password = os.Getenv("CHATBOT_REMINDME_PASSWORD")
	fs, err := filesystem.New("remindme")
	if err != nil {
		print.Badln(err)
	}
	if password == "" {
		print.Warningln("Missing CHATBOT_REMINDME_PASSWORD environment var")
		password = "Something you shouldn't use"
	}
	saltFile := fs.GetPasswordSaltFile()
	if _, err := os.Stat(saltFile); os.IsNotExist(err) {
		print.Warningf("%s does not exist creating a random salt\n", saltFile)
		if err := symmetric.KeyFromPassword([]byte(password), nil); err != nil {
			print.Badln(err)
		}
		tmpSalt := symmetric.GetPasswordSalt()
		copy(salt[:], tmpSalt[:])
		if err := fs.WriteToFile(salt[:], saltFile); err != nil {
			print.Badln(err)
		}
	} else {
		tmpSalt, err := filesystem.LoadFile(saltFile)
		if err != nil {
			print.Badln(err)
		}
		copy(salt[:], tmpSalt)
		if err := symmetric.KeyFromPassword([]byte(password), &salt); err != nil {
			print.Badln(err)
		}
	}
	return symmetric
}

func stateSetup() {
	fs, err := filesystem.New("remindme")
	if err != nil {
		print.Badln(err)
	}
	stateFile := fs.GetStateFile()
	ourState.file = stateFile
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		print.Warningln("No state file, creating one")
		_, err = os.Create(stateFile)
		if err != nil {
			print.Badln(err)
		}
		return
	}
	f, err := os.OpenFile(stateFile, os.O_RDONLY, 0600)
	if err != nil {
		print.Badln(err)
	}
	defer f.Close()
	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		print.Badln(err)
	}
	if len(buffer) == 0 {
		return
	}
	if err := json.Unmarshal(buffer, &reminders); err != nil {
		print.Badln(err)
	}
}

func init() {
	debugPrintf = func(format string, v ...interface{}) {
		return
	}
	s := cryptoSetup()
	ourState.symmetric = s
	stateSetup()
	go getReminder()
}

func main() {}
