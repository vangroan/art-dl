package common

// Config holds the values passed into the application from the command line
// or from files.
type Config struct {
	Directory        string
	SeedURLs         []string
	GalleryFile      string
	ConcurrencyLevel int
}
