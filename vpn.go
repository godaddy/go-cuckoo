package cuckoo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// VPNStatus Returns VPN status.
// For now this decodes to a blank interface.
func (c *Client) VPNStatus(ctx context.Context) (interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/vpn/status", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.MakeRequest(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		break
	case 404:
		return nil, fmt.Errorf("not available")
	default:
		message := struct {
			Message string `json:"message"`
		}{}
		json.NewDecoder(resp.Body).Decode(&message)
		return nil, fmt.Errorf("bad response code: %d, message: %s", resp.StatusCode, message)
	}

	// TODO change this to real structure
	var status interface{}
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		return nil, fmt.Errorf("cuckoo: status response marshalling error: %w", err)
	}

	return status, nil
}
