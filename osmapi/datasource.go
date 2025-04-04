package osmapi

import (
	"context"
	"net/http"
)

// RateLimiter waits until the next allowed request.
// This interface is met by `golang.org/x/time/rate.Limiter`
// and is meant to be used with it.
// For example:
//
//	// 10 qps
//	osmapi.DefaultDatasource.Limiter = rate.NewLimiter(10, 1)
type RateLimiter interface {
	Wait(context.Context) error
}

// Datasource defines context about the http client to use to make requests.
type Datasource struct {
	// If Limiter is non-nil.
	// The datasource will wait/block until the request is allowed by the rate limiter.
	// It is good practice to use this parameter when making may
	// concurrent requests against the prod osm api.
	Limiter RateLimiter
	BaseURL string
	Client  *http.Client
}

// NewDatasource creates a Datasource using the given client.
func NewDatasource(client *http.Client) *Datasource {
	return &Datasource{
		Client: client,
	}
}
