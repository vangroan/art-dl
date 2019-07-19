package common

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const commentSymbol string = "#"

// LoadGalleryFile reads a list of gallery
// URLs from a text file.
func LoadGalleryFile(filepath string) ([]string, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("Failed to gallery file: %s", err)
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
