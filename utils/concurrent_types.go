package utils

import "sync"

// ConcurrentIndexEntry stores document IDs and their frequencies with thread-safe access
type ConcurrentIndexEntry struct {
	sync.RWMutex
	DocIDs []int
	Freqs  []float64
}
