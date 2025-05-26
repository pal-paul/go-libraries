package bigquery

import (
	"fmt"
)

var (
	errProjectIdBlank       = fmt.Errorf("project id can't be blank")
	errFailedToCreateClient = fmt.Errorf("failed to create a new bigquery client")
	errInvalidClient        = fmt.Errorf("invalid client, bigquery client is not initialized")
	errInvalidDataset       = fmt.Errorf("invalid dataset, dataset id can't be blank")
	errInvalidTable         = fmt.Errorf("invalid table, table id can't be blank")
)
