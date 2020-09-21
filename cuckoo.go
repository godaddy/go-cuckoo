package cuckoo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Status is a collection of information on the status of cuckoo
type Status struct {
	Tasks           Tasks     `json:"tasks"`
	Diskspace       Diskspace `json:"diskspace"`
	Version         string    `json:"version"`
	ProtocolVersion int64     `json:"protocol_version"`
	Hostname        string    `json:"hostname"`
	Machines        Machines  `json:"machines"`
}

// Diskspace reported by cuckoo
type Diskspace struct {
	Analyses  Analyses `json:"analyses"`
	Binaries  Analyses `json:"binaries"`
	Temporary Analyses `json:"temporary"`
}

// Analyses reported by cuckoo
type Analyses struct {
	Total int64 `json:"total"`
	Free  int64 `json:"free"`
	Used  int64 `json:"used"`
}

// Machines reported by cuckoo
type Machines struct {
	Available int64 `json:"available"`
	Total     int64 `json:"total"`
}

// Tasks reported by cuckoo
type Tasks struct {
	Reported  int64 `json:"reported"`
	Running   int64 `json:"running"`
	Total     int64 `json:"total"`
	Completed int64 `json:"completed"`
	Pending   int64 `json:"pending"`
}

// CuckooStatus Returns status of the cuckoo server.
// In version 1.3 the diskspace entry was added.
// The diskspace entry shows the used, free, and total diskspace at the disk where the respective directories can be found.
// The diskspace entry allows monitoring of a Cuckoo node through the Cuckoo API.
// Note that each directory is checked separately as one may create a symlink for $CUCKOO/storage/analyses to a separate harddisk, but keep $CUCKOO/storage/binaries as-is.
// (This feature is only available under Unix!)
//
// In version 1.3 the cpuload entry was also added - the cpuload entry shows the CPU load for the past minute, the past 5 minutes, and the past 15 minutes, respectively.
// (This feature is only available under Unix!)
func (c *Client) CuckooStatus(ctx context.Context) (*Status, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/cuckoo/status", c.BaseURL), nil)
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
		return nil, fmt.Errorf("machine not found")
	default:
		return nil, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	status := &Status{}
	if err := json.NewDecoder(resp.Body).Decode(status); err != nil {
	    return nil, fmt.Errorf("cuckoo: status response marshalling error: %w", err)
	}

	return status, nil
}

// Exit Shuts down the server if in debug mode and using the werkzeug server.
func (c *Client) Exit(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/exit", c.BaseURL), nil)
	if err != nil {
		return err
	}

	resp, err := c.MakeRequest(req)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 200:
		break
	case 403:
		return fmt.Errorf("this call can only be used in debug mode")
	case 500:
		return fmt.Errorf("generic 500 error")
	default:
		return fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	return nil
}
