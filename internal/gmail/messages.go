package gmail

import (
	"encoding/base64"
	"fmt"
	"strings"

	"google.golang.org/api/gmail/v1"
)

// Message represents a simplified email message
type Message struct {
	ID       string
	ThreadId string
	Subject  string
	From     string
	To       string
	Date     string
	Snippet  string
	Body     string
}

// SearchMessages searches for messages matching the query
func (c *Client) SearchMessages(query string, maxResults int64) ([]*Message, error) {
	call := c.Service.Users.Messages.List(c.UserID).Q(query)
	if maxResults > 0 {
		call = call.MaxResults(maxResults)
	}

	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search messages: %w", err)
	}

	var messages []*Message
	for _, msg := range resp.Messages {
		m, err := c.GetMessage(msg.Id, false)
		if err != nil {
			continue // Skip messages that fail to fetch
		}
		messages = append(messages, m)
	}

	return messages, nil
}

// GetMessage retrieves a single message by ID
func (c *Client) GetMessage(messageID string, includeBody bool) (*Message, error) {
	format := "metadata"
	if includeBody {
		format = "full"
	}

	msg, err := c.Service.Users.Messages.Get(c.UserID, messageID).Format(format).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return parseMessage(msg, includeBody), nil
}

// GetThread retrieves all messages in a thread.
// The id parameter can be either a thread ID or a message ID.
// If a message ID is provided, the thread ID is resolved automatically.
func (c *Client) GetThread(id string) ([]*Message, error) {
	thread, err := c.Service.Users.Threads.Get(c.UserID, id).Format("full").Do()
	if err != nil {
		// If the ID wasn't found as a thread ID, try treating it as a message ID
		msg, msgErr := c.Service.Users.Messages.Get(c.UserID, id).Format("minimal").Do()
		if msgErr != nil {
			// Return the original thread error if message lookup also fails
			return nil, fmt.Errorf("failed to get thread: %w", err)
		}
		// Use the thread ID from the message
		thread, err = c.Service.Users.Threads.Get(c.UserID, msg.ThreadId).Format("full").Do()
		if err != nil {
			return nil, fmt.Errorf("failed to get thread: %w", err)
		}
	}

	var messages []*Message
	for _, msg := range thread.Messages {
		messages = append(messages, parseMessage(msg, true))
	}

	return messages, nil
}

func parseMessage(msg *gmail.Message, includeBody bool) *Message {
	m := &Message{
		ID:       msg.Id,
		ThreadId: msg.ThreadId,
		Snippet:  msg.Snippet,
	}

	// Extract headers
	for _, header := range msg.Payload.Headers {
		switch strings.ToLower(header.Name) {
		case "subject":
			m.Subject = header.Value
		case "from":
			m.From = header.Value
		case "to":
			m.To = header.Value
		case "date":
			m.Date = header.Value
		}
	}

	if includeBody && msg.Payload != nil {
		m.Body = extractBody(msg.Payload)
	}

	return m
}

func extractBody(payload *gmail.MessagePart) string {
	// Try to get plain text body first
	if payload.MimeType == "text/plain" && payload.Body != nil && payload.Body.Data != "" {
		decoded, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			return string(decoded)
		}
	}

	// Check parts for multipart messages
	for _, part := range payload.Parts {
		if part.MimeType == "text/plain" && part.Body != nil && part.Body.Data != "" {
			decoded, err := base64.URLEncoding.DecodeString(part.Body.Data)
			if err == nil {
				return string(decoded)
			}
		}
		// Recursively check nested parts
		if len(part.Parts) > 0 {
			body := extractBody(part)
			if body != "" {
				return body
			}
		}
	}

	// Fallback to HTML if no plain text found
	if payload.MimeType == "text/html" && payload.Body != nil && payload.Body.Data != "" {
		decoded, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			return string(decoded)
		}
	}

	for _, part := range payload.Parts {
		if part.MimeType == "text/html" && part.Body != nil && part.Body.Data != "" {
			decoded, err := base64.URLEncoding.DecodeString(part.Body.Data)
			if err == nil {
				return string(decoded)
			}
		}
	}

	return ""
}
