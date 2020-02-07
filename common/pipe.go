package common

import (
	"sync"

	log "github.com/sirupsen/logrus"
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

// SeedGalleries is a pipeline producer takes
// matched rules and generates a stream of
// gallery usernames.
func SeedGalleries(logger *log.Entry, matches ...RuleMatch) <-chan string {
	out := make(chan string)
	l := logger.WithFields(log.Fields{"stage": "SeedGalleries"})

	go func() {
		defer close(out)
		for _, match := range matches {
			if match.UserInfo == "" {
				l.WithFields(log.Fields{"userinfo": match.UserInfo}).
					Fatal("Scraper was instantiated with rules containing no user names")
			}

			l.WithFields(log.Fields{"userinfo": match.UserInfo}).
				Debugf("Seeding gallery")

			out <- match.UserInfo
		}
	}()

	return out
}
