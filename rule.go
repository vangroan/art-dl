package artdl

import (
	"regexp"
)

const (
	UserInfo string = "userinfo"
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
		ruleMatches := make([]RuleMatch, 0)

		for _, url := range seedURLs {
			groups := rule.pattern.SubexpNames()
			matches := rule.pattern.FindStringSubmatch(url)

			for groupIdx, group := range groups {
				
			}

			if rule.pattern.MatchString(url) {
				matches = append(matches, RuleMatch{
					OrigURI: url,
					UserInfo: 
				})
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

type RuleMatch struct {
	// OrigURI is the original URI parameter that was passed in
	OrigURI string

	// UserInfo is the username used to identify a gallery in the URI
	UserInfo string
}
