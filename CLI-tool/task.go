package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}

const tasksFile = "tasks.json"

func loadTasks() ([]Task, error) {
	if _, err := os.Stat(tasksFile); os.IsNotExist(err) {
		return []Task{}, nil
	}

	data, err := os.ReadFile(tasksFile)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tasksFile, data, 0644)
}

func AddTask(description string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	id := 1
	if len(tasks) > 0 {
		id = tasks[len(tasks)-1].ID + 1
	}

	newTask := Task{
		ID:          id,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}

	tasks = append(tasks, newTask)

	err = saveTasks(tasks)
	if err != nil {
		return err
	}

	fmt.Printf("Task added successfully (ID: %d)\n", newTask.ID)
	return nil
}

func ListTasks() error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	fmt.Println("ID  Done  Description")
	fmt.Println("---------------------")
	for _, task := range tasks {
		status := " "
		if task.Completed {
			status = "x"
		}
		fmt.Printf("%-3d [%s]   %s\n", task.ID, status, task.Description)
	}
	return nil
}

func CompleteTask(id int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	found := false
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Completed = true
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}

	err = saveTasks(tasks)
	if err != nil {
		return err
	}

	fmt.Printf("Task %d marked as completed\n", id)
	return nil
}

func DeleteTask(id int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	newTasks := []Task{}
	found := false
	for _, task := range tasks {
		if task.ID == id {
			found = true
			continue
		}
		newTasks = append(newTasks, task)
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}

	err = saveTasks(newTasks)
	if err != nil {
		return err
	}

	fmt.Printf("Task %d deleted\n", id)
	return nil
}
