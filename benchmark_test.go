package adaptiveratelimit

import (
	"testing"
	"time"
)

func BenchmarkAllow(b *testing.B) {
	cfg := AdaptiveConfig{
		TargetLatency: 200 * time.Millisecond,
		MaxErrorRate:  0.05,
		IncreaseStep:  1,
		DecreaseStep:  2,
		MinLimit:      1,
		MaxLimit:      100,
		Cooldown:      time.Second,
	}

	limiter := NewAdaptivePerSecond(1000, cfg)
	defer limiter.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Allow()
	}
}

func BenchmarkRecord(b *testing.B) {
	cfg := AdaptiveConfig{
		TargetLatency: 200 * time.Millisecond,
		MaxErrorRate:  0.05,
		IncreaseStep:  1,
		DecreaseStep:  2,
		MinLimit:      1,
		MaxLimit:      100,
		Cooldown:      time.Second,
	}

	limiter := NewAdaptivePerSecond(1000, cfg)
	defer limiter.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Record(100*time.Millisecond, nil)
	}
}
