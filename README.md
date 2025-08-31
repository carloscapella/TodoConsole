
# To-Do List Console App (Go)

![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Build](https://img.shields.io/badge/build-passing-brightgreen)

## Features
- Full CRUD for tasks (create, list, edit, complete, delete)
- Priority support: low, medium (default), high
- Tags support: add, edit, and list tags for each task
- Edit tasks: change title, completed status, priority, and tags
- Choose between JSON file or SQLite database for persistence
- Clean Architecture: modular, testable, and maintainable code
- Simple and intuitive CLI with flags
- Easily extensible for new features

## Project Structure
```
cmd/todo/main.go              # CLI entry point
internal/domain/task.go       # Task entity definition and business rules
internal/usecase/task.go      # Use cases (business logic)
internal/repository/json.go   # JSON persistence implementation
internal/repository/sqlite.go # SQLite persistence implementation
internal/handler/cli.go       # CLI flag handling and user interaction
```

## Requirements
- Go 1.22 or higher
- (For SQLite) gcc or build-essential (for cgo)

## Getting Started
Clone the repository and install dependencies:
```sh
git clone https://github.com/carloscapella/TodoConsole.git
cd TodoConsole
go mod tidy
```

## Usage

### Add a Task
```sh
# Add a task with default priority (medium)
go run cmd/todo/main.go --json --add "Read a book"
# Add a task with high priority and tags
go run cmd/todo/main.go --sqlite tasks.db --add "Finish report" --priority high --tags work,urgent
```

### List Tasks
```sh
go run cmd/todo/main.go --json --list
go run cmd/todo/main.go --sqlite tasks.db --list
```

### Edit Tasks
```sh
# Change the title of task with ID 2
go run cmd/todo/main.go --edit 2 --title "New title for task 2"
# Mark task 2 as completed
go run cmd/todo/main.go --edit 2 --set-completed true
# Change both title and completed status
go run cmd/todo/main.go --edit 2 --title "Read Go book" --set-completed true
# Change priority and tags
go run cmd/todo/main.go --edit 2 --priority low --tags "home,reading"
```

### Complete or Delete Tasks
```sh
go run cmd/todo/main.go --complete 1
go run cmd/todo/main.go --delete 1
```

### Use a custom JSON file
```sh
go run cmd/todo/main.go --json --file mytasks.json --add "Plan vacation"
go run cmd/todo/main.go --json --file mytasks.json --list
```

### Use a custom SQLite database
```sh
go run cmd/todo/main.go --sqlite /tmp/mytasks.db --add "Finish report"
go run cmd/todo/main.go --sqlite /tmp/mytasks.db --list
```

### Run as a compiled binary
```sh
go build -o todo cmd/todo/main.go
./todo --sqlite tasks.db --add "Backup files"
./todo --sqlite tasks.db --list
```

## Available Flags
- `--json`                  Use JSON file (tasks.json) for persistence
- `--sqlite <file>`         Use SQLite database for persistence
- `--add <task>`            Add a new task
- `--list`                  List all tasks
- `--edit <id>`             Edit a task by ID (use with --title, --set-completed, --priority, --tags)
- `--title <title>`         New title for the task (used with --edit)
- `--set-completed <bool>`  Set completed status: true or false (used with --edit)
- `--priority <level>`      Set task priority: low, medium, high
- `--tags <tags>`           Set task tags (comma-separated list)
- `--complete <id>`         Mark a task as completed
- `--delete <id>`           Delete a task by ID
- `--file <file>`           JSON file name (used with --json)

## Best Practices
- Follows Clean Architecture for separation of concerns
- Uses interfaces for repository abstraction
- Handles errors gracefully and provides user feedback
- Modular code for easy testing and extension

**Performance Note:**
If you use the JSON file for persistence, keep in mind that every create, update, or delete operation rewrites the entire file. This is fine for a small number of tasks, but performance may degrade with thousands of tasks. For better scalability, use the SQLite backend.

## Possible Extensions
- Unit and integration tests
- Import/export tasks
- REST API or web interface
- User authentication

## Contributing
Contributions, issues, and feature requests are welcome! Feel free to fork the repo and submit a pull request.

## Author
Carlos Capella

---
Open to feedback and improvements!
