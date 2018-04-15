package media

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/sha3"
)

func shaFileName(fileName string) string {
	digest := make([]byte, 32)
	sha3.ShakeSum256(digest, []byte(fileName))
	return hex.EncodeToString(digest)
}

//Setup takes a query, downloads the file, and returns the path
func Setup(query string, f func(string) ([]byte, error)) (string, error) {
	tmpfn := filepath.Join("/tmp", shaFileName(query))

	// Create the file
	out, err := os.Create(tmpfn)
	if err != nil {
		return "", fmt.Errorf("[Media Error] Unable to create file %v", err)
	}
	output, err := f(query)
	if err != nil {
		return "", fmt.Errorf("[Media Error] Query error %v", err)
	}
	// Write the body to file
	_, err = out.Write(output)
	if err != nil {
		return "", fmt.Errorf("[Media Error] Unable to write gif to file, %v", err)
	}
	err = out.Close()
	if err != nil {
		return "", fmt.Errorf("[Media Error] Error closing file %v", err)
	}
	return tmpfn, nil
}
