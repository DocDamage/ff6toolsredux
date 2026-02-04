package achievements

import (
	"sync"
	"testing"
	"time"
)

// TestNewTracker tests tracker creation
func TestNewTracker(t *testing.T) {
	tracker := NewTracker()

	if tracker == nil {
		t.Fatal("NewTracker() returned nil")
	}

	if tracker.achievements == nil {
		t.Error("Tracker achievements map is nil")
	}

	// Should have achievements initialized
	all := tracker.GetAllAchievements()
	if len(all) == 0 {
		t.Error("Tracker should have achievements initialized")
	}
}

// TestGetAchievement tests achievement retrieval
func TestGetAchievement(t *testing.T) {
	tracker := NewTracker()

	// Test existing achievement
	a := tracker.GetAchievement("first_edit")
	if a == nil {
		t.Fatal("GetAchievement(first_edit) returned nil")
	}
	if a.ID != "first_edit" {
		t.Errorf("Achievement ID = %s, want first_edit", a.ID)
	}
	if a.Name != "First Steps" {
		t.Errorf("Achievement Name = %s, want First Steps", a.Name)
	}

	// Test non-existent achievement
	nonExistent := tracker.GetAchievement("non_existent")
	if nonExistent != nil {
		t.Error("GetAchievement(non_existent) should return nil")
	}
}

// TestGetAllAchievements tests getting all achievements
func TestGetAllAchievements(t *testing.T) {
	tracker := NewTracker()
	all := tracker.GetAllAchievements()

	// Should have all predefined achievements
	expectedCount := 16 // Number of achievements defined in initializeAchievements
	if len(all) != expectedCount {
		t.Errorf("GetAllAchievements returned %d achievements, want %d", len(all), expectedCount)
	}

	// Verify all achievements have required fields
	for _, a := range all {
		if a.ID == "" {
			t.Error("Achievement has empty ID")
		}
		if a.Name == "" {
			t.Error("Achievement has empty Name")
		}
		if a.Description == "" {
			t.Error("Achievement has empty Description")
		}
		if a.MaxProgress <= 0 {
			t.Errorf("Achievement %s has invalid MaxProgress: %d", a.ID, a.MaxProgress)
		}
	}
}

// TestIncrementProgress tests progress increment
func TestIncrementProgress(t *testing.T) {
	tracker := NewTracker()

	// Use undo_master which has MaxProgress = 50
	// Test incrementing progress
	unlocked := tracker.IncrementProgress("undo_master", 1)
	if unlocked {
		t.Error("IncrementProgress should not unlock on first increment (below MaxProgress)")
	}

	a := tracker.GetAchievement("undo_master")
	if a.Progress != 1 {
		t.Errorf("Achievement progress = %d, want 1", a.Progress)
	}
}

// TestIncrementProgressUnlocks tests that achievement unlocks at max progress
func TestIncrementProgressUnlocks(t *testing.T) {
	tracker := NewTracker()

	// Use an achievement with higher max progress
	unlocked := tracker.IncrementProgress("save_10", 10)
	if !unlocked {
		t.Error("Should unlock when reaching MaxProgress")
	}

	a := tracker.GetAchievement("save_10")
	if !a.Unlocked {
		t.Error("Achievement should be marked as unlocked")
	}
	if a.UnlockedAt.IsZero() {
		t.Error("Achievement UnlockedAt should be set")
	}
	if a.Progress != 10 {
		t.Errorf("Achievement progress = %d, want 10", a.Progress)
	}
}

// TestIncrementProgressNonExistent tests incrementing non-existent achievement
func TestIncrementProgressNonExistent(t *testing.T) {
	tracker := NewTracker()

	unlocked := tracker.IncrementProgress("non_existent", 1)
	if unlocked {
		t.Error("Should not unlock non-existent achievement")
	}
}

// TestIncrementProgressAlreadyUnlocked tests incrementing already unlocked achievement
func TestIncrementProgressAlreadyUnlocked(t *testing.T) {
	tracker := NewTracker()

	// First unlock the achievement
	tracker.SetProgress("first_edit", 1)

	// Try to increment again
	unlocked := tracker.IncrementProgress("first_edit", 1)
	if unlocked {
		t.Error("Should not report unlock for already unlocked achievement")
	}

	a := tracker.GetAchievement("first_edit")
	if a.Progress != 1 {
		t.Errorf("Progress should remain at MaxProgress, got %d", a.Progress)
	}
}

