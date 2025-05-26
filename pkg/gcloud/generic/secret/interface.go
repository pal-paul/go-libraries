// Package secret provides a type-safe interface to Google Cloud Secret Manager.
// It supports generic typing for secret data, allowing strongly typed secret retrieval
// while maintaining compatibility with the underlying Secret Manager API.
//
// The package provides functionality for:
// - Creating and managing secrets
// - Adding new secret versions
// - Retrieving secrets with type safety
// - Pattern-based secret retrieval
// - Version-specific secret access
package secret

//go:generate mockgen -source=interface.go -destination=mocks/mock-secret.go -package=mocks
import (
	"regexp"

	sm "cloud.google.com/go/secretmanager/apiv1"
)

// secret implements the SecretInterface for a specific type T.
// It maintains the configuration and client connection to Google Cloud Secret Manager.
type secret[T any] struct {
	conf   *Config
	client *sm.Client
}

// SecretInterface defines the operations available for managing secrets in Google Cloud Secret Manager.
// The type parameter T determines the structure that secrets will be unmarshaled into when using
// typed retrieval methods.
type SecretInterface[T any] interface {
	// GetBytes retrieves a secret's latest version as raw bytes.
	// This method is useful when you need the raw secret data or when working with
	// non-JSON secrets.
	//
	// Parameters:
	//   - name: The name of the secret to retrieve
	//
	// Returns:
	//   - []byte: The secret data
	//   - error: An error if the operation fails
	//
	// The error will be of type:
	//   - ErrInvalidSecretName: If the name is empty
	//   - ErrFailedToCreateClient: If the client is not initialized
	GetBytes(name string) ([]byte, error)

	// Get retrieves and unmarshals a secret's latest version into type T.
	// This method provides type-safe secret retrieval by automatically unmarshaling
	// the secret data into the specified type.
	//
	// Parameters:
	//   - name: The name of the secret to retrieve
	//
	// Returns:
	//   - T: The unmarshaled secret data
	//   - error: An error if the operation fails
	//
	// The error will be of type:
	//   - ErrInvalidSecretName: If the name is empty
	//   - ErrFailedToCreateClient: If the client is not initialized
	//   - json.UnmarshalError: If the secret data cannot be unmarshaled into type T
	Get(name string) (T, error)

	// GetVersion retrieves a specific version of a secret as raw bytes.
	// This method allows access to historical versions of secrets when needed.
	//
	// Parameters:
	//   - name: The name of the secret to retrieve
	//   - version: The version identifier (e.g., "1", "2", "latest")
	//
	// Returns:
	//   - []byte: The secret data
	//   - error: An error if the operation fails
	//
	// The error will be of type:
	//   - ErrInvalidSecretName: If the name is empty
	//   - ErrInvalidSecretVersion: If the version is empty
	//   - ErrFailedToCreateClient: If the client is not initialized
	GetVersion(name string, version string) ([]byte, error)

	// GetSecrets retrieves all secrets matching a regular expression pattern.
	// This method is useful for retrieving groups of related secrets.
	//
	// Parameters:
	//   - secretsRegexp: A regular expression pattern to match secret names
	//
	// Returns:
	//   - []SecretData: A list of matching secrets with their data
	//   - error: An error if the operation fails
	//
	// The error will be of type:
	//   - ErrInvalidPattern: If the pattern is nil
	//   - ErrFailedToCreateClient: If the client is not initialized
	//   - ErrFailedToListSecrets: If listing secrets fails
	GetSecrets(secretsRegexp *regexp.Regexp) ([]SecretData, error)

	// CreateSecret creates a new secret in Secret Manager.
	// The secret is created with automatic replication by default.
	//
	// Parameters:
	//   - secretName: The name for the new secret
	//
	// Returns:
	//   - error: An error if the operation fails
	//
	// The error will be of type:
	//   - ErrInvalidSecretName: If the name is empty
	//   - ErrFailedToCreateClient: If the client is not initialized
	//   - ErrFailedToCreateSecret: If secret creation fails
	CreateSecret(secretName string) error

	// AddSecretVersion adds a new version to an existing secret.
	//
	// Parameters:
	//   - secretName: The name of the secret to version
	//   - payload: The secret data to store
	//
	// Returns:
	//   - error: An error if the operation fails
	//
	// The error will be of type:
	//   - ErrInvalidSecretName: If the name is empty
	//   - ErrInvalidSecretData: If the payload is nil or empty
	//   - ErrFailedToCreateClient: If the client is not initialized
	AddSecretVersion(secretName string, payload []byte) error
}
