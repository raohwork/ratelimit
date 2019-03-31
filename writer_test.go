package ratelimit

import (
	"io/ioutil"
	"testing"
	"time"
)

func TestWriter(t *testing.T) {
	cases := []struct {
		name    string
		rate    int64
		burst   int64
		prefill bool // whether to enable burst
		expect  time.Duration
	}{
		{
			name:   "R1K-B1K",
			rate:   1 * KB,
			burst:  1 * KB,
			expect: 2 * time.Second,
		},
		{
			name:   "R1K-B5K",
			rate:   1 * KB,
			burst:  5 * KB,
			expect: 2 * time.Second,
		},
		{
			name:    "R1K-B5K*burst",
			rate:    1 * KB,
			burst:   5 * KB,
			expect:  0 * time.Second,
			prefill: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := NewFromRate(c.rate, c.burst, 0)
			w := b.NewWriter(ioutil.Discard)

			if c.prefill {
				b.Return(c.burst)
			}

			begin := time.Now().UnixNano()
			data := make([]byte, c.rate)
			for _ = range []int{1, 2} {
				n, err := w.Write(data)
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if int64(n) != c.rate {
					t.Fatalf("written != rate: %d", n)
				}
			}
			end := time.Now().UnixNano()
			dur := time.Duration(end - begin)
			if dur < c.expect || dur > c.expect+time.Second {
				t.Errorf("unexpected time: %d", dur)
			}
		})
	}
}