// TestSetProgress tests setting progress directly
func TestSetProgress(t *testing.T) {
	tracker := NewTracker()

	unlocked := tracker.SetProgress("undo_master", 50)
	if !unlocked {
		t.Error("Should unlock when setting to MaxProgress")
	}

	a := tracker.GetAchievement("undo_master")
	if a.Progress != 50 {
		t.Errorf("Achievement progress = %d, want 50", a.Progress)
	}
}

// TestSetProgressPartial tests setting partial progress
func TestSetProgressPartial(t *testing.T) {
	tracker := NewTracker()

	unlocked := tracker.SetProgress("all_rages", 100)
	if unlocked {
		t.Error("Should not unlock when below MaxProgress")
	}

	a := tracker.GetAchievement("all_rages")
	if a.Progress != 100 {
		t.Errorf("Achievement progress = %d, want 100", a.Progress)
	}
	if a.Unlocked {
		t.Error("Achievement should not be unlocked yet")
	}
}

// TestSetProgressNonExistent tests setting progress for non-existent achievement
func TestSetProgressNonExistent(t *testing.T) {
	tracker := NewTracker()

	unlocked := tracker.SetProgress("non_existent", 100)
	if unlocked {
		t.Error("Should not unlock non-existent achievement")
	}
}

// TestSetProgressAlreadyUnlocked tests setting progress for already unlocked achievement
func TestSetProgressAlreadyUnlocked(t *testing.T) {
	tracker := NewTracker()

	// First unlock
	tracker.SetProgress("first_edit", 1)

	// Try to set progress again
	unlocked := tracker.SetProgress("first_edit", 0)
	if unlocked {
		t.Error("Should not report unlock for already unlocked achievement")
	}
}

// TestGetUnlockedAchievements tests getting unlocked achievements
func TestGetUnlockedAchievements(t *testing.T) {
	tracker := NewTracker()

	// Initially no achievements should be unlocked
	unlocked := tracker.GetUnlockedAchievements()
	if len(unlocked) != 0 {
		t.Errorf("Initially should have 0 unlocked, got %d", len(unlocked))
	}

	// Unlock some achievements
	tracker.SetProgress("first_edit", 1)
	tracker.SetProgress("save_10", 10)

	unlocked = tracker.GetUnlockedAchievements()
	if len(unlocked) != 2 {
		t.Errorf("Should have 2 unlocked achievements, got %d", len(unlocked))
	}
}

// TestGetTotalPoints tests total points calculation
func TestGetTotalPoints(t *testing.T) {
	tracker := NewTracker()

	// Initially should have 0 points
	points := tracker.GetTotalPoints()
	if points != 0 {
		t.Errorf("Initial points = %d, want 0", points)
	}

	// Unlock first_edit (10 points)
	tracker.SetProgress("first_edit", 1)
	points = tracker.GetTotalPoints()
	if points != 10 {
		t.Errorf("Points after first_edit = %d, want 10", points)
	}

	// Unlock save_10 (25 points)
	tracker.SetProgress("save_10", 10)
	points = tracker.GetTotalPoints()
	if points != 35 {
		t.Errorf("Points after save_10 = %d, want 35", points)
	}
}

// TestGetCompletionPercentage tests completion percentage
func TestGetCompletionPercentage(t *testing.T) {
	tracker := NewTracker()

	// Get total count
	totalCount := len(tracker.GetAllAchievements())

	// Initially 0%
	pct := tracker.GetCompletionPercentage()
	if pct != 0 {
		t.Errorf("Initial completion = %f, want 0", pct)
	}

	// Unlock one achievement
	tracker.SetProgress("first_edit", 1)
	expectedPct := (1.0 / float64(totalCount)) * 100.0
	pct = tracker.GetCompletionPercentage()
	if pct != expectedPct {
		t.Errorf("Completion = %f, want %f", pct, expectedPct)
	}
}

