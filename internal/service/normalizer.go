package service

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Removes case, diacritics, and non-alphanumeric characters.
func NormalizeText(text string) (string, error) {
	// Lowercase input to avoid case-sensitive inconsistencies.
	text = strings.ToLower(text)

	// Compose a transformer that:
	// - Decomposes Unicode characters (NFD),
	// - Removes non-spacing marks (diacritics),
	// - Recomposes characters to NFC form.
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)

	// Apply the chained transformer to the input string.
	// Returns transformed text, number of bytes read, and error.
	text, _, err := transform.String(t, text)
	if err != nil {
		return "", err
	}

	// Replace non-letter, non-digit, non-space characters with a single space.
	reSpecialChars, err := regexp.Compile(`[^\p{L}\p{N}\s]+`)
	if err != nil {
		return "", err
	}
	text = reSpecialChars.ReplaceAllString(text, " ")

	// Collapse multiple spaces into a single space.
	reSpaces, err := regexp.Compile(`\s+`)
	if err != nil {
		return "", err
	}
	text = reSpaces.ReplaceAllString(text, " ")

	// Trim leading/trailing spaces.
	return strings.TrimSpace(text), nil
}

func CleanStringJellyfin(text string) (string, error) {
	// Jellyfin points at this chars <, >, :, ", /, \, |, ?, *
	// as "known to cause issues"
	reForbiddenChars, err := regexp.Compile(`[<>:"/\\|?*]`)
	if err != nil {
		return "", err
	}
	cleanedText := reForbiddenChars.ReplaceAllString(text, "")
	return cleanedText, nil
}
