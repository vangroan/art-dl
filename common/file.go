package common

import (
	"fmt"
	"html"
	"io/ioutil"
	"regexp"
	"strings"
)

const commentSymbol string = "#"

// LoadGalleryFile reads a list of gallery
// URLs from a text file.
func LoadGalleryFile(filepath string) ([]string, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to gallery file: %s", err)
	}

	return parse(string(b)), nil
}

func parse(data string) []string {
	result := make([]string, 0)

	lines := strings.Split(strings.Replace(data, "\r\n", "\n", -1), "\n")
	for _, line := range lines {
		line = strings.Trim(line, " ")

		// Ignore comments
		if strings.HasPrefix(line, commentSymbol) {
			continue
		}

		// Ignore empty lines
		if line == "" {
			continue
		}

		result = append(result, line)
	}

	return result
}

// SanitizeFilename removes characters from the given filename
// which are reserved by file systems.
func SanitizeFilename(filename string) string {
	// List of reserved characters from Wikipedia
	// See: https://en.wikipedia.org/wiki/Filename#In_Windows
	reserved := regexp.MustCompile(`[\\/\?\*\:\|\<\>\,\;\=]+`)
	sanitized := reserved.ReplaceAllString(filename, "-")
	// Double quotes are reserved, but single quotes are permitted.
	return strings.ReplaceAll(sanitized, "\"", "'")
}

// SanitizeDirname removes characters from the given directory
// name which are reserved by file systems.
func SanitizeDirname(dirname string) string {
	// We expect names from web sources.
	unescaped := html.UnescapeString(dirname)
	// List of reserved characters from Wikipedia
	// See: https://en.wikipedia.org/wiki/Filename#In_Windows
	reserved := regexp.MustCompile(`[\\/\?\*\:\|\<\>\,\;\=]+`)
	sanitized := reserved.ReplaceAllString(unescaped, "-")
	// Trailing periods are not allowed in Windows,
	// but Unix allows a leading period to indicate
	// hidden folder.
	sanitized = strings.TrimRight(sanitized, ".")
	// Double quotes are reserved, but single quotes are permitted.
	return strings.ReplaceAll(sanitized, "\"", "'")
}
