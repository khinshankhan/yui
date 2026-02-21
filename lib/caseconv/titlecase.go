package caseconv

import (
	"strings"
	"unicode"
)

// TitleStyle represents different title case style guides.
type TitleStyle string

const (
	StyleAPA       TitleStyle = "apa"
	StyleChicago   TitleStyle = "chicago"
	StyleMLA       TitleStyle = "mla"
	StyleAP        TitleStyle = "ap"
	StyleBluebook  TitleStyle = "bluebook"
	StyleAMA       TitleStyle = "ama"
	StyleNYTimes   TitleStyle = "nytimes"
	StyleWikipedia TitleStyle = "wikipedia"
)

// AvailableTitleStyles returns all available title case styles.
func AvailableTitleStyles() []TitleStyle {
	return []TitleStyle{
		StyleAPA, StyleChicago, StyleMLA, StyleAP,
		StyleBluebook, StyleAMA, StyleNYTimes, StyleWikipedia,
	}
}

var (
	articles = map[string]bool{
		"a":   true,
		"an":  true,
		"the": true,
	}

	coordinatingConjunctions = map[string]bool{
		"for": true,
		"and": true,
		"nor": true,
		"but": true,
		"or":  true,
		"yet": true,
		"so":  true,
	}

	shortPrepositions = map[string]bool{
		"as":   true,
		"at":   true,
		"but":  true,
		"by":   true,
		"down": true,
		"for":  true,
		"from": true,
		"in":   true,
		"into": true,
		"like": true,
		"near": true,
		"of":   true,
		"off":  true,
		"on":   true,
		"onto": true,
		"out":  true,
		"over": true,
		"past": true,
		"per":  true,
		"plus": true,
		"than": true,
		"to":   true,
		"up":   true,
		"upon": true,
		"via":  true,
		"with": true,
	}

	allPrepositions = map[string]bool{
		"aboard":      true,
		"about":       true,
		"above":       true,
		"across":      true,
		"after":       true,
		"against":     true,
		"along":       true,
		"amid":        true,
		"among":       true,
		"around":      true,
		"as":          true,
		"at":          true,
		"before":      true,
		"behind":      true,
		"below":       true,
		"beneath":     true,
		"beside":      true,
		"besides":     true,
		"between":     true,
		"beyond":      true,
		"but":         true,
		"by":          true,
		"concerning":  true,
		"considering": true,
		"despite":     true,
		"down":        true,
		"during":      true,
		"except":      true,
		"for":         true,
		"from":        true,
		"in":          true,
		"inside":      true,
		"into":        true,
		"like":        true,
		"near":        true,
		"of":          true,
		"off":         true,
		"on":          true,
		"onto":        true,
		"out":         true,
		"outside":     true,
		"over":        true,
		"past":        true,
		"per":         true,
		"plus":        true,
		"regarding":   true,
		"round":       true,
		"since":       true,
		"than":        true,
		"through":     true,
		"throughout":  true,
		"till":        true,
		"to":          true,
		"toward":      true,
		"towards":     true,
		"under":       true,
		"underneath":  true,
		"unlike":      true,
		"until":       true,
		"up":          true,
		"upon":        true,
		"via":         true,
		"with":        true,
		"within":      true,
		"without":     true,
	}

	numberWords = map[string]bool{
		"one":   true,
		"two":   true,
		"three": true,
		"four":  true,
		"five":  true,
		"six":   true,
		"seven": true,
		"eight": true,
		"nine":  true,
		"ten":   true,
	}
)

