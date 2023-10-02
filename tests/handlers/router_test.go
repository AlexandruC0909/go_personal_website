package tests

import (
	"go_api/handlers"
	"net/http"
	"testing"
	"time"
)

func TestApiRouter_Run(t *testing.T) {
	apiRouter := &handlers.ApiRouter{}

	done := make(chan bool)

	go func() {
		apiRouter.Run()
		done <- true
	}()

	req, err := http.NewRequest("GET", "http://localhost/", nil)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	done <- true

	<-done
}
