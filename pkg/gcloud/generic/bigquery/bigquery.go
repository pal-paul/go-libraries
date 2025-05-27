package bigquery

//go:generate mockgen -source=interface.go -destination=mocks/mock-bigquery.go -package=mocks
import (
	"fmt"
	"strings"

	bq "cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type bigQuery[T any] struct {
	cfg    *Config
	client *bq.Client
}

// New returns a new BigQuery
//
// Parameters:
//   - opts: []Option [The options to configure the BigQuery]
//
// Returns:
//   - BigQueryInterface[T]: A new BigQueryInterface[T].
//   - error: An error if one occurs.
func New[T any](opts ...Option) (IBigQuery[T], error) {
	n := &bigQuery[T]{cfg: defaultConfig()}
	for _, opt := range opts {
		opt(n.cfg)
	}
	if n.cfg.ProjectId == "" {
		return nil, ErrProjectIdBlank{Value: "project ID is required"}
	}

	client, err := bq.NewClient(n.cfg.Context, n.cfg.ProjectId)
	if err != nil {
		return nil, ErrFailedToCreateClient{Value: fmt.Sprintf("failed to create client: %v", err)}
	}
	n.client = client
	return n, nil
}

// AppendMany appends a list of rows to a BigQuery table
// Parameters:
//   - dataSet: string [The dataset ID]
//   - table: string [The table ID]
//   - data: T [The data to append]
//
// Returns:
//   - error: An error if one occurs.
func (b *bigQuery[T]) AppendMany(dataSet string, table string, data []T) error {
	if dataSet == "" {
		return ErrInvalidDataset{Value: "dataset ID is required"}
	}
	if table == "" {
		return ErrInvalidTable{Value: "table ID is required"}
	}
	if b.client == nil {
		return ErrInvalidClient{Value: "client not initialized"}
	}

	ins := b.client.Dataset(dataSet).Table(table).Inserter()
	err := ins.Put(b.cfg.Context, data)
	if err != nil {
		return ErrFailedToAppend{Value: fmt.Sprintf("failed to append multiple rows: %v", err)}
	}
	return nil
}

// Append adds a single row of JSON data to a BigQuery table
// Parameters:
//   - dataSet: string [The dataset ID]
//   - table: string [The table ID]
//   - data: []byte [The JSON data to add]
//
// Returns:
//   - error: An error if one occurs.
func (b *bigQuery[T]) Append(dataSet string, table string, data T) error {
	if dataSet == "" {
		return ErrInvalidDataset{Value: "dataset ID is required"}
	}
	if table == "" {
		return ErrInvalidTable{Value: "table ID is required"}
	}
	if b.client == nil {
		return ErrInvalidClient{Value: "client not initialized"}
	}
	ins := b.client.Dataset(dataSet).Table(table).Inserter()
	err := ins.Put(b.cfg.Context, data)
	if err != nil {
		return ErrFailedToAppend{Value: fmt.Sprintf("failed to append row: %v", err)}
	}
	return nil
}

