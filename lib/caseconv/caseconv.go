package caseconv

import (
	"strings"
	"unicode"
)

func splitWords(str string) []string {
	cleaned := strings.ToLower(strings.TrimSpace(str))
	cleaned = strings.ReplaceAll(cleaned, "-", " ")
	cleaned = strings.ReplaceAll(cleaned, "_", " ")
	return strings.Fields(cleaned)
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func Convert(input, mode string) string {
	switch mode {
	case "upper":
		return strings.ToUpper(input)
	case "lower":
		return strings.ToLower(input)
	case "kebab":
		return strings.ToLower(strings.Join(splitWords(input), "-"))
	case "snake":
		return strings.ToLower(strings.Join(splitWords(input), "_"))
	case "pascal":
		parts := splitWords(input)
		for i, w := range parts {
			parts[i] = capitalize(w)
		}
		return strings.Join(parts, "")
	case "camel":
		parts := splitWords(input)
		if len(parts) == 0 {
			return ""
		}
		for i := range parts {
			if i == 0 {
				continue
			}
			parts[i] = capitalize(parts[i])
		}
		return parts[0] + strings.Join(parts[1:], "")
	default:
		return input
	}
}
