package utils

import (
	"strings"
	"unicode"

	snowballeng "github.com/kljensen/snowball/english"
)

// lowercaseFilter returns a slice of tokens normalized to lower case.
func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}
	return r
}

// characterFilter removes unwanted characters from tokens
func characterFilter(tokens []string) []string {
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		// Remove non-alphanumeric characters from start and end
		token = strings.TrimFunc(token, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})

		// Skip empty tokens or those that are too short
		if len(token) < 2 {
			continue
		}

		r = append(r, token)
	}
	return r
}

// stopwordFilter returns a slice of tokens with stop words removed.
func stopwordFilter(tokens []string) []string {
	var stopwords = map[string]struct{}{
		"a": {}, "about": {}, "above": {}, "after": {}, "again": {}, "against": {}, "all": {},
		"am": {}, "an": {}, "and": {}, "any": {}, "are": {}, "aren't": {}, "as": {}, "at": {},
		"be": {}, "because": {}, "been": {}, "before": {}, "being": {}, "below": {}, "between": {},
		"both": {}, "but": {}, "by": {}, "can": {}, "can't": {}, "cannot": {}, "could": {},
		"couldn't": {}, "did": {}, "didn't": {}, "do": {}, "does": {}, "doesn't": {}, "doing": {},
		"don't": {}, "down": {}, "during": {}, "each": {}, "few": {}, "for": {}, "from": {},
		"further": {}, "had": {}, "hadn't": {}, "has": {}, "hasn't": {}, "have": {}, "haven't": {},
		"having": {}, "he": {}, "he'd": {}, "he'll": {}, "he's": {}, "her": {}, "here": {},
		"here's": {}, "hers": {}, "herself": {}, "him": {}, "himself": {}, "his": {}, "how": {},
		"how's": {}, "i": {}, "i'd": {}, "i'll": {}, "i'm": {}, "i've": {}, "if": {}, "in": {},
		"into": {}, "is": {}, "isn't": {}, "it": {}, "it's": {}, "its": {}, "itself": {},
		"let's": {}, "me": {}, "more": {}, "most": {}, "mustn't": {}, "my": {}, "myself": {},
		"no": {}, "nor": {}, "not": {}, "of": {}, "off": {}, "on": {}, "once": {}, "only": {},
		"or": {}, "other": {}, "ought": {}, "our": {}, "ours": {}, "ourselves": {}, "out": {},
		"over": {}, "own": {}, "same": {}, "shan't": {}, "she": {}, "she'd": {}, "she'll": {},
		"she's": {}, "should": {}, "shouldn't": {}, "so": {}, "some": {}, "such": {}, "than": {},
		"that": {}, "that's": {}, "the": {}, "their": {}, "theirs": {}, "them": {}, "themselves": {},
		"then": {}, "there": {}, "there's": {}, "these": {}, "they": {}, "they'd": {}, "they'll": {},
		"they're": {}, "they've": {}, "this": {}, "those": {}, "through": {}, "to": {}, "too": {},
		"under": {}, "until": {}, "up": {}, "very": {}, "was": {}, "wasn't": {}, "we": {},
		"we'd": {}, "we'll": {}, "we're": {}, "we've": {}, "were": {}, "weren't": {}, "what": {},
		"what's": {}, "when": {}, "when's": {}, "where": {}, "where's": {}, "which": {},
		"while": {}, "who": {}, "who's": {}, "whom": {}, "why": {}, "why's": {}, "with": {},
		"won't": {}, "would": {}, "wouldn't": {}, "you": {}, "you'd": {}, "you'll": {},
		"you're": {}, "you've": {}, "your": {}, "yours": {}, "yourself": {}, "yourselves": {},
	}

	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := stopwords[token]; !ok {
			r = append(r, token)
		}
	}
	return r
}

// stemmerFilter returns a slice of stemmed tokens.
// Stemming is the process of reducing a word to its base or root form, which helps normalize words for text analysis.
// For example, "running," "runner," and "runs" might all be reduced to the root form "run".
func stemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = snowballeng.Stem(token, false)
	}
	return r
}
