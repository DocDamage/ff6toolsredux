package presets

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"ffvi_editor/models"
)

// loadResult represents the result of loading a preset
type loadResult struct {
	data interface{}
	err  error
}

// Manager handles party preset CRUD operations and storage
type Manager struct {
	presetsDir string
	presets    map[string]*models.PartyPreset
	mu         sync.RWMutex
}

// NewManager creates a new party preset manager
func NewManager(presetsDir string) (*Manager, error) {
	// Ensure directory exists
	if err := os.MkdirAll(presetsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create presets directory: %w", err)
	}

	m := &Manager{
		presetsDir: presetsDir,
		presets:    make(map[string]*models.PartyPreset),
	}

	// Load existing presets
	if err := m.loadAllPresets(); err != nil {
		return nil, fmt.Errorf("failed to load presets: %w", err)
	}

	return m, nil
}

// CreatePreset saves a new preset
func (m *Manager) CreatePreset(preset *models.PartyPreset) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if preset == nil {
		return fmt.Errorf("preset cannot be nil")
	}

	if err := preset.Validate(); err != nil {
		return fmt.Errorf("invalid preset: %w", err)
	}

	// Check for duplicate name
	for _, p := range m.presets {
		if p.Name == preset.Name {
			return fmt.Errorf("preset with name '%s' already exists", preset.Name)
		}
	}

	// Generate ID if not set
	if preset.ID == "" {
		preset.ID = generatePresetID()
	}

	preset.CreatedAt = time.Now()
	preset.UpdatedAt = time.Now()
	preset.Version = 1

	// Save to file
	if err := m.savePresetToFile(preset); err != nil {
		return err
	}

	m.presets[preset.ID] = preset
	return nil
}

// GetPreset retrieves a preset by ID
func (m *Manager) GetPreset(id string) (*models.PartyPreset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	preset, exists := m.presets[id]
	if !exists {
		return nil, fmt.Errorf("preset not found")
	}
	return preset, nil
}

// GetPresetByName retrieves a preset by name
func (m *Manager) GetPresetByName(name string) (*models.PartyPreset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, p := range m.presets {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("preset not found")
}

// ListPresets returns all presets, optionally filtered by tag
func (m *Manager) ListPresets(filterTag string) []*models.PartyPreset {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*models.PartyPreset, 0)
	for _, p := range m.presets {
		if filterTag == "" {
			result = append(result, p)
		} else {
			for _, tag := range p.Tags {
				if tag == filterTag {
					result = append(result, p)
					break
				}
			}
		}
	}

	// Sort by name
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// UpdatePreset updates an existing preset
func (m *Manager) UpdatePreset(preset *models.PartyPreset) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if preset == nil {
		return fmt.Errorf("preset cannot be nil")
	}

	if err := preset.Validate(); err != nil {
		return fmt.Errorf("invalid preset: %w", err)
	}

	if _, exists := m.presets[preset.ID]; !exists {
		return fmt.Errorf("preset not found")
	}

	preset.UpdatedAt = time.Now()
	preset.Version++

	// Save to file
	if err := m.savePresetToFile(preset); err != nil {
		return err
	}

	m.presets[preset.ID] = preset
	return nil
}

// DeletePreset removes a preset
func (m *Manager) DeletePreset(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	preset, exists := m.presets[id]
	if !exists {
		return fmt.Errorf("preset not found")
	}

	// Delete file
	presetPath := filepath.Join(m.presetsDir, preset.ID+".json")
	if err := os.Remove(presetPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete preset file: %w", err)
	}

	delete(m.presets, id)
	return nil
}

// GetFavorites returns all favorite presets
func (m *Manager) GetFavorites() []*models.PartyPreset {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*models.PartyPreset, 0)
	for _, p := range m.presets {
		if p.Favorite {
			result = append(result, p)
		}
	}

	// Sort by name
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// ToggleFavorite marks/unmarks a preset as favorite
func (m *Manager) ToggleFavorite(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	preset, exists := m.presets[id]
	if !exists {
		return fmt.Errorf("preset not found")
	}

	preset.Favorite = !preset.Favorite
	preset.UpdatedAt = time.Now()

	// Save updated preset
	if err := m.savePresetToFile(preset); err != nil {
		return err
	}

	return nil
}

// GetAllTags returns all unique tags used in presets
func (m *Manager) GetAllTags() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tagMap := make(map[string]bool)
	for _, p := range m.presets {
		for _, tag := range p.Tags {
			tagMap[tag] = true
		}
	}

	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}

	sort.Strings(tags)
	return tags
}

// ExportPreset exports a preset as JSON
func (m *Manager) ExportPreset(id string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	preset, exists := m.presets[id]
	if !exists {
		return nil, fmt.Errorf("preset not found")
	}

	return json.MarshalIndent(preset, "", "  ")
}

// ImportPreset imports a preset from JSON
func (m *Manager) ImportPreset(data []byte) (*models.PartyPreset, error) {
	var preset models.PartyPreset
	if err := json.Unmarshal(data, &preset); err != nil {
		return nil, fmt.Errorf("failed to parse preset JSON: %w", err)
	}

	if err := preset.Validate(); err != nil {
		return nil, fmt.Errorf("invalid preset: %w", err)
	}

	// Generate new ID to avoid conflicts
	preset.ID = generatePresetID()
	preset.CreatedAt = time.Now()
	preset.UpdatedAt = time.Now()
	preset.Version = 1

	return &preset, nil
}

// Private helper functions

func (m *Manager) savePresetToFile(preset *models.PartyPreset) error {
	presetPath := filepath.Join(m.presetsDir, preset.ID+".json")

	data, err := json.MarshalIndent(preset, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal preset: %w", err)
	}

	if err := os.WriteFile(presetPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write preset file: %w", err)
	}

	return nil
}

func (m *Manager) loadAllPresets() error {
	files, err := os.ReadDir(m.presetsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Directory doesn't exist yet, that's OK
		}
		return err
	}

	// Filter JSON files and batch load with worker pool
	var jsonFiles []fs.DirEntry
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			jsonFiles = append(jsonFiles, file)
		}
	}

	if len(jsonFiles) == 0 {
		return nil
	}

	// Use worker pool for concurrent loading (4 workers for file I/O)
	const maxWorkers = 4
	semaphore := make(chan struct{}, maxWorkers)
	results := make(chan loadResult, len(jsonFiles))
	var wg sync.WaitGroup

	for _, file := range jsonFiles {
		wg.Add(1)
		go func(f fs.DirEntry) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire worker slot
			defer func() { <-semaphore }() // Release worker slot

			presetPath := filepath.Join(m.presetsDir, f.Name())
			data, err := os.ReadFile(presetPath)
			if err != nil {
				results <- loadResult{nil, err}
				return
			}

			var preset models.PartyPreset
			if err := json.Unmarshal(data, &preset); err != nil {
				results <- loadResult{nil, err}
				return
			}

			results <- loadResult{&preset, nil}
		}(file)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for result := range results {
		if result.err != nil {
			continue // Skip on error
		}
		if result.data != nil {
			if preset, ok := result.data.(*models.PartyPreset); ok {
				m.presets[preset.ID] = preset
			}
		}
	}

	return nil
}

func generatePresetID() string {
	return time.Now().Format("20060102_150405_") + generateRandomSuffix()
}

func generateRandomSuffix() string {
	return time.Now().Format("000")
}
