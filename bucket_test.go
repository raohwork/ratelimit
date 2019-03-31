package ratelimit

import (
	"testing"
	"time"
)

func TestNewFromRate(t *testing.T) {
	cases := []struct {
		name       string
		expectCap  int64
		expectDur  time.Duration
		expectUnit int64
		rate       int64
		burst      int64
		unit       int64
	}{
		{
			name:       "U0-R1K-B1K",
			rate:       1 * KB,
			burst:      1 * KB,
			expectCap:  1 * KB,
			expectDur:  time.Second / (1 * KB),
			expectUnit: 102,
		},
		{
			name:       "U1K-R1K-B1K",
			rate:       1 * KB,
			burst:      1 * KB,
			unit:       1 * KB,
			expectCap:  1 * KB,
			expectDur:  time.Second / (1 * KB),
			expectUnit: 1 * KB,
		},
		{
			name:       "U0-R1K-B5K",
			rate:       1 * KB,
			burst:      5 * KB,
			expectCap:  5 * KB,
			expectDur:  time.Second / (1 * KB),
			expectUnit: 102,
		},
		{
			name:       "U100B-R1K-B5K",
			rate:       1 * KB,
			burst:      5 * KB,
			unit:       100,
			expectCap:  5 * KB,
			expectDur:  time.Second / (1 * KB),
			expectUnit: 100,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := NewFromRate(c.rate, c.burst, c.unit)

			if b.capacity != c.expectCap {
				t.Errorf(
					"expect capacity: %d, actual: %d",
					c.expectCap, b.capacity,
				)
			}
			if b.fillInterval != c.expectDur {
				t.Errorf(
					"expect duration: %d, actual: %d",
					c.expectDur, b.fillInterval,
				)
			}
			if b.transferUnit != c.expectUnit {
				t.Errorf(
					"expect unit: %d, actual: %d",
					c.expectUnit, b.transferUnit,
				)
			}
		})
	}
}
