package search

// SearchResult represents a single search result
type SearchResult struct {
	Type        ResultType
	Name        string
	Description string
	Details     map[string]interface{}
	ID          string // For filtering/jumping
}

// ResultType defines the type of result
type ResultType string

const (
	ResultTypeCharacter ResultType = "Character"
	ResultTypeItem      ResultType = "Item"
	ResultTypeSpell     ResultType = "Spell"
	ResultTypeEsper     ResultType = "Esper"
	ResultTypeEquipment ResultType = "Equipment"
)

// Index provides searchable access to save data
type Index struct {
	characters []SearchResult
	items      []SearchResult
	spells     []SearchResult
	espers     []SearchResult
	equipment  []SearchResult
}

// NewIndex creates a new search index
func NewIndex() *Index {
	return &Index{
		characters: make([]SearchResult, 0),
		items:      make([]SearchResult, 0),
		spells:     make([]SearchResult, 0),
		espers:     make([]SearchResult, 0),
		equipment:  make([]SearchResult, 0),
	}
}

// AddCharacter adds a character to the index
func (idx *Index) AddCharacter(id string, name string, level int) {
	idx.characters = append(idx.characters, SearchResult{
		Type:        ResultTypeCharacter,
		Name:        name,
		Description: "Character",
		Details: map[string]interface{}{
			"level": level,
		},
		ID: id,
	})
}

// AddItem adds an item to the index
func (idx *Index) AddItem(id string, name string, itemType string) {
	idx.items = append(idx.items, SearchResult{
		Type:        ResultTypeItem,
		Name:        name,
		Description: itemType,
		Details: map[string]interface{}{
			"type": itemType,
		},
		ID: id,
	})
}

// AddSpell adds a spell to the index
func (idx *Index) AddSpell(id string, name string, magicType string) {
	idx.spells = append(idx.spells, SearchResult{
		Type:        ResultTypeSpell,
		Name:        name,
		Description: magicType,
		Details: map[string]interface{}{
			"type": magicType,
		},
		ID: id,
	})
}

// AddEsper adds an esper to the index
func (idx *Index) AddEsper(id string, name string) {
	idx.espers = append(idx.espers, SearchResult{
		Type:        ResultTypeEsper,
		Name:        name,
		Description: "Esper",
		Details:     make(map[string]interface{}),
		ID:          id,
	})
}

// AddEquipment adds equipment to the index
func (idx *Index) AddEquipment(id string, name string, equipType string) {
	idx.equipment = append(idx.equipment, SearchResult{
		Type:        ResultTypeEquipment,
		Name:        name,
		Description: equipType,
		Details: map[string]interface{}{
			"type": equipType,
		},
		ID: id,
	})
}

// Search performs a case-insensitive search across all indexed items
func (idx *Index) Search(query string) []SearchResult {
	results := make([]SearchResult, 0)

	query = lowerString(query)

	// Search characters
	for _, result := range idx.characters {
		if matchesQuery(result.Name, query) {
			results = append(results, result)
		}
	}

	// Search items
	for _, result := range idx.items {
		if matchesQuery(result.Name, query) {
			results = append(results, result)
		}
	}

	// Search spells
	for _, result := range idx.spells {
		if matchesQuery(result.Name, query) {
			results = append(results, result)
		}
	}

	// Search espers
	for _, result := range idx.espers {
		if matchesQuery(result.Name, query) {
			results = append(results, result)
		}
	}

	// Search equipment
	for _, result := range idx.equipment {
		if matchesQuery(result.Name, query) {
			results = append(results, result)
		}
	}

	return results
}

// SearchByType performs a search filtered by result type
func (idx *Index) SearchByType(query string, resultType ResultType) []SearchResult {
	results := make([]SearchResult, 0)
	query = lowerString(query)

	switch resultType {
	case ResultTypeCharacter:
		for _, result := range idx.characters {
			if matchesQuery(result.Name, query) {
				results = append(results, result)
			}
		}
	case ResultTypeItem:
		for _, result := range idx.items {
			if matchesQuery(result.Name, query) {
				results = append(results, result)
			}
		}
	case ResultTypeSpell:
		for _, result := range idx.spells {
			if matchesQuery(result.Name, query) {
				results = append(results, result)
			}
		}
	case ResultTypeEsper:
		for _, result := range idx.espers {
			if matchesQuery(result.Name, query) {
				results = append(results, result)
			}
		}
	case ResultTypeEquipment:
		for _, result := range idx.equipment {
			if matchesQuery(result.Name, query) {
				results = append(results, result)
			}
		}
	}

	return results
}

// Clear resets the index
func (idx *Index) Clear() {
	idx.characters = make([]SearchResult, 0)
	idx.items = make([]SearchResult, 0)
	idx.spells = make([]SearchResult, 0)
	idx.espers = make([]SearchResult, 0)
	idx.equipment = make([]SearchResult, 0)
}

// Helper functions

// matchesQuery checks if a name contains the query using fuzzy matching
func matchesQuery(name, query string) bool {
	name = lowerString(name)

	// Exact match
	if name == query {
		return true
	}

	// Contains
	if contains(name, query) {
		return true
	}

	// Fuzzy match: check if query characters appear in order
	return fuzzyMatch(name, query)
}

// fuzzyMatch implements simple fuzzy matching
func fuzzyMatch(text, pattern string) bool {
	patternIdx := 0
	for _, char := range text {
		if patternIdx < len(pattern) && char == rune(pattern[patternIdx]) {
			patternIdx++
		}
	}
	return patternIdx == len(pattern)
}

// contains checks if text contains substr
func contains(text, substr string) bool {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// lowerString converts a string to lowercase
func lowerString(s string) string {
	result := ""
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			result += string(char + 32)
		} else {
			result += string(char)
		}
	}
	return result
}
