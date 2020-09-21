package cuckoo

import (
	"context"
	"os"
	"testing"
)

func getTestingClient() *Client {
	return New(&Config{APIKey: os.Getenv("API_KEY"), BaseURL: os.Getenv("BASEURL")})
}

func TestClient(t *testing.T) {
	c := getTestingClient()

	// Check we have valid credentials
	err := c.CheckAuth(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
}
