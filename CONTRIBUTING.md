# Contributing to runx

Thanks for your interest in contributing!

## Getting started

1. Fork the repo and clone your fork.
2. Install Go 1.21+.
3. Install dependencies:

```bash
go mod tidy
```

4. Build and verify:

```bash
go build ./...
go test ./...
```

## Making changes

- Open an issue before starting significant work.
- Keep changes focused — one feature or fix per PR.
- Add or update tests for any changed behaviour.
- Run the linter before submitting:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run
```

## Submitting a PR

1. Push your branch and open a pull request against `main`.
2. Fill in the PR template.
3. CI must pass before merging.

## Project layout

```
cmd/          CLI commands (Cobra)
internal/
  store/      JSON project storage
  detect/     Auto-detection of run commands
```

## Commit style

Use short, present-tense messages:

```
add shell auto-detection for Bun projects
fix list command sorting when names share a prefix
```

## License

By contributing you agree your changes will be licensed under the [MIT License](LICENSE).
