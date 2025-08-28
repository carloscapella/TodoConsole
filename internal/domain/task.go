package domain

import "time"

// Task represents a to-do item
// Use struct tags for JSON serialization
// Use pointers for optional fields if needed
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
