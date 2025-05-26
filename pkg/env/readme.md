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
    
    // Required field example
    APIKey      string        `env:"API_KEY,required"`
    
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
