package adaptiveratelimit

import (
	"sync"
	"time"
)

// AdaptiveConfig
type AdaptiveConfig struct {
	TargetLatency time.Duration
	MaxErrorRate  float64

	IncreaseStep int
	DecreaseStep int

	MinLimit int
	MaxLimit int
	Cooldown time.Duration
}

// Limiter allows a fixed number of requests per second.
type Limiter struct {
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

// NewAdaptivePerSecond creates a new limiter with a given requests-per-second limit.
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

// Allow returns true if the request is allowed, false otherwise.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.count >= l.currentLimit {
		return false
	}

	l.count++
	return true
}

// Read Current Limit
func (l *Limiter) CurrentLimit() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.currentLimit
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

func (l *Limiter) Stop() {
	close(l.stopCh)
}

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
