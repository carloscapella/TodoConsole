package usecase

import (
	"time"
	"todo/internal/domain"
)

// Edit updates the title and/or completed status of a task by id.
func (uc *TaskUseCase) Edit(id int, newTitle string, completed *bool) error {
	task, err := uc.getTask(id)
	if err != nil {
		return err
	}
	if task == nil {
		return domain.ErrTaskNotFound
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

func (uc *TaskUseCase) Add(title string, priority domain.Priority, tags []string) error {
	// If the priority is invalid or empty, use the default priority
	if priority == "" || !domain.ValidatePriority(string(priority)) {
		priority = domain.GetDefaultPriority()
	}

	task := &domain.Task{
		Title:     title,
		Completed: false,
		Priority:  priority,
		Tags:      domain.NormalizeTags(tags),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return uc.repo.Create(task)
}

func (uc *TaskUseCase) getTask(id int) (*domain.Task, error) {
	task, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, domain.ErrTaskNotFound
	}
	return task, nil
}

func (uc *TaskUseCase) updateTask(task *domain.Task) error {
	task.UpdatedAt = time.Now()
	return uc.repo.Update(task)
}

func (uc *TaskUseCase) Complete(id int) error {
	task, err := uc.getTask(id)
	if err != nil {
		return err
	}
	task.Completed = true
	return uc.updateTask(task)
}

func (uc *TaskUseCase) UpdatePriority(id int, priority domain.Priority) error {
	if !domain.ValidatePriority(string(priority)) {
		return domain.ErrInvalidPriority
	}

	task, err := uc.getTask(id)
	if err != nil {
		return err
	}

	task.Priority = priority
	return uc.updateTask(task)
}

func (uc *TaskUseCase) UpdateTags(id int, tags []string) error {
	task, err := uc.getTask(id)
	if err != nil {
		return err
	}

	task.Tags = domain.NormalizeTags(tags)
	return uc.updateTask(task)
}

func (uc *TaskUseCase) Delete(id int) error {
	return uc.repo.Delete(id)
}
