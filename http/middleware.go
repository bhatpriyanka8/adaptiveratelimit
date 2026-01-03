package http

import (
	"net/http"
	"time"

	"github.com/bhatpriyanka8/adaptiveratelimit"
)

// Middleware returns an HTTP middleware that applies adaptive
// rate limiting to incoming requests.
//
// Requests that exceed the current limit are rejected with
// HTTP status 429 (Too Many Requests).
func Middleware(l *adaptiveratelimit.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !l.Allow() {
				http.Error(w, "rate limited", http.StatusTooManyRequests)
				return
			}

			start := time.Now()
			err := func() error {
				next.ServeHTTP(w, r)
				return nil
			}()

			l.Record(time.Since(start), err)
		})
	}
}
