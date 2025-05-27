package bigquery_test

import (
	"context"
	"testing"

	bq "cloud.google.com/go/bigquery"
	"github.com/pal-paul/go-libraries/pkg/gcloud/generic/bigquery"
	"github.com/stretchr/testify/assert"
)

type TestData struct {
	Name string `bigquery:"name"`
	Age  int    `bigquery:"age"`
}

func TestBigQueryNew(t *testing.T) {
	tests := []struct {
		name      string
		opts      []bigquery.Option
		wantError bool
	}{
		{
			name: "success with valid project ID",
			opts: []bigquery.Option{
				bigquery.WithProjectId("test-project"),
				bigquery.WithContext(context.Background()),
			},
			wantError: false,
		},
		{
			name:      "error with empty project ID",
			opts:      []bigquery.Option{},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := bigquery.New[TestData](tt.opts...)
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

func TestBigQueryAppend(t *testing.T) {
	tests := []struct {
		name      string
		dataset   string
		table     string
		data      TestData
		wantError bool
		errorType error
	}{
		{
			name:    "error with empty dataset",
			dataset: "",
			table:   "test-table",
			data: TestData{
				Name: "John",
				Age:  30,
			},
			wantError: true,
			errorType: bigquery.ErrInvalidDataset{},
		},
		{
			name:    "error with empty table",
			dataset: "test-dataset",
			table:   "",
			data: TestData{
				Name: "John",
				Age:  30,
			},
			wantError: true,
			errorType: bigquery.ErrInvalidTable{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := bigquery.New[TestData](
				bigquery.WithProjectId("test-project"),
				bigquery.WithContext(context.Background()),
			)
			assert.NoError(t, err)

			err = client.Append(tt.dataset, tt.table, tt.data)
			if tt.wantError {
				assert.Error(t, err)
				assert.IsType(t, tt.errorType, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBigQueryAppendMany(t *testing.T) {
	tests := []struct {
		name      string
		dataset   string
		table     string
		data      []TestData
		wantError bool
		errorType error
	}{
		{
			name:    "error with empty dataset",
			dataset: "",
			table:   "test-table",
			data: []TestData{
				{Name: "John", Age: 30},
				{Name: "Jane", Age: 25},
			},
			wantError: true,
			errorType: bigquery.ErrInvalidDataset{},
		},
		{
			name:    "error with empty table",
			dataset: "test-dataset",
			table:   "",
			data: []TestData{
				{Name: "John", Age: 30},
				{Name: "Jane", Age: 25},
			},
			wantError: true,
			errorType: bigquery.ErrInvalidTable{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := bigquery.New[TestData](
				bigquery.WithProjectId("test-project"),
				bigquery.WithContext(context.Background()),
			)
			assert.NoError(t, err)

			err = client.AppendMany(tt.dataset, tt.table, tt.data)
			if tt.wantError {
				assert.Error(t, err)
				assert.IsType(t, tt.errorType, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBigQueryImportJsonFile(t *testing.T) {
	tests := []struct {
		name             string
		dataset          string
		table            string
		gcsFile          string
		schema           bq.Schema
		writeDisposition bq.TableWriteDisposition
		wantError        bool
		errorType        error
	}{
		{
			name:    "error with empty dataset",
			dataset: "",
			table:   "test-table",
			gcsFile: "gs://bucket/file.json",
			schema: bq.Schema{
				{Name: "name", Type: bq.StringFieldType},
				{Name: "age", Type: bq.IntegerFieldType},
			},
			writeDisposition: bq.WriteAppend,
			wantError:        true,
			errorType:        bigquery.ErrInvalidDataset{},
		},
		{
			name:    "error with empty table",
			dataset: "test-dataset",
			table:   "",
			gcsFile: "gs://bucket/file.json",
			schema: bq.Schema{
				{Name: "name", Type: bq.StringFieldType},
				{Name: "age", Type: bq.IntegerFieldType},
			},
			writeDisposition: bq.WriteAppend,
			wantError:        true,
			errorType:        bigquery.ErrInvalidTable{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := bigquery.New[TestData](
				bigquery.WithProjectId("test-project"),
				bigquery.WithContext(context.Background()),
			)
			assert.NoError(t, err)

			err = client.ImportJsonFile(tt.dataset, tt.table, tt.gcsFile, tt.schema, tt.writeDisposition)
			if tt.wantError {
				assert.Error(t, err)
				assert.IsType(t, tt.errorType, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBigQueryImportJsonFiles(t *testing.T) {
	schema := bq.Schema{
		{Name: "name", Type: bq.StringFieldType},
		{Name: "age", Type: bq.IntegerFieldType},
	}

	tests := []struct {
		name             string
		dataset          string
		table            string
		gcsFiles         []string
		schema           bq.Schema
		writeDisposition bq.TableWriteDisposition
		wantError        bool
		errorType        error
	}{
		{
			name:             "error with empty dataset",
			dataset:          "",
			table:            "test-table",
			gcsFiles:         []string{"gs://bucket/file1.json", "gs://bucket/file2.json"},
			schema:           schema,
			writeDisposition: bq.WriteAppend,
			wantError:        true,
			errorType:        bigquery.ErrInvalidDataset{},
		},
		{
			name:             "error with empty table",
			dataset:          "test-dataset",
			table:            "",
			gcsFiles:         []string{"gs://bucket/file1.json", "gs://bucket/file2.json"},
			schema:           schema,
			writeDisposition: bq.WriteAppend,
			wantError:        true,
			errorType:        bigquery.ErrInvalidTable{},
		},
		{
			name:             "error with empty files list",
			dataset:          "test-dataset",
			table:            "test-table",
			gcsFiles:         []string{},
			schema:           schema,
			writeDisposition: bq.WriteAppend,
			wantError:        true,
			errorType:        bigquery.ErrInvalidTable{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := bigquery.New[TestData](
				bigquery.WithProjectId("test-project"),
				bigquery.WithContext(context.Background()),
			)
			assert.NoError(t, err)

			err = client.ImportJsonFiles(tt.dataset, tt.table, tt.gcsFiles, tt.schema, tt.writeDisposition)
			if tt.wantError {
				assert.Error(t, err)
				assert.IsType(t, tt.errorType, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBigQueryExecuteQuery(t *testing.T) {
	t.Run("error with empty query", func(t *testing.T) {
		client, err := bigquery.New[TestData](
			bigquery.WithProjectId("test-project"),
			bigquery.WithContext(context.Background()),
		)
		assert.NoError(t, err)

		results, err := client.ExecuteQuery("")
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.IsType(t, bigquery.ErrInvalidQuery{}, err)
	})
}
