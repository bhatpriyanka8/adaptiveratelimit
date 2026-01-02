package adaptiveratelimit

import "testing"

func TestEWMAConverges(t *testing.T) {
	ewma := NewEWMA(0.5)

	ewma.Update(100)
	ewma.Update(100)
	ewma.Update(100)

	if ewma.Value() < 90 || ewma.Value() > 110 {
		t.Fatalf("expected EWMA to converge near 100, got %f", ewma.Value())
	}
}

func TestEWMAReactsToChange(t *testing.T) {
	ewma := NewEWMA(0.5)

	ewma.Update(100)
	ewma.Update(100)
	ewma.Update(300)

	if ewma.Value() <= 100 {
		t.Fatalf("expected EWMA to increase after spike, got %f", ewma.Value())
	}
}
