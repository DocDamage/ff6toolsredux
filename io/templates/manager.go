package templates

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	templatesModel "ffvi_editor/models/templates"
)

// loadResult represents the result of loading a template
type loadResult struct {
	data interface{}
	err  error
}

// Manager handles template CRUD operations and storage
type Manager struct {
	templatesDir string
	templates    map[string]*templatesModel.CharacterTemplate
	mu           sync.RWMutex
	builtInTags  []string
}

// NewManager creates a new template manager
func NewManager(templatesDir string) (*Manager, error) {
	// Ensure directory exists
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create templates directory: %w", err)
	}

	m := &Manager{
		templatesDir: templatesDir,
		templates:    make(map[string]*templatesModel.CharacterTemplate),
		builtInTags:  []string{"speedrun", "challenge", "casual", "optimal", "dps", "tank", "healer"},
	}

	// Load existing templates
	if err := m.loadAllTemplates(); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	return m, nil
}

// CreateTemplate saves a new template
func (m *Manager) CreateTemplate(template *templatesModel.CharacterTemplate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if template == nil {
		return templatesModel.ErrInvalidTemplate
	}

	// Check for duplicate name
	for _, t := range m.templates {
		if t.Name == template.Name {
			return templatesModel.ErrDuplicateName
		}
	}

	// Save to file
	if err := m.saveTemplateToFile(template); err != nil {
		return err
	}

	m.templates[template.ID] = template
	return nil
}

// GetTemplate retrieves a template by ID
func (m *Manager) GetTemplate(id string) (*templatesModel.CharacterTemplate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	template, exists := m.templates[id]
	if !exists {
		return nil, templatesModel.ErrTemplateNotFound
	}
	return template, nil
}

// GetTemplateByName retrieves a template by name
func (m *Manager) GetTemplateByName(name string) (*templatesModel.CharacterTemplate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, t := range m.templates {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, templatesModel.ErrTemplateNotFound
}

// ListTemplates returns all templates, optionally filtered by tag
func (m *Manager) ListTemplates(filterTag string) []*templatesModel.CharacterTemplate {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*templatesModel.CharacterTemplate, 0)
	for _, t := range m.templates {
		if filterTag == "" {
			result = append(result, t)
		} else {
			for _, tag := range t.Tags {
				if tag == filterTag {
					result = append(result, t)
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

// UpdateTemplate updates an existing template
func (m *Manager) UpdateTemplate(template *templatesModel.CharacterTemplate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if template == nil {
		return templatesModel.ErrInvalidTemplate
	}

	if _, exists := m.templates[template.ID]; !exists {
		return templatesModel.ErrTemplateNotFound
	}

	template.UpdatedAt = time.Now()
	template.Version++

	// Save to file
	if err := m.saveTemplateToFile(template); err != nil {
		return err
	}

	m.templates[template.ID] = template
	return nil
}

// DeleteTemplate removes a template
func (m *Manager) DeleteTemplate(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	template, exists := m.templates[id]
	if !exists {
		return templatesModel.ErrTemplateNotFound
	}

	// Delete file
	templatePath := filepath.Join(m.templatesDir, template.ID+".json")
	if err := os.Remove(templatePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete template file: %w", err)
	}

	delete(m.templates, id)
	return nil
}

// GetFavorites returns all favorite templates
func (m *Manager) GetFavorites() []*templatesModel.CharacterTemplate {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*templatesModel.CharacterTemplate, 0)
	for _, t := range m.templates {
		if t.Favorite {
			result = append(result, t)
		}
	}

	// Sort by name
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// ToggleFavorite marks/unmarks a template as favorite
func (m *Manager) ToggleFavorite(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	template, exists := m.templates[id]
	if !exists {
		return templatesModel.ErrTemplateNotFound
	}

	template.Favorite = !template.Favorite
	template.UpdatedAt = time.Now()

	// Save updated template
	if err := m.saveTemplateToFile(template); err != nil {
		return err
	}

	return nil
}

// GetAllTags returns all unique tags used in templates
func (m *Manager) GetAllTags() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tagMap := make(map[string]bool)
	for _, t := range m.templates {
		for _, tag := range t.Tags {
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

// ExportTemplate exports a template as JSON
func (m *Manager) ExportTemplate(id string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	template, exists := m.templates[id]
	if !exists {
		return nil, templatesModel.ErrTemplateNotFound
	}

	return json.MarshalIndent(template, "", "  ")
}

// ImportTemplate imports a template from JSON
func (m *Manager) ImportTemplate(data []byte) (*templatesModel.CharacterTemplate, error) {
	var template templatesModel.CharacterTemplate
	if err := json.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("failed to parse template JSON: %w", err)
	}

	if template.Character == nil {
		return nil, templatesModel.ErrInvalidTemplate
	}

	// Generate new ID to avoid conflicts
	template.ID = generateTemplateID()
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	template.Version = 1

	return &template, nil
}

// Private helper functions

func (m *Manager) saveTemplateToFile(template *templatesModel.CharacterTemplate) error {
	templatePath := filepath.Join(m.templatesDir, template.ID+".json")

	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	if err := os.WriteFile(templatePath, data, 0644); err != nil {
		return templatesModel.ErrIOError
	}

	return nil
}

func (m *Manager) loadAllTemplates() error {
	files, err := os.ReadDir(m.templatesDir)
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

			templatePath := filepath.Join(m.templatesDir, f.Name())
			data, err := os.ReadFile(templatePath)
			if err != nil {
				results <- loadResult{nil, err}
				return
			}

			var template templatesModel.CharacterTemplate
			if err := json.Unmarshal(data, &template); err != nil {
				results <- loadResult{nil, err}
				return
			}

			results <- loadResult{&template, nil}
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
			if template, ok := result.data.(*templatesModel.CharacterTemplate); ok {
				m.templates[template.ID] = template
			}
		}
	}

	return nil
}

func generateTemplateID() string {
	return time.Now().Format("20060102_150405_") + generateRandomSuffix()
}

func generateRandomSuffix() string {
	return time.Now().Format("000")
}
