// Package secret provides configuration options for the Google Cloud Secret Manager client.
package secret

import "context"

// Config holds the configuration for the Secret Manager client.
// It includes project identification and context for operations.
type Config struct {
	// ProjectId is the Google Cloud project ID where secrets are stored.
	// If not provided, it will be inferred from the environment.
	ProjectId string

	// Context is used for operation timeouts and cancellation.
	// If not provided, context.Background() will be used.
	Context context.Context
}

// Option is a function that modifies the client configuration.
// It's used to provide a clean, flexible API for configuration.
type Option func(conf *Config)

// WithOptions combines multiple configuration options into a single option.
// This is useful when you want to apply multiple configuration changes together.
//
// Example:
//
//	client := secret.New(
//	    secret.WithOptions(
//	        secret.WithProjectId("my-project"),
//	        secret.WithContext(ctx),
//	    ),
//	)
func WithOptions(opts ...Option) Option {
	return func(conf *Config) {
		for _, opt := range opts {
			opt(conf)
		}
	}
}

// WithProjectId sets the Google Cloud project ID for the client.
// The project ID determines where secrets will be stored and retrieved from.
//
// Parameters:
//   - projectId: The Google Cloud project ID
//
// Example:
//
//	client := secret.New(secret.WithProjectId("my-project"))
func WithProjectId(projectId string) Option {
	return func(conf *Config) {
		conf.ProjectId = projectId
	}
}

// WithContext sets the context for Secret Manager operations.
// The context can be used to set timeouts or cancel operations.
//
// Parameters:
//   - ctx: The context to use for operations
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	client := secret.New(secret.WithContext(ctx))
func WithContext(ctx context.Context) Option {
	return func(conf *Config) {
		conf.Context = ctx
	}
}

// defaultConfig creates a default configuration with:
// - Background context
// - Project ID from environment (via Application Default Credentials)
func defaultConfig() *Config {
	conf := &Config{
		Context: context.Background(),
	}
	return conf
}
