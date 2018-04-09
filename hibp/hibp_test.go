package hibp

import "testing"

func TestGet(t *testing.T) {
	res, err := Get("http://google.com", func(input string) string {
		return input
	})
	if err != nil {
		t.Fatalf("Error doing get request %v", err)
	}
	if len(res) <= 0 {
		t.Fatalf("Results length is less than or equal to zero")
	}
}

func TestSetupGetRequest(t *testing.T) {
	expectedHeader := map[string]string{"User-Agent": "KeyBase Chatbot"}
	expectedURL := "https://google.com"
	req, err := setupGetRequest(expectedURL)
	if err != nil {
		t.Fatalf("Error during setup of get request %v", err)
	}
	if req.URL.String() != expectedURL {
		t.Fatalf("Expected %s Got %s", req.URL.String(), expectedURL)
	}
	if req.Header.Get("User-Agent") != expectedHeader["User-Agent"] {
		t.Fatalf("Expected %s Got %s", expectedHeader["User-Agent"], req.Header.Get("User-Agent"))
	}
}
