package cuckoo

import (
	"context"
	"testing"
)

func TestFilesView(t *testing.T) {
	c := getTestingClient()

	// This test is expecting there to be a file with ID 1 on the server
	file, err := c.FilesView(context.Background(), &FileID{ID: 1})
	if err != nil {
		t.Error(err)
		return
	}
	if file.ID != 1 || file.Md5 == "" || file.Sha1 == "" || file.Sha256 == "" || file.Sha512 == "" {
		t.Errorf("File did not return all information")
	}

	// Try getting that file
	_, err = c.FilesGet(context.Background(), file.Sha256)
	if err != nil {
		t.Error(err)
		return
	}

	// Try getting a file that doesn't exist
	_, err = c.FilesGet(context.Background(), "notexisting")
	if err == nil {
		t.Errorf("should have errored on getting file")
		return
	}
}
