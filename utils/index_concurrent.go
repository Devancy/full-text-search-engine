package utils

import (
	"math"
	"runtime"
	"sort"
	"sync"
)

// ConcurrentIndex is an inverted index with concurrent processing capabilities.
// It maps tokens to document IDs and their frequencies.
type ConcurrentIndex struct {
	sync.RWMutex
	entries  sync.Map // map[string]*ConcurrentIndexEntry
	docCount int
}

// NewConcurrentIndex creates a new ConcurrentIndex instance
func NewConcurrentIndex() *ConcurrentIndex {
	return &ConcurrentIndex{}
}

func (idx *ConcurrentIndex) Clear() {
	idx.entries.Range(func(key, value any) bool {
		idx.entries.Delete(key)
		return true
	})
	idx.docCount = 0
}

func (idx *ConcurrentIndex) Stats() IndexStats {
	termCount := 0
	idx.entries.Range(func(key, value any) bool {
		termCount++
		return true
	})
	return IndexStats{
		DocumentCount: idx.docCount,
		TermCount:     termCount,
	}
}

// Add adds documents to the ConcurrentIndex with TF-IDF scoring using parallel processing
func (idx *ConcurrentIndex) Add(docs []*Document) {
	if len(docs) == 0 {
		return
	}

	// Update document count for IDF calculation
	idx.Lock()
	idx.docCount += len(docs)
	idx.Unlock()

	// Process documents in parallel
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	docChan := make(chan *Document, numWorkers*2)

	// Start worker goroutines
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for doc := range docChan {

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
					// Ensure ConcurrentIndexEntry is created with float32 slice
					entry, _ := idx.entries.LoadOrStore(token, &ConcurrentIndexEntry{
						DocIDs: make([]int, 0, 64),
						Freqs:  make([]float32, 0, 64), // Use float32
					})
					indexEntry := entry.(*ConcurrentIndexEntry)

					// Lock only this entry while updating it
					indexEntry.Lock()
					indexEntry.DocIDs = append(indexEntry.DocIDs, doc.ID)
					// Calculate TF as frequency / total tokens in document
					// Cast result to float32 before appending
					tf := float32(float64(freq) / float64(totalTokens))
					indexEntry.Freqs = append(indexEntry.Freqs, tf) // Append float32
					indexEntry.Unlock()
				}
			}
		}()
	}

	for _, doc := range docs {
		docChan <- doc
	}
	close(docChan)
	wg.Wait()

	// TF is stored directly, IDF calculated during Search
	// idx.calculateIDF()
}

// Search queries the ConcurrentIndex for the given text and returns scored results
func (idx *ConcurrentIndex) Search(text string) []SearchResult {
	tokens := analyze(text)
	if len(tokens) == 0 {
		return nil
	}

	scores := make(map[int]float32)
	var scoresMutex sync.RWMutex

	// Calculate scores for each token
	for _, token := range tokens {
		if entry, ok := idx.entries.Load(token); ok {
			indexEntry := entry.(*ConcurrentIndexEntry)
			indexEntry.RLock()

			// Calculate IDF for the current term
			// Must read docCount within the lock to ensure consistency if Add is running concurrently
			// Use RLock on the main index to safely read docCount
			idx.RLock()
			docCount := idx.docCount
			idx.RUnlock()

			// IDF = log(N/(df + 1)) + 1
			idf := float32(math.Log(float64(docCount)/(float64(len(indexEntry.DocIDs))+1.0)) + 1.0)

			for i, docID := range indexEntry.DocIDs {
				scoresMutex.Lock()
				// Score is TF (from entry.Freqs) * IDF (calculated now)
				scores[docID] += indexEntry.Freqs[i] * idf
				scoresMutex.Unlock()
			}
			indexEntry.RUnlock()
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
