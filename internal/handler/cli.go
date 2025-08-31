package handler

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"todo/internal/domain"
	"todo/internal/usecase"
)

// CLIFlags defines all available CLI flags
type CLIFlags struct {
	// Persistence options
	UseJSON    bool
	SQLitePath string
	FilePath   string

	// Task operations
	Add          string
	List         bool
	Complete     int
	Delete       int
	Edit         int
	NewTitle     string
	SetCompleted string
	Priority     string
	Tags         string
}

// Command represents the CLI command handler
type Command struct {
	Flags   CLIFlags
	useCase *usecase.TaskUseCase
}

// NewCommand creates a new Command instance
func NewCommand() *Command {
	cmd := &Command{}
	cmd.setupFlags()
	flag.Parse() // Parse flags immediately after setup
	return cmd
}

// SetUseCase sets the use case for the command
func (c *Command) SetUseCase(uc *usecase.TaskUseCase) {
	c.useCase = uc
}

// setupFlags configures all CLI flags
func (c *Command) setupFlags() {
	// Persistence flags
	flag.BoolVar(&c.Flags.UseJSON, "json", false, "Use JSON file for persistence")
	flag.StringVar(&c.Flags.SQLitePath, "sqlite", "", "Path to SQLite file for persistence")
	flag.StringVar(&c.Flags.FilePath, "file", "tasks.json", "JSON file name (used with --json)")

	// Task operation flags
	flag.StringVar(&c.Flags.Add, "add", "", "Add a new task")
	flag.BoolVar(&c.Flags.List, "list", false, "List all tasks")
	flag.IntVar(&c.Flags.Complete, "complete", 0, "Mark task as complete by ID")
	flag.IntVar(&c.Flags.Delete, "delete", 0, "Delete task by ID")
	flag.IntVar(&c.Flags.Edit, "edit", 0, "Edit a task by ID")
	flag.StringVar(&c.Flags.NewTitle, "title", "", "New title for the task (used with --edit)")
	flag.StringVar(&c.Flags.SetCompleted, "set-completed", "", "Set completed status: true or false (used with --edit)")
	flag.StringVar(&c.Flags.Priority, "priority", "", fmt.Sprintf("Set task priority: %s, %s, %s",
		domain.PriorityLow, domain.PriorityMedium, domain.PriorityHigh))
	flag.StringVar(&c.Flags.Tags, "tags", "", "Set task tags (comma-separated list)")
}

// Execute executes the command based on provided flags
func (c *Command) Execute() error {
	switch {
	case c.Flags.Edit > 0:
		return c.executeEdit()
	case c.Flags.Add != "":
		return c.executeAdd()
	case c.Flags.List:
		return c.executeList()
	case c.Flags.Complete > 0:
		return c.executeComplete()
	case c.Flags.Delete > 0:
		return c.executeDelete()
	default:
		c.printUsage()
		return nil
	}
}

func (c *Command) executeEdit() error {
	// Update completion status if specified
	var completedPtr *bool
	if c.Flags.SetCompleted != "" {
		val := c.Flags.SetCompleted == "true"
		completedPtr = &val
	}

	if err := c.useCase.Edit(c.Flags.Edit, c.Flags.NewTitle, completedPtr); err != nil {
		return err
	}

	// Update priority if specified
	if c.Flags.Priority != "" {
		priority := domain.Priority(c.Flags.Priority)
		if !domain.ValidatePriority(string(priority)) {
			return domain.ErrInvalidPriority
		}
		if err := c.useCase.UpdatePriority(c.Flags.Edit, priority); err != nil {
			return fmt.Errorf("error updating priority: %v", err)
		}
	}

	// Update tags if specified
	if c.Flags.Tags != "" {
		tags := parseTags(c.Flags.Tags)
		if err := c.useCase.UpdateTags(c.Flags.Edit, tags); err != nil {
			return fmt.Errorf("error updating tags: %v", err)
		}
	}

	fmt.Printf("Task %d edited successfully\n", c.Flags.Edit)
	return nil
}

func (c *Command) executeAdd() error {
	// Handle priority
	priority := domain.Priority(c.Flags.Priority)
	if c.Flags.Priority == "" {
		priority = domain.GetDefaultPriority()
	} else if !domain.ValidatePriority(string(priority)) {
		return domain.ErrInvalidPriority
	}

	// Handle tags
	tags := parseTags(c.Flags.Tags)

	if err := c.useCase.Add(c.Flags.Add, priority, tags); err != nil {
		return err
	}
	fmt.Printf("Task added successfully\n")
	return nil
}

func (c *Command) executeList() error {
	tasks, err := c.useCase.List()
	if err != nil {
		return err
	}
	for _, t := range tasks {
		tags := "no tags"
		if len(t.Tags) > 0 {
			tags = strings.Join(t.Tags, ", ")
		}
		fmt.Printf("[%d] %s\n\tStatus: %s\n\tPriority: %s\n\tTags: %s\n",
			t.ID, t.Title, t.GetStatus(), t.Priority, tags)
	}
	return nil
}

func (c *Command) executeComplete() error {
	if err := c.useCase.Complete(c.Flags.Complete); err != nil {
		return err
	}
	fmt.Printf("Task %d marked as completed\n", c.Flags.Complete)
	return nil
}

func (c *Command) executeDelete() error {
	if err := c.useCase.Delete(c.Flags.Delete); err != nil {
		return err
	}
	fmt.Printf("Task %d deleted successfully\n", c.Flags.Delete)
	return nil
}

func (c *Command) printUsage() {
	fmt.Println("Usage:")
	flag.PrintDefaults()
}

// parseTags splits and normalizes a comma-separated list of tags
func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return nil
	}

	var tags []string
	for _, tag := range strings.Split(tagsStr, ",") {
		if trimmed := strings.TrimSpace(tag); trimmed != "" {
			tags = append(tags, trimmed)
		}
	}
	return tags
}

// RunCLI is the main entry point for the CLI
func RunCLI(uc *usecase.TaskUseCase) {
	cmd := NewCommand()
	cmd.SetUseCase(uc)

	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
