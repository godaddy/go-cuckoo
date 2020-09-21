package cuckoo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// MemoryList Returns a list of memory dump files or one memory dump file associated with the specified task ID.
//
// Returns a []string{} of dump file names
func (c *Client) MemoryList(ctx context.Context, taskID int) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/memory/list/%d", c.BaseURL, taskID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.MakeRequest(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 404:
		return nil, fmt.Errorf("file or folder not found")
	case 200:
		break
	default:
		return nil, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	dumpFiles := struct {
		Files []string `json:"dump_files"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&dumpFiles)
	if err != nil {
		return nil, fmt.Errorf("cuckoo: status response marshalling error: %w", err)
	}

	return dumpFiles.Files, nil
}

// MemoryGet Returns one memory dump file associated with the specified task ID.
//
// pid - numerical identifier (pid) of a single memory dump file (e.g. 205, 1908).
//
// Note that depending on your cuckoo setup, sometimes you won't be able to download memory
// dumps. See this issue:
// https://github.com/cuckoosandbox/cuckoo/issues/2327
//
// This function returns the direct reader from the cuckoo api.
func (c *Client) MemoryGet(ctx context.Context, taskID int, pID int) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/memory/get/%d/%d", c.BaseURL, taskID, pID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.MakeRequest(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 404:
		return nil, fmt.Errorf("Memory dump not found")
	case 200:
		break
	default:
		return nil, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
