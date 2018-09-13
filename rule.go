package artdl

import (
	"regexp"
)

type RuleResolver struct {
	entries []RuleEntry
}

func NewRuleResolver() *RuleResolver {
	return &RuleResolver{
		entries: make([]RuleEntry, 0),
	}
}

func (resolver *RuleResolver) SetMappings(entries ...RuleEntry) {
	resolver.entries = entries
}

func (resolver *RuleResolver) Resolve(seedURLs []string) []Scraper {
	scrapers := make([]Scraper, 0)
	for _, rule := range resolver.entries {
		matches := make([]string, 0)

		for _, url := range seedURLs {
			if rule.pattern.MatchString(url) {
				matches = append(matches, url)
			}
		}

		if len(matches) > 0 {
			scrapers = append(scrapers, rule.factory(matches, nil))
		}
	}
	return scrapers
}

// RuleEntry maps a regex pattern to a scraper factory.
type RuleEntry struct {
	pattern *regexp.Regexp
	factory RuleFactoryFunc
}

type RuleFactoryFunc func(matches []string, config *Config) Scraper

func MapRule(pattern string, factory RuleFactoryFunc) RuleEntry {
	return RuleEntry{
		pattern: regexp.MustCompile(pattern),
		factory: factory,
	}
}
