package handler

import (
	"fmt"
	"os"
	"todo/internal/usecase"
)

func RunCLI(uc *usecase.TaskUseCase, add string, list bool, complete int, deleteTask int, edit int, newTitle string, setCompleted string) {
	switch {
	case edit > 0:
		var completedPtr *bool
		if setCompleted != "" {
			val := false
			if setCompleted == "true" {
				val = true
			}
			completedPtr = &val
		}
		err := uc.Edit(edit, newTitle, completedPtr)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Task edited!")
	case add != "":
		err := uc.Add(add)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Task added!")
	case list:
		tasks, err := uc.List()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		for _, t := range tasks {
			fmt.Printf("[%d] %s (done: %v)\n", t.ID, t.Title, t.Completed)
		}
	case complete > 0:
		err := uc.Complete(complete)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Task completed!")
	case deleteTask > 0:
		err := uc.Delete(deleteTask)
		if err != nil {
			fmt.Println("Errors:", err)
			os.Exit(1)
		}
		fmt.Println("Task deleted!")
	default:
		fmt.Println("Usage:")
		fmt.Println("  -add string\n\t\tAdd a new task")
		fmt.Println("  -list\n\t\tList all tasks")
		fmt.Println("  -complete int\n\t\tMark task as complete by ID")
		fmt.Println("  -delete int\n\t\tDelete task by ID")
		fmt.Println("  -edit int\n\t\tEdit a task by ID")
		fmt.Println("  -title string\n\t\tNew title for the task (used with --edit)")
		fmt.Println("  -set-completed string\n\t\tSet completed status: true or false (used with --edit)")
	}
}
