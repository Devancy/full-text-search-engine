package utils

import (
	"math"
	"sort"
)

// IndexEntry stores document IDs and their frequencies
type IndexEntry struct {
	DocIDs []int
	Freqs  []float64
}

// Index is an inverted index. It maps tokens to document IDs and their frequencies.
type Index struct {
	entries  map[string]*IndexEntry
	docCount int
}

// NewIndex creates a new Index instance
func NewIndex() *Index {
	return &Index{
		entries: make(map[string]*IndexEntry),
	}
}

func (idx *Index) Clear() {
	idx.entries = make(map[string]*IndexEntry)
	idx.docCount = 0
}

func (idx *Index) Stats() IndexStats {
	return IndexStats{
		DocumentCount: idx.docCount,
		TermCount:     len(idx.entries),
	}
}

// Add adds documents to the Index with TF-IDF scoring
func (idx *Index) Add(docs []*Document) {
	if len(docs) == 0 {
		return
	}

	// Update document count for IDF calculation
	idx.docCount += len(docs)

	for _, doc := range docs {

		// Count token frequencies in document
		tokenFreq := make(map[string]int)
		tokens := analyze(doc.Text)
		totalTokens := len(tokens)
		if totalTokens == 0 {
			continue
		}

		// Calculate term frequencies
		for _, token := range tokens {
			tokenFreq[token]++
		}

		// Update index with document frequencies
		for token, freq := range tokenFreq {
			if idx.entries[token] == nil {
				idx.entries[token] = &IndexEntry{
					DocIDs: make([]int, 0, 64),
					Freqs:  make([]float64, 0, 64),
				}
			}
			entry := idx.entries[token]

			entry.DocIDs = append(entry.DocIDs, doc.ID)
			// Calculate TF as frequency / total tokens in document
			tf := float64(freq) / float64(totalTokens)
			entry.Freqs = append(entry.Freqs, tf)
		}
	}

	idx.calculateIDF()
}

// calculateIDF updates term frequencies with IDF (Inverse Document Frequency) scores
func (idx *Index) calculateIDF() {
	for _, entry := range idx.entries {
		// IDF = log(N/(df + 1)) + 1  // Adding 1 to avoid division by zero and negative values
		idf := math.Log(float64(idx.docCount)/(float64(len(entry.DocIDs))+1.0)) + 1.0

		// Update frequencies with TF-IDF score
		for i := range entry.Freqs {
			entry.Freqs[i] *= idf
		}
	}
}

// SearchResult represents a scored search result
type SearchResult struct {
	DocID int
	Score float64
}

// Search queries the Index for the given text and returns scored results
func (idx *Index) Search(text string) []SearchResult {
	tokens := analyze(text)
	if len(tokens) == 0 {
		return nil
	}

	// Calculate scores for each matching document
	scores := make(map[int]float64)
	for _, token := range tokens {
		if entry, ok := idx.entries[token]; ok {
			for i, docID := range entry.DocIDs {
				scores[docID] += entry.Freqs[i]
			}
		}
	}

	if len(scores) == 0 {
		return nil
	}

	results := make([]SearchResult, 0, len(scores))
	for docID, score := range scores {
		results = append(results, SearchResult{
			DocID: docID,
			Score: score,
		})
	}

	// Sort results by score (highest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}
