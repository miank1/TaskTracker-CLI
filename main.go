package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type TaskDetails struct {
	ID          int       `json:"id,omitempty"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

var taskDetail []TaskDetails

func loadTasks() {
	data, err := os.ReadFile("todolist.json")
	if err == nil {
		_ = json.Unmarshal(data, &taskDetail)
	}
}

func saveTasks() {
	data, err := json.MarshalIndent(taskDetail, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}
	// Writing to JSON File
	err = os.WriteFile("todolist.json", data, 0644)
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <command> [arguments]")
		return
	}

	cmd := os.Args[1]
	loadTasks()

	switch cmd {
	case "add":
		if len(os.Args) < 4 {
			fmt.Println("Usage: add <description> <status>")
			return
		}
		id := 1
		if len(taskDetail) > 0 {
			id = taskDetail[len(taskDetail)-1].ID + 1
		}
		task := TaskDetails{
			ID:          id,
			Description: os.Args[2],
			Status:      os.Args[3],
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		taskDetail = append(taskDetail, task)
		saveTasks()
		fmt.Println("Task added.")

	case "update":
		if len(os.Args) < 5 {
			fmt.Println("Usage: update <id> <new_description> <new_status>")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		updated := false
		for i := range taskDetail {
			if taskDetail[i].ID == id {
				taskDetail[i].Description = os.Args[3]
				taskDetail[i].Status = os.Args[4]
				taskDetail[i].UpdatedAt = time.Now()
				updated = true
				break
			}
		}
		if updated {
			saveTasks()
			fmt.Println("Task updated.")
		} else {
			fmt.Println("Task not found.")
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: delete <id>")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		for i, task := range taskDetail {
			if task.ID == id {
				taskDetail = append(taskDetail[:i], taskDetail[i+1:]...)
				saveTasks()
				fmt.Println("Task deleted.")
				return
			}
		}
		fmt.Println("Task not found.")

	case "list":
		if len(taskDetail) == 0 {
			fmt.Println("No tasks found.")
			return
		}
		for _, task := range taskDetail {
			fmt.Printf("%d: %s [%s] (Created: %s)\n", task.ID, task.Description, task.Status, task.CreatedAt.Format(time.RFC1123))
		}

	case "mark-in-progress", "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Usage:", cmd, "<id>")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		status := "in-progress"
		if cmd == "mark-done" {
			status = "done"
		}
		updated := false
		for i := range taskDetail {
			if taskDetail[i].ID == id {
				taskDetail[i].Status = status
				taskDetail[i].UpdatedAt = time.Now()
				updated = true
				break
			}
		}
		if updated {
			saveTasks()
			fmt.Println("Task status updated.")
		} else {
			fmt.Println("Task not found.")
		}
	default:
		fmt.Println("Unknown command.")
	}
}
