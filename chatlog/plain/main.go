package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jamesbcook/chat-bot-plugins/chatlog"
)

var (
	//Name that keybase will use for background plugins
	Name = "log"
)

type getting string

//Logger export symbol
var Logger getting

var (
	err error
	l   = &logger{}
)

type logger struct {
	f *os.File
}

//Write data to a log file.
func (g getting) Write(p []byte) (int, error) {
	return l.write(p)
}

//Start logging and return file handle
func start() (*logger, error) {
	f, err := os.OpenFile(chatlog.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error opening file %v", err)
	}
	l.f = f
	return l, nil
}

//Write input to log file and sync
func (l *logger) write(p []byte) (int, error) {
	formated := fmt.Sprintf(chatlog.StrFMT,
		time.Now().Format(chatlog.TimeFMT), p)
	return l.f.Write([]byte(formated))
}

func init() {
	l, err = start()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {}
