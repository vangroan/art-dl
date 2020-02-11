package common

import (
	"os"
	"path/filepath"
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
	l := logger.WithField("stage", "SeedGalleries")

	go func() {
		defer close(out)
		for _, match := range matches {
			if match.UserInfo == "" {
				l.WithField("userinfo", match.UserInfo).
					Fatal("Scraper was instantiated with rules containing no user names")
			}

			l.WithField("userinfo", match.UserInfo).
				Debugf("Seeding gallery")

			out <- match.UserInfo
		}
	}()

	return out
}

// EnsureExists creates an empty directory for
// the user gallery if it doesn't exist.
func EnsureExists(logger *log.Entry, directory string, cancel <-chan struct{}, usernames <-chan string) <-chan string {
	out := make(chan string)
	l := log.WithField("stage", "EnsureExists")

	go func() {
		defer close(out)

		for username := range usernames {
			l.Debugf(filepath.Join(directory, username))
			_ = os.MkdirAll(filepath.Join(directory, username), os.ModePerm)

			select {
			// Forward command
			case out <- username:
			case <-cancel:
				return
			}
		}
	}()

	return out
}
