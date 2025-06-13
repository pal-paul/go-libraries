# HTTP Client Package

A lightweight, easy-to-use HTTP client package for Go applications. This package provides a simple interface to perform common HTTP operations with support for custom headers, request bodies, and error handling.

## Installation

```bash
go get github.com/pal-paul/go-libraries/pkg/http-client
```

## Features

- Simple interface for common HTTP methods (GET, POST, PUT, DELETE)
- Custom header support
- Automatic response body reading
- Status code handling
- Error handling with custom error types
- Easy to mock for testing

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/pal-paul/go-libraries/pkg/http-client"
)

func main() {
    // Create a new HTTP client
    client := http_client.New()

    // Make a GET request
    headers := map[string]string{
        "Accept": "application/json",
    }
    
    body, statusCode, err := client.Get("https://api.example.com/data", headers)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Status: %d\nBody: %s\n", statusCode, string(body))
}
```

## API Reference

### GET Request

```go
Get(url string, headers map[string]string) ([]byte, int, error)
```

Performs an HTTP GET request.

- **Parameters**:
  - `url`: The target URL
  - `headers`: Map of request headers
- **Returns**:
  - `[]byte`: Response body
  - `int`: HTTP status code
  - `error`: Any error that occurred

Example:

```go
headers := map[string]string{
    "Authorization": "Bearer token123",
}
body, status, err := client.Get("https://api.example.com/users", headers)
```

### POST Request

```go
Post(url string, postBody []byte, headers map[string]string) ([]byte, int, error)
```

Performs an HTTP POST request.

- **Parameters**:
  - `url`: The target URL
  - `postBody`: Request body as bytes
  - `headers`: Map of request headers
- **Returns**:
  - `[]byte`: Response body
  - `int`: HTTP status code
  - `error`: Any error that occurred

Example:

```go
requestBody := []byte(`{"name": "John Doe"}`)
headers := map[string]string{
    "Content-Type": "application/json",
}
body, status, err := client.Post("https://api.example.com/users", requestBody, headers)
```

### PUT Request

```go
Put(url string, putBody []byte, headers map[string]string) ([]byte, int, error)
```

Performs an HTTP PUT request.

- **Parameters**:
  - `url`: The target URL
  - `putBody`: Request body as bytes
  - `headers`: Map of request headers
- **Returns**:
  - `[]byte`: Response body
  - `int`: HTTP status code
  - `error`: Any error that occurred

Example:

```go
requestBody := []byte(`{"name": "John Doe", "age": 30}`)
headers := map[string]string{
    "Content-Type": "application/json",
}
body, status, err := client.Put("https://api.example.com/users/123", requestBody, headers)
```

### DELETE Request

```go
Delete(url string, headers map[string]string) ([]byte, int, error)
```

Performs an HTTP DELETE request.

- **Parameters**:
  - `url`: The target URL
  - `headers`: Map of request headers
- **Returns**:
  - `[]byte`: Response body
  - `int`: HTTP status code
  - `error`: Any error that occurred

Example:

```go
headers := map[string]string{
    "Authorization": "Bearer token123",
}
body, status, err := client.Delete("https://api.example.com/users/123", headers)
```

## Error Handling

The package provides specific error types for common HTTP-related issues:

```go
if err != nil {
    switch err.(type) {
    case *http_client.RequestError:
        // Handle request formation errors
    case *http_client.ResponseError:
        // Handle response reading errors
    default:
        // Handle other errors
    }
}
```

Common error scenarios:

- Invalid URL
- Network connectivity issues
- Response reading errors
- Invalid request bodies
- Server errors (500s)
- Client errors (400s)

## Best Practices

1. **Header Management**:

   ```go
   headers := map[string]string{
       "Content-Type": "application/json",
       "Accept": "application/json",
       "User-Agent": "MyApp/1.0",
   }
   ```

2. **Error Handling**:

   ```go
   body, status, err := client.Get(url, headers)
   if err != nil {
       // Handle error
       return err
   }
   if status >= 400 {
       // Handle HTTP error status
       return fmt.Errorf("server returned status %d", status)
   }
   ```

3. **Request Body Handling**:

   ```go
   // For structured data, use encoding/json
   data := map[string]interface{}{
       "key": "value",
   }
   requestBody, err := json.Marshal(data)
   if err != nil {
       return err
   }
   ```

## Testing

The package includes a mock client for testing:

```go
import "github.com/pal-paul/go-libraries/pkg/http-client/mocks"

func TestMyFunction(t *testing.T) {
    mockClient := mocks.NewMockIHttpClient(ctrl)
    mockClient.EXPECT().
        Get("https://api.example.com", gomock.Any()).
        Return([]byte(`{"status": "ok"}`), 200, nil)
        
    // Use mockClient in your tests
}
```

To run tests:

```bash
go test -v ./...
```

## Performance Considerations

- The client reuses HTTP connections by default
- Response bodies are always fully read and closed
- Large responses should be handled with care
- Consider timeouts for production use

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This package is released under the MIT License.
