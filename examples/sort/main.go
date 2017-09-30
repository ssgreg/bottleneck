package main

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/ssgreg/bottleneck"
)

func main() {

	for i := 1000; i < 1000000000; i *= 10 {
		// Our bc has 2 bottleneck entries:
		// 0 - array creation and filling with random numbers
		// 1 - sorting
		bc := bottleneck.NewCalculator()

		// Will add array creation and filling time slice to the
		// bottleneck entry (0).
		bc.TimeSlice(bottleneck.Index0)

		ra := make([]int, i)
		for j := 0; j < i; j++ {
			ra[j] = rand.Int()
		}

		// Will add sorting time slice to the bottleneck entry (1).
		bc.TimeSlice(bottleneck.Index1)

		sort.Ints(ra)

		entries := bc.Stats()
		fmt.Printf("array len is %v, creating new random array takes %v, sorting takes %v\n", i, entries[0].Duration, entries[1].Duration)
	}

}
