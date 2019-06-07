package common

import (
	"regexp"
)

const (
	userInfo string = "userinfo"
)

// RuleResolver maps URLs to factory functions for scrapers.
type RuleResolver struct {
	entries []RuleEntry
}

// NewRuleResolver creates a new `RuleResolver`
func NewRuleResolver() *RuleResolver {
	return &RuleResolver{
		entries: make([]RuleEntry, 0),
	}
}

// SetMappings reads multiple entries into the resolver's mapping
func (resolver *RuleResolver) SetMappings(entries ...RuleEntry) {
	resolver.entries = entries
}

// Resolve takes multiples URLs and matches them with its rule
// mappings. Each matched rule results in an instance of a scraper.
//
// Returned scrapers are instantiated using the factory functions
// given in the resolver's mapping.
func (resolver *RuleResolver) Resolve(seedURLs []string) []Scraper {
	scrapers := make([]Scraper, 0)
	for _, rule := range resolver.entries {
		result := make([]RuleMatch, 0)

		for _, url := range seedURLs {
			if rule.pattern.MatchString(url) {
				// Map capture groups for easier use
				groups := rule.pattern.SubexpNames()
				captures := rule.pattern.FindStringSubmatch(url)
				matches := make(map[string]string)
				for idx, group := range groups {
					if group == userInfo {
						matches[group] = captures[idx]
					}
				}

				ruleMatch := RuleMatch{
					OrigURI: url,
				}

				if ui, ok := matches[userInfo]; ok {
					ruleMatch.UserInfo = ui
				}

				result = append(result, ruleMatch)
			}
		}

		if len(result) > 0 {
			scrapers = append(scrapers, rule.factory(result, nil))
		}
	}
	return scrapers
}

// RuleEntry maps a regex pattern to a scraper factory.
type RuleEntry struct {
	pattern *regexp.Regexp
	factory RuleFactoryFunc
}

// RuleFactoryFunc is factory function that is expected to
// create an instance of a `Scraper`, given the parameters
// in the rule match and config.
type RuleFactoryFunc func(matches []RuleMatch, config *Config) Scraper

// MapRule is a helper for creating a `RuleEntry`.
func MapRule(pattern string, factory RuleFactoryFunc) RuleEntry {
	return RuleEntry{
		pattern: regexp.MustCompile(pattern),
		factory: factory,
	}
}

// RuleMatch contains parameters extracted from a URL when
// a rule successfully matches.
type RuleMatch struct {
	// OrigURI is the original URI parameter that was passed in
	OrigURI string

	// UserInfo is the username used to identify a gallery in the URI
	UserInfo string
}
