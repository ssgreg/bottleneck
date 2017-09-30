package bottleneck

import (
	"sync/atomic"
	"time"
)

// Entry describes a bottleneck.
type Entry struct {
	Duration   time.Duration
	CallCount  int64
	Percentage float64
	timestamp  time.Duration
}

// Index describes a concrete bottleneck entry.
type Index int

// Possible Index values.
const (
	Index0 Index = iota
	Index1
	Index2
	Index3
	Index4
	Index5
	Index6
	Index7
	IndexCount
)

// Calculator allows to find bottlenecks in your code.
type Calculator struct {
	entries []Entry
	last    uint64
}

// NewCalculator creates an instance of Calculator.
func NewCalculator() *Calculator {
	return &Calculator{make([]Entry, int(IndexCount)), 0}
}

// TimeSlice starts a new time slice and adds current time slice to
// the bottleneck specified by the index in the previous call.
// Must be always called from the same goroutine.
func (c *Calculator) TimeSlice(index Index) {
	n := uint64(time.Now().UnixNano())
	last := atomic.SwapUint64(&c.last, n|uint64(index)<<61)
	if last != 0 {
		lastIndex := last >> 61
		last &^= 7 << 61
		atomic.AddInt64((*int64)(&c.entries[lastIndex].Duration), int64(n-last))
		atomic.AddInt64(&c.entries[lastIndex].CallCount, 1)
	}
}

// Stats returns bottleneck entries.
// Could be called from the different goroutines.
func (c *Calculator) Stats() []Entry {
	n := uint64(time.Now().UnixNano())
	for {
		last := atomic.LoadUint64(&c.last)
		if last == 0 {
			break
		}
		index := last >> 61
		if atomic.CompareAndSwapUint64(&c.last, last, n|index<<61) {
			atomic.AddInt64((*int64)(&c.entries[index].Duration), int64(n-last&^(7<<61)))
			break
		}
	}

	var total time.Duration
	result := make([]Entry, len(c.entries))
	for i := 0; i < len(c.entries); i++ {
		result[i].Duration = time.Duration(atomic.LoadInt64((*int64)(&c.entries[i].Duration)))
		result[i].CallCount = atomic.LoadInt64(&c.entries[i].CallCount)
		total += result[i].Duration
	}
	for i := 0; i < len(c.entries); i++ {
		result[i].Percentage = float64(result[i].Duration) / float64(total)
	}
	return result
}
