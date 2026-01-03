package adaptiveratelimit

import (
	"testing"
	"time"
)

var cfg = AdaptiveConfig{
	TargetLatency: 200 * time.Millisecond,
	MaxErrorRate:  0.05,
	IncreaseStep:  1,
	DecreaseStep:  2,
	MinLimit:      1,
	MaxLimit:      100,
}

func TestLimiterAllowsUpToLimit(t *testing.T) {
	limiter := NewAdaptivePerSecond(2, cfg)
	defer limiter.Stop()

	if !limiter.Allow() {
		t.Fatal("expected first request to be allowed")
	}

	if !limiter.Allow() {
		t.Fatal("expected second request to be allowed")
	}

	if limiter.Allow() {
		t.Fatal("expected third request to be rate-limited")
	}
}

func TestLimiterResetsAfterOneSecond(t *testing.T) {
	limiter := NewAdaptivePerSecond(1, cfg)
	defer limiter.Stop()

	if !limiter.Allow() {
		t.Fatal("expected request to be allowed")
	}

	time.Sleep(1100 * time.Millisecond)

	if !limiter.Allow() {
		t.Fatal("expected limiter to reset after one second")
	}
}

func TestLimiterDecreasesLimitOnHighLatency(t *testing.T) {
	limiter := NewAdaptivePerSecond(10, cfg)
	defer limiter.Stop()

	// simulate bad latency
	for i := 0; i < 20; i++ {
		limiter.Record(500*time.Millisecond, nil)
	}

	time.Sleep(1100 * time.Millisecond)

	if limiter.CurrentLimit() >= 10 {
		t.Fatal("expected limit to decrease due to high latency")
	}
}
