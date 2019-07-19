package common

import (
	"sync"
)

// Scraper is the interface used by the application to
// start up a scarping job.
type Scraper interface {
	GetName() string
	Run(wg *sync.WaitGroup, matches []RuleMatch) error
}
