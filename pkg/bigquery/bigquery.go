package bigquery

//go:generate mockgen -source=interface.go -destination=mocks/mock-bigquery.go -package=mocks
import (
	"context"
	"fmt"
	"strings"

	bq "cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// New returns a new BigQuery
//
// Parameters:
//   - opts: []Option [The options to configure the BigQuery]
//
// Returns:
//   - BigQueryInterface[T]: A new BigQueryInterface[T].
//   - error: An error if one occurs.
func New[T any](opts ...Option) (IBigQuery[T], error) {
	n := &bigQuery[T]{}
	for _, opt := range opts {
		opt(n.cfg)
	}
	if n.cfg.ProjectId == "" {
		return nil, errProjectIdBlank
	}
	ctx := context.Background()
	if n.cfg.Context != nil {
		ctx = n.cfg.Context
	}
	client, err := bq.NewClient(ctx, n.cfg.ProjectId)
	if err != nil {
		return nil, errFailedToCreateClient
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
		return errInvalidDataset
	}
	if table == "" {
		return errInvalidTable
	}
	if b.client == nil {
		return errInvalidClient
	}

	ins := b.client.Dataset(dataSet).Table(table).Inserter()
	err := ins.Put(b.cfg.Context, data)
	if err != nil {
		return err
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
		return errInvalidDataset
	}
	if table == "" {
		return errInvalidTable
	}
	if b.client == nil {
		return errInvalidClient
	}
	ins := b.client.Dataset(dataSet).Table(table).Inserter()
	err := ins.Put(b.cfg.Context, data)
	if err != nil {
		return err
	}
	return nil
}

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
func (b *bigQuery[T]) ImportJsonFile(
	dataSet string,
	table string,
	gcsFile string,
	schema bq.Schema,
	writeDisposition bq.TableWriteDisposition,
) error {
	return b.ImportJsonFiles(dataSet, table, []string{gcsFile}, schema, writeDisposition)
}

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
func (b *bigQuery[T]) ImportJsonFiles(
	dataSet string,
	table string,
	gcsFile []string,
	schema bq.Schema,
	writeDisposition bq.TableWriteDisposition,
) error {
	if dataSet == "" {
		return errInvalidDataset
	}
	if table == "" {
		return errInvalidTable
	}
	if b.client == nil {
		return errInvalidClient
	}

	gcsRef := bq.NewGCSReference(gcsFile...)
	gcsRef.SourceFormat = bq.JSON
	gcsRef.Schema = schema
	loader := b.client.Dataset(dataSet).Table(table).LoaderFrom(gcsRef)
	loader.WriteDisposition = writeDisposition

	job, err := loader.Run(b.cfg.Context)
	if err != nil {
		return err
	}
	status, err := job.Wait(b.cfg.Context)
	if err != nil {
		return err
	}

	if status.Err() != nil {
		var errs []string
		for _, errDetail := range status.Errors {
			errs = append(errs, fmt.Sprintf("%s: %s", errDetail.Reason, errDetail.Message))
		}
		return fmt.Errorf("job completed with error: %v", strings.Join(errs, ","))
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
	var t []T
	if b.client == nil {
		client, err := bq.NewClient(b.cfg.Context, b.cfg.ProjectId)
		if err != nil {
			return nil, errFailedToCreateClient
		}
		b.client = client
		defer client.Close()
	}
	query := b.client.Query(sql)
	rowIterator, err := query.Read(b.cfg.Context)
	if err != nil {
		return t, err
	}
	for {
		var row T
		err := rowIterator.Next(&row)
		if err == iterator.Done {
			return t, nil
		}
		if err != nil {
			return t, fmt.Errorf("error iterating through results: %v", err)
		}
		t = append(t, row)
	}
}
