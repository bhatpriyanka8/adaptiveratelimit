package main

import (
	"net/http"
	"time"

	"github.com/bhatpriyanka8/adaptiveratelimit"
	adapthttp "github.com/bhatpriyanka8/adaptiveratelimit/http"
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

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Write([]byte("ok"))
	})

	http.ListenAndServe(
		":8080",
		adapthttp.Middleware(limiter)(handler),
	)
}
