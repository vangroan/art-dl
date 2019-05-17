package common

import (
	"sync"
	"testing"
)

type noopScraper struct{}

func (s *noopScraper) GetName() string { return "NoOp" }

func (s *noopScraper) Run(wg *sync.WaitGroup) {}

func TestResolve(t *testing.T) {
	// Arrange
	results := make([]RuleMatch, 0)
	var f RuleFactoryFunc = func(matches []RuleMatch, config *Config) Scraper {
		results = append(results, matches...)
		return &noopScraper{}
	}
	resolver := NewRuleResolver()
	resolver.SetMappings(
		MapRule(`(?P<userinfo>[a-zA-Z0-9_-]+)\.deviantart\.com`, f),
		MapRule(`(?P<userinfo>[a-zA-Z0-9_-]+)\.artstation\.com`, f),
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

	assert := func(expected, actual string) {
		if expected != actual {
			t.Fatalf("Expected %s, actual %s", expected, actual)
		}
	}
	expected := []string{"one", "two", "three"}

	for idx, match := range results {
		assert(expected[idx], match.UserInfo)
	}
}
