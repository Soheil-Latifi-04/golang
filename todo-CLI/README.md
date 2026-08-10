# todo-CLI

A simple command-line to-do list manager written in Go. Todos are stored locally as JSON by default, with an optional PostgreSQL storage backend available in the code.

## Features

- Add, edit, toggle, delete, and list todos from the command line
- Tracks creation and completion timestamps for each todo
- Pretty table output in the terminal
- Pluggable storage layer (`Storage[T]` interface) with two implementations:
  - **JSON file storage** (used by default, `todos.json`)
  - **PostgreSQL storage** (via GORM, currently commented out in `main.go`)

## Requirements

- Go 1.26.5 or later

## Installation

Clone the repository and move into the project directory:

```bash
git clone https://github.com/Soheil-Latifi-04/golang.git
cd golang/todo-CLI
```

Download dependencies:

```bash
go mod tidy
```

Build the binary:

```bash
go build -o todo
```

## Usage

Run directly with `go run .`, or use the built binary (`./todo`) after building.

```bash
go run . -add "Buy groceries"
go run . -list
go run . -toggle 0
go run . -edit "0:Buy groceries and cook dinner"
go run . -del 0
```

### Flags

| Flag      | Description                                              | Example                          |
|-----------|------------------------------------------------------------|-----------------------------------|
| `-add`    | Add a new todo with the given title                        | `-add "Write README"`            |
| `-list`   | List all todos in a table                                  | `-list`                          |
| `-toggle` | Toggle the completed status of the todo at the given index | `-toggle 2`                      |
| `-edit`   | Edit the title of the todo at the given index (`id:title`) | `-edit "1:New title"`            |
| `-del`    | Delete the todo at the given index                          | `-del 3`                         |

Indexes are zero-based and correspond to the row numbers shown by `-list`.

Todos are automatically saved to `todos.json` in the working directory after every run.

## Storage backends

By default, the app uses `JSONStorage`, which persists todos to a local `todos.json` file — no setup required.

A `PostgresStorage` implementation is also included (see `postgresql-storage.go`) but is disabled by default in `main.go`. To use it instead:

1. Set a `DATABASE_DSN` environment variable with your PostgreSQL connection string (a `.env` file is supported via `godotenv`).
2. In `main.go`, swap the `JSONStorage` initialization for the commented-out `PostgresStorage` block.

## Project structure

```
todo-CLI/
├── main.go                  # Entry point: loads storage, parses flags, saves state
├── command.go                # CLI flag definitions and command dispatch
├── todo.go                   # Todo/Todos types and core operations (add, edit, toggle, delete, print)
├── storage.go                 # Generic Storage[T] interface
├── json-storage.go            # JSON file storage implementation
├── postgresql-storage.go      # PostgreSQL storage implementation (GORM)
├── cofig.go                   # Loads config (e.g. DATABASE_DSN) from environment/.env
├── go.mod / go.sum            # Go module definition and dependencies
```

## Dependencies

- [aquasecurity/table](https://github.com/aquasecurity/table) – terminal table rendering
- [joho/godotenv](https://github.com/joho/godotenv) – `.env` file loading
- [jackc/pgx](https://github.com/jackc/pgx) – PostgreSQL driver (used via GORM)
- [gorm.io/gorm](https://gorm.io/) + [gorm.io/driver/postgres](https://gorm.io/) – ORM for the PostgreSQL storage backend
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

## License

No license file is currently included in this repository. Check with the repository owner before reuse.
