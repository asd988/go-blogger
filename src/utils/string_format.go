package utils

import (
	"regexp"
	"strings"
)

func Slugify(title string) string {
	slug := strings.ToLower(title)

	// Remove non-alphanumeric and non-hyphen characters
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}
