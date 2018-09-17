package main

import (
	"flag"
	"log"
	"os"
	"sync"

	artdl "github.com/vangroan/art-dl"
	scrapers "github.com/vangroan/art-dl/scrapers"
)

type seedURLFlags []string

func (urls *seedURLFlags) String() string {
	return "Seed URL Flags"
}

func (urls *seedURLFlags) Set(value string) error {
	*urls = append(*urls, value)
	return nil
}

func parseFlags() artdl.Config {
	config := artdl.Config{}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var seeds seedURLFlags

	flag.StringVar(&config.Directory, "directory", cwd, "The target directory to save downloaded images. Default is current working directory.")
	flag.Var(&seeds, "gallery", "Gallery URL")

	flag.Parse()

	config.SeedURLs = seeds

	return config
}

func main() {
	log.Println("Starting...")

	// Gather configuration from command line
	config := parseFlags()
	if len(config.SeedURLs) == 0 {
		log.Println("No galleries provided!")
		return
	}
	log.Printf("Config: %+v \n", config)

	// Resolve rules
	resolver := artdl.NewRuleResolver()
	resolver.SetMappings(
		artdl.MapRule(`(?P<userinfo>[a-zA-Z0-9_-]+)\.deviantart\.com`, scrapers.NewDeviantArtScraper),
	)
	scrapers := resolver.Resolve(config.SeedURLs)

	if len(scrapers) == 0 {
		log.Println("No rules matched provided galleries!")
		return
	}

	// Run scrapers
	var wg sync.WaitGroup
	for _, scraper := range scrapers {
		log.Println("Starting up scraper")
		wg.Add(1)
		go scraper.Run(&wg)
	}
	wg.Wait()
	log.Println("Shutting down...")
}
