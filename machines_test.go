package cuckoo

import (
	"context"
	"testing"
)

func TestMachinesList(t *testing.T) {
	c := getTestingClient()

	if err := c.CheckAuth(context.Background()); err != nil {
		t.Error(err)
		return
	}

	// This test expects machines to be available on the cuckoo server
	machines, err := c.MachinesList(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if len(machines) == 0 {
		t.Errorf("no machines found")
		return
	}

	machine, err := c.MachinesView(context.Background(), machines[0].Name)
	if err != nil {
		t.Error(err)
		return
	}
	if machine.Name != machines[0].Name {
		t.Errorf("Inconsistency in machine names %s - %s", machine.Name, machines[0].Name)
	}
	if machine.Name == "" {
		t.Errorf("Empty machine name")
	}

	// Make sure getting random machine fails
	_, err = c.MachinesView(context.Background(), "doesnotexist")
	if err == nil {
		t.Errorf("should have errored")
		return
	}
}
