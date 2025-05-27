package bigquery

import "context"

type Config struct {
	// ProjectId is the Google Cloud Project ID
	ProjectId string

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

func WithProjectId(projectId string) Option {
	if projectId == "" {
		panic("projectId is empty")
	}
	return func(cfg *Config) {
		cfg.ProjectId = projectId
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

func defaultConfig() *Config {
	return &Config{
		Context: context.Background(),
	}
}
