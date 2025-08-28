package handler

import (
	"flag"
	"fmt"
	"os"
	"todo/internal/usecase"
)

func RunCLI(uc *usecase.TaskUseCase) {
	add := flag.String("add", "", "Add a new task")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Mark task as complete by ID")
	delete := flag.Int("delete", 0, "Delete task by ID")
	flag.Parse()

	switch {
	case *add != "":
		err := uc.Add(*add)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Task added!")
	case *list:
		tasks, err := uc.List()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		for _, t := range tasks {
			fmt.Printf("[%d] %s (done: %v)\n", t.ID, t.Title, t.Completed)
		}
	case *complete > 0:
		err := uc.Complete(*complete)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Task completed!")
	case *delete > 0:
		err := uc.Delete(*delete)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Task deleted!")
	default:
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}
}
