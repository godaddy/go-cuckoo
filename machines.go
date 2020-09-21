package cuckoo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Machine returned by cuckoo
type Machine struct {
	Status           interface{} `json:"status"`
	Locked           bool        `json:"locked"`
	Name             string      `json:"name"`
	ResultserverIP   string      `json:"resultserver_ip"`
	IP               string      `json:"ip"`
	Tags             []string    `json:"tags"`
	Label            string      `json:"label"`
	LockedChangedOn  interface{} `json:"locked_changed_on"`
	Platform         string      `json:"platform"`
	Snapshot         interface{} `json:"snapshot"`
	Interface        interface{} `json:"interface"`
	StatusChangedOn  interface{} `json:"status_changed_on"`
	ID               int64       `json:"id"`
	ResultserverPort string      `json:"resultserver_port"`
}

// MachinesList Returns a list with details on the analysis machines available to Cuckoo.
func (c *Client) MachinesList(ctx context.Context) ([]*Machine, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/machines/list", c.BaseURL), nil)
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
	default:
		return nil, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	machines := struct {
		Machines []*Machine `json:"machines"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&machines)
	if err != nil {
		return nil, fmt.Errorf("cuckoo: status response marshalling error: %w", err)
	}

	return machines.Machines, nil
}

// MachinesView Returns details on the analysis machine associated with the given name.
func (c *Client) MachinesView(ctx context.Context, machineName string) (*Machine, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/machines/view/%s", c.BaseURL, machineName), nil)
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

	machine := struct {
		Machine *Machine `json:"machine"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&machine)
	if err != nil {
		return nil, fmt.Errorf("cuckoo: status response marshalling error: %w", err)
	}

	return machine.Machine, nil
}
