package fan_out_fan_in

import (
	"concurrency-patterns/pipeline"
	"fmt"
	mathRand "math/rand"
	"runtime"
	"sync"
	"time"
)

func exampleExec() {
	rand := func() any { return mathRand.Intn(50_000_000) }

	done := make(chan any)
	defer close(done)

	start := time.Now()

	randIntStream := pipeline.ToInt(done, pipeline.RepeatFn(done, rand))

	fmt.Println("Primes:")
	for prime := range pipeline.Take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took %v\n", time.Since(start))
}

func primeFinder(done chan any, stream <-chan int) <-chan any {
	// todo implement a long running prime find procedure
	return nil
}

func FanOutFanInExec() {
	done := make(chan any)
	defer close(done)

	start := time.Now()

	rand := func() any { return mathRand.Intn(50_000_000) }

	randIntStream := pipeline.ToInt(done, pipeline.RepeatFn(done, rand))

	// returns the number of logical CPUs usable by the current process.
	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders\n", numFinders)

	finders := make([]<-chan any, numFinders)

	fmt.Println("Primes:")

	// FANNING OUT
	//
	// "fanning-out" means starting multiple goroutines to handle input from a pipeline
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}

	// FANNING IN
	for prime := range pipeline.Take(done, FanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took %v\n", time.Since(start))
}

// FanIn joins multiple streams of data into a single stream
//
// it does so by leveraging the fanning-in pattern
//
// "fanning-in" means multiplexing or joining together multiple streams of data into a single stream
func FanIn(
	done chan any,
	channels ...<-chan any,
) <-chan any {
	var wg sync.WaitGroup
	multiplexedStream := make(chan any)

	multiplex := func(c <-chan any) {
		defer wg.Done()

		for i := range c {
			select {
			case <-done:
			case multiplexedStream <- i:
			}
		}
	}

	// select from all the channels
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	// wait for all the reads to complete
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}