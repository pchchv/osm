package osmapi

import "context"

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
