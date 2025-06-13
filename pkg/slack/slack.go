package slack

//go:generate mockgen -source=interface.go -destination=mocks/mock-slack.go -package=mocks
import (
	"net/http"
)

const (
	baseUrl = "https://slack.com/api"
)

type slack struct {
	cfg        *Config
	httpClient *http.Client
}

// New creates a new Slack client with the provided options.
func New(opts ...Option) ISlack {
	s := &slack{
		cfg:        defaultConfig(),
		httpClient: &http.Client{},
	}

	for _, opt := range opts {
		opt(s.cfg)
	}

	return s
}
