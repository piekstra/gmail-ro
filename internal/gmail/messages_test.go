package gmail

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/gmail/v1"
)

func TestParseMessage(t *testing.T) {
	t.Run("extracts headers correctly", func(t *testing.T) {
		msg := &gmail.Message{
			Id:       "msg123",
			ThreadId: "thread456",
			Snippet:  "This is a test...",
			Payload: &gmail.MessagePart{
				Headers: []*gmail.MessagePartHeader{
					{Name: "Subject", Value: "Test Subject"},
					{Name: "From", Value: "alice@example.com"},
					{Name: "To", Value: "bob@example.com"},
					{Name: "Date", Value: "Mon, 1 Jan 2024 12:00:00 +0000"},
				},
			},
		}

		result := parseMessage(msg, false)

		assert.Equal(t, "msg123", result.ID)
		assert.Equal(t, "thread456", result.ThreadId)
		assert.Equal(t, "Test Subject", result.Subject)
		assert.Equal(t, "alice@example.com", result.From)
		assert.Equal(t, "bob@example.com", result.To)
		assert.Equal(t, "Mon, 1 Jan 2024 12:00:00 +0000", result.Date)
		assert.Equal(t, "This is a test...", result.Snippet)
	})

	t.Run("extracts thread ID", func(t *testing.T) {
		msg := &gmail.Message{
			Id:       "msg123",
			ThreadId: "thread789",
			Payload: &gmail.MessagePart{
				Headers: []*gmail.MessagePartHeader{},
			},
		}

		result := parseMessage(msg, false)

		assert.Equal(t, "msg123", result.ID)
		assert.Equal(t, "thread789", result.ThreadId)
	})

	t.Run("handles case-insensitive headers", func(t *testing.T) {
		msg := &gmail.Message{
			Id: "msg123",
			Payload: &gmail.MessagePart{
				Headers: []*gmail.MessagePartHeader{
					{Name: "SUBJECT", Value: "Upper Case"},
					{Name: "from", Value: "lower@example.com"},
					{Name: "To", Value: "mixed@example.com"},
				},
			},
		}

		result := parseMessage(msg, false)

		assert.Equal(t, "Upper Case", result.Subject)
		assert.Equal(t, "lower@example.com", result.From)
		assert.Equal(t, "mixed@example.com", result.To)
	})

	t.Run("handles missing headers gracefully", func(t *testing.T) {
		msg := &gmail.Message{
			Id: "msg123",
			Payload: &gmail.MessagePart{
				Headers: []*gmail.MessagePartHeader{},
			},
		}

		result := parseMessage(msg, false)

		assert.Equal(t, "msg123", result.ID)
		assert.Empty(t, result.Subject)
		assert.Empty(t, result.From)
		assert.Empty(t, result.To)
		assert.Empty(t, result.Date)
	})
}

