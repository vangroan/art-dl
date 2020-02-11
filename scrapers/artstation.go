package scrapers

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	artdl "github.com/vangroan/art-dl/common"
)

// ArtStationScraper scrapes galleries on artstation.com
type ArtStationScraper struct {
	baseScraper
}

// NewArtStationScraper creates a new artstation scraper.
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
	logger := log.WithFields(log.Fields{
		"scraper": s.GetName(),
	})
	if logger == nil {
		return fmt.Errorf("Failed to create logger")
	}

	// TODO: Move the cancellation token out into main function
	cancel := make(chan struct{})
	defer close(cancel)

	seeds := artdl.SeedGalleries(logger, matches...)
	usernames := artdl.EnsureExists(logger, "artstation", cancel, seeds)
	for userinfo := range usernames {
		logger.Infof("User: %s\n", userinfo)
	}

	return nil
}
