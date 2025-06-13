package slack

import "context"

type Config struct {
	Token   string
	BaseURL string
	Context context.Context
}

// Option is a function that configures a Config.
type Option func(*Config)

// WithToken sets the Slack API token.
func WithToken(token string) Option {
	return func(cfg *Config) {
		cfg.Token = token
	}
}

// WithContext sets the context for API requests.
func WithContext(ctx context.Context) Option {
	return func(cfg *Config) {
		cfg.Context = ctx
	}
}

// WithBaseURL sets a custom base URL for API requests.
func WithBaseURL(url string) Option {
	return func(cfg *Config) {
		cfg.BaseURL = url
	}
}

func defaultConfig() *Config {
	return &Config{
		BaseURL: baseUrl,
		Context: context.Background(),
	}
}
