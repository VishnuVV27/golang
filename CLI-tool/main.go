package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli add <description>")
			return
		}
		description := os.Args[2]
		err := AddTask(description)
		if err != nil {
			fmt.Printf("Error adding task: %v\n", err)
		}

	case "list":
		err := ListTasks()
		if err != nil {
			fmt.Printf("Error listing tasks: %v\n", err)
		}

	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli complete <task_id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}
		err = CompleteTask(id)
		if err != nil {
			fmt.Printf("Error completing task: %v\n", err)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli delete <task_id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}
		err = DeleteTask(id)
		if err != nil {
			fmt.Printf("Error deleting task: %v\n", err)
		}

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Task CLI Tool")
	fmt.Println("Usage:")
	fmt.Println("  task-cli add <description>   Add a new task")
	fmt.Println("  task-cli list                List all tasks")
	fmt.Println("  task-cli complete <id>       Mark a task as completed")
	fmt.Println("  task-cli delete <id>         Delete a task")
}
