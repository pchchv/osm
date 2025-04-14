package replication

import (
	"fmt"
	"net/http"
	"time"
)

// BaseURL defines the default planet server to hit.
const BaseURL = "https://planet.osm.org"

var (
	// DefaultDatasource is the Datasource used by the package level convenience functions.
	DefaultDatasource = &Datasource{
		Client: &http.Client{
			Timeout: 30 * time.Minute,
		},
	}
	// timeFormats contains the set of different formats we've see the time in.
	timeFormats = []string{
		"2006-01-02 15:04:05.999999999 Z",
		"2006-01-02 15:04:05.999999999 +00:00",
		"2006-01-02T15\\:04\\:05Z",
	}
)

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

func decodeTime(s string) (t time.Time, err error) {
	for _, format := range timeFormats {
		if t, err = time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return t, err
}
