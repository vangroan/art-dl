package scrapers

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
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
	seedURLs := make([]string, 0)
	for _, ruleMatch := range ruleMatches {
		if ruleMatch.UserInfo == "" {
			panic("DeviantArt scraper was instantiated with rules containing no user names")
		}
		// seedURLs = append(seedURLs, fmt.Sprintf(galleryURLFmt, ruleMatch.UserInfo))

		u, err := makeRssURL(ruleMatch.UserInfo)
		if err != nil {
			log.Fatal("Error : ", err)
		}
		seedURLs = append(seedURLs, u.String())
	}

	return &DeviantArtScraper{
		baseScraper: baseScraper{
			config: config,
		},

		seeds: seedURLs,
	}
}

// GetName returns a descriptive name for the scraper.
func (s *DeviantArtScraper) GetName() string {
	return "DeviantArt"
}

// Run starts the scraper
func (s *DeviantArtScraper) Run(wg *sync.WaitGroup) error {
	defer wg.Done()

	toRssFetch := make(chan string, 128)
	toDownload := make(chan string, 128)

	defer close(toRssFetch)
	defer close(toDownload)

	// Copy URLs in goroutine so it keeps feeding fetch even if the channel is full
	go func() {
		for _, seedURL := range s.seeds {
			log.Printf("Seeding: %s\n", seedURL)

			toRssFetch <- seedURL
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case url := <-toRssFetch:
			s.fetchRss(url, toDownload)
		case url := <-toDownload:
			s.download(url)
		case t := <-ticker.C:
			log.Println("Current time: ", t)
		}
		log.Println("Looping...")
	}
}

// fetchRss retrieves the RSS XML from the url.
func (s *DeviantArtScraper) fetchRss(u string, toDownload chan string) {
	log.Println("Fetching RSS : ", u)

	// Retrieve RSS feed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(u)
	if err != nil {
		log.Println("Error : ", err)
		return
	}

	log.Println("Feed : ", feed.Title)

	// Schedule Image Downloads
	for _, item := range feed.Items {
		toDownload <- item.Content
	}
}

func (s *DeviantArtScraper) fetch(url string, toDownload chan string) {
	log.Println("Fetching : ", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Printf("%+v", resp)

	toDownload <- url
}

func (s *DeviantArtScraper) download(url string) {
	log.Println("Downloading : ", url)
}

// makeRssURL creates a URL with the appropriate query parameters
// for retrieving a user's gallery.
func makeRssURL(username string) (*url.URL, error) {
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
