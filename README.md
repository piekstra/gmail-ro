# gmail-ro

A read-only command-line interface for Gmail. Search, read, and view email threads without any ability to modify, send, or delete messages.

## Features

- **Read-only access** - Uses `gmail.readonly` OAuth scope exclusively
- **Search messages** - Full Gmail search syntax support
- **Read messages** - View complete message content
- **View threads** - Read entire conversation threads
- **JSON output** - Machine-readable output for scripting

## Installation

### Homebrew (macOS/Linux)

```bash
brew install piekstra/tap/gmail-ro
```

### Download Binary

Download the latest release for your platform from the [Releases page](https://github.com/piekstra/gmail-ro/releases).

### Build from Source

```bash
go install github.com/piekstra/gmail-ro@latest
```

## Setup

### 1. Create Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Gmail API:
   - Go to **APIs & Services** > **Library**
   - Search for "Gmail API"
   - Click **Enable**

### 2. Create OAuth Credentials

1. Go to **APIs & Services** > **Credentials**
2. Click **Create Credentials** > **OAuth client ID**
3. If prompted, configure the OAuth consent screen:
   - Choose **External** user type
   - Fill in required fields (app name, support email)
   - Add scope: `https://www.googleapis.com/auth/gmail.readonly`
   - Add your email as a test user
4. For Application type, select **Desktop app**
5. Click **Create**
6. Download the JSON file

### 3. Configure gmail-ro

1. Create the config directory:
   ```bash
   mkdir -p ~/.config/gmail-ro
   ```

2. Move the downloaded credentials file:
   ```bash
   mv ~/Downloads/client_secret_*.json ~/.config/gmail-ro/credentials.json
   ```

### 4. Authenticate

Run any command to trigger the OAuth flow:

```bash
gmail-ro search "is:unread"
```

1. A URL will be displayed - open it in your browser
2. Sign in with your Google account
3. Grant read-only access to Gmail
4. Copy the authorization code
5. Paste it back into the terminal

Your token will be saved to `~/.config/gmail-ro/token.json` for future use.

## Usage

### Search Messages

```bash
# Basic search
gmail-ro search "from:someone@example.com"

# Limit results
gmail-ro search "subject:meeting" --max 20

# JSON output
gmail-ro search "is:unread" --json

# Date range
gmail-ro search "after:2024/01/01 before:2024/02/01"

# Combined queries
gmail-ro search "from:alice@example.com subject:project has:attachment"
```

### Read a Message

```bash
# Get message by ID (from search results)
gmail-ro read 18abc123def456

# JSON output
gmail-ro read 18abc123def456 --json
```

### View Thread

```bash
# View entire conversation
gmail-ro thread 18abc123def456

# JSON output
gmail-ro thread 18abc123def456 --json
```

### Search Query Reference

gmail-ro supports all Gmail search operators:

| Operator | Example | Description |
|----------|---------|-------------|
| `from:` | `from:alice@example.com` | Messages from sender |
| `to:` | `to:bob@example.com` | Messages to recipient |
| `subject:` | `subject:meeting` | Subject contains word |
| `is:` | `is:unread`, `is:starred` | Message state |
| `has:` | `has:attachment` | Has attachment |
| `after:` | `after:2024/01/01` | After date |
| `before:` | `before:2024/02/01` | Before date |
| `label:` | `label:work` | Has label |
| `in:` | `in:inbox`, `in:sent` | In folder |
| `larger:` | `larger:5M` | Larger than size |
| `smaller:` | `smaller:1M` | Smaller than size |

See [Gmail search operators](https://support.google.com/mail/answer/7190) for the complete list.

## Shell Completion

gmail-ro supports tab completion for bash, zsh, fish, and PowerShell.

### Bash

```bash
# Load in current session
source <(gmail-ro completion bash)

# Install permanently (Linux)
gmail-ro completion bash | sudo tee /etc/bash_completion.d/gmail-ro > /dev/null

# Install permanently (macOS with Homebrew)
gmail-ro completion bash > $(brew --prefix)/etc/bash_completion.d/gmail-ro
```

### Zsh

```bash
# Load in current session
source <(gmail-ro completion zsh)

# Install permanently
mkdir -p ~/.zsh/completions
gmail-ro completion zsh > ~/.zsh/completions/_gmail-ro

# Add to ~/.zshrc if not already present:
# fpath=(~/.zsh/completions $fpath)
# autoload -Uz compinit && compinit
```

### Fish

```bash
# Load in current session
gmail-ro completion fish | source

# Install permanently
gmail-ro completion fish > ~/.config/fish/completions/gmail-ro.fish
```

### PowerShell

```powershell
# Load in current session
gmail-ro completion powershell | Out-String | Invoke-Expression

# Install permanently (add to $PROFILE)
gmail-ro completion powershell >> $PROFILE
```

## Configuration

Configuration files are stored in `~/.config/gmail-ro/`:

| File | Description |
|------|-------------|
| `credentials.json` | OAuth client credentials (from Google Cloud Console) |
| `token.json` | OAuth access/refresh token (created automatically) |

## Security

- This tool only requests **read-only** access to Gmail
- No write, send, or delete operations are possible
- OAuth tokens are stored locally with `0600` permissions
- Credentials never leave your machine

## Troubleshooting

### "Unable to read credentials file"

Ensure `credentials.json` exists:
```bash
ls -la ~/.config/gmail-ro/credentials.json
```

### "Token has been expired or revoked"

Delete the token file and re-authenticate:
```bash
rm ~/.config/gmail-ro/token.json
gmail-ro search "test"
```

### "Access blocked: This app's request is invalid"

Your OAuth consent screen may not be properly configured. Ensure:
1. The Gmail API is enabled
2. Your email is added as a test user (for apps in testing mode)
3. The `gmail.readonly` scope is added

## License

MIT License - see [LICENSE](LICENSE) for details.
