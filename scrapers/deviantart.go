package scrapers

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	artdl "github.com/vangroan/art-dl/common"
)

const (
	directory     string = "deviantart"
	galleryURLFmt string = "https://%s.deviantart.com/gallery"
	rssURL        string = "http://backend.deviantart.com/rss.xml"
)

// DeviantArtScraper scrapes galleries on deviantart.com
type DeviantArtScraper struct {
	baseScraper

	seeds []string
}

// NewDeviantArtScraper creates a new deviantart scraper
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

// makeRssURL creates a URL with the appropriate query parameters
// for retrieving a user's gallery.
func (s *DeviantArtScraper) makeRssURL(username string) (*url.URL, error) {
	u, err := url.Parse(rssURL)
	if err != nil {
		return nil, err
	}

	// Create RSS Query
	rssQuery := "gallery:" + username

	q := u.Query()
	q.Set("type", "deviation")
	q.Set("q", rssQuery)
	u.RawQuery = q.Encode()

	return u, nil
}

// GetName returns a descriptive name for the scraper.
func (s *DeviantArtScraper) GetName() string {
	return "DeviantArt"
}

// Run starts the scraper
func (s *DeviantArtScraper) Run(wg *sync.WaitGroup) error {
	defer wg.Done()

	seedErrors := make(chan error)
	toRssFetch := make(chan url.URL, 128)
	toFetch := make(chan string, 128)
	toDownload := make(chan string, 128)

	defer close(seedErrors)
	defer close(toRssFetch)
	defer close(toFetch)
	defer close(toDownload)

	// Copy URLs in goroutine so it keeps feeding fetch even if the channel is full
	go func() {
		for i, seed := range s.seeds {
			log.Printf("Seeding: %s\n", s.seeds[i])

			if u, err := url.Parse(seed); err == nil {
				toRssFetch <- *u
			} else {
				seedErrors <- err
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case err := <-seedErrors:
			return err
		case url := <-toRssFetch:
			s.fetchRss(url, toFetch)
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

// fetchRss retireves the RSS XML from the url.
func (s *DeviantArtScraper) fetchRss(u url.URL, toFetch chan string) {
	log.Println("Fetching RSS: ", u.String())

	res, err := http.Get(u.String())
	if err != nil {
		log.Println("Error : ", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("Error : ", err)
			return
		}

		var model artdl.Channel
		err = xml.Unmarshal(bytes, &model)
		if err != nil {
			log.Println("Error : ", err)
			return
		}

		for _, item := range model.Items {
			toFetch <- item.Link
		}
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
