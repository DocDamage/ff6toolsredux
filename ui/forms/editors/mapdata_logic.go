package editors

import (
	"sort"
	"strings"

	"ffvi_editor/models/consts"
	"ffvi_editor/settings"
)

// mapLocationItem represents a location in the world map
type mapLocationItem struct {
	Name        string
	HasFallback bool
	X           float64
	Y           float64
}

// mapDataLogic handles the business logic for map data operations
type mapDataLogic struct {
	world             int
	filteredLocations []mapLocationItem
}

// newMapDataLogic creates a new map data logic handler
func newMapDataLogic() *mapDataLogic {
	return &mapDataLogic{
		world:             1,
		filteredLocations: make([]mapLocationItem, 0),
	}
}

// defaultWorldLocationNames returns the default location names for a world
func defaultWorldLocationNames(world int) []string {
	if world == 2 {
		return []string{
			"Albrook",
			"Cave on the Veldt",
			"Colosseum",
			"Darill's Tomb",
			"Doma Castle",
			"Duncan's House",
			"Ebot's Rock",
			"Fanatics' Tower",
			"Figaro Castle",
			"Figaro Cave",
			"Gau's Father's House",
			"Jidoor",
			"Kefka's Tower",
			"Kohlingen",
			"Maranda",
			"Mobliz",
			"Narshe",
			"Nikeah",
			"Opera House",
			"Phoenix Cave",
			"Solitary Island",
			"South Figaro",
			"Thamasa",
			"The Ancient Castle",
			"The Veldt",
			"Triangle Island",
			"Tzen",
			"Zozo",
		}
	}
	// World of Balance - only include locations that have coordinates
	return []string{
		"Narshe",
		"Kohlingen",
		"Zozo",
		"Jidoor",
		"Opera House",
		"Maranda",
		"Vector",
		"Albrook",
		"Sealed Gate",
		"Tzen",
		"South Figaro",
		"Figaro Castle",
		"Figaro Cave",
		"Mt. Kolts",
		"Returners' Hideout",
		"Doma Castle",
		"Crescent Mountain",
		"Thamasa",
		"Espers' Gathering Place",
		"House in the Veldt",
		"Sabin's Cabin",
		"The Veldt",
		"Barren Falls",
		"Phantom Forest",
		"Imperial Base",
		"Triangle Island",
	}
}

// rebuildLocations rebuilds the filtered locations list based on search query
func (m *mapDataLogic) rebuildLocations(world int, searchQuery string) []mapLocationItem {
	m.world = world
	q := strings.ToLower(strings.TrimSpace(searchQuery))

	byName := make(map[string]mapLocationItem)

	// Add defaults from the map artwork
	for _, name := range defaultWorldLocationNames(world) {
		byName[name] = mapLocationItem{Name: name}
	}

	// Start with any user-added location names
	sm := settings.NewManager("")
	_ = sm.Load()
	for _, name := range sm.GetMapLocations(world) {
		byName[name] = mapLocationItem{Name: name}
	}

	// Include any names already calibrated
	for name := range sm.GetAllMapPoints(world) {
		if name == "" {
			continue
		}
		if _, ok := byName[name]; !ok {
			byName[name] = mapLocationItem{Name: name}
		}
	}

	// Merge in built-in fallbacks (where we have coordinates)
	for _, lm := range consts.Landmarks {
		if lm.World != world {
			continue
		}
		item, ok := byName[lm.Name]
		if !ok {
			item = mapLocationItem{Name: lm.Name}
		}
		item.HasFallback = true
		item.X = lm.X
		item.Y = lm.Y
		byName[lm.Name] = item
	}

	m.filteredLocations = m.filteredLocations[:0]
	for _, item := range byName {
		if q != "" && !strings.Contains(strings.ToLower(item.Name), q) {
			continue
		}
		m.filteredLocations = append(m.filteredLocations, item)
	}

	sort.Slice(m.filteredLocations, func(i, j int) bool {
		return strings.ToLower(m.filteredLocations[i].Name) < strings.ToLower(m.filteredLocations[j].Name)
	})

	return m.filteredLocations
}

// getFilteredLocations returns the current filtered locations
func (m *mapDataLogic) getFilteredLocations() []mapLocationItem {
	return m.filteredLocations
}

// findNearestLandmark finds the nearest landmark to the given coordinates
func findNearestLandmark(world int, x, y float64) (name string, lx, ly float64, found bool) {
	minDist := 1e9
	for _, lm := range consts.Landmarks {
		if lm.World != world {
			continue
		}
		dx := lm.X - x
		dy := lm.Y - y
		dist := dx*dx + dy*dy
		if dist < minDist {
			minDist = dist
			lx = lm.X
			ly = lm.Y
			name = lm.Name
			found = true
		}
	}
	return
}

// getMapImagePath returns the map image path for a world
func getMapImagePath(world int) string {
	if world == 2 {
		return "ui/embedded/resources/ff6-world-map_world-of-ruin.png"
	}
	return "ui/embedded/resources/ff6-world-map_world-of-balance.png"
}
