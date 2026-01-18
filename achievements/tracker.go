package achievements

import (
	"sync"
	"time"
)

// Achievement represents an unlockable achievement
type Achievement struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Category    string    `json:"category"` // "editor", "game", "community"
	Points      int       `json:"points"`
	Hidden      bool      `json:"hidden"`
	UnlockedAt  time.Time `json:"unlockedAt,omitempty"`
	Unlocked    bool      `json:"unlocked"`
	Progress    int       `json:"progress"`    // Current progress
	MaxProgress int       `json:"maxProgress"` // Required for unlock
}

// Category types
const (
	CategoryEditor    = "editor"
	CategoryGame      = "game"
	CategoryCommunity = "community"
)

// Tracker tracks achievement progress
type Tracker struct {
	achievements map[string]*Achievement
	mu           sync.RWMutex
	onUnlock     func(*Achievement)
}

// NewTracker creates a new achievement tracker
func NewTracker() *Tracker {
	t := &Tracker{
		achievements: make(map[string]*Achievement),
	}
	t.initializeAchievements()
	return t
}

// initializeAchievements creates all achievements
func (t *Tracker) initializeAchievements() {
	achievements := []*Achievement{
		// Editor Achievements
		{
			ID:          "first_edit",
			Name:        "First Steps",
			Description: "Make your first edit to a save file",
			Category:    CategoryEditor,
			Points:      10,
			MaxProgress: 1,
		},
		{
			ID:          "save_10",
			Name:        "Frequent Editor",
			Description: "Save 10 different files",
			Category:    CategoryEditor,
			Points:      25,
			MaxProgress: 10,
		},
		{
			ID:          "undo_master",
			Name:        "Undo Master",
			Description: "Use undo/redo 50 times",
			Category:    CategoryEditor,
			Points:      20,
			MaxProgress: 50,
		},
		{
			ID:          "template_creator",
			Name:        "Template Creator",
			Description: "Create and save 5 templates",
			Category:    CategoryEditor,
			Points:      30,
			MaxProgress: 5,
		},
		{
			ID:          "batch_operator",
			Name:        "Batch Operator",
			Description: "Use batch operations 10 times",
			Category:    CategoryEditor,
			Points:      25,
			MaxProgress: 10,
		},
		
		// Game Achievements
		{
			ID:          "max_level",
			Name:        "Maximum Power",
			Description: "Set a character to level 99",
			Category:    CategoryGame,
			Points:      15,
			MaxProgress: 1,
		},
		{
			ID:          "all_magic",
			Name:        "Master Mage",
			Description: "Give a character all magic spells",
			Category:    CategoryGame,
			Points:      30,
			MaxProgress: 1,
		},
		{
			ID:          "all_rages",
			Name:        "Rage Collector",
			Description: "Unlock all 254 Gau rages",
			Category:    CategoryGame,
			Points:      50,
			MaxProgress: 254,
		},
		{
			ID:          "speedrun_setup",
			Name:        "Speedrunner",
			Description: "Apply a speedrun preset configuration",
			Category:    CategoryGame,
			Points:      20,
			MaxProgress: 1,
		},
		{
			ID:          "esper_optimizer",
			Name:        "Esper Optimizer",
			Description: "Use the esper bonus calculator 10 times",
			Category:    CategoryGame,
			Points:      25,
			MaxProgress: 10,
		},
		
		// Community Achievements
		{
			ID:          "first_share",
			Name:        "Community Member",
			Description: "Share your first build",
			Category:    CategoryCommunity,
			Points:      15,
			MaxProgress: 1,
		},
		{
			ID:          "template_downloaded",
			Name:        "Template User",
			Description: "Download and use a community template",
			Category:    CategoryCommunity,
			Points:      10,
			MaxProgress: 1,
		},
		{
			ID:          "preset_uploader",
			Name:        "Preset Creator",
			Description: "Upload 3 presets to the marketplace",
			Category:    CategoryCommunity,
			Points:      30,
			MaxProgress: 3,
		},
		{
			ID:          "popular_creator",
			Name:        "Popular Creator",
			Description: "Have your preset downloaded 100 times",
			Category:    CategoryCommunity,
			Points:      50,
			Hidden:      true,
			MaxProgress: 100,
		},
		{
			ID:          "scripter",
			Name:        "Script Master",
			Description: "Execute 5 custom Lua scripts",
			Category:    CategoryEditor,
			Points:      35,
			MaxProgress: 5,
		},
		{
			ID:          "plugin_user",
			Name:        "Plugin Enthusiast",
			Description: "Install and use 3 plugins",
			Category:    CategoryEditor,
			Points:      25,
			MaxProgress: 3,
		},
	}

	for _, achievement := range achievements {
		t.achievements[achievement.ID] = achievement
	}
}

// IncrementProgress increments progress for an achievement
func (t *Tracker) IncrementProgress(achievementID string, amount int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	achievement, exists := t.achievements[achievementID]
	if !exists || achievement.Unlocked {
		return false
	}

	achievement.Progress += amount
	if achievement.Progress >= achievement.MaxProgress {
		achievement.Progress = achievement.MaxProgress
		achievement.Unlocked = true
		achievement.UnlockedAt = time.Now()
		
		// Call unlock callback
		if t.onUnlock != nil {
			go t.onUnlock(achievement)
		}
		
		return true // Achievement unlocked
	}

	return false // Progress updated but not unlocked
}

// SetProgress sets progress for an achievement
func (t *Tracker) SetProgress(achievementID string, progress int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	achievement, exists := t.achievements[achievementID]
	if !exists || achievement.Unlocked {
		return false
	}

	achievement.Progress = progress
	if achievement.Progress >= achievement.MaxProgress {
		achievement.Progress = achievement.MaxProgress
		achievement.Unlocked = true
		achievement.UnlockedAt = time.Now()
		
		if t.onUnlock != nil {
			go t.onUnlock(achievement)
		}
		
		return true
	}

	return false
}

// GetAchievement retrieves an achievement
func (t *Tracker) GetAchievement(id string) *Achievement {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.achievements[id]
}

// GetAllAchievements returns all achievements
func (t *Tracker) GetAllAchievements() []*Achievement {
	t.mu.RLock()
	defer t.mu.RUnlock()

	achievements := make([]*Achievement, 0, len(t.achievements))
	for _, achievement := range t.achievements {
		achievements = append(achievements, achievement)
	}
	return achievements
}

// GetUnlockedAchievements returns all unlocked achievements
func (t *Tracker) GetUnlockedAchievements() []*Achievement {
	t.mu.RLock()
	defer t.mu.RUnlock()

	achievements := make([]*Achievement, 0)
	for _, achievement := range t.achievements {
		if achievement.Unlocked {
			achievements = append(achievements, achievement)
		}
	}
	return achievements
}

// GetTotalPoints calculates total points earned
func (t *Tracker) GetTotalPoints() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	total := 0
	for _, achievement := range t.achievements {
		if achievement.Unlocked {
			total += achievement.Points
		}
	}
	return total
}

// GetCompletionPercentage calculates completion percentage
func (t *Tracker) GetCompletionPercentage() float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	total := len(t.achievements)
	unlocked := 0
	for _, achievement := range t.achievements {
		if achievement.Unlocked {
			unlocked++
		}
	}
	
	if total == 0 {
		return 0
	}
	
	return (float64(unlocked) / float64(total)) * 100.0
}

// SetUnlockCallback sets the callback for when achievements are unlocked
func (t *Tracker) SetUnlockCallback(callback func(*Achievement)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.onUnlock = callback
}
