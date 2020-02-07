package scrapers

import (
	"sync"

	artdl "github.com/vangroan/art-dl/common"
)

// ArtStationScraper scrapers gallerie on artstation.com
type ArtStationScraper struct {
	baseScraper
}

// NewArtStationScraper creates a new deviantart scraper
func NewArtStationScraper(id int, config *artdl.Config) artdl.Scraper {
	return &ArtStationScraper{
		baseScraper: baseScraper{
			id:     id,
			config: config,
		},
	}
}

// GetName returns a descriptive name for the scraper.
func (s *ArtStationScraper) GetName() string {
	return "ArtStation"
}

// Run starts the scraper
func (s *ArtStationScraper) Run(wg *sync.WaitGroup, matches []artdl.RuleMatch) error {
	defer wg.Done()

	return nil
}
