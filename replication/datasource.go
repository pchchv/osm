package replication

import (
	"net/http"
	"time"
)

// DefaultDatasource is the Datasource used by the package level convenience functions.
var DefaultDatasource = &Datasource{
	Client: &http.Client{
		Timeout: 30 * time.Minute,
	},
}

// Datasource defines context around replication data requests.
type Datasource struct {
	BaseURL string // will use package level BaseURL if empty
	Client  *http.Client
}

// NewDatasource creates a Datasource using the given client.
func NewDatasource(client *http.Client) *Datasource {
	return &Datasource{
		Client: client,
	}
}
