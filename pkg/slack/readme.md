# Slack Package

A Go package for interacting with the Slack API. This package provides a simple interface to perform common Slack operations like sending messages, uploading files, and managing reactions.

## Installation

```bash
go get github.com/pal-paul/go-libraries/pkg/slack
```

## Features

- Send formatted messages to channels
- Upload files with content
- Add and remove reactions
- Thread support
- Configurable client options
- Error handling

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "github.com/pal-paul/go-libraries/pkg/slack"
)

func main() {
    // Create a new Slack client
    client := slack.New(
        slack.WithToken("your-slack-token"),
        slack.WithContext(context.Background()),
    )

    // Send a message
    message := slack.Message{
        Text: "Hello from Go!",
        Blocks: []slack.Block{
            {
                Type: slack.SectionBlock,
                Text: &slack.Text{
                    Type: slack.Mrkdwn,
                    Text: "Hello *world*!",
                },
            },
        },
    }

    messageRef, err := client.AddFormattedMessage("your-channel", message)
    if err != nil {
        panic(err)
    }

    // Add a reaction to the message
    err = client.AddReaction("thumbsup", messageRef)
    if err != nil {
        panic(err)
    }
}
```

## Configuration

The package can be configured using various options:

```go
// Available options
WithToken(token string)       // Set Slack API token
WithContext(ctx context.Context) // Set context for API requests
WithBaseURL(url string)       // Set custom API base URL
```

## API Reference

### Message Operations

#### AddFormattedMessage

```go
AddFormattedMessage(channel string, message Message) (MessageRef, error)
```

Sends a formatted message to a Slack channel.

#### UploadFileWithContent

```go
UploadFileWithContent(fileType, fileName, title, content string, messageRef MessageRef) error
```

Uploads a file with content to Slack.

### Reaction Operations

#### AddReaction

```go
AddReaction(name string, item MessageRef) error
```

Adds a reaction emoji to a message.

#### RemoveReaction

```go
RemoveReaction(name string, item MessageRef) error
```

Removes a reaction emoji from a message.

## Types

### Message

```go
type Message struct {
    Channel string  // Channel ID or name
    Thread  string  // Thread timestamp for replies
    Text    string  // Plain text message
    Blocks  []Block // Message blocks for rich formatting
}
```

### Block Types

```go
const (
    SectionBlock  BlockType = "section"
    HeaderBlock   BlockType = "header"
    ActionsBlock  BlockType = "actions"
    RichTextBlock BlockType = "rich_text"
)
```

## Error Handling

The package returns meaningful errors for various scenarios:

- Invalid token
- Channel not found
- Rate limiting
- API errors
- Network issues

Example error handling:

```go
messageRef, err := client.AddFormattedMessage(channel, message)
if err != nil {
    switch {
    case err.Error() == slack.ErrInvalidToken:
        // Handle invalid token
    case err.Error() == slack.ErrInvalidChannel:
        // Handle invalid channel
    default:
        // Handle other errors
    }
}
```

## Testing

The package includes comprehensive tests. To run them:

```bash
go test -v ./...
```

For mock generation (requires `mockgen`):

```bash
go generate ./...
```

## Best Practices

1. **Token Security**: Never hardcode Slack tokens. Use environment variables or secure configuration management.
2. **Rate Limiting**: Be mindful of Slack API rate limits in production applications.
3. **Context Usage**: Use context with timeout for long-running operations.
4. **Error Handling**: Always check for errors and handle them appropriately.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This package is released under the MIT License.
