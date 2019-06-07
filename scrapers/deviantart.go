package scrapers

import (
	"log"
	"net/http"
	"net/url"
	"sync"

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

	// TODO: Move the cancellation token out into main function
	cancel := make(chan struct{})
	defer close(cancel)

	galleryURLs := artdl.IterateStrings(s.seeds...)
	galleryItems := fetchRssStage(cancel, galleryURLs)

	for i := range galleryItems {
		log.Println("Sink: ", i)
	}

	return nil
}

// fetchRssStage is a pipeline stage that retrieves RSS documents
//  and feeds them into an output channel.
func fetchRssStage(cancel <-chan struct{}, urls <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		for u := range urls {
			items, err := fetchRss(u)

			if err != nil {
				log.Println("Error: ", err)
				continue
			}

			for _, item := range items {
				out <- item
			}
		}
	}()

	return out
}

// fetchRss retrieves the RSS XML document from the url.
//
// Returns the image URLs conatined in the feed.
func fetchRss(u string) ([]string, error) {
	log.Println("Fetching RSS : ", u)

	// Retrieve RSS feed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(u)
	if err != nil {
		return nil, err
	}

	log.Println("Feed : ", feed.Title)

	result := make([]string, 0)

	// Schedule Image Downloads
	for _, item := range feed.Items {
		log.Println(item)
		if media, ok := item.Extensions["media"]; ok {
			if content, ok := media["content"]; ok {
				if len(content) > 0 {
					if contentURL, ok := content[0].Attrs["url"]; ok {
						result = append(result, contentURL)
					} else {
						log.Println("Warning: RSS feed item 'media:content' has no child URL", feed.Title)
					}
				} else {
					log.Println("Warning: RSS feed item has no 'media:content' children", feed.Title)
				}
			} else {
				log.Println("Warning: RSS feed item 'media:content' not found for feed ", feed.Title)
			}
		} else {
			log.Println("Warning: RSS feed item 'media' not found for feed ", feed.Title)
		}
	}

	return result, nil
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
