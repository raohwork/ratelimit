// Package ratelimit helps you limit the transfer rate using Token-Bucket algorithm.
// It is rewrote from and inspired by http://github.com/juju/ratelimit
package ratelimit

import (
	"sync"
	"time"
)

// Bucket is a thread-safe rate limiter.
// It uses Token-Bucket algorithm to limit the transfer rate.
type Bucket struct {
	lastTime     time.Time
	capacity     int64
	fillInterval time.Duration
	avail        int64
	lock         sync.Mutex
	transferUnit int64
	leastTime    time.Duration
}

func (b *Bucket) fill() (ret int64) {
	now := time.Now()
	waited := now.Sub(b.lastTime)

	if waited < b.leastTime {
		time.Sleep(b.leastTime - waited)
		now = time.Now()
		waited = now.Sub(b.lastTime)
	}
	b.lastTime = now

	tokens := int64(waited / b.fillInterval)
	b.avail += tokens
	if b.avail > b.capacity {
		b.avail = b.capacity
	}

	return b.avail
}

// Take will accquire at most n tokens from bucket.
// It returns the number of tokens accquired, not more than n or capacity.
//
// Take will block until (at least) a number of tokens (transferUnit) available,
// even if n < transferUnit.
func (b *Bucket) Take(n int64) (ret int64) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if n > b.capacity {
		n = b.capacity
	}

	if n > b.avail {
		if n > b.fill() {
			n = b.avail
		}
	}

	b.avail -= n
	return n
}

// Return releases n unused tokens.
func (b *Bucket) Return(n int64) {
	if n <= 0 {
		return
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	b.avail += n
	if b.avail > b.capacity {
		b.avail = b.capacity
	}
}

// Capacity returns capacity of this bucket.
func (b *Bucket) Capacity() (ret int64) {
	return b.capacity
}

// New creates a Bucket by specifying intervals to fill a token.
//
// Parameters
//
//   - fillInterval: duration between refillings (1 token a time)
//   - capacity: bucket capacity, worked as burst speed
//   - transferUnit: allocate/refill this amount of tokens each time
//
// You should use NewFromRate() in most case.
func New(fillInterval time.Duration, capacity, transferUnit int64) (ret *Bucket) {
	if capacity < 2 {
		capacity = 2
	}
	if transferUnit <= 0 {
		transferUnit = capacity / 2
		if transferUnit < 1 {
			transferUnit = 1
		}
	}

	return &Bucket{
		lastTime:     time.Now(),
		capacity:     capacity,
		fillInterval: fillInterval,
		avail:        0,
		lock:         sync.Mutex{},
		transferUnit: transferUnit,
		leastTime:    time.Duration(transferUnit) * fillInterval,
	}
}

// basic units
const (
	KB  = 1024
	MB  = KB * 1024
	KiB = 1000
	MiB = KiB * 1000
)

// NewFromRate creates a Bucket by specifying transfer rate
//
// Parameters
//
//   - rate: transfer rate in bytes per second
//   - burst: burst rate in bytes per second
//   - transferUnit: transfer unit (see New() for detail), <= 0 will be forced to rate/10
//
// The rate is capped at 1,000,000,000 bytes/s, which fills a token every 1ns.
func NewFromRate(rate, burst, transferUnit int64) (ret *Bucket) {
	if rate > int64(time.Second) {
		rate = int64(time.Second)
	}
	if transferUnit <= 0 {
		transferUnit = rate / 10
	}

	dur := time.Second / time.Duration(rate)
	return New(dur, burst, transferUnit)
}
