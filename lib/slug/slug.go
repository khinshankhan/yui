package slug

import (
	"strconv"
	"strings"
	"unicode"
)

// Make converts text into a URL-friendly slug.
func Make(input string) string {
	input = insertWordBoundaries(input)

	var b strings.Builder
	lastDash := false

	for _, r := range input {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(unicode.ToLower(r))
			lastDash = false
		case r == '-' || r == '_' || unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r):
			if b.Len() == 0 || lastDash {
				continue
			}
			b.WriteByte('-')
			lastDash = true
		}
	}

	return strings.Trim(b.String(), "-")
}

func insertWordBoundaries(input string) string {
	var b strings.Builder
	runes := []rune(strings.TrimSpace(input))

	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			if unicode.IsLower(prev) || unicode.IsDigit(prev) {
				b.WriteByte(' ')
			} else if unicode.IsUpper(prev) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				b.WriteByte(' ')
			}
		}

		b.WriteRune(r)
	}

	return b.String()
}

// NextAvailable returns base unless it is reserved, then appends a numeric suffix.
func NextAvailable(base string, reserved []string) string {
	if base == "" {
		return ""
	}

	used := make(map[string]struct{}, len(reserved))
	for _, item := range reserved {
		used[strings.TrimSpace(item)] = struct{}{}
	}

	if _, ok := used[base]; !ok {
		return base
	}

	for i := 1; ; i++ {
		candidate := base + "-" + strconv.Itoa(i)
		if _, ok := used[candidate]; !ok {
			return candidate
		}
	}
}
