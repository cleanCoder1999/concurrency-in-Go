package leak_prevention

import (
	"fmt"
	"math/rand"
	"time"
)

func doWork(
	done <-chan any, // by convention, the done channel is the first param of the function
	strings <-chan string,
) <-chan any {
	terminated := make(chan any)

	go func() {
		defer fmt.Println("doWork exited")
		defer close(terminated)

		for {
			select {
			case s := <-strings:
				// do something interesting
				fmt.Println(s)

			// this case checks whether the done channel has been signaled. If it has, we return from the goroutine
			case <-done:
				return
			}
		}
	}()

	return terminated
}

func leakPreventionExec_blockedOnAttemptingToWrite() {
	done := make(chan any)
	terminated := doWork(done, nil)

	go func() {
		// cancel the operation after 1 second
		time.Sleep(1 * time.Second)
		fmt.Println("cancelling do work goroutine...")

		// when closing the done channel, the goroutine in the doWork function is cancelled
		close(done)
	}()

	<-terminated // this call blocks the execution until the done channel is closed
	fmt.Println("Done")
}

func newRandStream(
	done <-chan any, // by convention, the done channel is the first param of the function
) <-chan int {
	randStream := make(chan int)

	go func() {
		defer fmt.Println("doWork exited")
		defer close(randStream)

		for {
			select {
			// this case would forever write random ints into the channel
			case randStream <- rand.Int():

			// this case checks whether the done channel has been signaled. If it has, we return from the goroutine
			case <-done:
				return
			}
		}
	}()

	return randStream
}

func leakPreventionExec_blockedOnAttemptingToRead() {
	done := make(chan any)
	randStream := newRandStream(done)

	fmt.Println("3 random ints")

	for i := 0; i < 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}

	// when closing the done channel, the goroutine in the newRandStream function is cancelled
	close(done)

	// simulate ongoing work
	time.Sleep(1 * time.Second)
}
