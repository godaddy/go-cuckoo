package cuckoo

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
)

func ExampleClient_TasksScreenshots() {
	// Create your client with your API key and URL
	c := getTestingClient()

	tasks := make(chan *Task)
	go c.ListAllTasks(context.Background(), tasks)

	// Find a task with screenshots
	for task := range tasks {
		// Get all screenshots
		screenshots, err := c.TasksScreenshots(context.Background(), task.ID, -1)
		if err != nil {
			continue
		}

		// Read all data and check ZIP contents
		allData, err := ioutil.ReadAll(screenshots)
		if err != nil {
			continue
		}

		zipReader, err := zip.NewReader(bytes.NewReader(allData), int64(len(allData)))
		if err != nil {
			continue
		}

		for _, file := range zipReader.File {
			fmt.Println("We found a file!")
			_ = file.Name
			return
		}
	}

	// Output: We found a file!
}
