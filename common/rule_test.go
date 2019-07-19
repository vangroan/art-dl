package common

import (
	"fmt"
	"sync"
	"testing"
)

type noopScraper struct {
	runCallback func(matches []RuleMatch)
}

func (s *noopScraper) GetName() string { return "NoOp" }

func (s *noopScraper) Run(wg *sync.WaitGroup, matches []RuleMatch) error {
	s.runCallback(matches)

	return nil
}

type result struct {
	ID    int
	Match RuleMatch
}

func TestResolve(t *testing.T) {
	// Arrange
	results := make([]result, 0)
	var f RuleFactoryFunc = func(id int, config *Config) Scraper {
		return &noopScraper{
			runCallback: func(matches []RuleMatch) {
				for _, match := range matches {
					fmt.Println(id)
					results = append(results, result{
						ID:    id,
						Match: match,
					})
				}
			},
		}
	}
	resolver := NewRuleResolver()
	resolver.SetMappings(
		MapRule(`(?P<userinfo>[a-zA-Z0-9_-]+)\.deviantart\.com`, "deviantart", f),
		MapRule(`(?P<userinfo>[a-zA-Z0-9_-]+)\.artstation\.com`, "artstation", f),
	)
	urls := []string{
		"https://one.deviantart.com",
		"https://two.deviantart.com",
		"https://three.artstation.com",
	}

	// Act
	scrapers := resolver.Resolve(urls)

	for _, entry := range scrapers {
		entry.Scraper.Run(nil, entry.Seeds)
	}

	// Assert
	if len(scrapers) != 2 {
		t.Fatalf("Expected %d, actual %d", 2, len(scrapers))
	}

	assertInt := func(expected, actual int) {
		if expected != actual {
			t.Fatalf("Expected %d, actual %d", expected, actual)
		}
	}

	assertStr := func(expected, actual string) {
		if expected != actual {
			t.Fatalf("Expected %s, actual %s", expected, actual)
		}
	}

	// deviantart
	assertInt(1, results[0].ID)
	assertStr("one", results[0].Match.UserInfo)
	assertStr("two", results[1].Match.UserInfo)

	// artstation
	assertInt(2, results[2].ID)
	assertStr("three", results[2].Match.UserInfo)
}
