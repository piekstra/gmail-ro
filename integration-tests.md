# Integration Tests

Comprehensive integration test suite for gmail-ro. Tests are designed to work against any active Gmail account with standard inbox content.

## Test Environment Setup

### Prerequisites
- Valid OAuth credentials configured (`~/.config/gmail-ro/credentials.json`, `~/.config/gmail-ro/token.json`)
- Access to a Gmail account with:
  - At least some messages in the inbox
  - At least one email with attachments (for attachment tests)
  - At least one email thread with multiple messages

### Verification
```bash
ls ~/.config/gmail-ro/credentials.json
ls ~/.config/gmail-ro/token.json
gmail-ro search "is:inbox" --max 1  # Quick connectivity check
```

---

## Version Command

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Print version | `gmail-ro version` | Shows "gmail-ro <version>" |

---

## Search Operations

### Basic Search

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Search inbox | `gmail-ro search "is:inbox" --max 5` | Returns messages with ID, ThreadID, From, Subject, Date, Snippet |
| Search with default limit | `gmail-ro search "is:inbox"` | Returns up to 10 messages (default) |
| Custom result limit | `gmail-ro search "is:inbox" --max 3` | Returns exactly 3 messages |
| JSON output | `gmail-ro search "is:inbox" --max 2 --json` | Valid JSON array with message objects |
| No results | `gmail-ro search "xyznonexistent12345uniquequery67890"` | "No messages found." |
| Search unread | `gmail-ro search "is:unread" --max 5` | Returns unread messages (if any) |
| Search starred | `gmail-ro search "is:starred" --max 5` | Returns starred messages (if any) |

### Query Operators

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| From filter | `gmail-ro search "from:noreply" --max 3` | Messages from addresses containing "noreply" |
| Subject filter | `gmail-ro search "subject:welcome" --max 3` | Messages with "welcome" in subject |
| Has attachment | `gmail-ro search "has:attachment" --max 3` | Messages with attachments |
| Date range | `gmail-ro search "after:2024/01/01" --max 3` | Messages after date |
| Combined query | `gmail-ro search "is:inbox has:attachment" --max 3` | Inbox messages with attachments |
| Label filter | `gmail-ro search "label:inbox" --max 3` | Messages in inbox |

### JSON Validation

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| JSON has required fields | `gmail-ro search "is:inbox" --max 1 --json \| jq '.[0] \| keys'` | Contains: id, threadId, from, subject, date, snippet |
| JSON ID is string | `gmail-ro search "is:inbox" --max 1 --json \| jq -e '.[0].id \| type == "string"'` | Returns true |
| JSON ThreadID present | `gmail-ro search "is:inbox" --max 1 --json \| jq -e '.[0].threadId != null'` | Returns true |

---

## Read Operations

### Read Single Message

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Read by ID | `MSG_ID=$(gmail-ro search "is:inbox" --max 1 --json \| jq -r '.[0].id'); gmail-ro read "$MSG_ID"` | Shows ID, From, To, Subject, Date, Body |
| Read JSON output | `MSG_ID=$(gmail-ro search "is:inbox" --max 1 --json \| jq -r '.[0].id'); gmail-ro read "$MSG_ID" --json` | Valid JSON with body content |
| Non-existent message | `gmail-ro read "0000000000000000"` | Error: 404 or "not found" |
| Invalid message ID | `gmail-ro read "invalid-id-format"` | Error message |

### Read Content Verification

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Body included | `MSG_ID=$(gmail-ro search "is:inbox" --max 1 --json \| jq -r '.[0].id'); gmail-ro read "$MSG_ID" --json \| jq -e '.body != null'` | Returns true |
| Headers present | `MSG_ID=$(gmail-ro search "is:inbox" --max 1 --json \| jq -r '.[0].id'); gmail-ro read "$MSG_ID"` | Output contains "From:", "To:", "Subject:", "Date:" |

---

## Thread Operations

### Thread by Thread ID

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| View thread | `THREAD_ID=$(gmail-ro search "is:inbox" --max 1 --json \| jq -r '.[0].threadId'); gmail-ro thread "$THREAD_ID"` | Shows "Thread contains N message(s)" and all messages |
| Thread JSON | `THREAD_ID=$(gmail-ro search "is:inbox" --max 1 --json \| jq -r '.[0].threadId'); gmail-ro thread "$THREAD_ID" --json` | Valid JSON array of messages |

