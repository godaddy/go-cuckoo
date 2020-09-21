package cuckoo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	// Standard number of results to get for paginated requests
	resultsPerPage = 10
)

// Task Statuses
var (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusReported  TaskStatus = "reported"
)

// ErrTaskNotFound is returned when the task is not found
var ErrTaskNotFound = fmt.Errorf("task not found")

// TaskStatus is a possible task status from cuckoo (pending, running, completed, reported)
type TaskStatus string

// Task is a task in cuckoo
type Task struct {
	Category       string        `json:"category"`
	Machine        interface{}   `json:"machine"`
	Errors         []interface{} `json:"errors"`
	Target         string        `json:"target"`
	Package        interface{}   `json:"package"`
	SampleID       interface{}   `json:"sample_id"`
	Guest          interface{}   `json:"guest"`
	Custom         interface{}   `json:"custom"`
	Owner          string        `json:"owner"`
	Priority       int64         `json:"priority"`
	Platform       interface{}   `json:"platform"`
	Options        interface{}   `json:"options"`
	Status         TaskStatus    `json:"status"`
	EnforceTimeout bool          `json:"enforce_timeout"`
	Timeout        int64         `json:"timeout"`
	Memory         bool          `json:"memory"`
	Tags           []string      `json:"tags"`
	ID             int           `json:"id"`
	AddedOn        string        `json:"added_on"`
	CompletedOn    interface{}   `json:"completed_on"`
}

// ListAllTasks Sends all tasks on cuckoo to the provided tasks channel.  It will
// close the channel once it completes or errors
//
// It will loop through all pages of the api until no more results are found
func (c *Client) ListAllTasks(ctx context.Context, tasksChan chan *Task) error {
	defer close(tasksChan)
	offset := 0

	// Keep looping until we get all tasks
	retryCount := 0
	for {
		tasks, err := c.ListTasks(ctx, resultsPerPage, offset)
		if err != nil {
			if strings.Contains(err.Error(), "500") {
				if retryCount >= 3 {
					return fmt.Errorf("max retries exceeded: %w", err)
				}
				// This is likely a temporary error, wait a bit and try again
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Second * 2):
				}
				retryCount++
				continue
			}
			return fmt.Errorf("error listing task page: %w", err)
		}
		if len(tasks) == 0 {
			// No more tasks, return
			return nil
		}

		retryCount = 0

		// Send all tasks to channel
		for _, task := range tasks {
			task := task
			select {
			case <-ctx.Done():
				return ctx.Err()
			case tasksChan <- task:
			}
		}

		offset += len(tasks)
	}
}

// ListTasks returns list of tasks.
//
// limit (optional) (int) - maximum number of returned tasks.
// offset (optional) (int) - data offset.
func (c *Client) ListTasks(ctx context.Context, limit, offset int) ([]*Task, error) {
	URL := fmt.Sprintf("%s/tasks/list/%d/%d", c.BaseURL, limit, offset)
	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.MakeRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	tasks := struct {
		Tasks []*Task `json:"tasks"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		return nil, err
	}

	return tasks.Tasks, nil
}

// ListTasksSample Returns list of tasks for sample.
func (c *Client) ListTasksSample(ctx context.Context, sampleID int) ([]*Task, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/tasks/sample/%d", c.BaseURL, sampleID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.MakeRequest(req)
	if err != nil {
		return nil, err
	}

	tasks := []*Task{}
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	return tasks, nil
}

// TasksView Returns details on the task associated with the specified ID.
func (c *Client) TasksView(ctx context.Context, taskID int) (*Task, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/tasks/view/%d", c.BaseURL, taskID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.MakeRequest(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 404:
		return nil, ErrTaskNotFound
	case 200:
		break
	default:
		return nil, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	task := struct {
		Task Task `json:"task"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		return nil, err
	}

	return &task.Task, nil
}

// TasksReschedule Reschedule a task with the specified ID and priority (default priority is 1 if -1 is passed in).
func (c *Client) TasksReschedule(ctx context.Context, taskID, priority int) (err error) {
	if priority == -1 {
		priority = 1
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/tasks/reschedule/%d/%d", c.BaseURL, taskID, priority), nil)
	if err != nil {
		return err
	}

	resp, err := c.MakeRequest(req)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 404:
		return ErrTaskNotFound
	case 200:
		break
	default:
		return fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	task := struct {
		Status string `json:"status"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		return err
	}
	if task.Status != "OK" {
		return fmt.Errorf("bad returned status: %s", task.Status)
	}

	return nil
}

// TasksDelete Removes the given task from the database and deletes the results.
func (c *Client) TasksDelete(ctx context.Context, taskID int) (err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/tasks/delete/%d", c.BaseURL, taskID), nil)
	if err != nil {
		return err
	}

	resp, err := c.MakeRequest(req)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 404:
		return ErrTaskNotFound
	case 500:
		body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1024*1024)) // Limit reading incase response is massive for some reason
		if err != nil {
			body = []byte{}
		}
		return fmt.Errorf("unable to delete the task, body: %s", body)
	case 200:
		break
	default:
		return fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	return nil
}

// TasksReport Returns the report associated with the specified task ID.
//
// It gets the reports in JSON format by default.  The report is very large and dynamic so it returns the http reader
func (c *Client) TasksReport(ctx context.Context, taskID int) (report io.ReadCloser, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/tasks/report/%d", c.BaseURL, taskID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.MakeRequest(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 404:
		return nil, fmt.Errorf("report not found")
	case 400:
		return nil, fmt.Errorf("invalid report format")
	case 200:
		break
	default:
		return nil, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// TasksScreenshots Returns one or all screenshots associated with the specified task ID.
//
// If screenshotNumber is -1, all screenshots are returned
//
// It will return a reader from the API reading the ZIP data of the screenshot(s).  You can use the zip package to read the files
func (c *Client) TasksScreenshots(ctx context.Context, taskID, screenshotNumber int) (zippedData io.ReadCloser, err error) {
	var req *http.Request
	if screenshotNumber == -1 {
		req, err = http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/tasks/screenshots/%d", c.BaseURL, taskID), nil)
	} else {
		req, err = http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/tasks/screenshots/%d/%d", c.BaseURL, taskID, screenshotNumber), nil)
	}

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

	return resp.Body, nil
}

// TasksReReport Re-run reporting for task associated with the specified task ID.
func (c *Client) TasksReReport(ctx context.Context, taskID int) (err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/tasks/rereport/%d", c.BaseURL, taskID), nil)

	if err != nil {
		return err
	}

	resp, err := c.MakeRequest(req)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 404:
		return fmt.Errorf("file or folder not found")
	case 200:
		break
	default:
		return fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	response := struct {
		Success bool `json:"success"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	if !response.Success {
		return fmt.Errorf("cuckoo returned non success")
	}

	return nil
}

// TasksReboot Add a reboot task to database from an existing analysis ID.
func (c *Client) TasksReboot(ctx context.Context, taskID int) (ID, rebootID int, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/tasks/reboot/%d", c.BaseURL, taskID), nil)

	if err != nil {
		return -1, -1, err
	}

	resp, err := c.MakeRequest(req)
	if err != nil {
		return -1, -1, err
	}
	switch resp.StatusCode {
	case 404:
		return -1, -1, fmt.Errorf("error creating reboot task")
	case 200:
		break
	default:
		return -1, -1, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	response := struct {
		TaskID   int `json:"task_id"`
		RebootID int `json:"reboot_id"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return -1, -1, err
	}

	return response.TaskID, response.RebootID, nil
}
