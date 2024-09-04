package bq

import (
	"context"

	bq "cloud.google.com/go/bigquery"
)

type Row map[string]bq.Value

type bigQuery struct {
	cfg      *Config
	client   *bq.Client
	datasets map[string]*Dataset
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

func (c *bigQuery) Dataset(name string) *Dataset {
	if c.datasets == nil {
		c.datasets = make(map[string]*Dataset)
	}
	if c.datasets[name] == nil {
		c.datasets[name] = &Dataset{
			Name:   name,
			client: c.client,
			ctx:    c.cfg.Context,
		}
	}
	return c.datasets[name]
}

type Dataset struct {
	Name   string
	table  map[string]*Table
	client *bq.Client
	ctx    context.Context
}

func (d *Dataset) Table(name string) *Table {
	if d.table == nil {
		d.table = make(map[string]*Table)
	}
	if d.table[name] == nil {
		d.table[name] = &Table{
			Name:    name,
			dataset: d,
		}
	}
	return d.table[name]
}

type IBigQuery interface{}
