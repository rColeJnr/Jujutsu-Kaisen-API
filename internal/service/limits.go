package service

import (
	"fmt"
	"net/url"
	"unicode/utf8"
)

// Maximum sizes for text fields, without them a client could store
// megabytes in a single "name" and every response would carry it.
const (
	maxNameLen  = 120 // names, terms, titles
	maxShortLen = 50
	maxTextLen  = 2000
	maxURLLen   = 2048
	maxListLen  = 500
)

type textRule struct {
	field string
	value string
	max   int
	isURL bool
}

func text(field, value string, max int) textRule {
	return textRule{field: field, value: value, max: max}
}

func imageURL(field, value string) textRule {
	return textRule{field: field, value: value, max: maxURLLen, isURL: true}
}

// str derefenres a string
func str(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// validateTexts checks for every rule, return error
// Length is counted in characters, not bytes.
func validateTexts(rules ...textRule) error {
	for _, r := range rules {
		if n := utf8.RuneCountInString(r.value); n > r.max {
			return fmt.Errorf("validation: %s must be at most %d characters, got: %d", r.field, r.max, n)
		}
		if r.isURL && r.value != "" {
			u, err := url.Parse(r.value)
			if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
				return fmt.Errorf("validation: %s must be a valid http or https URL", r.field)
			}
		}
	}
	return nil
}
