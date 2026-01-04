// Package adaptiveratelimit provides an adaptive, instance-local rate limiter
// for Go services.
//
// The limiter dynamically adjusts the allowed request rate based on observed
// latency and error signals, using EWMA smoothing and a background control loop.
// It is designed to provide safe backpressure while avoiding oscillation.
//
// This package is intended for single-process use and does not coordinate
// limits across multiple instances or nodes.
package adaptiveratelimit

import (
	"sync"
	"time"
)

// AdaptiveConfig defines the configuration parameters that control
// how the limiter adapts over time.
//
// The limiter increases capacity gradually when the system is healthy
// and backs off more aggressively when latency or error thresholds
// are exceeded.
type AdaptiveConfig struct {
	// TargetLatency is the desired average request latency.
	// Sustained latency above this value will cause the limiter to reduce capacity.
	TargetLatency time.Duration

	// MaxErrorRate is the maximum acceptable error rate (0.0–1.0).
	// Sustained error rates above this threshold will cause backoff.
	MaxErrorRate float64

	// IncreaseStep controls how much the limit is increased when the
	// system is healthy.
	IncreaseStep int

	// DecreaseStep controls how much the limit is reduced when the
	// system is under stress.
	DecreaseStep int

	// MinLimit is the lower bound on the allowed rate.
	MinLimit int

	// MaxLimit is the upper bound on the allowed rate.
	MaxLimit int

	// Cooldown specifies the minimum duration between consecutive
	// limit adjustments. This helps prevent oscillation.
	Cooldown time.Duration
}

// Limiter is an adaptive rate limiter that adjusts its throughput
// based on observed latency and error signals.
//
// Limiter is safe for concurrent use.
//
// Callers should invoke Allow before processing a request, and
// Record after the request completes to provide feedback to the
// control loop.
type Limiter struct {
	// unexported fields
	mu             sync.Mutex
	baseLimit      int
	currentLimit   int
	count          int
	lastReset      time.Time
	lastAdjustment time.Time

	latencyEWMA *EWMA
	errorEWMA   *EWMA

	cfg AdaptiveConfig

	stopCh chan struct{}
}

// NewAdaptivePerSecond creates a new adaptive rate limiter that
// starts at the given initial rate (requests per second) and
// adjusts over time using the provided configuration.
//
// The returned Limiter starts a background control loop and should
// be stopped by calling Stop when no longer needed.
func NewAdaptivePerSecond(limit int, cfg AdaptiveConfig) *Limiter {
	limiter := &Limiter{
		baseLimit:    limit,
		currentLimit: limit,
		lastReset:    time.Now(),
		cfg:          cfg,
		latencyEWMA:  NewEWMA(0.3),
		errorEWMA:    NewEWMA(0.2),
		stopCh:       make(chan struct{}),
	}
	limiter.startResetLoop()
	limiter.startAdaptiveLoop()
	return limiter
}

// Allow reports whether a request is allowed under the current rate limit.
//
// If Allow returns false, the caller should reject the request
// immediately (for example, by returning HTTP 429).
//
// Allow is safe to call concurrently and is designed to be lightweight.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.count >= l.currentLimit {
		return false
	}

	l.count++
	return true
}

func (l *Limiter) startResetLoop() {
	ticker := time.NewTicker(time.Second)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				l.mu.Lock()
				l.count = 0
				l.lastReset = time.Now()
				l.mu.Unlock()
			case <-l.stopCh:
				return
			}
		}
	}()
}

func (l *Limiter) startAdaptiveLoop() {
	ticker := time.NewTicker(time.Second)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				l.mu.Lock()

				now := time.Now()
				if now.Sub(l.lastAdjustment) < l.cfg.Cooldown {
					l.mu.Unlock()
					continue
				}

				avgLatency := time.Duration(l.latencyEWMA.Value()) * time.Millisecond
				errorRate := l.errorEWMA.Value()

				if avgLatency > l.cfg.TargetLatency || errorRate > l.cfg.MaxErrorRate {
					l.decreaseLimit()
				} else {
					l.increaseLimit()
				}

				l.lastAdjustment = now
				l.mu.Unlock()

			case <-l.stopCh:
				return
			}
		}
	}()
}

// Stop terminates the limiter's background control loop and releases
// associated resources.
//
// Stop should be called when the limiter is no longer needed.
// It is safe to call Stop multiple times.
func (l *Limiter) Stop() {
	close(l.stopCh)
}

// Record records the outcome of a completed request.
//
// The provided latency is used to update internal latency estimates.
// If err is non-nil, the request is treated as a failure and contributes
// to the error rate.
//
// Callers should invoke Record once per request after processing completes.
func (l *Limiter) Record(latency time.Duration, err error) {
	l.latencyEWMA.Update(float64(latency.Milliseconds()))

	if err != nil {
		l.errorEWMA.Update(1)
	} else {
		l.errorEWMA.Update(0)
	}
}

func (l *Limiter) increaseLimit() {
	l.currentLimit += l.cfg.IncreaseStep
	if l.currentLimit > l.cfg.MaxLimit {
		l.currentLimit = l.cfg.MaxLimit
	}
}

func (l *Limiter) decreaseLimit() {
	l.currentLimit -= l.cfg.DecreaseStep
	if l.currentLimit < l.cfg.MinLimit {
		l.currentLimit = l.cfg.MinLimit
	}
}

// CurrentLimit returns the current allowed rate.
func (l *Limiter) CurrentLimit() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.currentLimit
}

// ErrorRate returns the current smoothed error rate.
//
// The returned value is between 0.0 and 1.0.
func (l *Limiter) ErrorRate() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.errorEWMA.Value()
}

// AverageLatency returns the current smoothed average request latency.
func (l *Limiter) AverageLatency() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	return time.Duration(l.latencyEWMA.Value())
}
