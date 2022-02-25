package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	artdl "github.com/vangroan/art-dl/common"
	scrapers "github.com/vangroan/art-dl/scrapers"
)

const (
	version string = "2019.07.19"
)

type seedURLFlags []string

func (urls *seedURLFlags) String() string {
	return "Seed URL Flags"
}

func (urls *seedURLFlags) Set(value string) error {
	*urls = append(*urls, value)
	return nil
}

func parseFlags() (artdl.Config, bool) {
	config := artdl.Config{}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var printVersion bool
	var seeds seedURLFlags

	flag.BoolVar(&printVersion, "version", false, "Print art-dl version")
	flag.StringVar(&config.Directory, "directory", cwd, "The target directory to save downloaded images. Default is current working directory.")
	flag.Var(&seeds, "gallery", "Gallery URL")
	flag.StringVar(&config.GalleryFile, "file", "", "Gallery filename")

	flag.Parse()

	if printVersion {
		fmt.Printf("art-dl %s\n", version)
		return config, true
	}

	config.SeedURLs = seeds

	return config, false
}

func main() {
	// Gather configuration from command line

	config, close := parseFlags()
	if close {
		return
	}

	if config.GalleryFile != "" {
		urls, err := artdl.LoadGalleryFile(config.GalleryFile)
		if err != nil {
			log.Fatalln(err)
		}

		config.SeedURLs = append(config.SeedURLs, urls...)
	}

	if len(config.SeedURLs) == 0 {
		log.Fatalln("No galleries provided!")
	}

	log.Println("Starting...")

	log.Printf("Config: %+v \n", config)

	// Resolve rules
	resolver := artdl.NewRuleResolver()
	resolver.SetMappings(
		artdl.MapRule(`www\.deviantart\.com/(?P<userinfo>[a-zA-Z0-9_-]+)`, "deviantart", scrapers.NewDeviantArtScraper),
	)
	scrapers := resolver.Resolve(config.SeedURLs)

	if len(scrapers) == 0 {
		log.Println("No rules matched provided galleries!")
		return
	}

	// Run scrapers
	var wg sync.WaitGroup
	for _, entry := range scrapers {
		log.Println("Starting up scraper")
		wg.Add(1)
		go entry.Scraper.Run(&wg, entry.Seeds)
	}
	wg.Wait()
	log.Println("Shutting down...")
}
