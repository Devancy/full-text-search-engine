package utils

// Indexer defines the interface for full-text search index implementations
type Indexer interface {
	// Add adds documents to the index and updates TF-IDF scores
	Add(docs []*Document)

	// Search performs a full-text search and returns scored results
	Search(text string) []SearchResult

	// Stats returns statistics about the index
	Stats() IndexStats

	// Clear removes all documents from the index
	Clear()
}

// IndexStats contains statistics about the index
type IndexStats struct {
	DocumentCount int     // Total number of documents
	TermCount     int     // Total number of unique terms
	AvgDocLength  float64 // Average document length (in terms)
	MaxScore      float64 // Maximum score in the index
	MinScore      float64 // Minimum score in the index
	IndexSizeKB   int64   // Approximate size of the index in KB
}
