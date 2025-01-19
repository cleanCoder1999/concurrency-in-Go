package for_select

// ######## 1 sending iteration variables out on a channel
//
// Oftentimes you'll want to convert something that can be iterated over into values on a channel. This is nothing fancy, and usually looks something like this:

// for the sake of brevity and simplicity, done and intStream are package-level variables.
//
// AVOID package-level variables in your Go programs because it can be difficult to track the changes made to it,
// which makes it hard to understand how data is flowing through your program.
//
// you can find more details about that in Jon Bodner's "Learning Go"
var done = make(chan struct{})
var intStream = make(chan int, 5)

var sendingIterationVariablesOutAChannel = func() {
	for _, v := range []int{1, 2, 3, 4, 5} {
		select {
		case <-done:
			return
		case intStream <- v: // send the values somewhere
		}
	}
}

// ######## 2 looping infinitely waiting to be stopped
//
// It is very common to create goroutine that loop infinitely until they are stopped.
// There are a couple of variations of this one. Which one you choose is purely a stylistic preference.
//
// NOTE:
// Use default case in for-select loops with caution.
// Jon Bodner the author of "Learning Go", arguably the most comprehensive book about Go, states the following:
//
// "Having a default case inside in a for-select loop is almost always the wrong thing to do.
// It will be triggered every time through the loop when there's nothing to read or write for any of the cases.
// This makes your loop run constantly, which uses a great deal of CPU."
var loopingInfinitelyWaitingToBeStopped = func() {

	// 2.1 first variation keeps the select statement as short as possible
	//
	// if the done channel is not closed,
	// we will exit the select statement and continue on to the rest of our for loop's body
	for {
		select {
		case <-done:
			return
		default:
		}

		// do non-preemptable work
	}

	// 2.2 second variation embeds the work in a default clause of the select statement
	//
	// when we enter the select statement if the done channel has not been closed,
	// we'll execute the default clause instead
	for {
		select {
		case <-done:
			return
		default:
			// do non-preemptable work
		}
	}
}
