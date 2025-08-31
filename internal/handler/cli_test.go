package handler

import (
	"errors"
	"testing"
	"todo/internal/domain"
)

type mockUseCase struct {
	addErr, editErr, completeErr, delErr error
	listTasks                            []domain.Task
}

func (m *mockUseCase) Add(title string, priority domain.Priority, tags []string) error {
	return m.addErr
}
func (m *mockUseCase) Edit(id int, newTitle string, completed *bool) error   { return m.editErr }
func (m *mockUseCase) Complete(id int) error                                 { return m.completeErr }
func (m *mockUseCase) Delete(id int) error                                   { return m.delErr }
func (m *mockUseCase) List() ([]domain.Task, error)                          { return m.listTasks, nil }
func (m *mockUseCase) UpdatePriority(id int, priority domain.Priority) error { return nil }
func (m *mockUseCase) UpdateTags(id int, tags []string) error                { return nil }

func TestCommand_executeAdd_Error(t *testing.T) {
	cmd := &Command{useCase: &mockUseCase{addErr: errors.New("fail")}}
	cmd.Flags.Add = "test"
	err := cmd.executeAdd()
	if err == nil {
		t.Error("expected error")
	}
}

func TestCommand_executeList(t *testing.T) {
	cmd := &Command{useCase: &mockUseCase{listTasks: []domain.Task{{ID: 1, Title: "A"}}}}
	cmd.Flags.List = true
	err := cmd.executeList()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCommand_executeEdit_Error(t *testing.T) {
	cmd := &Command{useCase: &mockUseCase{editErr: errors.New("fail")}}
	cmd.Flags.Edit = 1
	err := cmd.executeEdit()
	if err == nil {
		t.Error("expected error")
	}
}

func TestCommand_executeComplete_Error(t *testing.T) {
	cmd := &Command{useCase: &mockUseCase{completeErr: errors.New("fail")}}
	cmd.Flags.Complete = 1
	err := cmd.executeComplete()
	if err == nil {
		t.Error("expected error")
	}
}

func TestCommand_executeDelete_Error(t *testing.T) {
	cmd := &Command{useCase: &mockUseCase{delErr: errors.New("fail")}}
	cmd.Flags.Delete = 1
	err := cmd.executeDelete()
	if err == nil {
		t.Error("expected error")
	}
}
