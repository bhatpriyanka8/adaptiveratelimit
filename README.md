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

## Quick Start

```go
import (
    "time"

    "github.com/bhatpriyanka8/adaptiveratelimit"
)

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

## Disclaimer

This project is developed and maintained in a personal capacity and
is not affiliated with or endorsed by any employer.
