# Integration Tests

Comprehensive integration test suite for gmro. Tests are designed to work against any active Gmail account with standard inbox content.

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
gmro search "is:inbox" --max 1  # Quick connectivity check
```

---

## Version Command

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Print version | `gmro version` | Shows "gmro <version>" |

---

## Search Operations

### Basic Search

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Search inbox | `gmro search "is:inbox" --max 5` | Returns messages with ID, ThreadID, From, Subject, Date, Snippet |
| Search with default limit | `gmro search "is:inbox"` | Returns up to 10 messages (default) |
| Custom result limit | `gmro search "is:inbox" --max 3` | Returns exactly 3 messages |
| JSON output | `gmro search "is:inbox" --max 2 --json` | Valid JSON array with message objects |
| No results | `gmro search "xyznonexistent12345uniquequery67890"` | "No messages found." |
| Search unread | `gmro search "is:unread" --max 5` | Returns unread messages (if any) |
| Search starred | `gmro search "is:starred" --max 5` | Returns starred messages (if any) |

### Query Operators

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| From filter | `gmro search "from:noreply" --max 3` | Messages from addresses containing "noreply" |
| Subject filter | `gmro search "subject:welcome" --max 3` | Messages with "welcome" in subject |
| Has attachment | `gmro search "has:attachment" --max 3` | Messages with attachments |
| Date range | `gmro search "after:2024/01/01" --max 3` | Messages after date |
| Combined query | `gmro search "is:inbox has:attachment" --max 3` | Inbox messages with attachments |
| Label filter | `gmro search "label:inbox" --max 3` | Messages in inbox |

### JSON Validation

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| JSON has required fields | `gmro search "is:inbox" --max 1 --json \| jq '.[0] \| keys'` | Contains: id, threadId, from, subject, date, snippet |
| JSON ID is string | `gmro search "is:inbox" --max 1 --json \| jq -e '.[0].id \| type == "string"'` | Returns true |
| JSON ThreadID present | `gmro search "is:inbox" --max 1 --json \| jq -e '.[0].threadId != null'` | Returns true |

---

## Read Operations

### Read Single Message

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Read by ID | `MSG_ID=$(gmro search "is:inbox" --max 1 --json \| jq -r '.[0].id'); gmro read "$MSG_ID"` | Shows ID, From, To, Subject, Date, Body |
| Read JSON output | `MSG_ID=$(gmro search "is:inbox" --max 1 --json \| jq -r '.[0].id'); gmro read "$MSG_ID" --json` | Valid JSON with body content |
| Non-existent message | `gmro read "0000000000000000"` | Error: 404 or "not found" |
| Invalid message ID | `gmro read "invalid-id-format"` | Error message |

### Read Content Verification

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Body included | `MSG_ID=$(gmro search "is:inbox" --max 1 --json \| jq -r '.[0].id'); gmro read "$MSG_ID" --json \| jq -e '.body != null'` | Returns true |
| Headers present | `MSG_ID=$(gmro search "is:inbox" --max 1 --json \| jq -r '.[0].id'); gmro read "$MSG_ID"` | Output contains "From:", "To:", "Subject:", "Date:" |

---

## Thread Operations

### Thread by Thread ID

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| View thread | `THREAD_ID=$(gmro search "is:inbox" --max 1 --json \| jq -r '.[0].threadId'); gmro thread "$THREAD_ID"` | Shows "Thread contains N message(s)" and all messages |
| Thread JSON | `THREAD_ID=$(gmro search "is:inbox" --max 1 --json \| jq -r '.[0].threadId'); gmro thread "$THREAD_ID" --json` | Valid JSON array of messages |

### Thread by Message ID

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Thread from message ID | `MSG_ID=$(gmro search "is:inbox" --max 1 --json \| jq -r '.[0].id'); gmro thread "$MSG_ID"` | Shows thread containing that message |
| Thread message count | `MSG_ID=$(gmro search "is:inbox" --max 1 --json \| jq -r '.[0].id'); gmro thread "$MSG_ID" --json \| jq 'length >= 1'` | Returns true (at least 1 message) |

### Error Cases

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Non-existent thread | `gmro thread "0000000000000000"` | Error: 404 or "not found" |

---

## Labels Operations

### List Labels

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| List all labels | `gmro labels` | Shows NAME, TYPE, TOTAL, UNREAD columns |
| Labels JSON output | `gmro labels --json` | Valid JSON array with label objects |
| Labels JSON has fields | `gmro labels --json \| jq -e '.[0] \| has("id", "name", "type")'` | Returns true |

### Label Display in Messages

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Search shows labels | `gmro search "is:inbox" --max 1` | Output may include "Labels:" line if message has user labels |
| Search shows categories | `gmro search "category:updates" --max 1` | Output may include "Categories: updates" |
| Search JSON has labels | `gmro search "is:inbox" --max 1 --json \| jq '.[0] \| has("labels", "categories")'` | Returns true |
| Read shows labels | `MSG_ID=$(gmro search "is:inbox" --max 1 --json \| jq -r '.[0].id'); gmro read "$MSG_ID"` | Output may include "Labels:" and "Categories:" lines |

### Label-Based Search

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Search by label | `gmro search "label:inbox" --max 3` | Returns inbox messages |
| Search by category | `gmro search "category:updates" --max 3` | Returns updates category messages |
| Exclude category | `gmro search "is:inbox -category:promotions" --max 3` | Returns inbox excluding promotions |
| Combined label search | `gmro search "is:inbox -category:social -category:promotions" --max 3` | Inbox excluding social and promotions |

---

## Attachment Operations

### Setup: Find Message with Attachments
```bash
# Store a message ID with attachments for subsequent tests
ATTACHMENT_MSG_ID=$(gmro search "has:attachment" --max 1 --json | jq -r '.[0].id')
```

### List Attachments

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| List attachments | `gmro attachments list "$ATTACHMENT_MSG_ID"` | Shows filename, type, size for each attachment |
| List JSON | `gmro attachments list "$ATTACHMENT_MSG_ID" --json` | Valid JSON array with attachment metadata |
| No attachments | `MSG_ID=$(gmro search "is:inbox -has:attachment" --max 1 --json \| jq -r '.[0].id'); gmro attachments list "$MSG_ID"` | "No attachments found for message." |

### JSON Validation

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Attachment has filename | `gmro attachments list "$ATTACHMENT_MSG_ID" --json \| jq -e '.[0].filename != null'` | Returns true |
| Attachment has mimeType | `gmro attachments list "$ATTACHMENT_MSG_ID" --json \| jq -e '.[0].mimeType != null'` | Returns true |
| Attachment has size | `gmro attachments list "$ATTACHMENT_MSG_ID" --json \| jq -e '.[0].size >= 0'` | Returns true |

### Download Attachments

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Download all (no flag) | `gmro attachments download "$ATTACHMENT_MSG_ID"` | Error: "must specify --filename or --all" |
| Download all | `gmro attachments download "$ATTACHMENT_MSG_ID" --all -o /tmp/gmail-test` | Downloads all attachments, shows "Downloaded: ..." |
| Download specific file | `FILENAME=$(gmro attachments list "$ATTACHMENT_MSG_ID" --json \| jq -r '.[0].filename'); gmro attachments download "$ATTACHMENT_MSG_ID" -f "$FILENAME" -o /tmp/gmail-test` | Downloads specific file |
| Non-existent filename | `gmro attachments download "$ATTACHMENT_MSG_ID" -f "nonexistent-file-12345.xyz"` | Error: "attachment not found" |
| Verify file created | `ls /tmp/gmail-test/` | Downloaded files exist |
| Verify file size > 0 | `stat -f%z /tmp/gmail-test/* \| head -1` (macOS) | Non-zero file size |

### Zip Extraction (if zip attachment available)

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Find zip attachment | `ZIP_MSG_ID=$(gmro search "has:attachment filename:zip" --max 1 --json \| jq -r '.[0].id')` | Message ID or null |
| Download and extract | `gmro attachments download "$ZIP_MSG_ID" -f "*.zip" --extract -o /tmp/gmail-zip-test` | Extracts to directory |
| Verify extraction | `ls /tmp/gmail-zip-test/*/` | Extracted files present |

---

## Error Handling

| Test Case | Command | Expected Result |
|-----------|---------|-----------------|
| Missing required arg (search) | `gmro search` | Error: accepts 1 arg(s), received 0 |
| Missing required arg (read) | `gmro read` | Error: accepts 1 arg(s), received 0 |
| Missing required arg (thread) | `gmro thread` | Error: accepts 1 arg(s), received 0 |
| Missing required arg (attachments list) | `gmro attachments list` | Error: accepts 1 arg(s), received 0 |
| Invalid subcommand | `gmro invalid` | Error: unknown command |
| Help flag | `gmro --help` | Shows usage information |
| Search help | `gmro search --help` | Shows search-specific help |

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
MSG_ID=$(gmro search "is:inbox" --max 1 --json | jq -r '.[0].id')

# 2. Read the full message
gmro read "$MSG_ID"

# 3. View the full thread
gmro thread "$MSG_ID"
```

