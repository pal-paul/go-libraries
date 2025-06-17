# env

Package env provides functionality for parsing environment variables into Go structs using struct tags. It supports various data types, custom unmarshaling, required fields, default values, and multiple environment variable names for a single field.

## Features

- **Basic Types Support**: Automatically parses strings, booleans, integers, floats, and time.Duration
- **Custom Unmarshaling**: Implement the `Unmarshaler` interface for custom type parsing
- **Required Fields**: Mark fields as required with the `required` tag
- **Default Values**: Specify default values with `default=value` tag
- **Multiple Environment Variables**: Specify multiple possible environment variable names for a field
- **Nested Structs**: Support for nested struct fields
- **Pointer Types**: Support for pointer fields
- **Environment Override**: Ability to override environment variables programmatically

## Installation

```bash
go get github.com/pal-paul/go-libraries/pkg/env
```

## Usage

### Basic Usage with Tags

```go

// Define your configuration struct with env tags
type Config struct {
    // Basic type examples
    Host        string        `env:"HOST"`
    Port        int           `env:"PORT"`
    Debug       bool          `env:"DEBUG"`
    Timeout     time.Duration `env:"TIMEOUT"`
    
    // Required field example (default behavior)
    APIKey      string        `env:"API_KEY,required"`
    
    // Explicitly optional field example
    OptionalKey string        `env:"OPTIONAL_KEY,required=false"`
    
    // Various ways to specify optional fields
    Feature1    string        `env:"FEATURE_1,required=false"`
    Feature2    string        `env:"FEATURE_2,required=0"`
    Feature3    string        `env:"FEATURE_3,required=no"`
    
    // Default value example
    CacheSize   int           `env:"CACHE_SIZE,default=100"`
    
    // Multiple env names example (will try REDIS_URL first, then CACHE_URL)
    CacheURL    string        `env:"REDIS_URL,CACHE_URL"`
    
    // Nested struct example
    Database struct {
        Host     string `env:"DB_HOST"`
        Port     int    `env:"DB_PORT"`
        Name     string `env:"DB_NAME"`
        User     string `env:"DB_USER"`
        Password string `env:"DB_PASSWORD,required"`
    }
    
    // Pointer field example
    MaxRetries  *int         `env:"MAX_RETRIES"`
}

func main() {
    // Create a new config instance
    cfg := &Config{}
    
    // Parse environment variables into the struct
    _, err := env.Unmarshal(cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    // Use the config values
    fmt.Printf("Host: %s\n", cfg.Host)
    fmt.Printf("Database: %s@%s:%d/%s\n", 
        cfg.Database.User, 
        cfg.Database.Host, 
        cfg.Database.Port, 
        cfg.Database.Name,
    )
    
    // Override example
    newPort := "8080"
    override := env.Override{
        "PORT": &newPort,
    }
    
    // Apply the override
    es := make(env.EnvSet)
    _, err = es.Apply(override, cfg)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Tag Options

### required

- `required`: Field is required (default behavior when specified)
- `required=true` or `required=1` or `required=yes`: Field is required
- `required=false` or `required=0` or `required=no`: Field is optional

### default

- `default=value`: Sets a default value if environment variable is not found

### Multiple Environment Variables

You can specify multiple environment variable names separated by commas. The first one found will be used:

```go
type Config struct {
    URL string `env:"PRIMARY_URL,SECONDARY_URL,FALLBACK_URL"`
}
```

## Custom Types

You can implement custom unmarshaling by implementing the `Unmarshaler` interface:

```go
type DatabaseConfig struct {
    ConnectionString string
}

func (d *DatabaseConfig) UnmarshalEnvironmentValue(data string) error {
    // Custom parsing logic
    d.ConnectionString = "parsed_" + data
    return nil
}

type Config struct {
    DB DatabaseConfig `env:"DATABASE_CONFIG"`
}
```

## Error Handling

The package provides specific error types for different scenarios:

### ErrMissingRequiredValue

Returned when a required field is not found in the environment:

```go
if err != nil {
    if missingErr, ok := err.(*env.ErrMissingRequiredValue); ok {
        fmt.Printf("Missing required environment variable: %s\n", missingErr.Value)
    }
}
```

### ErrInvalidValue

Returned when the input value is invalid (e.g., not a pointer to struct):

```go
cfg := Config{} // Should be &Config{}
_, err := env.Unmarshal(cfg) // This will return ErrInvalidValue
```

### ErrUnsupportedField

Returned when trying to set an unexported field:

```go
type Config struct {
    host string `env:"HOST"` // lowercase = unexported, will cause error
}
```

## Advanced Usage

### Environment Overrides

You can programmatically override environment variables:

```go
override := env.Override{
    "PORT": stringPtr("8080"),
    "DEBUG": stringPtr("true"),
}

es := make(env.EnvSet)
_, err := es.Apply(override, &config)
```

### Working with EnvSet

You can create and manipulate environment sets directly:

```go
// Create from current environment
es, err := env.EnvToEnvSet(os.Environ())

// Create manually
es := env.EnvSet{
    "HOST": "localhost",
    "PORT": "8080",
}

// Unmarshal from EnvSet
err := env.UnmarshalFromEnvSet(es, &config)
```

## Best Practices

1. **Use required fields for critical configuration**:
   ```go
   type Config struct {
       APIKey string `env:"API_KEY,required"`
       Host   string `env:"HOST,default=localhost"`
   }
   ```

2. **Provide sensible defaults**:
   ```go
   type Config struct {
       Port    int  `env:"PORT,default=8080"`
       Timeout time.Duration `env:"TIMEOUT,default=30s"`
   }
   ```

3. **Use multiple env names for backward compatibility**:
   ```go
   type Config struct {
       URL string `env:"NEW_URL,OLD_URL,LEGACY_URL"`
   }
   ```

4. **Group related configuration in nested structs**:
   ```go
   type Config struct {
       Server struct {
           Host string `env:"SERVER_HOST"`
           Port int    `env:"SERVER_PORT"`
       }
       Database struct {
           URL      string `env:"DB_URL,required"`
           MaxConns int    `env:"DB_MAX_CONNS,default=10"`
       }
   }
   ```

## Supported Types

- `string`
- `bool` (accepts: true/false, 1/0, yes/no, on/off)
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `time.Duration` (e.g., "1h30m", "5s", "100ms")
- Pointer types of above
- Custom types implementing `Unmarshaler` interface

## Testing

The package includes comprehensive tests. To run them:

```bash
go test -v ./pkg/env/...
```

For mock generation:

```bash
go generate ./pkg/env/...
```

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This package is released under the MIT License.
