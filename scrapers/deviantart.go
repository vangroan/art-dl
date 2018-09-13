package scrapers

import (
	"sync"

	artdl "github.com/vangroan/art-dl"
)

const (
	directory string = "deviantart"
)

type DeviantArtScraper struct {
	baseScraper
}

func NewDeviantArtScraper(seedURLs []string, config *artdl.Config) artdl.Scraper {
	return &DeviantArtScraper{
		baseScraper: baseScraper{
			config: config,
		},
	}
}

func (s *DeviantArtScraper) GetName() string {
	return "DeviantArt"
}

func (s *DeviantArtScraper) Run(wg *sync.WaitGroup) {
	defer wg.Done()
}
