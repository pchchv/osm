package replication

import "net/http"

// Datasource defines context around replication data requests.
type Datasource struct {
	BaseURL string // will use package level BaseURL if empty
	Client  *http.Client
}
