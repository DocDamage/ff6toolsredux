package docs

import (
	"fmt"
	"strings"
)

// HelpTopic represents a help topic
type HelpTopic struct {
	ID          string
	Title       string
	Content     string
	Category    string
	Keywords    []string
	RelatedTopics []string
}

// Category constants
const (
	CategoryGettingStarted = "getting-started"
	CategoryFeatures      = "features"
	CategoryTroubleshooting = "troubleshooting"
	CategoryAPI           = "api"
	CategoryAdvanced      = "advanced"
)

// HelpSystem manages help documentation
type HelpSystem struct {
	topics map[string]*HelpTopic
}

// NewHelpSystem creates a new help system
func NewHelpSystem() *HelpSystem {
	h := &HelpSystem{
		topics: make(map[string]*HelpTopic),
	}
	h.initializeTopics()
	return h
}

// initializeTopics creates all help topics
func (h *HelpSystem) initializeTopics() {
	topics := []*HelpTopic{
		// Getting Started
		{
			ID:       "quickstart",
			Title:    "Quick Start Guide",
			Category: CategoryGettingStarted,
			Content: `
## Quick Start Guide

Welcome to the Final Fantasy VI Save Editor!

### Opening a Save File

1. Click File → Open or press Ctrl+O
2. Navigate to your FF6 Pixel Remaster save file
3. The save file will load and display your character data

### Making Edits

- Navigate through tabs to edit different aspects of your save
- Changes are tracked and can be undone with Ctrl+Z
- Save your changes with Ctrl+S

### Safety Features

- Auto-backup is enabled by default (10 backups kept)
- Validation warnings will appear if something looks wrong
- Undo/redo support for all edits
`,
			Keywords: []string{"start", "begin", "open", "load"},
		},
		{
			ID:       "interface",
			Title:    "Understanding the Interface",
			Category: CategoryGettingStarted,
			Content: `
## Interface Overview

The editor is organized into several tabs:

### Validation Tab
- Shows any issues with your save file
- Auto-fix option for common problems
- Health checks and warnings

### Characters Tab
- Edit character stats, level, HP, MP
- Equipment management
- Command customization

### Inventory Tab
- Manage items and quantities
- Important items section
- Quick max/clear buttons

### Skills Tab
- Magic spells
- Abilities (Blitz, Lore, Rage, Dance, etc.)
- Esper assignments

### Party Tab
- Set active party members
- Configure formations

### Map Data Tab
- World map progress
- Locations unlocked
`,
			Keywords: []string{"interface", "tabs", "layout", "ui"},
		},
		
		// Features
		{
			ID:       "templates",
			Title:    "Using Templates",
			Category: CategoryFeatures,
			Content: `
## Character Templates

Templates let you save and reuse character builds.

### Creating a Template

1. Set up a character with desired stats/equipment
2. Click "Save as Template"
3. Give it a name and description
4. Template is saved for future use

### Applying a Template

1. Select a character
2. Click "Load Template"
3. Choose from your templates or community templates
4. Select "Replace" or "Merge" mode
5. Confirm to apply

### Sharing Templates

- Export templates to JSON files
- Share via share codes (Tools → Share Build)
- Upload to community marketplace
`,
			Keywords: []string{"template", "preset", "save", "reuse"},
		},
		{
			ID:       "batch_operations",
			Title:    "Batch Operations",
			Category: CategoryFeatures,
			Content: `
## Batch Operations

Perform operations on multiple items at once.

### Available Operations

- **Max All Stats**: Set all characters to level 99, max HP/MP, max stats
- **Max Items**: Set all items to 99
- **Max Magic**: Teach all spells to all characters
- **Equip Best**: Auto-equip best available gear
- **Clear Inventory**: Remove all items
- **Reset Characters**: Reset to default stats

### Using Batch Operations

1. Go to Tools → Batch Operations
2. Select operation
3. Preview changes
4. Confirm to apply

Note: Batch operations can be undone with Ctrl+Z
`,
			Keywords: []string{"batch", "bulk", "multiple", "max"},
		},
		{
			ID:       "cloud_sync",
			Title:    "Cloud Backup",
			Category: CategoryFeatures,
			Content: `
## Cloud Backup

Sync your save files to Google Drive or Dropbox.

### Setup

1. Go to Preferences → Cloud Sync
2. Choose provider (Google Drive or Dropbox)
3. Click "Authenticate" and sign in
4. Enable sync and configure settings

### Settings

- **Auto-sync**: Sync automatically on save
- **Sync Interval**: How often to sync (minutes)
- **Encryption**: Encrypt files before upload (recommended)
- **Conflict Resolution**: How to handle conflicts

### Conflict Resolution

When local and remote files differ:
- **Use Local**: Keep your local version
- **Use Remote**: Download remote version
- **Create Copy**: Save both versions
- **Prompt**: Ask you what to do
`,
			Keywords: []string{"cloud", "sync", "backup", "drive", "dropbox"},
		},
		{
			ID:       "scripting",
			Title:    "Lua Scripting",
			Category: CategoryAdvanced,
			Content: `
## Lua Scripting

Write custom scripts to automate save file edits.

### Available Functions

**Character Functions**
- getCharacter(id)
- setCharacterLevel(id, level)
- setCharacterHP(id, hp)
- setCharacterMP(id, mp)
- setCharacterStat(id, stat, value)

**Inventory Functions**
- getItemCount(itemId)
- setItemCount(itemId, count)
- addItem(itemId, count)

**Party Functions**
- getPartyMembers()
- setPartyMembers({id1, id2, ...})

**Magic Functions**
- hasMagic(charId, spellId)
- learnMagic(charId, spellId)

### Example Script

` + "```lua" + `
-- Set Terra to level 99 with max stats
setCharacterLevel(0, 99)
setCharacterHP(0, 9999)
setCharacterMP(0, 999)
setCharacterStat(0, "vigor", 128)
print("Terra maximized!")
` + "```" + `

### Running Scripts

1. Go to Tools → Script Editor
2. Write or load a script
3. Click "Run" to execute
4. Check output/log for results
`,
			Keywords: []string{"script", "lua", "automate", "programming"},
		},
		
		// Troubleshooting
		{
			ID:       "validation_errors",
			Title:    "Fixing Validation Errors",
			Category: CategoryTroubleshooting,
			Content: `
## Common Validation Errors

### Level Out of Range
**Error**: Character level must be between 1 and 99
**Fix**: Adjust level slider to valid range

### HP/MP Too High
**Error**: HP/MP exceeds maximum
**Fix**: Click "Auto-fix" or manually adjust values

### Invalid Equipment
**Error**: Equipment not compatible with character
**Fix**: Change equipment or use "Remove Invalid"

### Corrupted Save
**Error**: Save file structure invalid
**Fix**: Restore from backup or use last known good save

### Using Auto-fix

1. Go to Validation tab
2. Review errors
3. Click "Auto-fix All"
4. Verify changes
5. Save file
`,
			Keywords: []string{"error", "validation", "fix", "problem", "issue"},
		},
		{
			ID:       "backup_restore",
			Title:    "Backup and Restore",
			Category: CategoryTroubleshooting,
			Content: `
## Backup and Restore

### Auto-Backup

Auto-backup creates a copy before each save:
- Enabled by default
- Keeps last 10 backups
- Stored in ~/. ffvi-editor/backups/

### Manual Backup

1. File → Create Backup
2. Choose location
3. Save backup file

### Restoring from Backup

1. File → Restore from Backup
2. Select backup from list
3. Preview changes
4. Confirm to restore

### Backup Management

1. File → Manage Backups
2. View all backups with dates
3. Delete old backups
4. Export backup to file
`,
			Keywords: []string{"backup", "restore", "recover", "save"},
		},
		{
			ID:       "keyboard_shortcuts",
			Title:    "Keyboard Shortcuts",
			Category: CategoryFeatures,
			Content: `
## Keyboard Shortcuts

### File Operations
- **Ctrl+N**: New file
- **Ctrl+O**: Open file
- **Ctrl+S**: Save file
- **Ctrl+W**: Close file
- **Ctrl+Q**: Quit application

### Editing
- **Ctrl+Z**: Undo
- **Ctrl+Y**: Redo
- **Ctrl+F**: Find/Search
- **Ctrl+Shift+P**: Command palette

### Tools
- **Ctrl+B**: Create backup
- **Ctrl+Shift+V**: Validate save
- **Ctrl+,**: Preferences

### Customizing Shortcuts

1. Go to Preferences → Shortcuts
2. Click on shortcut to edit
3. Press new key combination
4. Click Save

Note: Some system shortcuts cannot be overridden.
`,
			Keywords: []string{"keyboard", "shortcuts", "hotkeys", "keys"},
		},
	}

	for _, topic := range topics {
		h.topics[topic.ID] = topic
	}
}

