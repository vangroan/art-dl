package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	artdl "github.com/vangroan/art-dl/common"
	scrapers "github.com/vangroan/art-dl/scrapers"
)

const (
	version string = "2019.06.07"
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

	if len(config.SeedURLs) == 0 {
		log.Println("No galleries provided!")
		return
	}

	log.Println("Starting...")

	log.Printf("Config: %+v \n", config)

	// Resolve rules
	resolver := artdl.NewRuleResolver()
	resolver.SetMappings(
		artdl.MapRule(`(?P<userinfo>[a-zA-Z0-9_-]+)\.deviantart\.com`, "deviantart", scrapers.NewDeviantArtScraper),
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
		go entry.Scraper.Run(&wg)
	}
	wg.Wait()
	log.Println("Shutting down...")
}
