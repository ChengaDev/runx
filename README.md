# runx

Run any saved project from anywhere on your system — no `cd` required.

## Quick Start

```bash
# 1. Install
go install github.com/ChengaDev/runx@latest

# 2. Save a project — interactive wizard
runx add

# Or if you're already inside the project directory
runx here

# 3. Run it from anywhere
runx run myapi

# 4. Pass extra flags to the underlying command
runx run myapi -- --port 9000

# 5. See all saved projects
runx list

# 6. Remove one
runx remove myapi
```

## Install

### Homebrew

```bash
brew install ChengaDev/tap/runx
```

### Go

```bash
go install github.com/ChengaDev/runx@latest
```

### Download

Grab a pre-built binary from [Releases](https://github.com/ChengaDev/runx/releases).

---

## Usage

### Add a project

Launch the interactive wizard — browse for the directory and confirm the auto-detected command:

```bash
runx add
```

Or pass arguments directly (supports `.` and relative paths):

```bash
runx add myapi --path ~/projects/myapi
runx add myapi --path .
runx add myapi --path ~/projects/myapi --cmd "go run ."
```

### Save the current directory

From inside a project, register it in one step — name defaults to the directory name:

```bash
cd ~/projects/myapi
runx here
```

**Auto-detection rules:**

| File found      | Command inferred             |
|-----------------|------------------------------|
| `package.json`  | `npm run dev`                |
| `manage.py`     | `python manage.py runserver` |
| `Cargo.toml`    | `cargo run`                  |
| `go.mod`        | `go run .`                   |
| `pyproject.toml`| `python main.py`             |
| `main.py`       | `python main.py`             |

### Run a project

```bash
runx run myapi
```

Pass extra arguments after `--`:

```bash
runx run myapi -- --port 9000
```

### List all projects

```bash
runx list
```

### Edit a project

```bash
runx edit myapi --cmd "go run ./cmd/server"
runx edit myapi --path ~/new/path
```

### Remove a project

```bash
runx remove myapi
```

---

## Storage

Projects are saved in `~/.runx/projects.json`.

Override the directory with the `RUNX_CONFIG_DIR` environment variable:

```bash
export RUNX_CONFIG_DIR=/custom/path
```

---

## Release

Releases are managed with [GoReleaser](https://goreleaser.com). Binaries are built for:

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

---

## License

MIT — see [LICENSE](LICENSE).
