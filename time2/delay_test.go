package time2

import (
	"math"
	"testing"
	"time"
)

func TestBackoffDelay(t *testing.T) {
	delay := NewBackoffDelay(500*time.Millisecond, 1*time.Second)

	start := time.Now()
	for i := 0; i < 2; i++ {
		time.Sleep(delay.NextDelay())
	}

	cost := time.Now().Sub(start)
	if math.Abs(float64(cost-1500*time.Millisecond)) > float64(100*time.Millisecond) {
		t.Fatalf("cost %s, but we want 1.5s", cost)
	}
}
