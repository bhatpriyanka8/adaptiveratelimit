# Adaptive Rate Limit

An adaptive rate limiter for Go that dynamically adjusts throughput
based on observed latency and error rates.

## Why

Traditional rate limiters use static limits and cannot react to
downstream health. This can cause cascading failures or underutilization.

This library adapts limits in real time using control-loop principles
similar to TCP congestion control.

## Features

- Adaptive request-per-second limits
- EWMA-based latency and error tracking
- Cooldown to prevent oscillation
- HTTP middleware and gRPC interceptor
- Clean goroutine lifecycle management

## How It Works

- Requests call Allow()
- Outcomes call Record()
- A background control loop evaluates health
- Limits are increased or decreased gradually

## Installation
```go
go get github.com/bhatpriyanka8/adaptiveratelimit
```
## Adaptive Configuration Explained

| Field            | Description |
|------------------|-------------|
| TargetLatency    | Desired average request latency. If exceeded, the limiter backs off. |
| MaxErrorRate     | Maximum acceptable error rate (0.0–1.0). |
| IncreaseStep     | How much to increase the limit when the system is healthy. |
| DecreaseStep     | How much to reduce the limit when the system is under stress. |
| MinLimit         | Lower bound on allowed requests per second. |
| MaxLimit         | Upper bound on allowed requests per second. |
| Cooldown         | Minimum duration between consecutive limit adjustments. |

The limiter increases capacity gradually when healthy and backs off faster under load.

## Quick Start

```go

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

if !limiter.Allow() {
    // reject request (e.g. HTTP 429)
    return
}

start := time.Now()
err := doWork()
limiter.Record(time.Since(start), err)
```
## Examples

For Runnable examples, refer
[HTTP Example](https://github.com/bhatpriyanka8/adaptiveratelimit/tree/main/examples/http)

[gRPC Example](https://github.com/bhatpriyanka8/adaptiveratelimit/tree/main/examples/grpc)

Go to any of these folders and just run main.go 
```
cd examples/http
go run main.go
```

go run main.go
## Disclaimer

This project is developed and maintained in a personal capacity and
is not affiliated with or endorsed by any employer.
