package repository

import (
	"os"
	"testing"
	"todo/internal/domain"
)

func TestSQLiteTaskRepository_CRUD(t *testing.T) {
	dbFile := "test_tasks.db"
	defer os.Remove(dbFile)
	repo, err := NewSQLiteTaskRepository(dbFile)
	if err != nil {
		t.Fatalf("failed to create repo: %v", err)
	}
	task := &domain.Task{Title: "Test", Priority: domain.PriorityMedium}
	if err := repo.Create(task); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	tasks, _ := repo.GetAll()
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
	fetched, err := repo.GetByID(task.ID)
	if err != nil || fetched.Title != "Test" {
		t.Errorf("get by id failed: %v", err)
	}
	task.Title = "Updated"
	if err := repo.Update(task); err != nil {
		t.Errorf("update failed: %v", err)
	}
	if err := repo.Delete(task.ID); err != nil {
		t.Errorf("delete failed: %v", err)
	}
	if _, err := repo.GetByID(task.ID); err == nil {
		t.Error("expected error for deleted task")
	}
}
