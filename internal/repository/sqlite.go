package repository

import (
	"database/sql"
	"encoding/json"
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
	// Try to add priority column if it doesn't exist (ignore error if already exists)
	db.Exec(`ALTER TABLE tasks ADD COLUMN priority TEXT DEFAULT 'medium'`)
	// Try to add tags column if it doesn't exist (ignore error if already exists)
	db.Exec(`ALTER TABLE tasks ADD COLUMN tags TEXT DEFAULT ''`)

	query := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL,
		priority TEXT NOT NULL DEFAULT 'medium',
		tags TEXT DEFAULT '',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
      )`
	_, err := db.Exec(query)
	return err
}

func (r *SQLiteTaskRepository) GetAll() ([]domain.Task, error) {
	rows, err := r.db.Query("SELECT id, title, completed, priority, tags, created_at, updated_at FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []domain.Task{}
	for rows.Next() {
		var t domain.Task
		var tagsStr string
		var createdAt, updatedAt string
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.Priority, &tagsStr, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		// Unmarshal tags from JSON string
		if tagsStr != "" {
			_ = unmarshalTags(tagsStr, &t.Tags)
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *SQLiteTaskRepository) GetByID(id int) (*domain.Task, error) {
	row := r.db.QueryRow("SELECT id, title, completed, priority, tags, created_at, updated_at FROM tasks WHERE id = ?", id)
	var t domain.Task
	var tagsStr string
	var createdAt, updatedAt string
	if err := row.Scan(&t.ID, &t.Title, &t.Completed, &t.Priority, &tagsStr, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task with id %d was not found", id)
		}
		return nil, err
	}
	t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	if tagsStr != "" {
		_ = unmarshalTags(tagsStr, &t.Tags)
	}
	return &t, nil
}

func (r *SQLiteTaskRepository) Create(task *domain.Task) error {
	now := time.Now().Format(time.RFC3339)
	tagsStr, err := marshalTags(task.Tags)
	if err != nil {
		return err
	}
	res, err := r.db.Exec(
		"INSERT INTO tasks (title, completed, priority, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		task.Title, task.Completed, task.Priority, tagsStr, now, now,
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
	tagsStr, err := marshalTags(task.Tags)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		"UPDATE tasks SET title = ?, completed = ?, priority = ?, tags = ?, updated_at = ? WHERE id = ?",
		task.Title, task.Completed, task.Priority, tagsStr, time.Now().Format(time.RFC3339), task.ID,
	)
	return err
}

// marshalTags serializa el slice de tags a JSON
func marshalTags(tags []string) (string, error) {
	if len(tags) == 0 {
		return "", nil
	}
	b, err := json.Marshal(tags)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// unmarshalTags deserializa el string JSON a slice de tags
func unmarshalTags(s string, tags *[]string) error {
	if s == "" {
		*tags = []string{}
		return nil
	}
	return json.Unmarshal([]byte(s), tags)
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
