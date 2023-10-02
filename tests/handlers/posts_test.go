package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type MockApiRouter struct{}

func (s *MockApiRouter) handleGetPosts(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func TestHandleGetPosts(t *testing.T) {
	apiRouter := &MockApiRouter{}

	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	os.Setenv("TEMPLATES_DIR", "/templates/posts")

	err = apiRouter.handleGetPosts(w, req)
	if err != nil {
		t.Fatalf("handleGetPosts returned an error: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

}