// ImportJsonFile loads newline-delimited JSON data from Cloud Storage to BigQuery
// Parameters:
//   - dataSet: string [The dataset ID]
//   - table: string [The table ID]
//   - gcsFile: string [The Cloud Storage file to load]
//   - schema: bq.Schema [The schema of the data]
//   - writeDisposition: bq.TableWriteDisposition [The write disposition]
//
// Returns:
//   - error: An error if one occurs.
func (b *bigQuery[T]) ImportJsonFile(
	dataSet string,
	table string,
	gcsFile string,
	schema bq.Schema,
	writeDisposition bq.TableWriteDisposition,
) error {
	if dataSet == "" {
		return ErrInvalidDataset{Value: "dataset ID is required"}
	}
	if table == "" {
		return ErrInvalidTable{Value: "table ID is required"}
	}
	if b.client == nil {
		return ErrInvalidClient{Value: "client not initialized"}
	}

	gcsRef := bq.NewGCSReference(gcsFile)
	gcsRef.SourceFormat = bq.JSON
	gcsRef.Schema = schema

	loader := b.client.Dataset(dataSet).Table(table).LoaderFrom(gcsRef)
	loader.WriteDisposition = writeDisposition

	job, err := loader.Run(b.cfg.Context)
	if err != nil {
		return ErrFailedToImport{Value: fmt.Sprintf("failed to start import job: %v", err)}
	}

	status, err := job.Wait(b.cfg.Context)
	if err != nil {
		return ErrFailedToImport{Value: fmt.Sprintf("failed while waiting for import job: %v", err)}
	}

	if status.Err() != nil {
		var errors []string
		for _, e := range status.Errors {
			errors = append(errors, e.Error())
		}
		return ErrFailedToImport{Value: fmt.Sprintf("import job failed: %s", strings.Join(errors, "; "))}
	}

	return nil
}

// ImportJsonFiles loads multiple JSON files from Cloud Storage to BigQuery
// Parameters:
//   - dataSet: string [The dataset ID]
//   - table: string [The table ID]
//   - gcsFiles: []string [The Cloud Storage files to load]
//   - schema: bq.Schema [The schema of the data]
//   - writeDisposition: bq.TableWriteDisposition [The write disposition]
//
// Returns:
//   - error: An error if one occurs.
func (b *bigQuery[T]) ImportJsonFiles(
	dataSet string,
	table string,
	gcsFiles []string,
	schema bq.Schema,
	writeDisposition bq.TableWriteDisposition,
) error {
	if dataSet == "" {
		return ErrInvalidDataset{Value: "dataset ID is required"}
	}
	if table == "" {
		return ErrInvalidTable{Value: "table ID is required"}
	}
	if len(gcsFiles) == 0 {
		return ErrInvalidTable{Value: "no files provided"}
	}
	if b.client == nil {
		return ErrInvalidClient{Value: "client not initialized"}
	}

	gcsRef := bq.NewGCSReference(gcsFiles...)
	gcsRef.SourceFormat = bq.JSON
	gcsRef.Schema = schema

	loader := b.client.Dataset(dataSet).Table(table).LoaderFrom(gcsRef)
	loader.WriteDisposition = writeDisposition

	job, err := loader.Run(b.cfg.Context)
	if err != nil {
		return ErrFailedToImport{Value: fmt.Sprintf("failed to start import job: %v", err)}
	}

	status, err := job.Wait(b.cfg.Context)
	if err != nil {
		return ErrFailedToImport{Value: fmt.Sprintf("failed while waiting for import job: %v", err)}
	}

	if status.Err() != nil {
		var errors []string
		for _, e := range status.Errors {
			errors = append(errors, e.Error())
		}
		return ErrFailedToImport{Value: fmt.Sprintf("import job failed: %s", strings.Join(errors, "; "))}
	}

	return nil
}

// ExecuteQuery executes a BigQuery query and returns the results as a list of rows
// Parameters:
//   - sql: string [The SQL query]
//
// Returns:
//   - []T: The results of the query
//   - error: An error if one occurs.
func (b *bigQuery[T]) ExecuteQuery(sql string) ([]T, error) {
	if sql == "" {
		return nil, ErrInvalidQuery{Value: "SQL query cannot be empty"}
	}
	if b.client == nil {
		return nil, ErrInvalidClient{Value: "client not initialized"}
	}

	query := b.client.Query(sql)
	it, err := query.Read(b.cfg.Context)
	if err != nil {
		return nil, ErrQueryExecution{Value: fmt.Sprintf("query execution failed: %v", err)}
	}

	var results []T
	for {
		var row T
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return results, ErrFailedToImport{Value: fmt.Sprintf("failed to read row: %v", err)}
		}
		results = append(results, row)
	}

	return results, nil
}
