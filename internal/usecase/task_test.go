package usecase

import (
	"errors"
	"reflect"
	"testing"
	"todo/internal/domain"
)

type mockRepo struct {
	tasks   []domain.Task
	saveErr error
}

func (m *mockRepo) GetAll() ([]domain.Task, error) {
	return m.tasks, nil
}
func (m *mockRepo) GetByID(id int) (*domain.Task, error) {
	for i := range m.tasks {
		if m.tasks[i].ID == id {
			return &m.tasks[i], nil
		}
	}
	return nil, errors.New("not found")
}
func (m *mockRepo) Create(task *domain.Task) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	task.ID = len(m.tasks) + 1
	m.tasks = append(m.tasks, *task)
	return nil
}
func (m *mockRepo) Update(task *domain.Task) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	for i := range m.tasks {
		if m.tasks[i].ID == task.ID {
			m.tasks[i] = *task
			return nil
		}
	}
	return errors.New("not found")
}
func (m *mockRepo) Delete(id int) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	for i := range m.tasks {
		if m.tasks[i].ID == id {
			m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func TestTaskUseCase_Add_And_List(t *testing.T) {
	repo := &mockRepo{}
	uc := NewTaskUseCase(repo)
	err := uc.Add("Test task", "high", []string{"go", "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tasks, err := uc.List()
	if err != nil || len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %v, err: %v", len(tasks), err)
	}
	if tasks[0].Title != "Test task" || tasks[0].Priority != "high" {
		t.Errorf("task fields not set correctly: %+v", tasks[0])
	}
}

func TestTaskUseCase_Complete(t *testing.T) {
	repo := &mockRepo{tasks: []domain.Task{{ID: 1, Title: "A", Completed: false}}}
	uc := NewTaskUseCase(repo)
	err := uc.Complete(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	task, _ := uc.getTask(1)
	if !task.Completed {
		t.Error("task should be completed")
	}
}

func TestTaskUseCase_UpdatePriority(t *testing.T) {
	repo := &mockRepo{tasks: []domain.Task{{ID: 1, Title: "A", Priority: "low"}}}
	uc := NewTaskUseCase(repo)
	err := uc.UpdatePriority(1, "high")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	task, _ := uc.getTask(1)
	if task.Priority != "high" {
		t.Errorf("expected priority high, got %s", task.Priority)
	}
}

func TestTaskUseCase_UpdateTags(t *testing.T) {
	repo := &mockRepo{tasks: []domain.Task{{ID: 1, Title: "A", Tags: []string{"old"}}}}
	uc := NewTaskUseCase(repo)
	err := uc.UpdateTags(1, []string{"new", "go"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	task, _ := uc.getTask(1)
	if !reflect.DeepEqual(task.Tags, []string{"new", "go"}) {
		t.Errorf("tags not updated: %+v", task.Tags)
	}
}

func TestTaskUseCase_Delete(t *testing.T) {
	repo := &mockRepo{tasks: []domain.Task{{ID: 1, Title: "A"}}}
	uc := NewTaskUseCase(repo)
	err := uc.Delete(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.tasks) != 0 {
		t.Error("task not deleted")
	}
}

func TestTaskUseCase_Edit(t *testing.T) {
	repo := &mockRepo{tasks: []domain.Task{{ID: 1, Title: "A", Completed: false}}}
	uc := NewTaskUseCase(repo)
	completed := true
	err := uc.Edit(1, "B", &completed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	task, _ := uc.getTask(1)
	if task.Title != "B" || !task.Completed {
		t.Errorf("edit failed: %+v", task)
	}
}

func TestTaskUseCase_Errors(t *testing.T) {
	repo := &mockRepo{}
	uc := NewTaskUseCase(repo)
	// Complete non-existent
	err := uc.Complete(99)
	if err == nil {
		t.Error("expected error for non-existent task")
	}
	// UpdatePriority invalid
	err = uc.UpdatePriority(99, "high")
	if err == nil {
		t.Error("expected error for non-existent task")
	}
	// UpdateTags invalid
	err = uc.UpdateTags(99, []string{"a"})
	if err == nil {
		t.Error("expected error for non-existent task")
	}
	// Delete invalid
	err = uc.Delete(99)
	if err == nil {
		t.Error("expected error for non-existent task")
	}
	// Edit invalid
	err = uc.Edit(99, "", nil)
	if err == nil {
		t.Error("expected error for non-existent task")
	}
}

func TestTaskUseCase_Add_InvalidPriority(t *testing.T) {
	repo := &mockRepo{}
	uc := NewTaskUseCase(repo)
	err := uc.Add("Test", "invalid", nil)
	if err != nil {
		t.Errorf("should fallback to default, got error: %v", err)
	}
	tasks, _ := uc.List()
	if tasks[0].Priority != domain.GetDefaultPriority() {
		t.Errorf("expected default priority, got %s", tasks[0].Priority)
	}
}

func TestTaskUseCase_Add_EmptyTitle(t *testing.T) {
	repo := &mockRepo{}
	uc := NewTaskUseCase(repo)
	err := uc.Add("", "medium", nil)
	if err != nil {
		t.Errorf("should allow empty title, got error: %v", err)
	}
}

func TestTaskUseCase_Create_SaveError(t *testing.T) {
	repo := &mockRepo{saveErr: errors.New("save failed")}
	uc := NewTaskUseCase(repo)
	err := uc.Add("fail", "medium", nil)
	if err == nil {
		t.Error("expected error on save failure")
	}
}
