package tee_channel

import (
	ordone "concurrency-patterns/or_done_channel"
	"concurrency-patterns/pipeline"
	"fmt"
)

// Tee splits values coming from a channel so that you can send them off into two separate areas of your codebase
//
// it's like the tee command of Unix-like systems. you can pass it a channel to read from,
// and it will return two separate channels that will get the same value
func Tee(
	done <-chan any,
	in <- chan any,
) (_, _ <-chan any) {
	out1 := make(chan any)
	out2 := make(chan any)

	go func() {
		defer close(out1)
		defer close(out2)

		for val := range ordone.OrDone(done, in) {
			var out1, out2 = out1, out2 // intentionally shadowed

			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case out1 <- val:
					out1 = nil // once read, ensure the out1 channel will be blocked so the other channel may continue
				case out2 <- val:
					out2 = nil // once read, ensure the out2 channel will be blocked so the other channel may continue
				}
			}
		}
	}()

	return out1, out2
}


func TeeChannelExec() {
	done := make(chan any)
	defer close(done)

	out1, out2 := Tee(done, pipeline.Take(done, pipeline.Repeat(done, 1,2), 4))

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}
