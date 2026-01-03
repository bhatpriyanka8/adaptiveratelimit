package main

import (
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/bhatpriyanka8/adaptiveratelimit"
	adaptgrpc "github.com/bhatpriyanka8/adaptiveratelimit/grpc"
)

func main() {
	cfg := adaptiveratelimit.AdaptiveConfig{
		TargetLatency: 200 * time.Millisecond,
		MaxErrorRate:  0.05,
		IncreaseStep:  1,
		DecreaseStep:  2,
		MinLimit:      1,
		MaxLimit:      100,
		Cooldown:      2 * time.Second,
	}

	limiter := adaptiveratelimit.NewAdaptivePerSecond(10, cfg)
	defer limiter.Stop()

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			adaptgrpc.UnaryServerInterceptor(limiter),
		),
	)

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("gRPC server listening on :50051")
	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
