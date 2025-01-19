package or_channel

import (
	"fmt"
	"time"
)

//var or = func(channels ...<-chan any) <-chan any

// or takes a variadic slice of channels and returns a single channel
func or(channels ...<-chan any) <-chan any {

	switch len(channels) {

	// since this is a recursive function, we must set up termination criteria:
	//
	// first, instead of the very attic slice is empty, we simply return a mill channel.
	//This is consistent with the idea of passing no channels; we wouldn't expect a composite channel to do anything.
	case 0:
		return nil
		// our second termination criteria states that if our variadic slice only contains one element, we just return that element
	case 1:
		return channels[0]
	}

	orDone := make(chan any)

	// main body of the function where recursion happens
	go func() {
		defer close(orDone)

		switch len(channels) {

		// because of how we are recursing, every recursive call to or will at least have two channels.
		// As an optimization to keep the number of goroutines constrained, we place a special case here for calls to or with only two channels.
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
				// here we recursively create an or-channel from all the channels in our slice after the third index,
				// and then select from this.
				//
				// this recurrence relation will be structure the rest of the slice into or-channels for a tree from which the first signal will return.
				//
				// we also pass in the orDone channel so that when the goroutines up the tree exit, goroutines down the tree also exit
			case <-or(append(channels[3:], orDone)...):
			}
		}
	}()

	return orDone
}

// orChannelExec is a brief example that takes channels that close after a set duration and uses the or function to combine these into a single channel that closes:
//
// notice that despite placing several generals in our call to or that takes various times to close,
// our channel that closes after one second causes the entire channel created by a call to or to close
// this is because - despite its place in a tree the or function builds - it will always close first, and thus the channels that depend on its closure will close as well
func orChannelExec() {
	sig := func(after time.Duration) <-chan any {
		c := make(chan any)

		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}
