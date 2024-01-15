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
		s := slug.Make(input)
		return strings.ToLower(s)
	} else if namingStyle == "name" {
		return input
	}

	log.Fatal().Msg("Invalid naming-scheme, must be one of: lowercase, slug, name")
	return input
}
