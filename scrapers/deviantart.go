package scrapers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	artdl "github.com/vangroan/art-dl"
)

const (
	directory     string = "deviantart"
	galleryURLFmt string = "https://%s.deviantart.com/gallery"
)

type DeviantArtScraper struct {
	baseScraper

	seeds []string
}

func NewDeviantArtScraper(ruleMatches []artdl.RuleMatch, config *artdl.Config) artdl.Scraper {
	seedURLs := make([]string, len(ruleMatches))
	for _, ruleMatch := range ruleMatches {
		if ruleMatch.UserInfo == "" {
			panic("DeviantArt scraper was instantiated with rules containing to user names")
		}
		seedURLs = append(seedURLs, fmt.Sprintf(galleryURLFmt, ruleMatch.UserInfo))
	}

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

	// Copy URLs in goroutine so it keeps feeding fetch even if the channel is full
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
		}
		log.Println("Looping...")
	}
}

func (s *DeviantArtScraper) fetch(url string, toDownload chan string) {
	log.Println("Fetching: ", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Printf("%+v", resp)

	toDownload <- url
}

func (s *DeviantArtScraper) download(url string) {
	log.Println("Downloading: ", url)
}
