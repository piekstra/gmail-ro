# gmro

A read-only command-line interface for Gmail. Search, read, and view email threads without any ability to modify, send, or delete messages.

## Features

- **Read-only access** - Uses `gmail.readonly` OAuth scope exclusively
- **Search messages** - Full Gmail search syntax support
- **Read messages** - View complete message content
- **View threads** - Read entire conversation threads
- **JSON output** - Machine-readable output for scripting

## Installation

### macOS

**Homebrew (recommended)**

```bash
brew install open-cli-collective/tap/gmail-readonly
```

> Note: This installs from our third-party tap.

---

### Windows

**Chocolatey**

```powershell
choco install gmail-readonly
```

**Winget**

```powershell
winget install OpenCLICollective.gmail-readonly
```

---

### Linux

**Snap**

```bash
sudo snap install ocli-gmail-ro
```

> Note: After installation, the command is available as `gmro`.

**APT (Debian/Ubuntu)**

```bash
# Add the GPG key
curl -fsSL https://open-cli-collective.github.io/linux-packages/keys/gpg.asc | sudo gpg --dearmor -o /usr/share/keyrings/open-cli-collective.gpg

# Add the repository
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/open-cli-collective.gpg] https://open-cli-collective.github.io/linux-packages/apt stable main" | sudo tee /etc/apt/sources.list.d/open-cli-collective.list

# Install
sudo apt update
sudo apt install gmro
```

> Note: This is our third-party APT repository, not official Debian/Ubuntu repos.

**DNF/YUM (Fedora/RHEL/CentOS)**

```bash
# Add the repository
sudo tee /etc/yum.repos.d/open-cli-collective.repo << 'EOF'
[open-cli-collective]
name=Open CLI Collective
baseurl=https://open-cli-collective.github.io/linux-packages/rpm
enabled=1
gpgcheck=1
gpgkey=https://open-cli-collective.github.io/linux-packages/keys/gpg.asc
EOF

# Install
sudo dnf install gmro
```

> Note: This is our third-party RPM repository, not official Fedora/RHEL repos.

**Binary download**

Download `.deb`, `.rpm`, or `.tar.gz` from the [Releases page](https://github.com/open-cli-collective/gmail-ro/releases).

---

### From Source

```bash
go install github.com/open-cli-collective/gmail-ro@latest
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

### 3. Configure gmro

1. Create the config directory:
   ```bash
   mkdir -p ~/.config/gmail-readonly
   ```

2. Move the downloaded credentials file:
   ```bash
   mv ~/Downloads/client_secret_*.json ~/.config/gmail-readonly/credentials.json
   ```

### 4. Authenticate

Run any command to trigger the OAuth flow:

```bash
gmro search "is:unread"
```

1. A URL will be displayed - open it in your browser
2. Sign in with your Google account
3. Grant read-only access to Gmail
4. Copy the authorization code
5. Paste it back into the terminal

Your token will be saved securely (system keychain on macOS/Linux, or `~/.config/gmail-readonly/token.json` as fallback).

## Usage

### Search Messages

```bash
# Basic search
gmro search "from:someone@example.com"

# Limit results
gmro search "subject:meeting" --max 20

# JSON output
gmro search "is:unread" --json

# Date range
gmro search "after:2024/01/01 before:2024/02/01"

# Combined queries
gmro search "from:alice@example.com subject:project has:attachment"
```

Search results include both `ID` (message ID) and `ThreadID` (thread ID). Either can be
used with the `thread` command.

### Read a Message

```bash
# Get message by ID (from search results)
gmro read 18abc123def456

# JSON output
gmro read 18abc123def456 --json
```

### View Thread

```bash
# View entire conversation (accepts message ID or thread ID)
gmro thread 18abc123def456

# JSON output
gmro thread 18abc123def456 --json
```

The thread command accepts either a message ID or thread ID. If you pass a message ID
(from search results), it automatically resolves to the containing thread.

### Version

```bash
gmro version
```

### Search Query Reference

gmro supports all Gmail search operators:

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

gmro supports tab completion for bash, zsh, fish, and PowerShell.

### Bash

```bash
# Load in current session
source <(gmro completion bash)

# Install permanently (Linux)
gmro completion bash | sudo tee /etc/bash_completion.d/gmro > /dev/null

# Install permanently (macOS with Homebrew)
gmro completion bash > $(brew --prefix)/etc/bash_completion.d/gmro
```

### Zsh

```bash
# Load in current session
source <(gmro completion zsh)

# Install permanently
mkdir -p ~/.zsh/completions
gmro completion zsh > ~/.zsh/completions/_gmro

# Add to ~/.zshrc if not already present:
# fpath=(~/.zsh/completions $fpath)
# autoload -Uz compinit && compinit
```

### Fish

```bash
# Load in current session
gmro completion fish | source

# Install permanently
gmro completion fish > ~/.config/fish/completions/gmro.fish
```

### PowerShell

```powershell
# Load in current session
gmro completion powershell | Out-String | Invoke-Expression

# Install permanently (add to $PROFILE)
gmro completion powershell >> $PROFILE
```

## Configuration

Configuration files are stored in `~/.config/gmail-readonly/`:

| File | Description |
|------|-------------|
| `credentials.json` | OAuth client credentials (from Google Cloud Console) |
| `token.json` | OAuth access/refresh token (fallback if keychain unavailable) |

## Security

- This tool only requests **read-only** access to Gmail
- No write, send, or delete operations are possible
- OAuth tokens are stored in system keychain (macOS Keychain / Linux secret-tool) when available
- File-based storage uses `0600` permissions
- Credentials never leave your machine

## Troubleshooting

### "Unable to read credentials file"

Ensure `credentials.json` exists:
```bash
ls -la ~/.config/gmail-readonly/credentials.json
```

### "Token has been expired or revoked"

Clear the token and re-authenticate:
```bash
gmro config clear
gmro search "test"
```

### "Access blocked: This app's request is invalid"

Your OAuth consent screen may not be properly configured. Ensure:
1. The Gmail API is enabled
2. Your email is added as a test user (for apps in testing mode)
3. The `gmail.readonly` scope is added

## License

MIT License - see [LICENSE](LICENSE) for details.
