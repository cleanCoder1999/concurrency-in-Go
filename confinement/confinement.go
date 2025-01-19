package confinement

import (
	"bytes"
	"fmt"
	"sync"
)

// ####### EXAMPLE 1 ####### - start

var chanOwner = func() <-chan int {

	// here we instantiate the channel within the lexical scope of the chanOwner function.
	// This limits the scope of the write aspect of the results channel to the closure defined below it.
	// In other words, it confines the write aspect of the channel to prevent other goroutines from writing to it.
	results := make(chan int, 5)

	go func() {
		defer close(results)
		for i := 0; i < 5; i++ {
			results <- i
		}
	}()

	return results
}

// here we receive a read-only copy of int channel.
// By declaring that the only usage we require is read access,
// we can confine usage of the channel within the consume function to only reads.
var consumer = func(results <-chan int) {
	for r := range results {
		fmt.Printf("Received: %d\n", r)
	}
	fmt.Println("Done receiving!")
}

func confinementExec_example1() {

	// here we receive the red aspect of the channel, and we are able to pass it into the consumer,
	// which can do nothing but read from it.
	// Once again, this confines the main goroutine to a read-only view of the channel.
	//
	// Note: suppose the function confinementExec will be executed in the main function
	results := chanOwner()

	consumer(results)
}

// ####### EXAMPLE 1 ####### - end

// ####### EXAMPLE 2 ####### - start

// let's take a look at an example of confinement that uses beta structure, which is not concurrent safe, and instance of bites.Buffer
//
// in this example, you can see that because printData does not close around the data slice, it cannot access it, and needs to take in a slice of bytes to operate on.
// We pass in different subsets of the slice, thus constraining the goroutines we start to only the part of the slice we are passing in.
// Because of the lexical scope, we've made it impossible(*) to do the wrong thing, and so we don't need to synchronize memory access or share data through communication.
//
// * Katherine: I'm ignoring the possibility of manually manipulating memory via the unsafe package. It's called unsafe for a reason!

var printData = func(wg *sync.WaitGroup, data []byte) {
	defer wg.Done()

	var buff bytes.Buffer
	for _, b := range data {
		fmt.Fprintf(&buff, "%c", b)
	}
	fmt.Println(buff.String())
}

func confinementExec_example2() {
	var wg sync.WaitGroup
	wg.Add(2)

	data := []byte("golang")

	go printData(&wg, data[:3])
	go printData(&wg, data[3:])

	wg.Wait()
}

// ####### EXAMPLE 2 ####### - end
