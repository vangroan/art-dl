package scrapers

import (
	artdl "github.com/vangroan/art-dl/common"
)

// baseScraper contains useful, commonly used operations.
type baseScraper struct {
	id int
	config *artdl.Config
}
