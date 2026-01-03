package adaptiveratelimit

import "sync"

// EWMA implements an exponentially weighted moving average.
//
// EWMA is used to smooth noisy signals such as latency and error rates.
// It is safe for concurrent use.
type EWMA struct {
	// unexported fields
	mu    sync.Mutex
	alpha float64
	value float64
	init  bool
}

// NewEWMA creates a new EWMA with the given smoothing factor alpha.
// Alpha must be between 0 and 1, where lower values result in
// heavier smoothing.
func NewEWMA(alpha float64) *EWMA {
	return &EWMA{
		alpha: alpha,
	}
}

// Update incorporates a new sample into the moving average.
func (e *EWMA) Update(sample float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.init {
		e.value = sample
		e.init = true
		return
	}

	e.value = e.alpha*sample + (1-e.alpha)*e.value
}

// Value returns the current EWMA value.
func (e *EWMA) Value() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.value
}
