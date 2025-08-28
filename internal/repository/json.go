package repository

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"todo/internal/domain"
)

type JSONTaskRepository struct {
	file string
	tasks []domain.Task
}

func NewJSONTaskRepository(file string) *JSONTaskRepository {
	r := &JSONTaskRepository{file: file}
	r.load()
	return r
}

func (r *JSONTaskRepository) load() {
	f, err := os.Open(r.file)
	if err != nil {
		r.tasks = []domain.Task{}
		return
	}
	defer f.Close()
	data, _ := ioutil.ReadAll(f)
	json.Unmarshal(data, &r.tasks)
}

func (r *JSONTaskRepository) save() error {
	data, err := json.MarshalIndent(r.tasks, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(r.file, data, 0644)
}

func (r *JSONTaskRepository) GetAll() ([]domain.Task, error) {
	return r.tasks, nil
}

func (r *JSONTaskRepository) GetByID(id int) (*domain.Task, error) {
	for i := range r.tasks {
		if r.tasks[i].ID == id {
			return &r.tasks[i], nil
		}
	}
	return nil, errors.New("not found")
}

func (r *JSONTaskRepository) Create(task *domain.Task) error {
	maxID := 0
	for _, t := range r.tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	task.ID = maxID + 1
	r.tasks = append(r.tasks, *task)
	return r.save()
}

func (r *JSONTaskRepository) Update(task *domain.Task) error {
	for i := range r.tasks {
		if r.tasks[i].ID == task.ID {
			r.tasks[i] = *task
			return r.save()
		}
	}
	return errors.New("not found")
}

func (r *JSONTaskRepository) Delete(id int) error {
	for i := range r.tasks {
		if r.tasks[i].ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return r.save()
		}
	}
	return errors.New("not found")
}
