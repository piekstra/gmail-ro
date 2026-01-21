# Integration Tests

Manual integration tests for gmail-ro. These require valid OAuth credentials.

## Setup

Ensure you have valid credentials configured:
```bash
ls ~/.config/gmail-ro/credentials.json
ls ~/.config/gmail-ro/token.json
```

## Tests

### Search returns ThreadId

Verify search results include thread IDs.

```bash
# Text output should show ThreadId
gmail-ro search "is:inbox" --max 1

# Expected output includes:
# ID: <message-id>
# ThreadId: <thread-id>
# ...

# JSON output should include ThreadId field
gmail-ro search "is:inbox" --max 1 --json | jq '.[0].ThreadId'
# Expected: non-empty string
```

### Thread command accepts message ID

Verify the thread command works with message IDs from search results.

```bash
# Get a message ID from search
MSG_ID=$(gmail-ro search "is:inbox" --max 1 --json | jq -r '.[0].ID')

# Thread command should work with message ID (not just thread ID)
gmail-ro thread "$MSG_ID"
# Expected: displays thread contents without 404 error
```

### Thread command accepts thread ID

Verify the thread command still works with direct thread IDs.

```bash
# Get a thread ID from search
THREAD_ID=$(gmail-ro search "is:inbox" --max 1 --json | jq -r '.[0].ThreadId')

# Thread command should work with thread ID
gmail-ro thread "$THREAD_ID"
# Expected: displays thread contents
```
