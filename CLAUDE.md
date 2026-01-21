# CLAUDE.md

This file provides guidance for AI agents working with the gmail-ro codebase.

## Project Overview

gmail-ro is a **read-only** command-line interface for Gmail written in Go. It uses OAuth2 for authentication and only requests the `gmail.readonly` scope - no write, send, or delete operations are possible.

## Quick Commands

```bash
# Build
make build

# Run tests
make test

# Run tests with coverage
make test-cover

# Lint
make lint

# Format code
make fmt

# All checks (format, lint, test)
make verify

# Install locally
make install

# Clean build artifacts
make clean
```

## Architecture

```
gmail-ro/
├── main.go                     # Entry point
├── cmd/
│   ├── root.go                 # Root command, version command
│   ├── search.go               # Search messages command
│   ├── read.go                 # Read single message command
│   ├── thread.go               # Read conversation thread command
│   ├── attachments.go          # Attachments parent command
│   ├── attachments_list.go     # List attachments subcommand
│   └── attachments_download.go # Download attachments subcommand
├── internal/
│   ├── gmail/
│   │   ├── client.go           # OAuth2 client, authentication
│   │   ├── client_test.go      # Client tests
│   │   ├── messages.go         # Message/Attachment structs, parsing
│   │   ├── messages_test.go    # Message parsing tests
│   │   ├── attachments.go      # Attachment download methods
│   │   └── attachments_test.go # Attachment tests
│   └── zip/
│       ├── extract.go          # Secure zip extraction
│       └── extract_test.go     # Zip extraction tests
├── .github/workflows/
│   ├── ci.yml                  # Lint and test on PR/push
│   ├── auto-release.yml        # Create tags on main push
│   └── release.yml             # Build and release binaries
├── Makefile                    # Build, test, lint targets
└── go.mod                      # Module: github.com/open-cli-collective/gmail-ro
```

## Key Patterns

### Read-Only by Design

This CLI intentionally only supports read operations:
- Uses `gmail.GmailReadonlyScope` exclusively
- Only calls `.List()` and `.Get()` Gmail API methods
- No `.Send()`, `.Delete()`, `.Modify()`, or `.Trash()` operations

### OAuth2 Configuration

Credentials are stored in `~/.config/gmail-ro/`:
- `credentials.json` - OAuth client credentials (from Google Cloud Console)
- `token.json` - OAuth access/refresh token (created on first auth)

### Command Patterns

All commands follow this structure:

```go
var searchCmd = &cobra.Command{
    Use:   "search <query>",
    Short: "Search for messages",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        ctx := context.Background()
        client, err := gmail.NewClient(ctx)
        if err != nil {
            return err
        }
        // ... use client
    },
}
```

### Output Formats

Commands support two output modes:
- **Text** (default): Human-readable formatted output
- **JSON** (`--json`): Machine-readable JSON for scripting

```go
if jsonOutput {
    enc := json.NewEncoder(os.Stdout)
    enc.SetIndent("", "  ")
    return enc.Encode(messages)
}
// ... text output
```

## Testing

Tests use `testify` for assertions and table-driven test patterns:

```go
func TestParseMessage(t *testing.T) {
    tests := []struct {
        name     string
        input    *gmail.Message
        expected *Message
    }{
        {"basic message", ...},
        {"multipart message", ...},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := parseMessage(tt.input, true)
            assert.Equal(t, tt.expected.Subject, result.Subject)
        })
    }
}
```

Run tests: `make test`

Coverage report: `make test-cover && open coverage.html`

## Adding a New Command

1. Create new file in `cmd/` (e.g., `cmd/labels.go`)
2. Define the command with `&cobra.Command{}`
3. Register in `init()` with `rootCmd.AddCommand()`
4. Add flags if needed
5. Write tests in `cmd/<name>_test.go`

Example:

```go
func init() {
    rootCmd.AddCommand(labelsCmd)
    labelsCmd.Flags().BoolVarP(&labelsJSON, "json", "j", false, "Output as JSON")
}

var labelsCmd = &cobra.Command{
    Use:   "labels",
    Short: "List Gmail labels",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation - READ ONLY operations
    },
}
```

