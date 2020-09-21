package cuckoo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Sample is a file sample returned by cuckoo
type Sample struct {
	Sha1     string `json:"sha1"`
	FileType string `json:"file_type"`
	FileSize int64  `json:"file_size"`
	Crc32    string `json:"crc32"`
	Ssdeep   string `json:"ssdeep"`
	Sha256   string `json:"sha256"`
	Sha512   string `json:"sha512"`
	ID       int64  `json:"id"`
	Md5      string `json:"md5"`
}

// FileID to look up in cuckoo.  You can set any of the
// fields and leave the others blank
type FileID struct {
	ID     int
	MD5    string
	SHA256 string
}

// FilesView Returns details on the file matching either the specified MD5 hash, SHA256 hash or ID.
func (c *Client) FilesView(ctx context.Context, fileID *FileID) (*Sample, error) {
	format := "id"
	id := fmt.Sprintf("%d", fileID.ID)
	switch {
	case fileID.MD5 != "":
		format = "md5"
		id = fileID.MD5
	case fileID.SHA256 != "":
		format = "sha256"
		id = fileID.SHA256
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/files/view/%s/%s", c.BaseURL, format, id), nil)
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
	case 400:
		return nil, fmt.Errorf("invalid lookup term")
	case 200:
		break
	default:
		return nil, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	sample := struct {
		Sample *Sample `json:"sample"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&sample)
	if err != nil {
		return nil, err
	}

	return sample.Sample, nil
}

// FilesGet Returns the binary content of the file matching the specified SHA256 hash.
func (c *Client) FilesGet(ctx context.Context, sha256 string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/files/get/%s", c.BaseURL, sha256), nil)
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
