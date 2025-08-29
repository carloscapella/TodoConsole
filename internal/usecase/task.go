package usecase

import (
	"errors"
	"time"
	"todo/internal/domain"
)

type TaskRepository interface {
	GetAll() ([]domain.Task, error)
	GetByID(id int) (*domain.Task, error)
	Create(task *domain.Task) error
	Update(task *domain.Task) error
	Delete(id int) error
}

type TaskUseCase struct {
	repo TaskRepository
}

func NewTaskUseCase(r TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: r}
}

func (uc *TaskUseCase) List() ([]domain.Task, error) {
	return uc.repo.GetAll()
}

func (uc *TaskUseCase) Add(title string) error {
	task := &domain.Task{
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return uc.repo.Create(task)
}

func (uc *TaskUseCase) Complete(id int) error {
	task, err := uc.repo.GetByID(id)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}
	task.Completed = true
	task.UpdatedAt = time.Now()
	return uc.repo.Update(task)
}

func (uc *TaskUseCase) Delete(id int) error {
	return uc.repo.Delete(id)
}

// Edit updates the title and/or completed status of a task by id.
func (uc *TaskUseCase) Edit(id int, newTitle string, completed *bool) error {
	task, err := uc.repo.GetByID(id)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}
	if newTitle != "" {
		task.Title = newTitle
	}
	if completed != nil {
		task.Completed = *completed
	}
	task.UpdatedAt = time.Now()
	return uc.repo.Update(task)
}
