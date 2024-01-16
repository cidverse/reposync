package util

import (
	"strings"

	"github.com/gosimple/slug"
	"github.com/rs/zerolog/log"
)

// Slugify transforms a string for generate better directory names
func Slugify(input string, namingStyle string) string {
	if namingStyle == "lowercase" {
		return strings.ToLower(input)
	} else if namingStyle == "slug" || namingStyle == "" {
		parts := strings.Split(input, "/")
		for i, part := range parts {
			parts[i] = slug.Make(part)
		}
		return strings.Join(parts, "/")
	} else if namingStyle == "name" {
		return input
	}

	log.Fatal().Msg("Invalid naming-scheme, must be one of: lowercase, slug, name")
	return input
}
