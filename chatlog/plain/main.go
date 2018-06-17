package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jamesbcook/chatbot-plugins/chatlog"
)

type logging string
type backgroundPlugin string

//Logger export symbol
var Logger logging

//BP for export
var BP backgroundPlugin

var (
	l            = &logger{}
	areDebugging = false
	debugWriter  *io.Writer
)

type logger struct {
	f *os.File
}

//Name that keybase will use for background plugins
func (b backgroundPlugin) Name() string {
	return "log"
}

//Debug output
func (b backgroundPlugin) Debug(set bool, writer *io.Writer) {
	areDebugging = set
	debugWriter = writer
}

func debug(input string) {
	if areDebugging && debugWriter != nil {
		output := fmt.Sprintf("[DEBUG] %s\n", input)
		(*debugWriter).Write([]byte(output))
	}
}

//Write data to a log file.
func (lo logging) Write(p []byte) (int, error) {
	return l.write(p)
}

//Start logging and return file handle
func start() (*logger, error) {
	f, err := os.OpenFile(chatlog.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("[Log Error] opening file %v", err)
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
	var err error
	l, err = start()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {}
