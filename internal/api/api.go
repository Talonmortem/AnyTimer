package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Task struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
}

func GetTasks(client *http.Client) ([]Task, error) {
	resp, err := client.Get("http://localhost:8082/tasks")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	if err := json.Unmarshal(body, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func CreateTask(client *http.Client, t Task) (Task, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return Task{}, err
	}

	resp, err := client.Post("http://localhost:8082/tasks/", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return Task{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Task{}, err
	}

	var task Task
	if err := json.Unmarshal(body, &task); err != nil {
		return Task{}, err
	}

	return task, nil
}
