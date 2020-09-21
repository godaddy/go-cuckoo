package cuckoo

import (
	"context"
	"fmt"
	"testing"
)

func TestCuckooStatus(t *testing.T) {
	c := getTestingClient()

	status, err := c.CuckooStatus(context.Background())
	if err != nil {
		t.Error(err)
	}

	fmt.Println(status.Version)

	// Check machine count
	machines, err := c.MachinesList(context.Background())
	if err != nil {
		t.Error(err)
	}
	if status.Machines.Total != int64(len(machines)) {
		t.Errorf("discrepency between number of machines %d vs %d", status.Machines.Total, len(machines))
	}

	// Check task count
	tasks := make(chan *Task)
	go func() {
		err := c.ListAllTasks(context.Background(), tasks)
		if err != nil {
			t.Error(err)
		}
	}()
	taskCount := 0
	for range tasks {
		taskCount++
	}

	if status.Tasks.Total != int64(taskCount) {
		t.Errorf("discrepency between number of tasks %d vs %d", status.Tasks.Total, taskCount)
	}
}
