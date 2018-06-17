package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jamesbcook/chatbot-plugins/chatlog"
	"golang.org/x/crypto/scrypt"
)

type logging string
type backgroundPlugin string

//Logger variable to be used as an export
var Logger logging

var (
	err          error
	l            = &logger{}
	aesgcm       cipher.AEAD
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

//Write encrypted data to a log file. Random 12 byte nonce is used, and put
//in front of the cipher text
func (lo logging) Write(p []byte) (int, error) {
	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, p, nil)
	saved := make([]byte, len(ciphertext)+12)
	copy(saved, nonce)
	copy(saved[12:], ciphertext)
	return l.write(saved)
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

func (lo logging) Decrypt(src []byte) ([]byte, error) {
	decoded := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(decoded, src)
	if err != nil {
		return nil, fmt.Errorf("[Log Error] decoding bytes %v", err)
	}
	nonce := make([]byte, 12)
	copy(nonce, decoded[:12])
	plaintext, err := aesgcm.Open(nil, nonce, decoded[12:], nil)
	if err != nil {
		return nil, fmt.Errorf("[Log Error] opening ciphertext %v", err)
	}
	return plaintext, nil
}

func expandKey(inputKey string) ([]byte, error) {
	salt := [32]byte{0x6f, 0x64, 0x0e, 0xc7, 0x7f, 0x9c, 0x7a, 0xb4, 0x5f, 0xb4, 0xcc, 0x74, 0xcd, 0x73, 0x91, 0x66, 0x90, 0xd7, 0x2e, 0xd1, 0xee, 0xa7, 0xa6, 0xcd, 0x2d, 0xb1, 0xab, 0xde, 0x9e, 0x77, 0x15, 0x0a}
	return scrypt.Key([]byte(inputKey), salt[:], 16384, 8, 1, 32)
}

func init() {
	l, err = start()
	if err != nil {
		log.Fatal(err)
	}
	var key []byte
	if res := os.Getenv("CHATBOT_LOG_KEY"); res == "" {
		key, err = expandKey("Some Default Password That you Shouldn't Use")
		if err != nil {
			panic(err)
		}
		log.Println("Missing CHATBOT_LOG_KEY using default key")
	} else {
		key, err = expandKey(os.Getenv("CHATBOT_LOG_KEY"))
		if err != nil {
			panic(err)
		}
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

}

func main() {}
