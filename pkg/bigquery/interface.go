package bigquery

import (
	"context"

	bq "cloud.google.com/go/bigquery"
)

type BigQueryRow map[string]bq.Value

type bigQuery[T any] struct {
	cfg    *Config
	client *bq.Client
}

type Config struct {
	// ProjectId is the Google Cloud Project ID
	ProjectId string

	// Context is the context to use for BigQuery operations
	Context context.Context
}
type Option func(cfg *Config)

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

type BigQueryInterface[T any] interface {
	// AppendMany appends a list of rows to a BigQuery table
	// Parameters:
	//   - dataSet: string [The dataset ID]
	//   - table: string [The table ID]
	//   - data: T [The data to append]
	//
	// Returns:
	//   - error: An error if one occurs.
	AppendMany(dataSet string, table string, data []T) error

	// Append adds a single row of JSON data to a BigQuery table
	// Parameters:
	//   - dataSet: string [The dataset ID]
	//   - table: string [The table ID]
	//   - data: []byte [The JSON data to add]
	//
	// Returns:
	//   - error: An error if one occurs.
	Append(dataSet string, table string, data T) error

	// ImportJsonFile loading newline-delimited JSON data from Cloud Storage to BigQuery
	// Parameters:
	//   - dataSet: string [The dataset ID]
	//   - table: string [The table ID]
	//   - gcsFile: string [The Cloud Storage file to load]
	//   - schema: bq.Schema [The schema of the data]
	//   - writeDisposition: bq.TableWriteDisposition [The write disposition]
	//
	// Returns:
	//   - error: An error if one occurs.
	ImportJsonFile(
		dataSet string,
		table string,
		gcsFile string,
		schema bq.Schema,
		writeDisposition bq.TableWriteDisposition,
	) error

	// ImportJsonFiles loading newline-delimited JSON data from Cloud Storage to BigQuery
	// Parameters:
	//   - dataSet: string [The dataset ID]
	//   - table: string [The table ID]
	//   - gcsFile: []string [The Cloud Storage files to load]
	//   - schema: bq.Schema [The schema of the data]
	//   - writeDisposition: bq.TableWriteDisposition [The write disposition]
	//
	// Returns:
	//   - error: An error if one occurs.
	ImportJsonFiles(
		dataSet string,
		table string,
		gcsFile []string,
		schema bq.Schema,
		writeDisposition bq.TableWriteDisposition,
	) error

	// ExecuteQuery executes a BigQuery query and returns the results as a list of rows
	// Parameters:
	//   - sql: string [The SQL query]
	//
	// Returns:
	//   - []T: The results of the query
	//   - error: An error if one occurs.
	ExecuteQuery(sql string) ([]T, error)
}
