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
					entry, _ := idx.entries.LoadOrStore(token, &ConcurrentIndexEntry{
						DocIDs: make([]int, 0, 64),
						Freqs:  make([]float64, 0, 64),
					})
					indexEntry := entry.(*ConcurrentIndexEntry)

					// Lock only this entry while updating it
					indexEntry.Lock()
					indexEntry.DocIDs = append(indexEntry.DocIDs, doc.ID)
					// Calculate TF as frequency / total tokens in document
					tf := float64(freq) / float64(totalTokens)
					indexEntry.Freqs = append(indexEntry.Freqs, tf)
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

	idx.calculateIDF()
}

// calculateIDF updates term frequencies with IDF scores
func (idx *ConcurrentIndex) calculateIDF() {
	idx.entries.Range(func(key, value interface{}) bool {
		entry := value.(*ConcurrentIndexEntry)
		entry.Lock()
		defer entry.Unlock()

		// IDF = log(N/(df + 1)) + 1  // Adding 1 to avoid division by zero and negative values
		idf := math.Log(float64(idx.docCount)/(float64(len(entry.DocIDs))+1.0)) + 1.0

		// Update frequencies with TF-IDF score
		for i := range entry.Freqs {
			entry.Freqs[i] *= idf
		}
		return true
	})
}

// Search queries the ConcurrentIndex for the given text and returns scored results
func (idx *ConcurrentIndex) Search(text string) []SearchResult {
	tokens := analyze(text)
	if len(tokens) == 0 {
		return nil
	}

	scores := make(map[int]float64)
	var scoresMutex sync.RWMutex

	// Calculate scores for each token
	for _, token := range tokens {
		if entry, ok := idx.entries.Load(token); ok {
			indexEntry := entry.(*ConcurrentIndexEntry)
			indexEntry.RLock()
			for i, docID := range indexEntry.DocIDs {
				scoresMutex.Lock()
				scores[docID] += indexEntry.Freqs[i]
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