### Thread by Message ID

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Thread from message ID | `MSG_ID=$(gmail-ro search "is:inbox" --max 1 --json \| jq -r '.[0].id'); gmail-ro thread "$MSG_ID"` | Shows thread containing that message |
| Thread message count | `MSG_ID=$(gmail-ro search "is:inbox" --max 1 --json \| jq -r '.[0].id'); gmail-ro thread "$MSG_ID" --json \| jq 'length >= 1'` | Returns true (at least 1 message) |

### Error Cases

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Non-existent thread | `gmail-ro thread "0000000000000000"` | Error: 404 or "not found" |

---

## Labels Operations

### List Labels

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| List all labels | `gmail-ro labels` | Shows NAME, TYPE, TOTAL, UNREAD columns |
| Labels JSON output | `gmail-ro labels --json` | Valid JSON array with label objects |
| Labels JSON has fields | `gmail-ro labels --json \| jq -e '.[0] \| has("id", "name", "type")'` | Returns true |

### Label Display in Messages

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Search shows labels | `gmail-ro search "is:inbox" --max 1` | Output may include "Labels:" line if message has user labels |
| Search shows categories | `gmail-ro search "category:updates" --max 1` | Output may include "Categories: updates" |
| Search JSON has labels | `gmail-ro search "is:inbox" --max 1 --json \| jq '.[0] \| has("labels", "categories")'` | Returns true |
| Read shows labels | `MSG_ID=$(gmail-ro search "is:inbox" --max 1 --json \| jq -r '.[0].id'); gmail-ro read "$MSG_ID"` | Output may include "Labels:" and "Categories:" lines |

### Label-Based Search

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Search by label | `gmail-ro search "label:inbox" --max 3` | Returns inbox messages |
| Search by category | `gmail-ro search "category:updates" --max 3` | Returns updates category messages |
| Exclude category | `gmail-ro search "is:inbox -category:promotions" --max 3` | Returns inbox excluding promotions |
| Combined label search | `gmail-ro search "is:inbox -category:social -category:promotions" --max 3` | Inbox excluding social and promotions |

---

## Attachment Operations

### Setup: Find Message with Attachments
```bash
# Store a message ID with attachments for subsequent tests
ATTACHMENT_MSG_ID=$(gmail-ro search "has:attachment" --max 1 --json | jq -r '.[0].id')
```

### List Attachments

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| List attachments | `gmail-ro attachments list "$ATTACHMENT_MSG_ID"` | Shows filename, type, size for each attachment |
| List JSON | `gmail-ro attachments list "$ATTACHMENT_MSG_ID" --json` | Valid JSON array with attachment metadata |
| No attachments | `MSG_ID=$(gmail-ro search "is:inbox -has:attachment" --max 1 --json \| jq -r '.[0].id'); gmail-ro attachments list "$MSG_ID"` | "No attachments found." |

### JSON Validation

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Attachment has filename | `gmail-ro attachments list "$ATTACHMENT_MSG_ID" --json \| jq -e '.[0].filename != null'` | Returns true |
| Attachment has mimeType | `gmail-ro attachments list "$ATTACHMENT_MSG_ID" --json \| jq -e '.[0].mimeType != null'` | Returns true |
| Attachment has size | `gmail-ro attachments list "$ATTACHMENT_MSG_ID" --json \| jq -e '.[0].size >= 0'` | Returns true |

### Download Attachments

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Download all (no flag) | `gmail-ro attachments download "$ATTACHMENT_MSG_ID"` | Error: "must specify --filename or --all" |
| Download all | `gmail-ro attachments download "$ATTACHMENT_MSG_ID" --all -o /tmp/gmail-test` | Downloads all attachments, shows "Downloaded: ..." |
| Download specific file | `FILENAME=$(gmail-ro attachments list "$ATTACHMENT_MSG_ID" --json \| jq -r '.[0].filename'); gmail-ro attachments download "$ATTACHMENT_MSG_ID" -f "$FILENAME" -o /tmp/gmail-test` | Downloads specific file |
| Non-existent filename | `gmail-ro attachments download "$ATTACHMENT_MSG_ID" -f "nonexistent-file-12345.xyz"` | Error: "attachment not found" |
| Verify file created | `ls /tmp/gmail-test/` | Downloaded files exist |
| Verify file size > 0 | `stat -f%z /tmp/gmail-test/* \| head -1` (macOS) | Non-zero file size |

### Zip Extraction (if zip attachment available)

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Find zip attachment | `ZIP_MSG_ID=$(gmail-ro search "has:attachment filename:zip" --max 1 --json \| jq -r '.[0].id')` | Message ID or null |
| Download and extract | `gmail-ro attachments download "$ZIP_MSG_ID" -f "*.zip" --extract -o /tmp/gmail-zip-test` | Extracts to directory |
| Verify extraction | `ls /tmp/gmail-zip-test/*/` | Extracted files present |

