package slack_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pal-paul/go-libraries/pkg/slack"
	"github.com/stretchr/testify/assert"
)

func setupMockServer(t *testing.T, expectedPath string, method string, status int, response []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedPath, r.URL.Path, "Path mismatch. Expected: %s, Got: %s", expectedPath, r.URL.Path)
		assert.Equal(t, method, r.Method, "Method mismatch. Expected: %s, Got: %s", method, r.Method)

		// Verify Slack token
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"), "Authorization header mismatch")

		// Set response headers
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(status)
		if response != nil {
			w.Write(response)
		}
	}))
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		opts      []slack.Option
		wantError bool
	}{
		{
			name: "success with valid options",
			opts: []slack.Option{
				slack.WithToken("test-token"),
				slack.WithContext(context.Background()),
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := slack.New(tt.opts...)
			assert.NotNil(t, client)
		})
	}
}

func TestAddFormattedMessage(t *testing.T) {
	tests := []struct {
		name      string
		channel   string
		message   slack.Message
		response  []byte
		status    int
		wantError bool
	}{
		{
			name:    "success",
			channel: "test-channel",
			message: slack.Message{
				Text: "Hello World",
				Blocks: []slack.Block{
					{
						Type: slack.SectionBlock,
						Text: &slack.Text{
							Type: slack.Mrkdwn,
							Text: "Test message",
						},
					},
				},
			},
			response: []byte(`{
				"ok": true,
				"channel": "test-channel",
				"ts": "1234567890.123456",
				"message": {
					"text": "Hello World"
				}
			}`),
			status:    http.StatusOK,
			wantError: false,
		},
		{
			name:      "invalid channel",
			channel:   "invalid-channel",
			message:   slack.Message{Text: "Test"},
			status:    http.StatusBadRequest,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupMockServer(
				t,
				"/api/chat.postMessage",
				http.MethodPost,
				tt.status,
				tt.response,
			)
			defer server.Close()

			client := slack.New(
				slack.WithToken("test-token"),
				slack.WithContext(context.Background()),
				slack.WithBaseURL(server.URL+"/api"),
			)

			messageRef, err := client.AddFormattedMessage(tt.channel, tt.message)
			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, messageRef)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.channel, messageRef.Channel)
				assert.NotEmpty(t, messageRef.Timestamp)
			}
		})
	}
}

func TestUploadFileWithContent(t *testing.T) {
	tests := []struct {
		name      string
		fileType  string
		fileName  string
		title     string
		content   string
		msgRef    slack.MessageRef
		response  []byte
		status    int
		wantError bool
	}{
		{
			name:     "success",
			fileType: "text",
			fileName: "test.txt",
			title:    "Test File",
			content:  "Hello World",
			msgRef: slack.MessageRef{
				Channel:   "test-channel",
				Timestamp: "1234567890.123456",
			},
			response: []byte(`{
				"ok": true,
				"file": {
					"id": "F12345678",
					"name": "test.txt"
				}
			}`),
			status:    http.StatusOK,
			wantError: false,
		},
		{
			name:     "upload failure",
			fileType: "text",
			fileName: "test.txt",
			title:    "Test File",
			content:  "Hello World",
			msgRef: slack.MessageRef{
				Channel:   "test-channel",
				Timestamp: "1234567890.123456",
			},
			status:    http.StatusBadRequest,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupMockServer(
				t,
				"/api/files.upload",
				http.MethodPost,
				tt.status,
				tt.response,
			)
			defer server.Close()

			client := slack.New(
				slack.WithToken("test-token"),
				slack.WithContext(context.Background()),
				slack.WithBaseURL(server.URL+"/api"),
			)

			err := client.UploadFileWithContent(tt.fileType, tt.fileName, tt.title, tt.content, tt.msgRef)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
