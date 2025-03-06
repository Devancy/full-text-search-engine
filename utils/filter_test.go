package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLowercaseFilter(t *testing.T) {
	var (
		in  = []string{"Cat", "DOG", "fish"}
		out = []string{"cat", "dog", "fish"}
	)
	assert.Equal(t, out, lowercaseFilter(in))
}

func TestStopwordFilter(t *testing.T) {
	var (
		in  = []string{"i", "am", "the", "cat"}
		out = []string{"am", "cat"}
	)
	assert.Equal(t, out, stopwordFilter(in))
}

func TestStemmerFilter(t *testing.T) {
	var (
		in  = []string{"cat", "cats", "fish", "fishing", "fished", "airline"}
		out = []string{"cat", "cat", "fish", "fish", "fish", "airlin"}
	)
	assert.Equal(t, out, stemmerFilter(in))
}

func TestCharacterFilter(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Remove punctuation from ends",
			input:    []string{"!hello!", ".world.", "?test?"},
			expected: []string{"hello", "world", "test"},
		},
		{
			name:     "Skip short tokens",
			input:    []string{"a", "ab", "abc"},
			expected: []string{"ab", "abc"},
		},
		{
			name:     "Preserve alphanumeric content",
			input:    []string{"hello123", "test42world", "123test"},
			expected: []string{"hello123", "test42world", "123test"},
		},
		{
			name:     "Mixed special characters",
			input:    []string{"@user#name", "$price100", "email@domain"},
			expected: []string{"user#name", "price100", "email@domain"},
		},
		{
			name:     "Empty and invalid tokens",
			input:    []string{"", "!", "@", "a", "#b#"},
			expected: []string{},
		},
		{
			name:     "Numbers and mixed content",
			input:    []string{"2023", "version2.0", "!2024!"},
			expected: []string{"2023", "version2.0", "2024"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := characterFilter(tt.input)
			assert.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}
