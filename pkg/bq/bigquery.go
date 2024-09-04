package bq

//go:generate mockgen -source=interface.go -destination=mocks/mock-bigquery.go -package=mocks
import (
	"context"

	bq "cloud.google.com/go/bigquery"
)

// New returns a new BigQuery
//
// Parameters:
//   - opts: []Option [The options to configure the BigQuery]
//
// Returns:
//   - BigQueryInterface[T]: A new BigQueryInterface[T].
//   - error: An error if one occurs.
func New(opts ...Option) (IBigQuery, error) {
	n := &bigQuery{}
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