---

## Error Handling

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Missing required arg (search) | `gmail-ro search` | Error: accepts 1 arg(s), received 0 |
| Missing required arg (read) | `gmail-ro read` | Error: accepts 1 arg(s), received 0 |
| Missing required arg (thread) | `gmail-ro thread` | Error: accepts 1 arg(s), received 0 |
| Missing required arg (attachments list) | `gmail-ro attachments list` | Error: accepts 1 arg(s), received 0 |
| Invalid subcommand | `gmail-ro invalid` | Error: unknown command |
| Help flag | `gmail-ro --help` | Shows usage information |
| Search help | `gmail-ro search --help` | Shows search-specific help |

---

## Output Format Consistency

### Text Output Structure

| Command Type | Expected Fields |
|--------------|-----------------|
| Search | ID, ThreadID, From, Subject, Date, Labels (if any), Categories (if any), Snippet, separator (---) |
| Read | ID, From, To, Subject, Date, Labels (if any), Categories (if any), "--- Body ---", body content |
| Thread | "Thread contains N message(s)", per-message: "=== Message X of Y ===", ID, From, To, Subject, Date, Labels, Categories, body |
| Labels | NAME, TYPE, TOTAL, UNREAD columns |
| Attachments List | "Found N attachment(s):", numbered list with filename, Type, Size |

### JSON Schema Validation

| Type | Required Fields |
|------|-----------------|
| Search result | id, threadId, from, subject, date, snippet, labels, categories |
| Message | id, threadId, from, to, subject, date, body, labels, categories |
| Attachment | filename, mimeType, size, partId |
| Label | id, name, type, messagesTotal, messagesUnread |

---

## End-to-End Workflows

### Workflow 1: Search -> Read -> Thread
```bash
# 1. Search for a message
MSG_ID=$(gmail-ro search "is:inbox" --max 1 --json | jq -r '.[0].id')

# 2. Read the full message
gmail-ro read "$MSG_ID"

# 3. View the full thread
gmail-ro thread "$MSG_ID"
```

### Workflow 2: Find and Download Attachments
```bash
# 1. Find message with attachments
ATTACHMENT_MSG_ID=$(gmail-ro search "has:attachment" --max 1 --json | jq -r '.[0].id')

# 2. List attachments
gmail-ro attachments list "$ATTACHMENT_MSG_ID"

# 3. Download all attachments
gmail-ro attachments download "$ATTACHMENT_MSG_ID" --all -o /tmp/gmail-attachments

# 4. Verify downloads
ls -la /tmp/gmail-attachments/
```

### Workflow 3: JSON Pipeline
```bash
# Extract all From addresses from recent inbox messages
gmail-ro search "is:inbox" --max 10 --json | jq -r '.[].from'

# Get message bodies from a thread
THREAD_ID=$(gmail-ro search "is:inbox" --max 1 --json | jq -r '.[0].threadId')
gmail-ro thread "$THREAD_ID" --json | jq -r '.[].body'
```

---

## Test Execution Checklist

### Setup
- [ ] Build latest: `make build`
- [ ] Verify credentials exist
- [ ] Quick connectivity test: `gmail-ro search "is:inbox" --max 1`

### Core Commands
- [ ] `gmail-ro version`
- [ ] `gmail-ro search` with various queries
- [ ] `gmail-ro read` by message ID
- [ ] `gmail-ro thread` by thread ID
- [ ] `gmail-ro thread` by message ID
- [ ] `gmail-ro labels` (list all labels)

### Labels
- [ ] `gmail-ro labels` text output
- [ ] `gmail-ro labels --json` JSON output
- [ ] Search by label/category
- [ ] Labels/categories in message output

### Attachments
- [ ] `gmail-ro attachments list`
- [ ] `gmail-ro attachments download --all`
- [ ] `gmail-ro attachments download --filename`
- [ ] Zip extraction with `--extract`

### Output Formats
- [ ] Text output for all commands
- [ ] JSON output for all commands
- [ ] JSON validates with jq

### Error Handling
- [ ] Missing arguments
- [ ] Invalid IDs
- [ ] Non-existent resources

### Cleanup
- [ ] Remove test downloads: `rm -rf /tmp/gmail-test /tmp/gmail-zip-test /tmp/gmail-attachments`

---

## Adding New Tests

When adding new features or fixing bugs:

1. Add test cases to the appropriate section above
2. Include both happy path and error cases
3. Document any known limitations or edge cases
4. Update the "Test Execution Checklist" if needed
