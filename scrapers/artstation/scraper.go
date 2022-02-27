package artstation

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/mmcdole/gofeed"
	artdl "github.com/vangroan/art-dl/common"
)

// https://www.artstation.com/666kart.rss?page=1
const (
	GalleryRule      string = `www\.artstation\.com/(?P<userinfo>[a-zA-Z0-9_-]+)`
	navigationLimit  int    = 9999
	directory        string = "artstation"
	concurrencyLevel int    = 8
	rssURL           string = "https://www.artstation.com/%s.rss?page=3"
	projectPattern   string = "https://www.artstation.com/artwork/(?P<projectid>[a-zA-Z0-9_-]+)"
	projectIDKey     string = "projectid"
	projectAPIURL    string = "https://www.artstation.com/projects/%s.json"
)

// ArtStationScraper scrapers gallerie on artstation.com
type ArtStationScraper struct {
	artdl.BaseScraper
}

// NewScraper creates a new deviantart scraper
func NewScraper(id int, config *artdl.Config) artdl.Scraper {
	return &ArtStationScraper{
		BaseScraper: artdl.BaseScraper{
			ID:     id,
			Config: config,
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

	// TODO: Move the cancellation token out into main function
	cancel := make(chan struct{})
	defer close(cancel)

	seeds := seedGalleries(matches...)
	usernames := ensureExistsStage(cancel, seeds)
	projectURLs := fetchRssStage(cancel, usernames)
	filenames := fetchProjectStage(cancel, projectURLs, 0)

	for filename := range filenames {
		log.Println("Done:", filename)
	}

	// filenames := make([]<-chan string, 0)
	// for i := 0; i < concurrencyLevel; i++ {
	// 	// Avoid conflicting IDs with other scrapers by offsetting
	// 	// download worker ID by scraper's ID and expected number
	// 	// of downloaders.
	// 	id := s.ID*concurrencyLevel + i
	// 	filenames = append(filenames, fetchProjectStage(cancel, projectURLs, id))
	// }

	// for filename := range artdl.MergeStrings(cancel, filenames...) {
	// 	log.Println("Done:", filename)
	// }

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
				log.Fatal("Artstation scraper was instantiated with rules containing no user names")
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
	USERS:
		for username := range usernames {

			// Artstation's RSS feed returns maximum 50 items per request
			offset := 1
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
// Returns the page URLs of projects.
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

	// Schedule Project JSON downloads
	for _, item := range feed.Items {
		if link := item.Link; link != "" {
			result = append(result, link)
		} else {
			log.Println("Warning: RSS feed item 'link' is empty ", feed.Title)
		}
	}

	return result, nil
}

// fetchProjectStage is a pipeline stage that will retrieve
// the HTML page of the project.
func fetchProjectStage(cancel <-chan struct{}, commands <-chan downloadCommand, id int) <-chan string {
	out := make(chan string)

	// Regex to extract project identifier from page URL.
	projectRegex := regexp.MustCompile(projectPattern)
	groups := projectRegex.SubexpNames()

	go func() {
		defer close(out)

		for cmd := range commands {
			// Wrap in function to call defers
			func() {
				// URL is for project HTML page, but we need to convert
				// it to a JSON URL to call the API.
				captures := projectRegex.FindStringSubmatch(cmd.url)
				var projectID string
				for idx, group := range groups {
					if group == projectIDKey {
						projectID = captures[idx]
						break
					}
				}

				if projectID == "" {
					log.Println("Warning: Failed to extract project ID from ", cmd.url)
					return
				}

				jsonURL := fmt.Sprintf(projectAPIURL, projectID)
				log.Println("Downloading JSON ", jsonURL)

				// Download json
				r, err := http.Get(jsonURL)
				if err != nil {
					log.Println("Warning: Failed to fetch project page JSON: ", err)
					return
				}
				defer r.Body.Close()

				var data ProjectData
				err = json.NewDecoder(r.Body).Decode(&data)
				if err != nil {
					log.Println("Warning: Failed to decode JSON: ", err)
					return
				}

				for _, asset := range data.Assets {
					// Wrap in function to call defers
					func() {
						if asset.ImageUrl != "" {
							log.Println("Downloading Image ", asset.ImageUrl)
							r, err := http.Get(asset.ImageUrl)
							if err != nil {
								log.Println("Warnining: Failed downloading image ", asset.ImageUrl)
								return
							}
							defer r.Body.Close()

							filepath := downloadProjectImage(cmd.username, projectID, asset.ImageUrl)

							select {
							case out <- filepath:
							case <-cancel:
								return
							}
						} else {
							log.Println("Warning: Asset image URL is empty")
						}
					}()
				}
			}()
		}
	}()

	return out
}

func downloadProjectImage(username string, project string, url string) string {
	dir := filepath.Join(directory, username, project)
	_ = os.MkdirAll(dir, os.ModePerm)

	filepath, err := downloadFile(url, dir)
	if err != nil {
		log.Printf("Worker [%d] Warning: %s", 0, err)
		return ""
	}

	return filepath
}

// downloadFile downloads a file to the target folder. If
// a file with same name exists, the download will not
// start.
//
// Returns the file path if the download was successful,
// an error if the file already exists, or the download
// failed.
func downloadFile(fileURL string, targetFolder string) (string, error) {
	// Determine filename
	u, err := url.Parse(fileURL)
	if err != nil {
		return "", err
	}

	filename := path.Base(u.Path)
	filepath := filepath.Join(targetFolder, filename)

	// Ensure file does not exist
	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		return "", fmt.Errorf("file '%s' exists", filepath)
	}

	// Start file download
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create new file
	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Stream download into file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

type downloadCommand struct {
	url      string
	username string
}

// makeRssURL creates a URL with the appropriate query parameters
// for retrieving a user's gallery.
func makeRssURL(username string, offset int) (*url.URL, error) {
	// Format URL with username
	a := fmt.Sprintf(rssURL, username)

	u, err := url.Parse(a)
	if err != nil {
		return nil, err
	}

	// Create RSS Query
	q := u.Query()
	q.Set("page", strconv.Itoa(offset))
	u.RawQuery = q.Encode()

	return u, nil
}

type ProjectData struct {
	Assets []AssetData `json:"assets"`
}

type AssetData struct {
	ImageUrl string `json:"image_url"`
}