## Gmail API Notes

### Search Query Syntax

The search command accepts standard Gmail search syntax:

| Operator | Example | Description |
|----------|---------|-------------|
| `from:` | `from:alice@example.com` | Messages from sender |
| `to:` | `to:bob@example.com` | Messages to recipient |
| `subject:` | `subject:meeting` | Subject contains word |
| `is:` | `is:unread` | Message state |
| `has:` | `has:attachment` | Has attachment |
| `after:` | `after:2024/01/01` | After date |
| `before:` | `before:2024/02/01` | Before date |

### Message Format

Gmail API returns messages in different formats:
- `metadata` - Headers only (faster, used for search results)
- `full` - Complete message with body (used for read/thread)

Body content is base64url encoded and may be nested in multipart structures.

### Attachment Handling

Attachments are identified by:
- `part.Filename != ""` - Primary detection
- `Content-Disposition: attachment` header - Secondary detection

For downloading:
- Small attachments may have inline data in `part.Body.Data`
- Large attachments require separate API call via `part.Body.AttachmentId`

```bash
# List attachments in a message
gmail-ro attachments list <message-id>
gmail-ro attachments list <message-id> --json

# Download attachments
gmail-ro attachments download <message-id> --all
gmail-ro attachments download <message-id> --filename report.pdf
gmail-ro attachments download <message-id> --all --output ~/Downloads

# Download and extract zip files
gmail-ro attachments download <message-id> --filename data.zip --extract
```

Zip extraction includes security safeguards:
- 100MB per-file limit, 500MB total limit
- Path traversal prevention
- Max 1000 files, 10 levels depth

## Common Issues

### "Unable to read credentials file"

Ensure OAuth credentials are set up:
```bash
mkdir -p ~/.config/gmail-ro
# Download credentials.json from Google Cloud Console
mv ~/Downloads/client_secret_*.json ~/.config/gmail-ro/credentials.json
```

### "Token has been expired or revoked"

Delete token and re-authenticate:
```bash
rm ~/.config/gmail-ro/token.json
gmail-ro search "test"  # Will prompt for re-auth
```

### "Access blocked: This app's request is invalid"

Check Google Cloud Console:
1. Gmail API is enabled
2. OAuth consent screen configured
3. Your email added as test user
4. `gmail.readonly` scope added

## Dependencies

Key dependencies:
- `github.com/spf13/cobra` - CLI framework
- `golang.org/x/oauth2` - OAuth2 client
- `google.golang.org/api/gmail/v1` - Gmail API client
- `github.com/stretchr/testify` - Testing assertions (dev)

## Security

- **Read-only scope**: Cannot modify, send, or delete emails
- **Local token storage**: OAuth tokens stored with 0600 permissions
- **No credential exposure**: Credentials never logged or transmitted

## Commit Conventions

Use conventional commits:

```
type(scope): description

feat(search): add date range filtering
fix(thread): handle missing message bodies
docs(readme): add installation instructions
```

| Prefix | Purpose | Triggers Release? |
|--------|---------|-------------------|
| `feat:` | New features | Yes |
| `fix:` | Bug fixes | Yes |
| `docs:` | Documentation only | No |
| `test:` | Adding/updating tests | No |
| `refactor:` | Code changes that don't fix bugs or add features | No |
| `chore:` | Maintenance tasks | No |
| `ci:` | CI/CD changes | No |

## CI & Release Workflow

Releases are automated with a dual-gate system to avoid unnecessary releases:

**Gate 1 - Path filter:** Only triggers when Go code changes (`**.go`, `go.mod`, `go.sum`)
**Gate 2 - Commit prefix:** Only `feat:` and `fix:` commits create releases

This means:
- `feat: add command` + Go files changed → release
- `fix: handle edge case` + Go files changed → release
- `docs:`, `ci:`, `test:`, `refactor:` → no release
- Changes only to docs, packaging, workflows → no release

**After merging a release-triggering PR:** The workflow creates a tag, which triggers GoReleaser to build binaries and publish to Homebrew. Chocolatey and Winget require manual workflow dispatch.
