package main

import (
	"fmt"
	"os"
	"todo/internal/handler"
	"todo/internal/repository"
	"todo/internal/usecase"
)

func main() {
	// Create and configure CLI command
	cmd := handler.NewCommand()

	// Initialize repository based on flags
	var repo usecase.TaskRepository
	if cmd.Flags.SQLitePath != "" {
		sqlRepo, err := repository.NewSQLiteTaskRepository(cmd.Flags.SQLitePath)
		if err != nil {
			fmt.Println("Error opening SQLite:", err)
			os.Exit(1)
		}
		repo = sqlRepo
	} else if cmd.Flags.UseJSON {
		repo = repository.NewJSONTaskRepository(cmd.Flags.FilePath)
	} else {
		fmt.Println("You must specify --json or --sqlite <file>")
		os.Exit(1)
	}

	// Initialize use case and execute command
	uc := usecase.NewTaskUseCase(repo)
	cmd.SetUseCase(uc)

	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
