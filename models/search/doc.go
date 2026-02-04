// Package search provides search indexing functionality.
//
// The search package handles:
//   - Full-text search indexing
//   - Fuzzy matching
//   - Search result ranking
//
// Searchable content:
//   - Items
//   - Spells
//   - Characters
//   - Enemies
//
// Example usage:
//
//	// Create search index
//	index := search.NewIndex()
//
//	// Add documents
//	index.Add(itemID, itemName, itemDescription)
//
//	// Search
//	results := index.Search("fire")
//	for _, result := range results {
//	    fmt.Printf("Found: %s (score: %f)\n", result.Name, result.Score)
//	}
package search
