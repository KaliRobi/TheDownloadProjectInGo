package server

import (
	"io"
	"net/http"
	"testing"
	"time"
)

func TestStart(t *testing.T) {

	testUrlContent := "Test content"

	downloadError := error(nil)

	go StartServer(testUrlContent, downloadError)

	time.Sleep(time.Second * 2)

	response, err := http.Get("http://localhost:8080/content")
	if err != nil {
		t.Fatalf("http.Get: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", response.Status)
	}

	responseBody, _ := io.ReadAll(response.Body)

	if len(responseBody) == 0 && err != nil {
		t.Fatalf("unexpected length: %v", err)
	}
}
