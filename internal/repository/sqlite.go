package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"todo/internal/domain"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteTaskRepository struct {
	db *sql.DB
}

func NewSQLiteTaskRepository(dbPath string) (*SQLiteTaskRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err := createTableIfNotExists(db); err != nil {
		return nil, err
	}
	return &SQLiteTaskRepository{db: db}, nil
}

func createTableIfNotExists(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)`
	_, err := db.Exec(query)
	return err
}

func (r *SQLiteTaskRepository) GetAll() ([]domain.Task, error) {
	rows, err := r.db.Query("SELECT id, title, completed, created_at, updated_at FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []domain.Task{}
	for rows.Next() {
		var t domain.Task
		var createdAt, updatedAt string
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *SQLiteTaskRepository) GetByID(id int) (*domain.Task, error) {
	row := r.db.QueryRow("SELECT id, title, completed, created_at, updated_at FROM tasks WHERE id = ?", id)
	var t domain.Task
	var createdAt, updatedAt string
	if err := row.Scan(&t.ID, &t.Title, &t.Completed, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task with id %d was not found", id)
		}
		return nil, err
	}
	t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &t, nil
}

func (r *SQLiteTaskRepository) Create(task *domain.Task) error {
	now := time.Now().Format(time.RFC3339)
	res, err := r.db.Exec(
		"INSERT INTO tasks (title, completed, created_at, updated_at) VALUES (?, ?, ?, ?)",
		task.Title, task.Completed, now, now,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		task.ID = int(id)
	}
	return nil
}

func (r *SQLiteTaskRepository) Update(task *domain.Task) error {
	_, err := r.db.Exec(
		"UPDATE tasks SET title = ?, completed = ?, updated_at = ? WHERE id = ?",
		task.Title, task.Completed, time.Now().Format(time.RFC3339), task.ID,
	)
	return err
}

func (r *SQLiteTaskRepository) Delete(id int) error {
	// Check if task exists first
	task, err := r.GetByID(id)
	if err != nil {
		return err
	}
	if task == nil {
		return fmt.Errorf("task with id %d was not found", id)
	}
	res, err := r.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d was not found", id)
	}
	fmt.Printf("Task with id %d was deleted successfully\n", id)
	return nil
}
