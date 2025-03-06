package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	utils "github.com/devancy/full-text-search-engine/utils"
	"github.com/eiannone/keyboard"
)

func main() {
	// Set up logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetPrefix("[Search Engine] ")

	// Parse command line flags
	var (
		dumpPath      = "enwiki-latest-abstract1.xml.gz"
		useConcurrent bool
		maxResults    int
	)
	flag.StringVar(&dumpPath, "p", "enwiki-latest-abstract1.xml.gz", "wiki abstract dump path")
	flag.BoolVar(&useConcurrent, "c", false, "use concurrent indexing")
	flag.IntVar(&maxResults, "n", 5, "maximum number of results to display")
	flag.Parse()

	log.Println("Running Full Text Search Engine")

	// Verify dump file exists
	if _, err := os.Stat(dumpPath); os.IsNotExist(err) {
		log.Fatalf("Dump file not found: %s", dumpPath)
	}

	// Load documents
	start := time.Now()
	log.Printf("Loading documents from %s...", dumpPath)
	docs, err := utils.LoadDocuments(dumpPath)
	if err != nil {
		log.Fatalf("Failed to load documents: %v", err)
	}
	log.Printf("Loaded %d documents in %v", len(docs), time.Since(start))

	// Create and initialize index
	start = time.Now()
	var idx utils.Indexer
	if useConcurrent {
		idx = utils.NewConcurrentIndex()
		log.Println("Using concurrent index")
	} else {
		idx = utils.NewIndex()
		log.Println("Using simple index")
	}

	log.Println("Indexing documents...")
	idx.Add(docs)
	log.Printf("Indexed %d documents in %v", len(docs), time.Since(start))

	// Set up interactive search
	if err := keyboard.Open(); err != nil {
		log.Fatalf("Failed to initialize keyboard: %v", err)
	}
	defer keyboard.Close()

	fmt.Println("\nEnter your search query (press Ctrl+C to exit, Enter to search):")
	fmt.Print("> ")

	var query strings.Builder
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Printf("Keyboard error: %v", err)
			continue
		}

		switch key {
		case keyboard.KeyCtrlC:
			fmt.Println("\nExiting...")
			return
		case keyboard.KeyEnter:
			if query.Len() > 0 {
				performSearch(idx, docs, query.String(), maxResults)
				query.Reset()
				fmt.Println("\nEnter your search query (press Ctrl+C to exit, Enter to search):")
				fmt.Print("> ")
			}
		case keyboard.KeyBackspace, keyboard.KeyBackspace2:
			if query.Len() > 0 {
				str := query.String()
				query.Reset()
				query.WriteString(str[:len(str)-1])
				fmt.Printf("\r%s", strings.Repeat(" ", 100)) // Clear line
				fmt.Printf("\r> %s", query.String())
			}
		default:
			if char != 0 {
				query.WriteRune(char)
				fmt.Printf("\r> %s", query.String())
			}
		}
	}
}

func performSearch(idx utils.Indexer, docs []*utils.Document, query string, maxResults int) {
	results := idx.Search(query)
	fmt.Printf("\nSearch Results for: %q\n", query)

	if len(results) == 0 {
		fmt.Println("No matches found.")
		return
	}

	fmt.Println("\nResults (sorted by relevance):")
	fmt.Println(strings.Repeat("-", 80))

	for i, result := range results {
		if i >= maxResults {
			fmt.Printf("\n... and %d more results\n", len(results)-maxResults)
			break
		}
		doc := docs[result.DocID]
		fmt.Printf("\n%d. %s\n", i+1, doc.Title)
		fmt.Printf("   Score: %.4f\n", result.Score)
		fmt.Printf("   URL: %s\n", doc.URL)
		fmt.Printf("   %s\n", doc.Text)
		fmt.Println(strings.Repeat("-", 80))
	}
}
