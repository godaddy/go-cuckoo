// Package cuckoo is a go client library to interact with the cuckoo REST api
package cuckoo

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultTimeout = time.Second * 30
)

// Client to interact with cuckoo
type Client struct {
	// Auth that fits in: "Authorization: Bearer %s"
	APIKey  string
	BaseURL string

	// Client used for requests
	Client *http.Client
}

// Config is the configuration required to create a client
type Config struct {
	APIKey  string
	BaseURL string
	// Optional, if nil a new client will be created
	// with a defaultTimeout
	Client *http.Client
}

// New Creates a new client based on the provided API Key
func New(c *Config) *Client {
	client := c.Client
	if client == nil {
		client = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	return &Client{
		APIKey:  c.APIKey,
		BaseURL: c.BaseURL,
		Client:  client,
	}
}

// CheckAuth returns an error if the APIKey is not valid, or no error if valid.
//
// Under the hood it simply makes a call to /tasks/list
func (c *Client) CheckAuth(ctx context.Context) error {
	_, err := c.ListTasks(ctx, 1, 0)
	if err != nil {
		return err
	}
	return nil
}