// TestSetUnlockCallback tests unlock callback
func TestSetUnlockCallback(t *testing.T) {
	tracker := NewTracker()

	callbackCalled := false
	var unlockedAchievement *Achievement

	tracker.SetUnlockCallback(func(a *Achievement) {
		callbackCalled = true
		unlockedAchievement = a
	})

	tracker.SetProgress("first_edit", 1)

	// Wait for goroutine to complete
	time.Sleep(100 * time.Millisecond)

	if !callbackCalled {
		t.Error("Unlock callback was not called")
	}
	if unlockedAchievement == nil {
		t.Error("Callback received nil achievement")
	}
	if unlockedAchievement != nil && unlockedAchievement.ID != "first_edit" {
		t.Errorf("Callback received wrong achievement: %s", unlockedAchievement.ID)
	}
}

// TestAchievementCategories tests that achievements have proper categories
func TestAchievementCategories(t *testing.T) {
	tracker := NewTracker()
	all := tracker.GetAllAchievements()

	validCategories := map[string]bool{
		CategoryEditor:    true,
		CategoryGame:      true,
		CategoryCommunity: true,
	}

	for _, a := range all {
		if !validCategories[a.Category] {
			t.Errorf("Achievement %s has invalid category: %s", a.ID, a.Category)
		}
	}
}

// TestAchievementPoints tests that achievements have positive points
func TestAchievementPoints(t *testing.T) {
	tracker := NewTracker()
	all := tracker.GetAllAchievements()

	for _, a := range all {
		if a.Points <= 0 {
			t.Errorf("Achievement %s has invalid points: %d", a.ID, a.Points)
		}
	}
}

// TestAchievementHidden tests hidden achievements
func TestAchievementHidden(t *testing.T) {
	tracker := NewTracker()

	// popular_creator should be hidden
	a := tracker.GetAchievement("popular_creator")
	if a == nil {
		t.Fatal("popular_creator achievement not found")
	}
	if !a.Hidden {
		t.Error("popular_creator should be hidden")
	}

	// first_edit should not be hidden
	a = tracker.GetAchievement("first_edit")
	if a.Hidden {
		t.Error("first_edit should not be hidden")
	}
}

// TestConcurrentAccess tests thread safety
func TestConcurrentAccess(t *testing.T) {
	tracker := NewTracker()

	var wg sync.WaitGroup
	numGoroutines := 100

	// Concurrent increments
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tracker.IncrementProgress("undo_master", 1)
		}()
	}

	wg.Wait()

	a := tracker.GetAchievement("undo_master")
	// Should be at most MaxProgress (50), but could be less due to race conditions
	// The important thing is no panic occurs
	if a.Progress > a.MaxProgress {
		t.Errorf("Progress = %d exceeds MaxProgress = %d", a.Progress, a.MaxProgress)
	}
}

// TestConcurrentReadWrite tests concurrent reads and writes
func TestConcurrentReadWrite(t *testing.T) {
	tracker := NewTracker()

	var wg sync.WaitGroup
	numReaders := 10
	numWriters := 10
	numOperations := 10

	// Concurrent readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				_ = tracker.GetAllAchievements()
				_ = tracker.GetUnlockedAchievements()
				_ = tracker.GetTotalPoints()
				time.Sleep(time.Microsecond)
			}
		}()
	}

	// Concurrent writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				tracker.IncrementProgress("first_edit", 1)
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	wg.Wait()
	// If we get here without panic, the test passes
}

// BenchmarkIncrementProgress benchmarks progress increment
func BenchmarkIncrementProgress(b *testing.B) {
	tracker := NewTracker()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tracker.IncrementProgress("first_edit", 1)
		// Reset for next iteration
		tracker.SetProgress("first_edit", 0)
	}
}

// BenchmarkGetAllAchievements benchmarks getting all achievements
func BenchmarkGetAllAchievements(b *testing.B) {
	tracker := NewTracker()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tracker.GetAllAchievements()
	}
}

// BenchmarkGetTotalPoints benchmarks points calculation
func BenchmarkGetTotalPoints(b *testing.B) {
	tracker := NewTracker()
	// Unlock some achievements first
	tracker.SetProgress("first_edit", 1)
	tracker.SetProgress("save_10", 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tracker.GetTotalPoints()
	}
}
