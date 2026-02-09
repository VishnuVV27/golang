package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		handleAdd(os.Args[2:])
	case "list":
		handleList(os.Args[2:])
	case "complete":
		handleComplete(os.Args[2:])
	case "delete":
		handleDelete(os.Args[2:])
	case "edit":
		handleEdit(os.Args[2:])
	case "help":
		printUsage()
	default:
		printUsage()
	}
}

// reorderArgs moves flags to the front so flag.Parse detects them.
// flagsWithValues is a list of flag names that expect a value (e.g. "priority").
func reorderArgs(args []string, flagsWithValues []string) []string {
	var flagArgs []string
	var nonFlagArgs []string

	skipNext := false
	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}

		if strings.HasPrefix(arg, "-") {
			flagArgs = append(flagArgs, arg)
			flagName := strings.TrimPrefix(arg, "-")
			// Check if this flag takes a value
			takesValue := false
			for _, f := range flagsWithValues {
				if f == flagName {
					takesValue = true
					break
				}
			}
			if takesValue && i+1 < len(args) {
				flagArgs = append(flagArgs, args[i+1])
				skipNext = true
			}
		} else {
			nonFlagArgs = append(nonFlagArgs, arg)
		}
	}
	return append(flagArgs, nonFlagArgs...)
}

func handleAdd(args []string) {
	// Reorder: flags first
	args = reorderArgs(args, []string{"priority", "due"})

	cmd := flag.NewFlagSet("add", flag.ExitOnError)
	priorityPtr := cmd.String("priority", "", "Priority (H, M, L)")
	duePtr := cmd.String("due", "", "Due date (YYYY-MM-DD)")

	cmd.Parse(args)

	if len(cmd.Args()) < 1 {
		fmt.Println("Usage: task-cli add <description> [-priority H|M|L] [-due YYYY-MM-DD]")
		return
	}

	description := strings.Join(cmd.Args(), " ")
	priority := parsePriority(*priorityPtr)
	dueDate := parseDueDate(*duePtr)

	err := AddTask(description, priority, dueDate)
	if err != nil {
		fmt.Printf("Error adding task: %v\n", err)
	}
}

func handleList(args []string) {
	// Reorder: flags first
	args = reorderArgs(args, []string{"priority"})

	cmd := flag.NewFlagSet("list", flag.ExitOnError)
	allPtr := cmd.Bool("all", false, "Show all tasks (pending and completed)")
	donePtr := cmd.Bool("done", false, "Show completed tasks only")
	priorityPtr := cmd.String("priority", "", "Filter by priority (H, M, L)")

	cmd.Parse(args)

	priority := parsePriority(*priorityPtr)

	opts := ListOptions{
		ShowAll:  *allPtr,
		ShowDone: *donePtr,
		Priority: priority,
	}

	err := ListTasks(opts)
	if err != nil {
		fmt.Printf("Error listing tasks: %v\n", err)
	}
}

func handleComplete(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: task-cli complete <task_id>")
		return
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid task ID")
		return
	}
	err = CompleteTask(id)
	if err != nil {
		fmt.Printf("Error completing task: %v\n", err)
	}
}

func handleDelete(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: task-cli delete <task_id>")
		return
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid task ID")
		return
	}
	err = DeleteTask(id)
	if err != nil {
		fmt.Printf("Error deleting task: %v\n", err)
	}
}

func handleEdit(args []string) {
	// Reorder: flags first
	args = reorderArgs(args, []string{"priority", "due"})

	cmd := flag.NewFlagSet("edit", flag.ExitOnError)
	priorityPtr := cmd.String("priority", "", "New priority (H, M, L)")
	duePtr := cmd.String("due", "", "New due date (YYYY-MM-DD)")

	cmd.Parse(args)

	if len(cmd.Args()) < 1 {
		fmt.Println("Usage: task-cli edit <task_id> [new description] [-priority H|M|L] [-due YYYY-MM-DD]")
		return
	}

	id, err := strconv.Atoi(cmd.Args()[0])
	if err != nil {
		fmt.Println("Invalid task ID")
		return
	}

	newDescription := ""
	if len(cmd.Args()) > 1 {
		newDescription = strings.Join(cmd.Args()[1:], " ")
	}

	newPriority := parsePriority(*priorityPtr)
	newDueDate := parseDueDate(*duePtr)

	err = EditTask(id, newDescription, newPriority, newDueDate)
	if err != nil {
		fmt.Printf("Error editing task: %v\n", err)
	}
}

func parsePriority(p string) int {
	switch strings.ToUpper(p) {
	case "H":
		return PriorityHigh
	case "M":
		return PriorityMedium
	case "L":
		return PriorityLow
	default:
		return 0 // None/Default
	}
}

func parseDueDate(d string) time.Time {
	if d == "" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02", d)
	if err != nil {
		fmt.Println("Warning: Invalid date format. Use YYYY-MM-DD. Ignoring due date.")
		return time.Time{}
	}
	return t
}

func printUsage() {
	fmt.Println("Task CLI Tool (v2.0)")
	fmt.Println("Usage:")
	fmt.Println("  add <desc> [-priority H|M|L] [-due YYYY-MM-DD]")
	fmt.Println("  list [-all] [-done] [-priority H|M|L]")
	fmt.Println("  edit <id> [desc] [-priority H|M|L] [-due YYYY-MM-DD]")
	fmt.Println("  complete <id>")
	fmt.Println("  delete <id>")
}
