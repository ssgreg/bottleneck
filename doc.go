/*
Package bottleneck allows do detect bottlenecks in your code. Package uses atomic to minimize impact on performance. Bottleneck info returning by `Stats` function is eventually consistent.

Which operation takes more time?

The example below proves that creating and filling an array with random numbers takes less time than sorting.

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

The output is:

	array len is 1000, creating new random array takes 40.059µs, sorting takes 213.93µs
	array len is 10000, creating new random array takes 455.646µs, sorting takes 2.163691ms
	array len is 100000, creating new random array takes 4.297062ms, sorting takes 26.020184ms
	array len is 1000000, creating new random array takes 40.96832ms, sorting takes 307.671125ms
	array len is 10000000, creating new random array takes 390.894645ms, sorting takes 3.574017192s
	array len is 100000000, creating new random array takes 3.911881876s, sorting takes 40.516029548s

Monitoring of a goroutine execution

The example below starts two goroutines. The first one is imitating a worker with two complex jobs. The second one monitors worker's executions. It shows percentage, number of calls and total execution time per each job.

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
				//

				// Next time slice will be added to the bottleneck with index 2.
				bc.TimeSlice(bottleneck.Index2)

				// Do complex job (2). It's 2 times longer than job (1).
				time.Sleep(time.Duration(rand.Int() % 1000 * 1000000))
				//

				// Next time slice will be added to the bottleneck with index 0.
				bc.TimeSlice(bottleneck.Index0)
			}
		}()

		go func() {
			// Monitor that 10 times per second prints bottleneck stats.
			for {
				time.Sleep(time.Millisecond * 100)
				bns := bc.Stats()
				fmt.Printf("j1: %0.3f%%, %v, %d calls | j2: %0.3f%%, %v, %d calls\n", bns[1].Percentage, bns[1].Duration, bns[1].CallCount, bns[2].Percentage, bns[2].Duration, bns[2].CallCount)
			}
		}()

		wg.Wait()
	}

A fragment of output:

	j1: 1.000%, 100.260332ms, 0 calls | j2: 0.000%, 0s, 0 calls
	j1: 1.000%, 201.013773ms, 0 calls | j2: 0.000%, 0s, 0 calls
	j1: 1.000%, 301.297164ms, 0 calls | j2: 0.000%, 0s, 0 calls
	j1: 1.000%, 402.127131ms, 0 calls | j2: 0.000%, 0s, 0 calls
	j1: 0.815%, 410.234664ms, 1 calls | j2: 0.185%, 92.8835350ms, 0 calls
	j1: 0.679%, 410.234664ms, 1 calls | j2: 0.321%, 193.718235ms, 0 calls
	j1: 0.582%, 410.234664ms, 1 calls | j2: 0.418%, 294.723694ms, 0 calls
	j1: 0.509%, 410.234664ms, 1 calls | j2: 0.491%, 395.139638ms, 0 calls
	j1: 0.453%, 410.234664ms, 1 calls | j2: 0.547%, 495.800695ms, 0 calls
	j1: 0.452%, 454.419565ms, 1 calls | j2: 0.548%, 551.816673ms, 1 calls
	j1: 0.502%, 555.145899ms, 1 calls | j2: 0.498%, 551.816673ms, 1 calls
	j1: 0.543%, 655.362978ms, 1 calls | j2: 0.457%, 551.816673ms, 1 calls
	j1: 0.559%, 731.396896ms, 2 calls | j2: 0.441%, 576.023745ms, 1 calls
	j1: 0.571%, 804.403279ms, 2 calls | j2: 0.429%, 603.695591ms, 2 calls
	j1: 0.600%, 905.306768ms, 2 calls | j2: 0.400%, 603.695591ms, 2 calls
	j1: 0.625%, 1.005623715s, 2 calls | j2: 0.375%, 603.695591ms, 2 calls
	j1: 0.647%, 1.106392719s, 2 calls | j2: 0.353%, 603.695591ms, 2 calls
	j1: 0.646%, 1.169336820s, 3 calls | j2: 0.354%, 641.748760ms, 2 calls
	j1: 0.612%, 1.169336820s, 3 calls | j2: 0.388%, 741.900221ms, 2 calls
	j1: 0.581%, 1.169336820s, 3 calls | j2: 0.419%, 842.655883ms, 2 calls
	j1: 0.563%, 1.188210892s, 3 calls | j2: 0.437%, 923.956385ms, 3 calls
	j1: 0.582%, 1.289025414s, 3 calls | j2: 0.418%, 923.956385ms, 3 calls
	j1: 0.601%, 1.390096321s, 3 calls | j2: 0.399%, 923.956385ms, 3 calls
	j1: 0.591%, 1.428142697s, 4 calls | j2: 0.409%, 986.964089ms, 3 calls
	j1: 0.574%, 1.443765872s, 4 calls | j2: 0.426%, 1.072230598s, 4 calls
	j1: 0.590%, 1.544728063s, 4 calls | j2: 0.410%, 1.072230598s, 4 calls
	j1: 0.605%, 1.644750924s, 5 calls | j2: 0.395%, 1.072519408s, 4 calls
	j1: 0.584%, 1.644750924s, 5 calls | j2: 0.416%, 1.173230030s, 4 calls
	j1: 0.564%, 1.644750924s, 5 calls | j2: 0.436%, 1.273398048s, 4 calls
	j1: 0.545%, 1.644750924s, 5 calls | j2: 0.455%, 1.374261959s, 4 calls
	j1: 0.527%, 1.644750924s, 5 calls | j2: 0.473%, 1.474492000s, 4 calls
	j1: 0.527%, 1.698508948s, 5 calls | j2: 0.473%, 1.521447279s, 5 calls
*/
package bottleneck
