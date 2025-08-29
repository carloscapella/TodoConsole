package handler

import (
	"flag"
	"fmt"
	"os"
	"todo/internal/usecase"
)

// Command representa las opciones y flags de la CLI
type Command struct {
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

	useCase *usecase.TaskUseCase
}

// NewCommand crea una nueva instancia de Command
func NewCommand() *Command {
	cmd := &Command{}
	cmd.setupFlags()
	flag.Parse() // Parse flags immediately after setup
	return cmd
}

// SetUseCase establece el caso de uso para el comando
func (c *Command) SetUseCase(uc *usecase.TaskUseCase) {
	c.useCase = uc
}

// setupFlags configura todos los flags de la CLI
func (c *Command) setupFlags() {
	// Persistence flags
	flag.BoolVar(&c.UseJSON, "json", false, "Use JSON file for persistence")
	flag.StringVar(&c.SQLitePath, "sqlite", "", "Path to SQLite file for persistence")
	flag.StringVar(&c.FilePath, "file", "tasks.json", "JSON file name (used with --json)")

	// Task operation flags
	flag.StringVar(&c.Add, "add", "", "Add a new task")
	flag.BoolVar(&c.List, "list", false, "List all tasks")
	flag.IntVar(&c.Complete, "complete", 0, "Mark task as complete by ID")
	flag.IntVar(&c.Delete, "delete", 0, "Delete task by ID")
	flag.IntVar(&c.Edit, "edit", 0, "Edit a task by ID")
	flag.StringVar(&c.NewTitle, "title", "", "New title for the task (used with --edit)")
	flag.StringVar(&c.SetCompleted, "set-completed", "", "Set completed status: true or false (used with --edit)")
}

// Execute ejecuta el comando según los flags proporcionados
func (c *Command) Execute() error {
	switch {
	case c.Edit > 0:
		return c.executeEdit()
	case c.Add != "":
		return c.executeAdd()
	case c.List:
		return c.executeList()
	case c.Complete > 0:
		return c.executeComplete()
	case c.Delete > 0:
		return c.executeDelete()
	default:
		c.printUsage()
		return nil
	}
}

func (c *Command) executeEdit() error {
	var completedPtr *bool
	if c.SetCompleted != "" {
		val := c.SetCompleted == "true"
		completedPtr = &val
	}

	if err := c.useCase.Edit(c.Edit, c.NewTitle, completedPtr); err != nil {
		return err
	}
	fmt.Printf("Task %d edited successfully\n", c.Edit)
	return nil
}

func (c *Command) executeAdd() error {
	if err := c.useCase.Add(c.Add); err != nil {
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
		fmt.Printf("[%d] %s (done: %v)\n", t.ID, t.Title, t.Completed)
	}
	return nil
}

func (c *Command) executeComplete() error {
	if err := c.useCase.Complete(c.Complete); err != nil {
		return err
	}
	fmt.Printf("Task %d marked as completed\n", c.Complete)
	return nil
}

func (c *Command) executeDelete() error {
	if err := c.useCase.Delete(c.Delete); err != nil {
		return err
	}
	fmt.Printf("Task %d deleted successfully\n", c.Delete)
	return nil
}

func (c *Command) printUsage() {
	fmt.Println("Usage:")
	flag.PrintDefaults()
}

// RunCLI es el punto de entrada principal para la CLI
func RunCLI(uc *usecase.TaskUseCase) {
	cmd := NewCommand()
	cmd.SetUseCase(uc)

	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
