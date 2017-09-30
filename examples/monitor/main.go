package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ssgreg/bottleneck"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Our bc has 3 bottleneck entries:
	// 0 - auxiliary: for loop, checks if break is needed
	// 1 - complex job (1)
	// 2 - complex job (2)
	bc := bottleneck.NewCalculator()

	go func() {
		defer wg.Done()
		finish := time.Now().Add(time.Second * 50)
		for {
			if time.Now().After(finish) {
				break
			}

			// Next time slice will be added to the bottleneck with index 1.
			bc.TimeSlice(bottleneck.Index1)
			// Do complex job (1).
			time.Sleep(time.Duration(rand.Int() % 500 * 1000000))

			// Next time slice will be added to the bottleneck with index 2.
			bc.TimeSlice(bottleneck.Index2)
			// Do complex job (2). It's 2 times longer than job (1).
			time.Sleep(time.Duration(rand.Int() % 1000 * 1000000))

			// Next time slice will be added to the bottleneck with index 0. d
			bc.TimeSlice(bottleneck.Index0)
		}
	}()

	go func() {
		// Monitor that 10 times per second prints bottleneck stats.
		for {
			time.Sleep(time.Millisecond * 100)
			bns := bc.Stats()
			fmt.Printf("w1: %0.3f%%, %v, %d calls | w2: %0.3f%%, %v, %d calls\n", bns[1].Percentage, bns[1].Duration, bns[1].CallCount, bns[2].Percentage, bns[2].Duration, bns[2].CallCount)
		}
	}()

	wg.Wait()
}