// ToTitleStyle converts a string to title case using the specified style guide.
func ToTitleStyle(s string, style TitleStyle) string {
	if s == "" {
		return s
	}

	tokens := tokenize(s)
	if len(tokens) == 0 {
		return s
	}

	firstWordIdx := -1
	lastWordIdx := -1
	for i, t := range tokens {
		if t.isWord {
			if firstWordIdx == -1 {
				firstWordIdx = i
			}
			lastWordIdx = i
		}
	}

	amaProperNoun := map[int]bool{}
	if style == StyleAMA {
		amaProperNoun = detectAMAProperNoun(tokens)
	}

	capitalizeNextWord := false
	for i := range tokens {
		if !tokens[i].isWord {
			if strings.ContainsAny(tokens[i].text, ":—–") {
				capitalizeNextWord = true
			}
			continue
		}

		isFirst := i == firstWordIdx || capitalizeNextWord || amaProperNoun[i]
		isLast := i == lastWordIdx

		tokens[i].text = capitalizeWord(tokens[i].text, style, isFirst, isLast)

		_, _, trailing := extractPunctuation(tokens[i].text)
		capitalizeNextWord = strings.ContainsAny(trailing, ":—–")
	}

	var result strings.Builder
	for _, t := range tokens {
		result.WriteString(t.text)
	}
	return result.String()
}

func detectAMAProperNoun(tokens []token) map[int]bool {
	result := map[int]bool{}
	for i := 0; i < len(tokens); i++ {
		if !tokens[i].isWord {
			continue
		}

		_, core, _ := extractPunctuation(tokens[i].text)
		if strings.ToLower(core) != "of" {
			continue
		}

		var span []int
		for j := i + 1; j < len(tokens); j++ {
			if !tokens[j].isWord {
				if strings.ContainsAny(tokens[j].text, ":—–") {
					break
				}
				continue
			}
			span = append(span, j)
		}

		if len(span) < 2 {
			continue
		}

		_, firstCore, _ := extractPunctuation(tokens[span[0]].text)
		firstLower := strings.ToLower(firstCore)
		if firstLower == "" || articles[firstLower] || coordinatingConjunctions[firstLower] || allPrepositions[firstLower] || numberWords[firstLower] {
			continue
		}
		if len([]rune(firstLower)) < 4 {
			continue
		}

		for _, idx := range span {
			_, c, _ := extractPunctuation(tokens[idx].text)
			lc := strings.ToLower(c)
			if lc == "" || articles[lc] || coordinatingConjunctions[lc] || allPrepositions[lc] {
				continue
			}
			result[idx] = true
		}
	}
	return result
}

// token represents a piece of text (word or delimiter)
type token struct {
	text   string
	isWord bool
}

// tokenize splits text into tokens, preserving spaces and handling hyphens within words
func tokenize(s string) []token {
	var tokens []token
	var current strings.Builder
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if unicode.IsSpace(r) {
			if current.Len() > 0 {
				tokens = append(tokens, token{text: current.String(), isWord: true})
				current.Reset()
			}
			tokens = append(tokens, token{text: string(r), isWord: false})
		} else if isWordSeparator(r) {
			if current.Len() > 0 {
				tokens = append(tokens, token{text: current.String(), isWord: true})
				current.Reset()
			}
			tokens = append(tokens, token{text: string(r), isWord: false})
		} else {
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, token{text: current.String(), isWord: true})
	}

	return tokens
}

func isWordSeparator(r rune) bool {
	switch r {
	case ':', '/', '&', '—', '–':
		return true
	default:
		return false
	}
}

// capitalizeWord applies the appropriate capitalization to a word based on style
func capitalizeWord(word string, style TitleStyle, isFirst, isLast bool) string {
	if strings.Contains(word, "-") {
		return capitalizeHyphenated(word, style, isFirst, isLast)
	}

	leading, core, trailing := extractPunctuation(word)
	if core == "" {
		return word
	}

	lowerCore := strings.ToLower(core)

	if isFirst {
		return leading + capitalize(core) + trailing
	}

	// Keep common contraction form "'n'" lowercase in titles.
	if lowerCore == "n" && strings.Contains(leading, "'") && strings.Contains(trailing, "'") {
		return leading + "n" + trailing
	}

	// AMA is sentence case, but preserve already-capitalized words
	// (eg proper nouns provided in input).
	if style == StyleAMA && hasUppercase(core) {
		return leading + core + trailing
	}

	if shouldBeLowercase(lowerCore, style, isLast) {
		return leading + lowerCore + trailing
	}

	return leading + capitalize(core) + trailing
}

func hasUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// extractPunctuation separates leading/trailing punctuation from the core word
func extractPunctuation(word string) (leading, core, trailing string) {
	runes := []rune(word)
	start := 0
	end := len(runes)

	for start < end && !unicode.IsLetter(runes[start]) && !unicode.IsDigit(runes[start]) {
		start++
	}

	for end > start && !unicode.IsLetter(runes[end-1]) && !unicode.IsDigit(runes[end-1]) {
		end--
	}

	leading = string(runes[:start])
	core = string(runes[start:end])
	trailing = string(runes[end:])
	return
}

// shouldBeLowercase determines if a word should remain lowercase based on style rules
func shouldBeLowercase(word string, style TitleStyle, isLast bool) bool {
	wordLen := len([]rune(word))

	switch style {
	case StyleAPA:
		if isLast {
			return false
		}
		if wordLen >= 4 {
			return false
		}
		return articles[word] || shortPrepositions[word] || coordinatingConjunctions[word]

	case StyleChicago:
		if isLast {
			return false
		}
		return articles[word] || coordinatingConjunctions[word] || allPrepositions[word]

	case StyleMLA:
		if isLast {
			return false
		}
		return articles[word] || coordinatingConjunctions[word] || allPrepositions[word]

	case StyleAP:
		if isLast {
			return false
		}
		return articles[word] || coordinatingConjunctions[word] || allPrepositions[word]

	case StyleBluebook:
		if isLast {
			return false
		}
		if articles[word] || coordinatingConjunctions[word] {
			return true
		}
		if allPrepositions[word] && wordLen <= 4 {
			return true
		}
		return false

	case StyleAMA:
		return true

	case StyleNYTimes:
		if isLast {
			return false
		}
		return articles[word] || coordinatingConjunctions[word] || allPrepositions[word]

	case StyleWikipedia:
		if isLast {
			return false
		}
		if articles[word] || coordinatingConjunctions[word] {
			return true
		}
		if allPrepositions[word] {
			return wordLen < 5
		}
		return false

	default:
		return false
	}
}

// capitalizeHyphenated handles hyphenated words according to style rules
func capitalizeHyphenated(word string, style TitleStyle, isFirst, isLast bool) string {
	parts := strings.Split(word, "-")
	if len(parts) == 0 {
		return word
	}

	leading, firstCore, _ := extractPunctuation(parts[0])
	if firstCore != "" {
		parts[0] = firstCore
	}

	_, lastCore, trailing := extractPunctuation(parts[len(parts)-1])
	if lastCore != "" {
		parts[len(parts)-1] = lastCore
	}

	for i := range parts {
		part := parts[i]
		if part == "" {
			continue
		}

		lowerPart := strings.ToLower(part)
		isFirstPart := i == 0

		switch style {
		case StyleAP:
			if isFirstPart {
				parts[i] = capitalize(part)
			} else if articles[lowerPart] || coordinatingConjunctions[lowerPart] || allPrepositions[lowerPart] {
				parts[i] = lowerPart
			} else {
				parts[i] = capitalize(part)
			}

		case StyleAPA:
			parts[i] = capitalize(part)

		case StyleMLA:
			if isFirstPart {
				parts[i] = capitalize(part)
			} else {
				if !articles[lowerPart] && !allPrepositions[lowerPart] && !coordinatingConjunctions[lowerPart] {
					parts[i] = capitalize(part)
				} else {
					parts[i] = lowerPart
				}
			}

		case StyleChicago:
			if isFirstPart {
				parts[i] = capitalize(part)
			} else if articles[lowerPart] || coordinatingConjunctions[lowerPart] || allPrepositions[lowerPart] {
				parts[i] = lowerPart
			} else {
				parts[i] = capitalize(part)
			}

		case StyleAMA:
			if isFirst && isFirstPart {
				parts[i] = capitalize(part)
			} else if hasUppercase(part) {
				parts[i] = part
			} else {
				parts[i] = lowerPart
			}

		default:
			if isFirstPart {
				parts[i] = capitalize(part)
			} else if shouldBeLowercase(lowerPart, style, false) {
				parts[i] = lowerPart
			} else {
				parts[i] = capitalize(part)
			}
		}
	}

	result := leading + strings.Join(parts, "-") + trailing
	return result
}
