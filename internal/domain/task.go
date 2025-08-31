package domain

import (
	"fmt"
	"strings"
	"time"
)

// Priority represents the priority level of a task
type Priority string

// TaskStatus represents the status of a task
type TaskStatus string

// Task priorities
const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Task statuses
const (
	StatusPending   TaskStatus = "pending"
	StatusCompleted TaskStatus = "completed"
)

// Task represents a to-do item
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	Priority  Priority  `json:"priority"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TaskError represents domain-specific errors
type TaskError struct {
	Code    string
	Message string
}

func (e TaskError) Error() string {
	return e.Message
}

// Common task errors
var (
	ErrInvalidPriority = TaskError{
		Code: "INVALID_PRIORITY",
		Message: fmt.Sprintf("invalid priority. Must be one of: %s, %s, %s",
			PriorityLow, PriorityMedium, PriorityHigh),
	}
	ErrTaskNotFound = TaskError{
		Code:    "TASK_NOT_FOUND",
		Message: "task not found",
	}
)

// GetStatus returns the current status of the task
func (t *Task) GetStatus() TaskStatus {
	if t.Completed {
		return StatusCompleted
	}
	return StatusPending
}

// ValidatePriority checks if a priority value is valid
func ValidatePriority(p string) bool {
	switch Priority(p) {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}

// GetDefaultPriority returns the default priority for new tasks
func GetDefaultPriority() Priority {
	return PriorityMedium
}

// NormalizeTags cleans and normalizes the tags
func NormalizeTags(tags []string) []string {
	if len(tags) == 0 {
		return []string{}
	}

	uniqueTags := make(map[string]bool)
	for _, tag := range tags {
		if cleaned := cleanTag(tag); cleaned != "" {
			uniqueTags[cleaned] = true
		}
	}

	return mapToSlice(uniqueTags)
}

// cleanTag removes spaces and converts to lowercase
func cleanTag(tag string) string {
	return strings.ToLower(strings.TrimSpace(tag))
}

// mapToSlice converts a map[string]bool to []string
func mapToSlice(m map[string]bool) []string {
	result := make([]string, 0, len(m))
	for key := range m {
		result = append(result, key)
	}
	return result
}
