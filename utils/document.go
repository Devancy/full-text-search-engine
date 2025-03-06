package utils

import (
	"compress/gzip"
	"encoding/xml"
	"os"
	"runtime"
	"sync"
)

// Document represents a Wikipedia abstract dump Document.
type Document struct {
	Title string `xml:"title"`
	URL   string `xml:"url"`
	Text  string `xml:"abstract"`
	ID    int
}

// LoadDocuments parses a Wikipedia abstract dump and returns a slice of documents.
// Dump example: https://dumps.wikimedia.your.org/enwiki/latest/enwiki-latest-abstract1.xml.gz
func LoadDocuments(path string) ([]*Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	dec := xml.NewDecoder(gz)
	dump := struct {
		Documents []*Document `xml:"doc"`
	}{}
	if err := dec.Decode(&dump); err != nil {
		return nil, err
	}

	// Use a worker pool to assign IDs concurrently
	numWorkers := runtime.NumCPU()
	docs := dump.Documents
	chunkSize := len(docs) / numWorkers
	var wg sync.WaitGroup

	for i := range numWorkers {
		wg.Add(1)
		start := i * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			end = len(docs)
		}

		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				docs[i].ID = i
			}
		}(start, end)
	}
	wg.Wait()

	return docs, nil
}
