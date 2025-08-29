package main

import (
	"flag"
	"fmt"
	"os"
	"todo/internal/handler"
	"todo/internal/repository"
	"todo/internal/usecase"
)

func main() {
	// Persistence flags
	jsonFlag := flag.Bool("json", false, "Use JSON file for persistence")
	sqliteFlag := flag.String("sqlite", "", "Path to SQLite file for persistence")
	fileFlag := flag.String("file", "tasks.json", "JSON file name (used with --json)")

	// Command flags
	add := flag.String("add", "", "Add a new task")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Mark task as complete by ID")
	deleteTask := flag.Int("delete", 0, "Delete task by ID")
	edit := flag.Int("edit", 0, "Edit a task by ID")
	title := flag.String("title", "", "New title for the task (used with --edit)")
	setCompleted := flag.String("set-completed", "", "Set completed status: true or false (used with --edit)")

	flag.Parse()

	var repo usecase.TaskRepository
	if *sqliteFlag != "" {
		sqlRepo, err := repository.NewSQLiteTaskRepository(*sqliteFlag)
		if err != nil {
			fmt.Println("Error opening SQLite:", err)
			os.Exit(1)
		}
		repo = sqlRepo
	} else if *jsonFlag {
		repo = repository.NewJSONTaskRepository(*fileFlag)
	} else {
		fmt.Println("You must specify --json or --sqlite <file>")
		os.Exit(1)
	}

	uc := usecase.NewTaskUseCase(repo)
	handler.RunCLI(uc, *add, *list, *complete, *deleteTask, *edit, *title, *setCompleted)
}
