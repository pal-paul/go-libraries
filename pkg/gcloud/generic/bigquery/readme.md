# BigQuery Package

A generic BigQuery client wrapper for Google Cloud BigQuery operations.

## Features

- Generic type support for strongly-typed data operations
- Simplified interface for common BigQuery operations
- Error handling with typed errors
- Support for both single and batch operations
- JSON file import capabilities
- Query execution with type-safe results

## Installation

```bash
go get github.com/pal-paul/go-libraries/pkg/gcloud/generic/bigquery
```

## Usage

### Initialize Client

```go
import (
    "context"
    "github.com/pal-paul/go-libraries/pkg/gcloud/generic/bigquery"
)

type Person struct {
    Name string `bigquery:"name"`
    Age  int    `bigquery:"age"`
}

// Create a new client
client, err := bigquery.New[Person](
    bigquery.WithProjectId("your-project-id"),
    bigquery.WithContext(context.Background()),
)
if err != nil {
    log.Fatal(err)
}
```

### Single Row Operations

```go
// Append a single row
person := Person{
    Name: "John Doe",
    Age:  30,
}
err = client.Append("dataset_id", "table_id", person)
```

### Batch Operations

```go
// Append multiple rows
people := []Person{
    {Name: "John Doe", Age: 30},
    {Name: "Jane Doe", Age: 25},
}
err = client.AppendMany("dataset_id", "table_id", people)
```

### Import JSON Files

```go
// Import a single JSON file
schema := bigquery.Schema{
    {Name: "name", Type: bigquery.StringFieldType},
    {Name: "age", Type: bigquery.IntegerFieldType},
}
err = client.ImportJsonFile(
    "dataset_id",
    "table_id",
    "gs://bucket-name/file.json",
    schema,
    bigquery.WriteAppend,
)

// Import multiple JSON files
files := []string{
    "gs://bucket-name/file1.json",
    "gs://bucket-name/file2.json",
}
err = client.ImportJsonFiles(
    "dataset_id",
    "table_id",
    files,
    schema,
    bigquery.WriteAppend,
)
```

### Execute Queries

```go
// Execute a query and get typed results
query := "SELECT name, age FROM dataset_id.table_id WHERE age > 25"
results, err := client.ExecuteQuery(query)
if err != nil {
    log.Fatal(err)
}
for _, person := range results {
    fmt.Printf("Name: %s, Age: %d\n", person.Name, person.Age)
}
```

## Error Handling

The package provides typed errors for better error handling:

```go
if err != nil {
    switch err.(type) {
    case bigquery.ErrProjectIdBlank:
        // Handle missing project ID
    case bigquery.ErrInvalidDataset:
        // Handle invalid dataset
    case bigquery.ErrInvalidTable:
        // Handle invalid table
    case bigquery.ErrInvalidQuery:
        // Handle invalid query
    case bigquery.ErrQueryExecution:
        // Handle query execution error
    case bigquery.ErrFailedToImport:
        // Handle import failure
    default:
        // Handle unknown error
    }
}
```

## Error Types

- `ErrProjectIdBlank`: Project ID is missing or empty
- `ErrFailedToCreateClient`: Failed to create BigQuery client
- `ErrInvalidClient`: Client is not properly initialized
- `ErrInvalidDataset`: Dataset ID is missing or invalid
- `ErrInvalidTable`: Table ID is missing or invalid
- `ErrFailedToImport`: Failed to import data
- `ErrFailedToAppend`: Failed to append data
- `ErrInvalidQuery`: Query is empty or invalid
- `ErrQueryExecution`: Error during query execution
- `ErrInvalidGCSFile`: Invalid Google Cloud Storage file path
- `ErrFailedToRead`: Failed to read data from BigQuery

## Best Practices

1. Always initialize the client with proper context and project ID
2. Use appropriate error handling for each operation
3. Close the client when done (uses context for cancellation)
4. Use appropriate schema definitions for imports
5. Consider batch operations for better performance
6. Use appropriate write disposition for import operations

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