### Workflow 2: Find and Download Attachments
```bash
# 1. Find message with attachments
ATTACHMENT_MSG_ID=$(gmro search "has:attachment" --max 1 --json | jq -r '.[0].id')

# 2. List attachments
gmro attachments list "$ATTACHMENT_MSG_ID"

# 3. Download all attachments
gmro attachments download "$ATTACHMENT_MSG_ID" --all -o /tmp/gmail-attachments

# 4. Verify downloads
ls -la /tmp/gmail-attachments/
```

### Workflow 3: JSON Pipeline
```bash
# Extract all From addresses from recent inbox messages
gmro search "is:inbox" --max 10 --json | jq -r '.[].from'

# Get message bodies from a thread
THREAD_ID=$(gmro search "is:inbox" --max 1 --json | jq -r '.[0].threadId')
gmro thread "$THREAD_ID" --json | jq -r '.[].body'
```

---

## Test Execution Checklist

### Setup
- [ ] Build latest: `make build`
- [ ] Verify credentials exist
- [ ] Quick connectivity test: `gmro search "is:inbox" --max 1`

### Core Commands
- [ ] `gmro version`
- [ ] `gmro search` with various queries
- [ ] `gmro read` by message ID
- [ ] `gmro thread` by thread ID
- [ ] `gmro thread` by message ID
- [ ] `gmro labels` (list all labels)

### Labels
- [ ] `gmro labels` text output
- [ ] `gmro labels --json` JSON output
- [ ] Search by label/category
- [ ] Labels/categories in message output

### Attachments
- [ ] `gmro attachments list`
- [ ] `gmro attachments download --all`
- [ ] `gmro attachments download --filename`
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
