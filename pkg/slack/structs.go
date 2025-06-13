package slack

// MessageRef represents a reference to a Slack message
type MessageRef struct {
	Channel   string `json:"channel"`
	Timestamp string `json:"ts"`
}

// BlockType represents the type of a Slack message block
type BlockType string

const (
	SectionBlock  BlockType = "section"
	HeaderBlock   BlockType = "header"
	ActionsBlock  BlockType = "actions"
	RichTextBlock BlockType = "rich_text"
)

// TextType represents the type of text in a Slack message
type TextType string

const (
	Mrkdwn               TextType = "mrkdwn"
	PlainText            TextType = "plain_text"
	RichTextSection      TextType = "rich_text_section"
	RichTextPreformatted TextType = "rich_text_preformatted"
)

// ActionType represents the type of action in a Slack message
type ActionType string

const (
	Button     ActionType = "button"
	UserSelect ActionType = "users_select"
)

// Message represents a Slack message
type Message struct {
	Channel string  `json:"channel,omitempty"`
	Thread  string  `json:"thread_ts,omitempty"`
	Text    string  `json:"text,omitempty"`
	Blocks  []Block `json:"blocks,omitempty"`
}

// Text represents text content in a Slack message
type Text struct {
	Type  TextType `json:"type,omitempty"`
	Text  string   `json:"text,omitempty"`
	Emoji bool     `json:"emoji,omitempty"`
}

// Field represents a field in a Slack message block
type Field struct {
	Type TextType `json:"type,omitempty"`
	Text string   `json:"text,omitempty"`
}

// Element represents an element in a Slack message block
type Element struct {
	Type     string `json:"type,omitempty"`
	Text     *Text  `json:"text,omitempty"`
	Style    string `json:"style,omitempty"`
	Value    string `json:"value,omitempty"`
	Elements []Text `json:"elements,omitempty"`
	ActionId string `json:"action_id,omitempty"`
}

// Block represents a block in a Slack message
type Block struct {
	Type     BlockType `json:"type,omitempty"`
	Text     *Text     `json:"text,omitempty"`
	Fields   []Field   `json:"fields,omitempty"`
	Elements []Element `json:"elements,omitempty"`
	BlockId  string    `json:"block_id,omitempty"`
}

// SlackResponse handles parsing out errors from the web api.
type SlackResponse struct {
	Ok               bool                  `json:"ok"`
	Error            string                `json:"error"`
	Channel          string                `json:"channel"`
	Ts               string                `json:"ts"`
	ResponseMetadata SlackResponseMetadata `json:"response_metadata"`
}

// SlackResponseMetadata contains metadata about the response
type SlackResponseMetadata struct {
	Cursor   string   `json:"next_cursor"`
	Messages []string `json:"messages"`
	Warnings []string `json:"warnings"`
}

// FileUploadResponse represents a response from Slack's files.upload API
type FileUploadResponse struct {
	SlackResponse
	File struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"file"`
}
