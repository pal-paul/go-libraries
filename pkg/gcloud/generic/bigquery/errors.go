package bigquery

import (
	"fmt"
)

type ErrProjectIdBlank struct {
	Value string
}

func (e ErrProjectIdBlank) Error() string {
	return fmt.Sprintf("project id can't be blank [%s]", e.Value)
}

type ErrFailedToCreateClient struct {
	Value string
}

func (e ErrFailedToCreateClient) Error() string {
	return fmt.Sprintf("failed to create a new bigquery client [%s]", e.Value)
}

type ErrInvalidClient struct {
	Value string
}

func (e ErrInvalidClient) Error() string {
	return fmt.Sprintf("invalid client, bigquery client is not initialized [%s]", e.Value)
}

type ErrInvalidDataset struct {
	Value string
}

func (e ErrInvalidDataset) Error() string {
	return fmt.Sprintf("invalid dataset, dataset id can't be blank [%s]", e.Value)
}

type ErrInvalidTable struct {
	Value string
}

func (e ErrInvalidTable) Error() string {
	return fmt.Sprintf("invalid table, table id can't be blank [%s]", e.Value)
}

type ErrFailedToImport struct {
	Value string
}

func (e ErrFailedToImport) Error() string {
	return fmt.Sprintf("failed to import data to bigquery [%s]", e.Value)
}

type ErrFailedToAppend struct {
	Value string
}

func (e ErrFailedToAppend) Error() string {
	return fmt.Sprintf("failed to append data to bigquery table [%s]", e.Value)
}

type ErrInvalidQuery struct {
	Value string
}

func (e ErrInvalidQuery) Error() string {
	return fmt.Sprintf("invalid query: %s", e.Value)
}

type ErrQueryExecution struct {
	Value string
}

func (e ErrQueryExecution) Error() string {
	return fmt.Sprintf("failed to execute query: %s", e.Value)
}

type ErrInvalidGCSFile struct {
	Value string
}

func (e ErrInvalidGCSFile) Error() string {
	return fmt.Sprintf("invalid GCS file path: %s", e.Value)
}

type ErrFailedToRead struct {
	Value string
}

func (e ErrFailedToRead) Error() string {
	return fmt.Sprintf("failed to read data: %s", e.Value)
}
