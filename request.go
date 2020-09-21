package cuckoo

import (
	"fmt"
	"net/http"
)

// ErrNotAuthorized is returned when cuckoo replies with HTTP 401
var ErrNotAuthorized = fmt.Errorf("not authorized")

// MakeRequest performs the provided request adding in the appropriate auth header
//
// It will return ErrNotAuthorized if the auth fails
func (c *Client) MakeRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// TODO: Check if this status code is accurate
	if resp.StatusCode == 401 {
		return resp, ErrNotAuthorized
	}

	return resp, err
}
