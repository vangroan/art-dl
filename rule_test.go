package artdl

import (
	"sync"
	"testing"
)

type noopScraper struct{}

func (s *noopScraper) GetName() string { return "NoOp" }

func (s *noopScraper) Run(wg *sync.WaitGroup) {}

func TestResolve(t *testing.T) {
	// Arrange
	var f RuleFactoryFunc = func(matches []string, config *Config) Scraper {
		return &noopScraper{}
	}
	resolver := NewRuleResolver()
	resolver.SetMappings(
		MapRule(`(?P<username>[a-zA-Z0-9_-]+)\.deviantart\.com`, f),
		MapRule(`(?P<username>[a-zA-Z0-9_-]+)\.artstation\.com`, f),
	)
	urls := []string{
		"https://one.deviantart.com",
		"https://two.deviantart.com",
		"https://three.artstation.com",
	}

	// Act
	scrapers := resolver.Resolve(urls)

	// Assert
	if len(scrapers) != 2 {
		t.Fatalf("Expected %d, actual %d", 2, len(scrapers))
	}
}
