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

// config holds the application configuration values derived from flags.
type config struct {
	dumpPath      string
	useConcurrent bool
	maxResults    int
}

func main() {
	setupLogging()
	cfg := parseFlags()

	log.Println("Running Full Text Search Engine")

	docs, err := loadDocuments(cfg.dumpPath)
	if err != nil {
		log.Fatalf("Initialization error: %v", err)
	}

	idx, err := createAndPopulateIndex(docs, cfg.useConcurrent)
	if err != nil {
		log.Fatalf("Initialization error: %v", err)
	}

	if err := runInteractiveSearch(idx, docs, cfg); err != nil {
		log.Fatalf("Runtime error: %v", err)
	}
}

// setupLogging configures the log output format.
func setupLogging() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetPrefix("[Search Engine] ")
}

// parseFlags parses command-line flags and returns a config struct.
func parseFlags() (cfg config) {
	flag.StringVar(&cfg.dumpPath, "p", "enwiki-latest-abstract1.xml.gz", "wiki abstract dump path")
	flag.BoolVar(&cfg.useConcurrent, "c", false, "use concurrent indexing")
	flag.IntVar(&cfg.maxResults, "n", 5, "maximum number of results to display")
	flag.Parse()
	return cfg
}

// loadDocuments loads documents from the specified path and validates the path.
func loadDocuments(dumpPath string) ([]*utils.Document, error) {
	if _, err := os.Stat(dumpPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("dump file not found: %s", dumpPath)
	}

	start := time.Now()
	log.Printf("Loading documents from %s...", dumpPath)
	docs, err := utils.LoadDocuments(dumpPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load documents: %w", err)
	}
	log.Printf("Loaded %d documents in %v", len(docs), time.Since(start))
	return docs, nil
}

// createAndPopulateIndex creates the appropriate indexer (concurrent or simple) and adds documents.
func createAndPopulateIndex(docs []*utils.Document, useConcurrent bool) (utils.Indexer, error) {
	start := time.Now()
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
	return idx, nil
}

// runInteractiveSearch handles the main user interaction loop for searching.
func runInteractiveSearch(idx utils.Indexer, docs []*utils.Document, cfg config) error {
	if err := keyboard.Open(); err != nil {
		return fmt.Errorf("failed to initialize keyboard: %w", err)
	}
	defer keyboard.Close()

	fmt.Println("\nEnter your search query (press Ctrl+C to exit, Enter to search):")
	fmt.Print("> ")

	var query strings.Builder
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			// Log non-critical keyboard errors and continue
			log.Printf("Keyboard error: %v", err)
			continue
		}

		switch key {
		case keyboard.KeyCtrlC:
			fmt.Println("\nExiting...")
			return nil // Normal exit
		case keyboard.KeyEnter:
			if query.Len() > 0 {
				queryString := query.String()
				results := performSearch(idx, queryString)
				fmt.Printf("\nSearch Results for: %q\n", queryString)
				displayResults(results, docs, cfg.maxResults)
				query.Reset()
				fmt.Println("\nEnter your search query (press Ctrl+C to exit, Enter to search):")
				fmt.Print("> ")
			} else {
				// Handle Enter press when query is empty (do nothing, just reprint prompt)
				fmt.Print("\r> ")
			}
		case keyboard.KeyBackspace, keyboard.KeyBackspace2:
			if query.Len() > 0 {
				str := query.String()
				query.Reset()
				query.WriteString(str[:len(str)-1])
				// Clear line and reprint query
				fmt.Printf("\r%s", strings.Repeat(" ", 100))
				fmt.Printf("\r> %s", query.String())
			}
		default:
			// Handle printable characters
			if char != 0 {
				query.WriteRune(char)
				fmt.Printf("\r> %s", query.String())
			}
		}
	}
}

// displayResults handles printing search results with pagination.
func displayResults(results []utils.SearchResult, docs []*utils.Document, pageSize int) {
	if len(results) == 0 {
		fmt.Println("No matches found.")
		return
	}

	startIndex := 0
displayLoop:
	for {
		endIndex := min(startIndex+pageSize, len(results))

		// Print header only for the first page
		if startIndex == 0 {
			fmt.Println("\nResults (sorted by relevance):")
			fmt.Println(strings.Repeat("-", 80))
		}

		// Print results for the current page
		for i := startIndex; i < endIndex; i++ {
			result := results[i]
			// Ensure DocID is within bounds
			if result.DocID >= 0 && result.DocID < len(docs) {
				doc := docs[result.DocID]
				fmt.Printf("\n%d. %s\n", i+1, doc.Title)
				fmt.Printf("   Score: %.4f\n", result.Score)
				fmt.Printf("   URL: %s\n", doc.URL)
				fmt.Printf("   %s\n", doc.Text)
				fmt.Println(strings.Repeat("-", 80))
			} else {
				log.Printf("Warning: Invalid DocID %d found in search results.", result.DocID)
			}
		}

		startIndex = endIndex

		// Check if more results are available
		if startIndex < len(results) {
			remaining := len(results) - startIndex
			nextCount := min(remaining, pageSize)
			fmt.Printf("\nPress Enter for next %d results (%d remaining), or any other key to return to query...\n", nextCount, remaining)

			// Wait for user input for pagination
			_, pageKey, pageErr := keyboard.GetKey()
			if pageErr != nil {
				log.Printf("Keyboard error during pagination: %v", pageErr)
				break displayLoop
			}

			switch pageKey {
			case keyboard.KeyEnter:
				continue displayLoop // Show next page
			case keyboard.KeyCtrlC:
				fmt.Println("\nExiting...")
				// We need a way to signal exit from here back to main or os.Exit
				// For simplicity now, just break the display loop. A channel or error return could be better.
				os.Exit(0) // Force exit if Ctrl+C pressed during pagination
			default:
				break displayLoop // Exit pagination loop, return to query prompt
			}
		} else { // No more results to display
			fmt.Println("\nEnd of results.")
			break displayLoop
		}
	}
}

// performSearch searches the index and returns all matching results sorted by relevance.
func performSearch(idx utils.Indexer, query string) []utils.SearchResult {
	start := time.Now()
	log.Printf("Searching for: %q", query)
	results := idx.Search(query)
	log.Printf("Search completed in %v, found %d results.", time.Since(start), len(results))
	return results
}
