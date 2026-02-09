package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	Priority    int       `json:"priority"` // 1: Low, 2: Medium, 3: High
	DueDate     time.Time `json:"due_date,omitempty"`
}

const (
	PriorityLow    = 1
	PriorityMedium = 2
	PriorityHigh   = 3
)

func getDBPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(homeDir, ".task-cli")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "tasks.json"), nil
}

func loadTasks() ([]Task, error) {
	dbPath, err := getDBPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return []Task{}, nil
	}

	data, err := os.ReadFile(dbPath)
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
	dbPath, err := getDBPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dbPath, data, 0644)
}

func AddTask(description string, priority int, dueDate time.Time) error {
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
		Priority:    priority,
		DueDate:     dueDate,
	}

	tasks = append(tasks, newTask)

	err = saveTasks(tasks)
	if err != nil {
		return err
	}

	fmt.Printf("Task added successfully (ID: %d)\n", newTask.ID)
	return nil
}

type ListOptions struct {
	ShowAll  bool
	ShowDone bool
	Priority int
}

func ListTasks(opts ListOptions) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	// Filter
	var filtered []Task
	for _, task := range tasks {
		if !opts.ShowAll && !opts.ShowDone && task.Completed {
			continue // Default: skip completed
		}
		if opts.ShowDone && !task.Completed {
			continue // Show only done
		}
		if opts.Priority > 0 && task.Priority != opts.Priority {
			continue
		}
		filtered = append(filtered, task)
	}
	tasks = filtered

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	// Sort: High Priority First, then Due Date (closest first), then CreatedAt
	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].Priority != tasks[j].Priority {
			return tasks[i].Priority > tasks[j].Priority // High to Low
		}
		if !tasks[i].DueDate.IsZero() && !tasks[j].DueDate.IsZero() {
			return tasks[i].DueDate.Before(tasks[j].DueDate) // Sooner first
		}
		// If one has due date and other doesn't, prioritize one with due date?
		// Let's say tasks with due dates come before tasks without.
		if !tasks[i].DueDate.IsZero() && tasks[j].DueDate.IsZero() {
			return true
		}
		if tasks[i].DueDate.IsZero() && !tasks[j].DueDate.IsZero() {
			return false
		}

		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	fmt.Println("ID  Done Pri Due         Description")
	fmt.Println("------------------------------------")
	for _, task := range tasks {
		status := " "
		if task.Completed {
			status = "x"
		}
		pri := " "
		switch task.Priority {
		case PriorityHigh:
			pri = "H"
		case PriorityMedium:
			pri = "M"
		case PriorityLow:
			pri = "L"
		}

		due := ""
		if !task.DueDate.IsZero() {
			due = task.DueDate.Format("2006-01-02")
		}

		fmt.Printf("%-3d [%s]  %s   %-10s  %s\n", task.ID, status, pri, due, task.Description)
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

func EditTask(id int, newDescription string, newPriority int, newDueDate time.Time) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	found := false
	for i, task := range tasks {
		if task.ID == id {
			if newDescription != "" {
				tasks[i].Description = newDescription
			}
			if newPriority > 0 {
				tasks[i].Priority = newPriority
			}
			if !newDueDate.IsZero() {
				tasks[i].DueDate = newDueDate
			}
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

	fmt.Printf("Task %d updated\n", id)
	return nil
}
