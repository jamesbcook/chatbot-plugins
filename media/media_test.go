package media

import (
	"bytes"
	"net/http"
	"os"
	"testing"
)

func TestSha(t *testing.T) {
	expected := "840d1ce81a4327840b54cb1d419907fd1f62359bad33656e058653d2e4172a43"
	res := shaFileName("Hello World")
	if res != expected {
		t.Fatalf("Sha doesn't match expected %s got %s", expected, res)
	}
}

func TestSetup(t *testing.T) {
	res, err := Setup("https://google.com", func(input string) ([]byte, error) {
		resp, err := http.Get(input)
		if err != nil {
			t.Fatalf("Error doing http get request %v", err)
		}
		defer resp.Body.Close()
		var buf bytes.Buffer
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			t.Fatalf("Error reading from body %v", err)
		}
		return buf.Bytes(), nil
	})
	if err != nil {
		t.Fatalf("Error during setup %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results less than or equal to 0")
	}
	if err := os.Remove(res); err != nil {
		t.Fatalf("Error removing file %v", err)
	}
}