func TestExtractBody(t *testing.T) {
	t.Run("extracts plain text body", func(t *testing.T) {
		bodyText := "Hello, this is the message body."
		encoded := base64.URLEncoding.EncodeToString([]byte(bodyText))

		payload := &gmail.MessagePart{
			MimeType: "text/plain",
			Body: &gmail.MessagePartBody{
				Data: encoded,
			},
		}

		result := extractBody(payload)
		assert.Equal(t, bodyText, result)
	})

	t.Run("extracts plain text from multipart message", func(t *testing.T) {
		bodyText := "Plain text content"
		encoded := base64.URLEncoding.EncodeToString([]byte(bodyText))

		payload := &gmail.MessagePart{
			MimeType: "multipart/alternative",
			Parts: []*gmail.MessagePart{
				{
					MimeType: "text/plain",
					Body: &gmail.MessagePartBody{
						Data: encoded,
					},
				},
				{
					MimeType: "text/html",
					Body: &gmail.MessagePartBody{
						Data: base64.URLEncoding.EncodeToString([]byte("<p>HTML content</p>")),
					},
				},
			},
		}

		result := extractBody(payload)
		assert.Equal(t, bodyText, result)
	})

	t.Run("falls back to HTML if no plain text", func(t *testing.T) {
		htmlContent := "<p>HTML only</p>"
		encoded := base64.URLEncoding.EncodeToString([]byte(htmlContent))

		payload := &gmail.MessagePart{
			MimeType: "text/html",
			Body: &gmail.MessagePartBody{
				Data: encoded,
			},
		}

		result := extractBody(payload)
		assert.Equal(t, htmlContent, result)
	})

	t.Run("handles nested multipart", func(t *testing.T) {
		bodyText := "Nested plain text"
		encoded := base64.URLEncoding.EncodeToString([]byte(bodyText))

		payload := &gmail.MessagePart{
			MimeType: "multipart/mixed",
			Parts: []*gmail.MessagePart{
				{
					MimeType: "multipart/alternative",
					Parts: []*gmail.MessagePart{
						{
							MimeType: "text/plain",
							Body: &gmail.MessagePartBody{
								Data: encoded,
							},
						},
					},
				},
			},
		}

		result := extractBody(payload)
		assert.Equal(t, bodyText, result)
	})

	t.Run("returns empty string for empty body", func(t *testing.T) {
		payload := &gmail.MessagePart{
			MimeType: "text/plain",
			Body:     &gmail.MessagePartBody{},
		}

		result := extractBody(payload)
		assert.Empty(t, result)
	})

	t.Run("returns empty string for nil body", func(t *testing.T) {
		payload := &gmail.MessagePart{
			MimeType: "text/plain",
		}

		result := extractBody(payload)
		assert.Empty(t, result)
	})

	t.Run("handles invalid base64 gracefully", func(t *testing.T) {
		payload := &gmail.MessagePart{
			MimeType: "text/plain",
			Body: &gmail.MessagePartBody{
				Data: "not-valid-base64!!!",
			},
		}

		result := extractBody(payload)
		assert.Empty(t, result)
	})
}

func TestMessageStruct(t *testing.T) {
	t.Run("message struct has all fields", func(t *testing.T) {
		msg := &Message{
			ID:       "test-id",
			ThreadId: "thread-id",
			Subject:  "Test Subject",
			From:     "from@example.com",
			To:       "to@example.com",
			Date:     "2024-01-01",
			Snippet:  "Preview...",
			Body:     "Full body content",
		}

		assert.Equal(t, "test-id", msg.ID)
		assert.Equal(t, "thread-id", msg.ThreadId)
		assert.Equal(t, "Test Subject", msg.Subject)
		assert.Equal(t, "from@example.com", msg.From)
		assert.Equal(t, "to@example.com", msg.To)
		assert.Equal(t, "2024-01-01", msg.Date)
		assert.Equal(t, "Preview...", msg.Snippet)
		assert.Equal(t, "Full body content", msg.Body)
	})
}

func TestParseMessageWithBody(t *testing.T) {
	t.Run("includes body when requested", func(t *testing.T) {
		bodyText := "This is the full body"
		encoded := base64.URLEncoding.EncodeToString([]byte(bodyText))

		msg := &gmail.Message{
			Id: "msg123",
			Payload: &gmail.MessagePart{
				MimeType: "text/plain",
				Headers: []*gmail.MessagePartHeader{
					{Name: "Subject", Value: "Test"},
				},
				Body: &gmail.MessagePartBody{
					Data: encoded,
				},
			},
		}

		result := parseMessage(msg, true)
		assert.Equal(t, bodyText, result.Body)
	})

	t.Run("excludes body when not requested", func(t *testing.T) {
		bodyText := "This should not appear"
		encoded := base64.URLEncoding.EncodeToString([]byte(bodyText))

		msg := &gmail.Message{
			Id: "msg123",
			Payload: &gmail.MessagePart{
				MimeType: "text/plain",
				Headers: []*gmail.MessagePartHeader{
					{Name: "Subject", Value: "Test"},
				},
				Body: &gmail.MessagePartBody{
					Data: encoded,
				},
			},
		}

		result := parseMessage(msg, false)
		assert.Empty(t, result.Body)
	})
}
