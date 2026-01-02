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

## Example

```go
limiter := adaptiveratelimit.NewAdaptivePerSecond(10, cfg)

if !limiter.Allow() {
    return errors.New("rate limited")
}

start := time.Now()
err := doWork()
limiter.Record(time.Since(start), err)


## Disclaimer

This project is developed and maintained in a personal capacity and
is not affiliated with or endorsed by any employer.
