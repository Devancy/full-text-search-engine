package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	idx := NewIndex()

	// Test empty index
	assert.Empty(t, idx.Search("foo"))
	assert.Empty(t, idx.Search("donut"))

	// Test single document indexing
	doc1 := Document{ID: 1, Text: "A donut on a glass plate. Only the donuts."}
	idx.Add([]*Document{&doc1})

	// Stopwords should return no results
	assert.Empty(t, idx.Search("a"))
	assert.Empty(t, idx.Search("the"))
	assert.Empty(t, idx.Search("on"))

	// Test basic search functionality with stemming
	// "donut" and "donuts" should stem to the same term
	results := idx.Search("donut")
	assert.Len(t, results, 1)
	assert.Equal(t, 1, results[0].DocID)
	assert.Greater(t, results[0].Score, 0.0)

	// Test case insensitivity and stemming
	results = idx.Search("DONUTS")
	assert.Len(t, results, 1)
	assert.Equal(t, 1, results[0].DocID)

	results = idx.Search("glass")
	assert.Len(t, results, 1)
	assert.Equal(t, 1, results[0].DocID)

	// Test multiple document indexing
	doc2 := Document{ID: 2, Text: "donut is a donut"}
	idx.Add([]*Document{&doc2})

	// Test document frequency affects scoring
	results = idx.Search("donut")
	assert.Len(t, results, 2)
	assert.Contains(t, []int{1, 2}, results[0].DocID)
	assert.Contains(t, []int{1, 2}, results[1].DocID)

	// Document 2 should have higher score due to higher term frequency
	assert.Equal(t, 2, results[0].DocID, "Doc 2 should rank higher due to higher term frequency")
	assert.Equal(t, 1, results[1].DocID)
	assert.Greater(t, results[0].Score, results[1].Score)

	// Test term appearing in only one document
	results = idx.Search("glass")
	assert.Len(t, results, 1)
	assert.Equal(t, 1, results[0].DocID)
}

// TestIndexScoring tests the TF-IDF scoring implementation
func TestIndexScoring(t *testing.T) {
	idx := NewIndex()

	// Add documents with different term frequencies
	docs := []*Document{
		{ID: 1, Text: "apple banana apple"},  // apple appears twice
		{ID: 2, Text: "apple banana cherry"}, // all terms appear once
	}
	idx.Add(docs)

	// Test that higher term frequency results in higher score
	results := idx.Search("apple")
	assert.Len(t, results, 2)
	assert.Equal(t, 1, results[0].DocID, "Doc 1 should rank higher due to higher term frequency")
	assert.Equal(t, 2, results[1].DocID)
	assert.Greater(t, results[0].Score, results[1].Score)

	// Test that terms appearing in fewer documents have higher IDF
	results = idx.Search("cherry")
	assert.Len(t, results, 1)
	assert.Equal(t, 2, results[0].DocID, "Cherry appears in Doc 2 only")
	cherryScore := results[0].Score

	results = idx.Search("banana")
	assert.Len(t, results, 2)
	bananaScore := results[0].Score

	// Cherry appears in 1 doc, banana in 2 docs, so cherry should have higher score
	assert.Greater(t, cherryScore, bananaScore, "Terms in fewer documents should have higher scores")
}

// TestEmptyAndEdgeCases tests edge cases
func TestEmptyAndEdgeCases(t *testing.T) {
	idx := NewIndex()

	// Test empty index
	assert.Empty(t, idx.Search(""))
	assert.Empty(t, idx.Search(" "))

	// Test single character and short words (should be filtered out)
	idx.Add([]*Document{{ID: 1, Text: "a b c test"}})
	assert.Empty(t, idx.Search("a"))
	assert.Empty(t, idx.Search("b"))
	assert.Empty(t, idx.Search("c"))

	results := idx.Search("test")
	assert.Len(t, results, 1)

	// Test documents with only stopwords
	idx.Add([]*Document{{ID: 2, Text: "the a in on"}})
	assert.Empty(t, idx.Search("the"))
	assert.Empty(t, idx.Search("in"))
}

func generateLargeDataset(n int) []*Document {
	docs := make([]*Document, n)
	texts := []string{
		"The quick brown fox jumps over the lazy dog",
		"Pack my box with five dozen liquor jugs",
		"How vexingly quick daft zebras jump",
		"The five boxing wizards jump quickly",
		"Sphinx of black quartz, judge my vow",
	}

	for i := range n {
		docs[i] = &Document{
			ID:    i,
			Title: fmt.Sprintf("Document %d", i),
			Text:  texts[i%len(texts)],
		}
	}
	return docs
}

func BenchmarkIndexAdd(b *testing.B) {
	docs := generateLargeDataset(1000)

	b.Run("SimpleIndex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			idx := NewIndex()
			idx.Add(docs)
		}
	})

	b.Run("ConcurrentIndex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			idx := NewConcurrentIndex()
			idx.Add(docs)
		}
	})
}

func BenchmarkIndexAddLarge(b *testing.B) {
	docs := generateLargeDataset(1000000)

	b.Run("SimpleIndex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			idx := NewIndex()
			idx.Add(docs)
		}
	})

	b.Run("ConcurrentIndex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			idx := NewConcurrentIndex()
			idx.Add(docs)
		}
	})
}
