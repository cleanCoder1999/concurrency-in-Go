package pipeline

import (
	"fmt"
	"math/rand"
)

// Abstract:
// this pipeline is similar to our pipeline utilizing functions in the previous example, but it is different in very important ways.
//
// 1 we're using channels. This is obvious but sign significant because it allows two things:
// 1.1 at the end of our pipeline, we can use a range statement to extract the values
// 1.2 each stage we can safely execute currently because our inputs and outputs are safe in a current contexts
//
// 2 each stage of the pipeline is executing concurrently. This means that any stage only need
// 2.1 to wait for its inputs, and
// 2.2 to be able to send it outputs
//
// 3 finally, we range over this pipeline and values are pulled through the system

// generator converts a discrete set of values into a stream of data on a channel
//
// takes in a variadic slice of integers, construct a buffered journal of integers with a length equal to the incoming integer slice, starts a goroutine, and returns the constructed channel
//
// you will see a generator function frequently when working with pipe plants because at the beginning of the pipeline, you'll always have some batch of data that you need to convert to a channel.
func Generator(done <-chan any, integers ...int) <-chan int {
	intStream := make(chan int)

	go func() {
		defer close(intStream)

		for _, i := range integers {
			select {
			case <-done:
				return
			case intStream <- i:
			}
		}
	}()

	return intStream
}

func MultiplyChannel(
	done <-chan any,
	intStream <-chan int,
	multiplier int,
) <-chan int {
	multipliedStream := make(chan int)
	go func() {
		defer close(multipliedStream)

		for i := range intStream {
			select {
			case <-done:
				return
			case multipliedStream <- i * multiplier:
			}
		}
	}()

	return multipliedStream
}

func AddChannel(
	done <-chan any,
	intStream <-chan int,
	additive int,
) <-chan int {
	addedStream := make(chan int)
	go func() {
		defer close(addedStream)

		for i := range intStream {
			select {
			case <-done:
				return
			case addedStream <- i + additive:
			}
		}
	}()

	return addedStream
}

func ChannelProcessingExec() {
	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	done := make(chan any)

	intStream := Generator(done, ints...)

	// concurrent pipeline
	pipeline := MultiplyChannel(done, AddChannel(done, MultiplyChannel(done, intStream, 2), 1), 1)

	for v := range pipeline {
		fmt.Println(v)
	}
}

// ########### some handy generators:

// Repeat will Repeat the values you passed to it infinitely until you tell it to stop
func Repeat(
	done <-chan any,
	values ...any,
) <-chan any {
	valueStream := make(chan any)

	go func() {
		defer close(valueStream)

		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()

	return valueStream
}

// RepeatFn will Repeat the call to function fn infinitely until you tell it to stop
func RepeatFn(
	done <-chan any,
	fn func() any,
	) <-chan any {
	valueStream := make(chan any)

	go func() {
		defer close(valueStream)

		for {
			select {
			case <-done: return
			case valueStream <- fn():
			}
		}
	}()

	return valueStream
}

// Take reads num numbers from valueStream and writes it into a newly created channel
func Take(
	done <-chan any,
	valueStream <-chan any,
	num int,
) <-chan any {
	takeStream := make(chan any)

	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return

			// The difference lies in how the value is received from the channel:
			//
			// 1. `case takeStream <- <-valueStream:`
			// This line receives a value from `valueStream` and then sends it to `takeStream`.
			// The double arrow `<-` is used to receive the value first, and then it's sent to `takeStream`.
			//
			// 2. `case takeStream <- valueStream:`:
			// This line is incorrect syntax in Go. The single arrow `<-` is used for receiving values from channels,
			// but it cannot be used directly in the `case` statement like this.
			//
			// The correct way to write it would be
			// `case value := <-valueStream` and then `takeStream <- value`.
			//
			// So, the correct and idiomatic way to write it would be:
			//
			//  case value := <-valueStream:
			//    takeStream <- value
			//
			// This way, you receive the value from `valueStream` into the variable `value`, and then send it to `takeStream`.
			case takeStream <- <-valueStream:
			}
		}

	}()

	return takeStream
}

// Repeat and Take can be very powerful together
func ChannelProcessingExec2() {
	done := make(chan any)
	defer close(done)

	for num := range Take(done, Repeat(done, 1), 10) {
		fmt.Println(num)
	}
}

// RepeatFn and Take can be very powerful together
func ChannelProcessingExec3() {
	done := make(chan any)
	defer close(done)

	fn := func() any {
		return rand.Int()
	}

	for num := range Take(done, RepeatFn(done, fn), 10) {
		fmt.Println(num)
	}
}

// ToString converts the values sent via valueStream to a string
func ToString(
	done <-chan any,
	valueStream <-chan any,
) <-chan string {
	stringStream := make(chan string)

	go func() {
		defer close(stringStream)

		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(string):
			}
		}
	}()

	return stringStream
}

// ToInt converts the values sent via valueStream to int
func ToInt(
	done <-chan any,
	valueStream <-chan any,
) <-chan int {
	stringStream := make(chan int)

	go func() {
		defer close(stringStream)

		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(int):
			}
		}
	}()

	return stringStream
}

// RepeatFn and Take can be very powerful together
func ChannelProcessingExec4() {
	done := make(chan any)
	defer close(done)

	var msg string
	for token := range ToString(done, Take(done, Repeat(done, "I", "am."), 5)) {
		msg += token
	}

	fmt.Println("message: ", msg)
}