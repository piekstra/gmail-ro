# Contributing to gmail-ro

Thank you for your interest in contributing to gmail-ro!

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/open-cli-collective/gmail-ro.git
   cd gmail-ro
   ```

2. Install dependencies:
   ```bash
   make deps
   ```

3. Build:
   ```bash
   make build
   ```

4. Run tests:
   ```bash
   make test
   ```

## Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover

# Run quick tests (skip slow ones)
make test-short
```

## Code Style

- Run `gofmt` and `goimports` before committing
- Run the linter: `make lint`
- Follow Go conventions and idioms

## Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new feature
fix: fix a bug
docs: update documentation
test: add tests
refactor: refactor code
ci: update CI configuration
chore: maintenance tasks
```

Examples:
```
feat: add labels list command
fix: handle expired OAuth token gracefully
docs: update installation instructions
```

## Important: Read-Only Design

This CLI is intentionally **read-only**. Do not add features that:
- Send emails
- Delete emails
- Modify emails or labels
- Perform any write operations

The `gmail.readonly` scope is fundamental to this project's design.

## Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make your changes
4. Run `make verify`
5. Commit with a conventional commit message
6. Push and create a pull request

## Project Structure

```
gmail-ro/
├── main.go               # Entry point
├── cmd/                  # Command implementations
│   ├── root.go           # Root command, version
│   ├── search.go         # Search messages
│   ├── read.go           # Read single message
│   └── thread.go         # Read conversation thread
├── internal/
│   └── gmail/            # Gmail API client
│       ├── client.go     # OAuth2 client, authentication
│       └── messages.go   # Message parsing, API operations
└── .github/              # GitHub workflows and templates
```

## Questions?

Open an issue or start a discussion on GitHub.
