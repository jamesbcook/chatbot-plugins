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
)

var (
	//Name of plugin module
	Name = "log"
)

type logging string

//Logger variable to be used as an export
var Logger logging

var (
	err    error
	l      = &logger{}
	aesgcm cipher.AEAD
)

type logger struct {
	f *os.File
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

func init() {
	l, err = start()
	if err != nil {
		log.Fatal(err)
	}
	var decodedKey []byte
	if res := os.Getenv("CHATBOT_LOG_KEY"); res == "" {
		decodedKey, err = hex.DecodeString("43f287ac2487750aaf4b3cafa3f4c979")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Missing CHATBOT_LOG_KEY using default key")
	} else {
		decodedKey, err = hex.DecodeString(os.Getenv("CHATBOT_LOG_KEY"))
		if err != nil {
			log.Fatal(err)
		}
	}

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		panic(err.Error())
	}
	aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

}

func main() {}
