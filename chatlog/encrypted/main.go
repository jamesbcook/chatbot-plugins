package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/jamesbcook/chatbot-external-api/crypto"
	"github.com/jamesbcook/chatbot-external-api/filesystem"
	"github.com/jamesbcook/chatbot-plugins/chatlog"
	"github.com/jamesbcook/print"
)

type logging string
type backgroundPlugin string

//Logger variable to be used as an export
var Logger logging

//BP for export
var BP backgroundPlugin

var (
	err          error
	l            = &logger{}
	ourState     = &state{}
	areDebugging = false
	debugPrintf  func(format string, v ...interface{})
)

type state struct {
	symmetric *crypto.Symmetric
	mutex     sync.RWMutex
	file      string
}

type logger struct {
	f *os.File
}

//Name that keybase will use for background plugins
func (b backgroundPlugin) Name() string {
	return "log"
}

//Debug output
func (b backgroundPlugin) Debug(set bool, writer *io.Writer) {
	var out io.Writer
	out = os.Stdout
	debugPrintf = print.Debugf(set, &out)
}

//Write encrypted data to a log file. Random 12 byte nonce is used, and put
//in front of the cipher text
func (lo logging) Write(p []byte) (int, error) {
	ourState.mutex.Lock()
	ciphertext := encrypt(p)
	ourState.mutex.Unlock()
	return l.write(ciphertext)
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
	encoded := hex.EncodeToString(p)
	formated := fmt.Sprintf(chatlog.StrFMT,
		time.Now().Format(chatlog.TimeFMT), encoded+"\n")
	return l.f.Write([]byte(formated))
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

func (lo logging) Decrypt(src []byte) ([]byte, error) {
	decoded := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(decoded, src)
	if err != nil {
		return nil, fmt.Errorf("[Log Error] decoding bytes %v", err)
	}
	plaintext, err := decrypt(decoded)
	if err != nil {
		return nil, fmt.Errorf("[Log Error] opening ciphertext %v", err)
	}
	return plaintext, nil
}

func cryptoSetup() *crypto.Symmetric {
	symmetric := &crypto.Symmetric{}
	var password string
	var salt [32]byte
	password = os.Getenv("CHATBOT_LOG_PASSWORD")
	fs, err := filesystem.New("log")
	if err != nil {
		print.Badln(err)
	}
	if password == "" {
		print.Warningln("Missing CHATBOT_LOG_PASSWORD environment var")
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

func init() {
	debugPrintf = func(format string, v ...interface{}) {
		return
	}
	l, err = start()
	if err != nil {
		print.Badln(err)
	}
	s := cryptoSetup()
	ourState.symmetric = s
}

func main() {}
