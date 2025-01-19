package error_handling

import (
	"fmt"
	"net/http"
)

// here we create a type that income passes both the *http.Response and the error
// possible from an iteration of the loop within a goroutine
type Result struct {
	Error    error
	Response *http.Response
}

// checkStatus returns the channel that can be read from to retrieve results of an iteration of our loop
func checkStatus(done <-chan any, urls ...string) <-chan Result {
	results := make(chan Result)
	go func() {
		defer close(results)

		for _, url := range urls {
			var result Result
			resp, err := http.Get(url)

			// here we create, the result instance with the Error and Response fields set
			result = Result{Error: err, Response: resp}

			select {
			case results <- result:
			case <-done:
			}
		}
	}()

	return results
}

func errorHandlingExec_example1() {
	done := make(chan any)
	defer close(done)

	urls := []string{"https://www.google.com", "https://www.badass"}

	for result := range checkStatus(done, urls...) {

		// in our "main" goroutine, we are able to deal with error coming out of the goroutine started by checkStatus intelligently,
		// and with the full context of the larger program
		if result.Error != nil {
			fmt.Printf("error: %v", result.Error)
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}

func errorHandlingExec_example2() {
	done := make(chan any)
	defer close(done)

	urls := []string{"https://www.google.com", "https://www.badass", "a", "b", "c"}

	var errCount int
	for result := range checkStatus(done, urls...) {

		// in our "main" goroutine, we are able to deal with error coming out of the goroutine started by checkStatus intelligently,
		// and with the full context of the larger program
		if result.Error != nil {
			fmt.Printf("error: %v", result.Error)
			errCount++

			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}
