package replication

import (
	"fmt"
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

// UnexpectedStatusCodeError is return for a non 200 or 404 status code.
type UnexpectedStatusCodeError struct {
	Code int
	URL  string
}

// Error returns an error message with some information.
func (e *UnexpectedStatusCodeError) Error() string {
	return fmt.Sprintf("replication: unexpected status code of %d for url %s", e.Code, e.URL)
}

// NotFound will return try if the error from one of the
// methods was due to the file not found on the remote host.
func NotFound(err error) bool {
	if err != nil {
		if e, ok := err.(*UnexpectedStatusCodeError); ok {
			return e.Code == http.StatusNotFound
		}
	}

	return false
}
