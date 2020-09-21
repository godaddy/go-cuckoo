package cuckoo_test

import (
	"context"
	"fmt"
	"os"

	cuckoo "github.com/godaddy/go-cukoo"
)

func ExampleClient() {
	c := cuckoo.New(&cuckoo.Config{APIKey: os.Getenv("API_KEY"), BaseURL: os.Getenv("BASEURL")})

	status, _ := c.CuckooStatus(context.Background())

	fmt.Println(status.Version)
}
