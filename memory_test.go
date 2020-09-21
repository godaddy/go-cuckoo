package cuckoo

import (
	"context"
	"testing"
)

func TestMemoryList(t *testing.T) {
	c := getTestingClient()

	if err := c.CheckAuth(context.Background()); err != nil {
		t.Error(err)
		return
	}

	tasks := make(chan *Task)
	go func() {
		err := c.ListAllTasks(context.Background(), tasks)
		if err != nil {
			t.Error(err)
		}
	}()

	// File a task that has a memory files
	for task := range tasks {
		memory, err := c.MemoryList(context.Background(), task.ID)
		if err != nil {
			if err.Error() == "file or folder not found" {
				continue
			}
			t.Error(err)
			return
		}

		if len(memory) == 0 {
			continue
		}

		// Try getting the raw content
		_, err = c.MemoryGet(context.Background(), task.ID, 3428)
		if err != nil {
			t.Error(err)
			return
		}

		// We found a sample that we could download
		return
	}

	t.Errorf("did not find a sample to test on")
}
