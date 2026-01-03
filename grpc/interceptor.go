package grpc

import (
	"context"
	"time"

	"github.com/bhatpriyanka8/adaptiveratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a gRPC unary interceptor that
// applies adaptive rate limiting to incoming RPCs.
//
// RPCs that exceed the current limit are rejected with a
// ResourceExhausted error.
func UnaryServerInterceptor(l *adaptiveratelimit.Limiter) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if !l.Allow() {
			return nil, status.Error(429, "rate limited")
		}

		start := time.Now()
		resp, err := handler(ctx, req)
		l.Record(time.Since(start), err)

		return resp, err
	}
}
