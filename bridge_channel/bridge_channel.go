package bridge_channel

import (
	ordone "concurrency-patterns/or_done_channel"
	"fmt"
)

// Bridge destructures a channel of channels into a simple, single channel
//
// this technique is called bridging the channels
func Bridge(
	done <-chan any,
	chanStream <-chan <-chan any, // channel of channels
) <-chan any {
	valStream := make(chan any)

	go func() {
		defer close(valStream)

		for {
			var stream <-chan any

			select {
			case maybeStream, ok := <-chanStream:
				if !ok {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}

			for val := range ordone.OrDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()

	return valStream
}

func BridgeChannelExec() {
	genVals := func() <-chan <-chan any {
		chanStream := make(chan (<-chan any))

		go func() {
			defer close(chanStream)

			for i := 0; i < 10; i++ {
				stream := make(chan any, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()

		return chanStream
	}()

	for v := range Bridge(nil, genVals) {
		fmt.Println(v)
	}
}
