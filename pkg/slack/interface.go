package slack

//go:generate mockgen -source=interface.go -destination=mocks/mock-slack.go -package=mocks

type ISlack interface {
	// UploadFileWithContent uploads a file with the given content to Slack.
	// Parameters:
	//   - fileType: The type of file (e.g., "text")
	//   - fileName: Name of the file
	//   - title: Title of the file
	//   - content: Content of the file
	//   - messageRef: Reference to a message if posting in a thread
	// Returns:
	//   - error: Any error that occurred during upload
	UploadFileWithContent(fileType, fileName, title, content string, messageRef MessageRef) error

	// AddFormattedMessage sends a formatted message to a Slack channel.
	// Parameters:
	//   - channel: The channel to send the message to
	//   - message: The message content and formatting
	// Returns:
	//   - MessageRef: Reference to the sent message
	//   - error: Any error that occurred while sending
	AddFormattedMessage(channel string, message Message) (MessageRef, error)

	// AddReaction adds a reaction emoji to a message.
	// Parameters:
	//   - name: Name of the reaction emoji
	//   - item: Reference to the message to react to
	// Returns:
	//   - error: Any error that occurred while adding the reaction
	AddReaction(name string, item MessageRef) error

	// RemoveReaction removes a reaction emoji from a message.
	// Parameters:
	//   - name: Name of the reaction emoji
	//   - item: Reference to the message to remove reaction from
	// Returns:
	//   - error: Any error that occurred while removing the reaction
	RemoveReaction(name string, item MessageRef) error
}
