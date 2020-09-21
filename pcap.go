package cuckoo

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// PcapGet Returns the content of the PCAP associated with the given task.
func (c *Client) PcapGet(ctx context.Context, taskID int) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/pcap/get/%d", c.BaseURL, taskID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.MakeRequest(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 404:
		return nil, fmt.Errorf("file not found")
	case 200:
		break
	default:
		return nil, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
