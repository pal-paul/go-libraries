# Google Cloud Secret Manager Package

A generic Go package that provides a simple, type-safe interface for interacting with Google Cloud Secret Manager. This package allows you to manage secrets with strongly typed responses using Go generics.

## Installation

```bash
go get github.com/pal-paul/go-libraries/pkg/gcloud/generic/secret
```

## Features

- Generic type-safe secret retrieval
- Create and manage secrets
- Add new secret versions
- Retrieve secrets by name or pattern matching
- Support for latest and specific versions
- Automatic client initialization and configuration

## Usage

### Initialize a Secret Client

```go
import (
    "github.com/pal-paul/go-libraries/pkg/gcloud/generic/secret"
)

// Create a client for secrets that will be unmarshaled into MySecretType
client, err := secret.New[MySecretType]()
if err != nil {
    log.Fatalf("Failed to create secret client: %v", err)
}
```

### Retrieve a Secret

```go
// Get a secret as bytes
data, err := client.GetBytes("my-secret")
if err != nil {
    log.Fatalf("Failed to get secret: %v", err)
}

// Get a typed secret (automatically unmarshaled)
secretData, err := client.Get("my-secret")
if err != nil {
    log.Fatalf("Failed to get secret: %v", err)
}
```

### Get a Specific Version

```go
data, err := client.GetVersion("my-secret", "1")
if err != nil {
    log.Fatalf("Failed to get secret version: %v", err)
}
```

### Create a New Secret

```go
err := client.CreateSecret("new-secret")
if err != nil {
    log.Fatalf("Failed to create secret: %v", err)
}
```

### Add a New Secret Version

```go
payload := []byte("my-secret-data")
err := client.AddSecretVersion("my-secret", payload)
if err != nil {
    log.Fatalf("Failed to add secret version: %v", err)
}
```

### Get Multiple Secrets by Pattern

```go
import "regexp"

pattern := regexp.MustCompile("^my-secret-.*$")
secrets, err := client.GetSecrets(pattern)
if err != nil {
    log.Fatalf("Failed to get secrets: %v", err)
}
```

## API Reference

### Types

#### `SecretInterface[T]`

The main interface for interacting with secrets. Type parameter `T` represents the type that secrets will be unmarshaled into.

#### `SecretData`

Structure containing secret data and metadata:

```go
type SecretData struct {
    Data []byte // The secret payload
    Name string // The secret name
}
```

### Functions

#### `New[T any](opts ...Option) (SecretInterface[T], error)`

Creates a new Secret client with optional configuration.

### Methods

#### `GetBytes(name string) ([]byte, error)`

Retrieves the latest version of a secret as raw bytes.

#### `Get(name string) (T, error)`

Retrieves and unmarshals the latest version of a secret into type T.

#### `GetVersion(name string, version string) ([]byte, error)`

Retrieves a specific version of a secret.

#### `GetSecrets(secretsRegexp *regexp.Regexp) ([]SecretData, error)`

Retrieves all secrets matching the provided regular expression pattern.

#### `CreateSecret(secretName string) error`

Creates a new secret.

#### `AddSecretVersion(secretName string, payload []byte) error`

Adds a new version to an existing secret.

## Error Handling

The package provides specific error types for common failure scenarios:

- `ErrFailedToCreateClient`: Client initialization failures
- `ErrInvalidSecretName`: Invalid secret name provided
- `ErrInvalidSecretVersion`: Invalid version specification

## Configuration

The client can be configured using Option functions. Default configuration includes:

- Automatic replication for new secrets
- Latest version retrieval by default
- Project ID from environment/application default credentials

## Best Practices

1. Always handle errors returned by the methods
2. Use strong typing with the generic parameter for type-safe secret handling
3. Use pattern matching carefully to avoid retrieving unnecessary secrets
4. Remember to close the client when done using it
5. Use version-specific retrieval when exact version control is needed

## Advanced Usage

### Using Context and Configuration Options

```go
import (
    "context"
    "time"
)

// Create a context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Create client with context and project ID
client, err := secret.New[MySecretType](
    secret.WithContext(ctx),
    secret.WithProjectID("my-project-id"),
)
```

### Error Handling Examples

```go
// Handle invalid secret name
_, err := client.GetBytes("")
if err != nil {
    var invalidName secret.ErrInvalidSecretName
    if errors.As(err, &invalidName) {
        log.Printf("Invalid secret name: %v", invalidName)
        return
    }
}

// Handle non-existent secret version
_, err = client.GetVersion("my-secret", "999")
if err != nil {
    // Check if it's a GCP API error
    if strings.Contains(err.Error(), "NOT_FOUND") {
        log.Printf("Secret version not found")
        return
    }
}
```

### Integration Testing

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    ctx := context.Background()
    client, err := secret.New[MySecretType](secret.WithContext(ctx))
    require.NoError(t, err)

    // Test full lifecycle
    secretName := "test-secret"
    
    // Create
    err = client.CreateSecret(secretName)
    require.NoError(t, err)

    // Add version
    payload := []byte(`{"key": "value"}`)
    err = client.AddSecretVersion(secretName, payload)
    require.NoError(t, err)

    // Retrieve
    data, err := client.GetBytes(secretName)
    require.NoError(t, err)
    assert.Equal(t, payload, data)
}
```

### Proper Cleanup and Resource Management

```go
// Create client with cleanup
client, err := secret.New[MySecretType]()
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// Use pattern matching for cleanup
cleanup := func() error {
    pattern := regexp.MustCompile("^test-.*$")
    secrets, err := client.GetSecrets(pattern)
    if err != nil {
        return fmt.Errorf("failed to list secrets: %v", err)
    }

    for _, s := range secrets {
        // Delete test secrets
        // Note: Implement deletion logic based on your needs
    }
    return nil
}
```

## Security Best Practices

1. **Access Control**
   - Use principle of least privilege with IAM roles
   - Grant minimal required permissions to service accounts
   - Regularly audit secret access

2. **Secret Management**
   - Never store secrets in version control
   - Rotate secrets regularly
   - Use version labels for secret rotation
   - Consider using customer-managed encryption keys (CMEK)

3. **Client Usage**
   - Always use contexts for timeouts and cancellation
   - Clean up resources properly
   - Validate input data before storing secrets
   - Use strong typing to prevent data corruption

4. **Monitoring and Auditing**
   - Enable audit logging for Secret Manager
   - Monitor secret access patterns
   - Set up alerts for suspicious activities
   - Review access logs regularly

## Known Limitations

1. Secret payload size limit (1MB)
2. Project-specific secrets (cross-project access requires additional setup)
3. IAM permission inheritance limitations
4. Version history limitations

## Troubleshooting

Common issues and solutions:

1. **Client Creation Fails**
   - Check GCP credentials are properly set
   - Verify Project ID is correct
   - Ensure Secret Manager API is enabled

2. **Permission Denied**
   - Verify IAM roles are properly assigned
   - Check service account permissions
   - Ensure project ID matches the secret location

3. **Secret Not Found**
   - Verify secret name is correct
   - Check if secret exists in the specified project
   - Ensure version exists if requesting specific version

4. **Context Deadline Exceeded**
   - Increase context timeout
   - Check network connectivity
   - Verify GCP API availability

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This package is distributed under the MIT License. See LICENSE file for details.
