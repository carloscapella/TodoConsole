package main

import (
	"todo/internal/handler"
	"todo/internal/repository"
	"todo/internal/usecase"
)

func main() {
	repo := repository.NewJSONTaskRepository("tasks.json")
	uc := usecase.NewTaskUseCase(repo)
	handler.RunCLI(uc)
}
