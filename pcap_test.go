package cuckoo

import (
	"context"
	"testing"
)

func TestPcapGet(t *testing.T) {
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

	for task := range tasks {
		_, err := c.PcapGet(context.Background(), task.ID)
		if err != nil {
			continue
		}
		return
	}
	t.Errorf("could not find pcap to test")
}
