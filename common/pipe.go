package common

import (
	"sync"
)

// IterateStrings takes multiple strings and
// outputs them as a channel.
func IterateStrings(strings ...string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for _, s := range strings {
			out <- s
		}
	}()

	return out
}

// MergeStrings takes multiple channels of strings and pipes them
// into a single output channel.
func MergeStrings(cancel <-chan struct{}, channels ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	output := func(c <-chan string) {
		defer wg.Done()

		for s := range c {
			out <- s
		}
	}

	// Adding to wait group must happen before spawing gorountine
	wg.Add(len(channels))

	// Spawn readers
	for _, c := range channels {
		go output(c)
	}

	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}
