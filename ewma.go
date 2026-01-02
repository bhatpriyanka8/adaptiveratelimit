package adaptiveratelimit

import "sync"

type EWMA struct {
	mu    sync.Mutex
	alpha float64
	value float64
	init  bool
}

func NewEWMA(alpha float64) *EWMA {
	return &EWMA{
		alpha: alpha,
	}
}

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

func (e *EWMA) Value() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.value
}
