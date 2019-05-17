package common

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
		result := make([]RuleMatch, 0)

		for _, url := range seedURLs {
			if rule.pattern.MatchString(url) {
				// Map capture groups for easier use
				groups := rule.pattern.SubexpNames()
				captures := rule.pattern.FindStringSubmatch(url)
				matches := make(map[string]string)
				for idx, group := range groups {
					if group == UserInfo {
						matches[group] = captures[idx]
					}
				}

				ruleMatch := RuleMatch{
					OrigURI: url,
				}

				if userInfo, ok := matches[UserInfo]; ok {
					ruleMatch.UserInfo = userInfo
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

type RuleFactoryFunc func(matches []RuleMatch, config *Config) Scraper

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
