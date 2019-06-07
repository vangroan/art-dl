package common

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