// GetTopic retrieves a help topic by ID
func (h *HelpSystem) GetTopic(id string) (*HelpTopic, error) {
	topic, exists := h.topics[id]
	if !exists {
		return nil, fmt.Errorf("topic not found: %s", id)
	}
	return topic, nil
}

// GetAllTopics returns all help topics
func (h *HelpSystem) GetAllTopics() []*HelpTopic {
	topics := make([]*HelpTopic, 0, len(h.topics))
	for _, topic := range h.topics {
		topics = append(topics, topic)
	}
	return topics
}

// GetTopicsByCategory returns topics in a category
func (h *HelpSystem) GetTopicsByCategory(category string) []*HelpTopic {
	topics := make([]*HelpTopic, 0)
	for _, topic := range h.topics {
		if topic.Category == category {
			topics = append(topics, topic)
		}
	}
	return topics
}

// Search searches help topics
func (h *HelpSystem) Search(query string) []*HelpTopic {
	query = strings.ToLower(query)
	results := make([]*HelpTopic, 0)

	for _, topic := range h.topics {
		// Search in title
		if strings.Contains(strings.ToLower(topic.Title), query) {
			results = append(results, topic)
			continue
		}

		// Search in keywords
		for _, keyword := range topic.Keywords {
			if strings.Contains(strings.ToLower(keyword), query) {
				results = append(results, topic)
				break
			}
		}

		// Search in content
		if strings.Contains(strings.ToLower(topic.Content), query) {
			// Only add if not already in results
			found := false
			for _, r := range results {
				if r.ID == topic.ID {
					found = true
					break
				}
			}
			if !found {
				results = append(results, topic)
			}
		}
	}

	return results
}

// GetCategories returns all available categories
func (h *HelpSystem) GetCategories() []string {
	return []string{
		CategoryGettingStarted,
		CategoryFeatures,
		CategoryTroubleshooting,
		CategoryAPI,
		CategoryAdvanced,
	}
}
