# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build ./...

# Run tests (with race detector)
go test ./... -race -count=1

# Run a single test
go test ./internal/store/... -run TestStoreName -race

# Lint (requires golangci-lint)
golangci-lint run

# Build and run the CLI locally
go run . <subcommand>
```

## Architecture

`runx` is a Go CLI tool built with [Cobra](https://github.com/spf13/cobra). Each subcommand (`add`, `run`, `list`, `edit`, `remove`) lives as its own file in `cmd/` and is registered in `cmd/root.go`.

**Data flow:**
1. A command handler in `cmd/` calls `store.New()` to load the project store
2. It reads/mutates the store via the `Store` methods (`Add`, `Get`, `Remove`, `List`)
3. Every mutation immediately persists to `~/.runx/projects.json` via `store.save()`

**Key packages:**
- `internal/store` — JSON-backed persistence. The store is a `map[string]Project` keyed by project name. Config directory defaults to `~/.runx` and can be overridden with `RUNX_CONFIG_DIR`.
- `internal/detect` — Heuristic run-command detection. Checks for indicator files (`go.mod`, `package.json`, etc.) in order of priority and returns the first match.
- `cmd/add.go` — Has two modes: interactive wizard (no args, uses `charmbracelet/huh`) and non-interactive (name arg + `--path` flag).

**Shell execution:** `cmd/run.go` runs commands via `sh -c` (or `cmd /c` on Windows). Extra args passed after `--` are shell-quoted via `shellQuote()` before being appended to the command string to prevent injection.

## Release

Releases are triggered by pushing a `v*` tag. GoReleaser (`.goreleaser.yaml`) builds cross-platform binaries and publishes a Homebrew formula to `ChengaDev/homebrew-tap`. Requires `HOMEBREW_TAP_GITHUB_TOKEN` secret in GitHub.
