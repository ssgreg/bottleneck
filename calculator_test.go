package bottleneck

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculator(t *testing.T) {
	c := NewCalculator()
	c.TimeSlice(0)
	time.Sleep(time.Millisecond * 100)
	c.TimeSlice(1)
	time.Sleep(time.Millisecond * 100)
	c.TimeSlice(0)
	time.Sleep(time.Millisecond * 100)

	entries := c.Stats()
	assert.EqualValues(t, 1, entries[0].CallCount)
	assert.EqualValues(t, 1, entries[1].CallCount)
	assert.InDelta(t, .666, entries[0].Percentage, .03)
	assert.InDelta(t, .333, entries[1].Percentage, .03)
	assert.InDelta(t, 200, int64(entries[0].Duration)/1000000, 10)
	assert.InDelta(t, 100, int64(entries[1].Duration)/1000000, 10)

	c.TimeSlice(0)

	entries = c.Stats()
	assert.EqualValues(t, 2, entries[0].CallCount)
	assert.EqualValues(t, 1, entries[1].CallCount)
	assert.InDelta(t, .666, entries[0].Percentage, .03)
	assert.InDelta(t, .333, entries[1].Percentage, .03)
	assert.InDelta(t, 200, int64(entries[0].Duration)/1000000, 10)
	assert.InDelta(t, 100, int64(entries[1].Duration)/1000000, 10)
}
