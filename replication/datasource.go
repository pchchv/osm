package replication

import (
	"net/http"
	"time"
)

// BaseURL defines the default planet server to hit.
const BaseURL = "https://planet.osm.org"

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

func (ds Datasource) client() *http.Client {
	if ds.Client != nil {
		return ds.Client
	}

	if DefaultDatasource.Client != nil {
		return DefaultDatasource.Client
	}

	return http.DefaultClient
}

func (ds Datasource) baseURL() string {
	if ds.BaseURL != "" {
		return ds.BaseURL
	}

	return BaseURL
}
