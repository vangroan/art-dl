package deviantart

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/mmcdole/gofeed"
	artdl "github.com/vangroan/art-dl/common"
)

const (
	GalleryRule      string = `www\.deviantart\.com/(?P<userinfo>[a-zA-Z0-9_-]+)`
	navigationLimit  int    = 9999
	directory        string = "deviantart"
	concurrencyLevel int    = 8
	galleryURLFmt    string = "https://www.deviantart.com/%s/gallery"
	rssURL           string = "http://backend.deviantart.com/rss.xml"
)

// DeviantArtScraper scrapes galleries on deviantart.com
type DeviantArtScraper struct {
	artdl.BaseScraper
}

// NewScraper creates a new deviantart scraper
func NewScraper(id int, config *artdl.Config) artdl.Scraper {
	return &DeviantArtScraper{
		BaseScraper: artdl.BaseScraper{
			ID:     id,
			Config: config,
		},
	}
}

// GetName returns a descriptive name for the scraper.
func (s *DeviantArtScraper) GetName() string {
	return "DeviantArt"
}

// Run starts the scraper
func (s *DeviantArtScraper) Run(wg *sync.WaitGroup, matches []artdl.RuleMatch) error {
	defer wg.Done()

	// TODO: Move the cancellation token out into main function
	cancel := make(chan struct{})
	defer close(cancel)

	seeds := seedGalleries(matches...)
	usernames := ensureExistsStage(cancel, seeds)
	downloadCommands := fetchRssStage(cancel, usernames)

	filenames := make([]<-chan string, 0)
	for i := 0; i < concurrencyLevel; i++ {
		// Avoid conflicting IDs with other scrapers by offsetting
		// download worker ID by scraper's ID and expected number
		// of downloaders.
		id := s.ID*concurrencyLevel + i
		filenames = append(filenames, downloadStage(cancel, downloadCommands, id))
	}

	for filename := range artdl.MergeStrings(cancel, filenames...) {
		log.Println("Done:", filename)
	}

	return nil
}

// seedGalleries takes the matched rules and generates
// a stream of gallery usernames.
func seedGalleries(matches ...artdl.RuleMatch) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for _, match := range matches {
			if match.UserInfo == "" {
				log.Fatal("DeviantArt scraper was instantiated with rules containing no user names")
			}

			out <- match.UserInfo
		}
	}()

	return out
}

// ensureExistsStage creates an empty directory for
// the user gallery if it doesn't exist.
func ensureExistsStage(cancel <-chan struct{}, usernames <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		for username := range usernames {
			_ = os.MkdirAll(filepath.Join(directory, username), os.ModePerm)

			select {
			// Forward command
			case out <- username:
			case <-cancel:
				return
			}
		}
	}()

	return out
}

// fetchRssStage is a pipeline stage that retrieves RSS documents
// and feeds them into an output channel.
func fetchRssStage(cancel <-chan struct{}, usernames <-chan string) <-chan downloadCommand {
	out := make(chan downloadCommand)

	go func() {
		defer close(out)

	USERS:
		for username := range usernames {

			// DeviantArt's RSS feed returns maximum 60 items per request
			offset := 0
		FETCHING:
			for offset < navigationLimit {
				log.Println("Offset:", offset)

				rssURL, err := makeRssURL(username, offset)
				if err != nil {
					log.Println("Error:", err)
				}

				items, err := fetchRss(rssURL.String())

				if err != nil {
					log.Println("Error:", err)
					continue USERS
				}

				for _, item := range items {
					select {
					case out <- downloadCommand{username: username, url: item}:
					case <-cancel:
						return
					}
				}

				if len(items) > 0 {
					// Continue navigating
					offset += len(items)
				} else {
					break FETCHING
				}

			}

		}
	}()

	return out
}

// fetchRss retrieves the RSS XML document from the url.
//
// Returns the image URLs conatined in the feed.
func fetchRss(u string) ([]string, error) {
	log.Println("Fetching RSS Feed:", u)

	// Retrieve RSS feed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(u)
	if err != nil {
		return nil, err
	}

	log.Println("Feed :", feed.Title)

	result := make([]string, 0)

	// Schedule Image Downloads
	for _, item := range feed.Items {
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

type downloadCommand struct {
	url      string
	username string
}

// downloadStage is a pipeline stage that takes a channel of download commands
// and downloads the images to the target directory.
//
// Returns a channel of filepaths to the downloaded files.
func downloadStage(cancel <-chan struct{}, commands <-chan downloadCommand, id int) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		for cmd := range commands {
			log.Printf("Worker [%d] Downloading %s", id, cmd.url)

			dir := filepath.Join(directory, cmd.username)
			filepath, err := artdl.DownloadFile(cmd.url, dir, true)
			if err != nil {
				log.Printf("Worker [%d] Warning: %s", id, err)
				continue
			}

			select {
			case out <- filepath:
			case <-cancel:
				return
			}
		}
	}()

	return out
}

// makeRssURL creates a URL with the appropriate query parameters
// for retrieving a user's gallery.
func makeRssURL(username string, offset int) (*url.URL, error) {
	u, err := url.Parse(rssURL)
	if err != nil {
		return nil, err
	}

	// Create RSS Query
	rssQuery := "gallery:" + username

	q := u.Query()
	q.Set("type", "deviation")
	q.Set("q", rssQuery)
	q.Set("offset", strconv.Itoa(offset))
	u.RawQuery = q.Encode()

	return u, nil
}
