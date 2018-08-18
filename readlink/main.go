package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
	"golang.org/x/net/html"
)

type extraPlugin string

//Extra for export interface
var Extra extraPlugin

var debugPrintf func(format string, v ...interface{})

//Name that keybase will use for background plugins
func (e extraPlugin) Name() string {
	return "linkreader"
}

//Debug output
func (e extraPlugin) Debug(set bool, writer *io.Writer) {
	debugPrintf = print.Debugf(set, writer)
}

func (e extraPlugin) Get(message string) (string, error) {
	fields := strings.Fields(message)
	debugPrintf("[READLINK] %v\n", fields)
	for _, field := range fields {
		if isValidURL(field) {
			res, err := http.Get(field)
			if err != nil {
				return "", err
			}
			defer res.Body.Close()
			n, err := html.Parse(res.Body)
			if err != nil {
				return "", err
			}
			title, found := transverseBody(n)
			if found {
				debugPrintf("FOUND! %s\n", title)
				retMessage := fmt.Sprintf("Link Title: %s\n", title)
				return retMessage, nil
			}
		}
	}
	return "", nil
}

func transverseBody(node *html.Node) (string, bool) {
	if node.Type == html.ElementNode && node.Data == "title" {
		return node.FirstChild.Data, true
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		res, ok := transverseBody(child)
		if ok {
			return res, true
		}
	}
	return "", false
}

func (e extraPlugin) Send(conversationID, message string) error {
	debugPrintf("Starting kbchat\n")
	w, err := kbchat.Start("chat")
	if err != nil {
		errMsg := fmt.Errorf("[Readlink Error] sending message %v", err)
		debugPrintf("%v\n", errMsg.Error())
		return errMsg
	}
	debugPrintf("Sending this message to messageID: %s\n", conversationID)
	if err := w.SendMessage(conversationID, message); err != nil {
		if err := w.Proc.Kill(); err != nil {
			return err
		}
		return err
	}
	debugPrintf("Killing child process\n")
	return w.Proc.Kill()
}

func isValidURL(query string) bool {
	_, err := url.ParseRequestURI(query)
	if err != nil {
		debugPrintf("%s is not a title\n", query)
		return false
	}
	debugPrintf("%s is a title\n", query)
	return true
}

func main() {}
