package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jamesbcook/chatbot-plugins/chatlog"
	"golang.org/x/crypto/hkdf"
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

func expandKey(inputKey string) io.Reader {
	salt := [32]byte{0x6f, 0x64, 0x0e, 0xc7, 0x7f, 0x9c, 0x7a, 0xb4, 0x5f, 0xb4, 0xcc, 0x74, 0xcd, 0x73, 0x91, 0x66, 0x90, 0xd7, 0x2e, 0xd1, 0xee, 0xa7, 0xa6, 0xcd, 0x2d, 0xb1, 0xab, 0xde, 0x9e, 0x77, 0x15, 0x0a}
	info := []byte{0x43, 0x68, 0x61, 0x74, 0x62, 0x6f, 0x74, 0x20, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x20, 0x4c, 0x6f, 0x67, 0x20, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e}
	return hkdf.New(sha256.New, []byte(inputKey), salt[:], info)
}

func init() {
	l, err = start()
	if err != nil {
		log.Fatal(err)
	}
	var dk io.Reader
	var key [32]byte
	if res := os.Getenv("CHATBOT_LOG_KEY"); res == "" {
		dk = expandKey("Some Default Password That you Shouldn't Use")
		log.Println("Missing CHATBOT_LOG_KEY using default key")
	} else {
		dk = expandKey(os.Getenv("CHATBOT_LOG_KEY"))
	}
	if _, err := io.ReadFull(dk, key[:]); err != nil {
		log.Fatal(err)
	}
	block, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err.Error())
	}
	aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

}

func main() {}
