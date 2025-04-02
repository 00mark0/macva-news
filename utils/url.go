package utils

import (
	"regexp"
	"strings"
)

// Slugify converts a string into a URL-friendly format
func Slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces and special characters with hyphens
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")

	// Trim hyphens from the start and end
	s = strings.Trim(s, "-")

	return s
}
