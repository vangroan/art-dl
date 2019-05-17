package common

import (
	"sync"
)

type Scraper interface {
	GetName() string
	Run(wg *sync.WaitGroup) error
}
