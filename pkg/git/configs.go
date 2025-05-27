package git

import "context"

type Config struct {
	Owner   string // Owner is the Git repository owner
	Repo    string // Repo is the Git repository name
	Token   string // Token is the Git access token
	BaseURL string // BaseURL is the base URL for the Git API

	// Context is the context to use for BigQuery operations
	Context context.Context
}
type Option func(cfg *Config)

func WithOptions(opts ...Option) Option {
	return func(conf *Config) {
		for _, opt := range opts {
			opt(conf)
		}
	}
}

func WithOwner(owner string) Option {
	if owner == "" {
		panic("owner is empty")
	}
	return func(cfg *Config) {
		cfg.Owner = owner
	}
}

func WithRepo(repo string) Option {
	if repo == "" {
		panic("repo is empty")
	}
	return func(cfg *Config) {
		cfg.Repo = repo
	}
}

func WithToken(token string) Option {
	if token == "" {
		panic("token is empty")
	}
	return func(cfg *Config) {
		cfg.Token = token
	}
}

func WithContext(ctx context.Context) Option {
	if ctx == nil {
		panic("context is nil")
	}
	return func(cfg *Config) {
		cfg.Context = ctx
	}
}

// WithBaseURL sets the base URL for the Git API
func WithBaseURL(url string) Option {
	return func(cfg *Config) {
		cfg.BaseURL = url
	}
}

func defaultConfig() *Config {
	return &Config{
		Context: context.Background(),
	}
}
