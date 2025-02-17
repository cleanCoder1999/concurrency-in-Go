package or_done_channel

// OrDone wraps the read from a channel c with a select statement that also selects from a done channel
//
// this approach allows to work with channels from disparate parts of a system
// since you can't make any assertions about how a channel will behave.
//
// In other words, you don't know if the fact that your goroutine was canceled means the channel
// you're reading from will have been canceled. OrDone addresses this problem.
func OrDone(
	done <-chan any,
	c <-chan any,
) <-chan any {
	valStream := make(chan any)

	go func() {
		defer close(valStream)

		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}

				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()

	return valStream
}