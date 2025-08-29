### Edit Tasks (Title and Completed Status)

You can edit a task's title and/or its completed status:

```sh
# Change the title of task with ID 2
go run cmd/todo/main.go --edit 2 --title "New title for task 2"

# Mark task 2 as completed
go run cmd/todo/main.go --edit 2 --set-completed true

# Change both title and completed status
go run cmd/todo/main.go --edit 2 --title "Read Go book" --set-completed true

# Mark task 2 as not completed
go run cmd/todo/main.go --edit 2 --set-completed false
```


# To-Do List Console App (Go)

![Go Version](https://img.shields.io/badge/Go-1.18%2B-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Build](https://img.shields.io/badge/build-passing-brightgreen)

## Advanced Usage Examples

### Use a custom JSON file
```sh
go run cmd/todo/main.go --json --file mytasks.json --add "Plan vacation"
go run cmd/todo/main.go --json --file mytasks.json --list
```

### Use a custom SQLite database and combine with other flags
```sh
go run cmd/todo/main.go --sqlite /tmp/mytasks.db --add "Finish report"
go run cmd/todo/main.go --sqlite /tmp/mytasks.db --list
go run cmd/todo/main.go --sqlite /tmp/mytasks.db --complete 2
go run cmd/todo/main.go --sqlite /tmp/mytasks.db --delete 2
```

### Chain commands (Unix style)
```sh
# Add several tasks and list them
go run cmd/todo/main.go --json --add "Task 1"
go run cmd/todo/main.go --json --add "Task 2"
go run cmd/todo/main.go --json --list | grep false
```

### Run as a compiled binary
```sh
go build -o todo cmd/todo/main.go
./todo --sqlite tasks.db --add "Backup files"
./todo --sqlite tasks.db --list
```

A clean architecture console application in Go for managing your tasks (To-Do List), supporting both JSON file and SQLite database persistence, with robust flag-based CLI interaction.

## Features
- Full CRUD for tasks (create, list, complete, delete)
- Choose between JSON file or SQLite database for persistence
- Clean Architecture: modular, testable, and maintainable code
- Simple and intuitive CLI with flags
- Easily extensible for new features

## Project Structure
```
cmd/todo/main.go              # CLI entry point
internal/domain/task.go       # Task entity definition
internal/usecase/task.go      # Use cases (business logic)
internal/repository/json.go   # JSON persistence implementation
internal/repository/sqlite.go # SQLite persistence implementation
internal/handler/cli.go       # CLI flag handling and user interaction
```

## Requirements
- Go 1.18 or higher
- (For SQLite) gcc or build-essential (for cgo)

## Getting Started
Clone the repository and install dependencies:
```sh
git clone https://github.com/carloscapella/TodoConsole.git
cd TodoConsole
go mod tidy
```

## Usage

### JSON Persistence
```sh
go run cmd/todo/main.go --json --add "Buy groceries"
go run cmd/todo/main.go --json --list
```

### SQLite Persistence
```sh
go run cmd/todo/main.go --sqlite tasks.db --add "Buy groceries"
go run cmd/todo/main.go --sqlite tasks.db --list
```

### Available Flags
- `--json`                  Use JSON file (tasks.json) for persistence
- `--sqlite <file>`         Use SQLite database for persistence
- `--add <task>`            Add a new task
- `--list`                  List all tasks
- `--complete <id>`         Mark a task as completed
- `--delete <id>`           Delete a task by ID

### Example Workflow
```sh
# Add a task
go run cmd/todo/main.go --json --add "Read a book"
# List tasks
go run cmd/todo/main.go --json --list
# Complete a task
go run cmd/todo/main.go --json --complete 1
# Delete a task
go run cmd/todo/main.go --json --delete 1
```

## Best Practices
- Follows Clean Architecture for separation of concerns
- Uses interfaces for repository abstraction
- Handles errors gracefully and provides user feedback
- Modular code for easy testing and extension

## Possible Extensions
- Unit and integration tests
- Edit/update task functionality
- Import/export tasks
- REST API or web interface
- User authentication

## Contributing
Contributions, issues, and feature requests are welcome! Feel free to fork the repo and submit a pull request.

## Author
Carlos Capella

---
Open to feedback and improvements!
