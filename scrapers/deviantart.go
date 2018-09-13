package scrapers

import (
	"log"
	"sync"
	"time"

	artdl "github.com/vangroan/art-dl"
)

const (
	directory string = "deviantart"
)

type DeviantArtScraper struct {
	baseScraper

	seeds []string
}

func NewDeviantArtScraper(seedURLs []string, config *artdl.Config) artdl.Scraper {
	return &DeviantArtScraper{
		baseScraper: baseScraper{
			config: config,
		},

		seeds: seedURLs,
	}
}

func (s *DeviantArtScraper) GetName() string {
	return "DeviantArt"
}

func (s *DeviantArtScraper) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	toFetch := make(chan string, 128)
	toDownload := make(chan string, 128)

	defer close(toFetch)
	defer close(toDownload)

	go func() {
		for i := 0; i < len(s.seeds); i++ {
			log.Printf("Seeding: %s\n", s.seeds[i])
			toFetch <- s.seeds[i]
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case url := <-toFetch:
			s.fetch(url, toDownload)
		case url := <-toDownload:
			s.download(url)
		case t := <-ticker.C:
			log.Println("Current time: ", t)
			return
		}
		log.Println("Looping...")
	}
}

func (s *DeviantArtScraper) fetch(url string, toDownload chan string) {
	log.Println("Fetching: ", url)
	toDownload <- url
}

func (s *DeviantArtScraper) download(url string) {
	log.Println("Downloading: ", url)
}
