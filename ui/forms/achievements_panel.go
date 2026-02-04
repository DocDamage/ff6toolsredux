package forms

import (
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/achievements"
)

// AchievementsPanel displays achievement progress
type AchievementsPanel struct {
	window  fyne.Window
	tracker *achievements.Tracker
}

// NewAchievementsPanel creates a new achievements panel
func NewAchievementsPanel(window fyne.Window, tracker *achievements.Tracker) *AchievementsPanel {
	return &AchievementsPanel{
		window:  window,
		tracker: tracker,
	}
}

// Show displays the achievements panel dialog
func (a *AchievementsPanel) Show() {
	// Get all achievements
	allAchievements := a.tracker.GetAllAchievements()

	// Sort by category and then by unlocked status
	sort.Slice(allAchievements, func(i, j int) bool {
		if allAchievements[i].Category != allAchievements[j].Category {
			return allAchievements[i].Category < allAchievements[j].Category
		}
		return allAchievements[i].Unlocked && !allAchievements[j].Unlocked
	})

	// Statistics
	totalPoints := a.tracker.GetTotalPoints()
	completion := a.tracker.GetCompletionPercentage()
	unlockedCount := len(a.tracker.GetUnlockedAchievements())
	totalCount := len(allAchievements)

	statsLabel := widget.NewLabel(fmt.Sprintf(
		"Unlocked: %d/%d (%.1f%%) | Points: %d",
		unlockedCount, totalCount, completion, totalPoints,
	))

	// Filter by category
	categorySelect := widget.NewSelect(
		[]string{"All", "Editor", "Game", "Community"},
		func(value string) {
			// Filter logic would refresh the list
		},
	)
	categorySelect.SetSelected("All")

	// Show only unlocked checkbox
	unlockedOnlyCheck := widget.NewCheck("Show only unlocked", func(checked bool) {
		// Filter logic would refresh the list
	})

	// Achievements list
	achievementList := widget.NewList(
		func() int { return len(allAchievements) },
		func() fyne.CanvasObject {
			return container.NewVBox(
				container.NewHBox(
					widget.NewLabel("🏆"),
					container.NewVBox(
						widget.NewLabel("Achievement Name"),
						widget.NewLabel("Description"),
						widget.NewProgressBar(),
					),
					widget.NewLabel("100 pts"),
				),
				widget.NewSeparator(),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= len(allAchievements) {
				return
			}

			achievement := allAchievements[id]
			box := item.(*fyne.Container).Objects[0].(*fyne.Container)

			// Icon
			icon := box.Objects[0].(*widget.Label)
			if achievement.Unlocked {
				icon.SetText("🏆")
			} else if achievement.Hidden && achievement.Progress == 0 {
				icon.SetText("❓")
			} else {
				icon.SetText("🔒")
			}

			// Details
			details := box.Objects[1].(*fyne.Container)
			nameLabel := details.Objects[0].(*widget.Label)
			descLabel := details.Objects[1].(*widget.Label)
			progressBar := details.Objects[2].(*widget.ProgressBar)

			// Set name (hide if hidden and not unlocked)
			if achievement.Hidden && !achievement.Unlocked {
				nameLabel.SetText("Hidden Achievement")
				descLabel.SetText("???")
			} else {
				nameLabel.SetText(achievement.Name)
				descLabel.SetText(achievement.Description)
			}

			// Progress bar
			if achievement.MaxProgress > 1 {
				progressBar.Show()
				progressBar.Min = 0
				progressBar.Max = float64(achievement.MaxProgress)
				progressBar.Value = float64(achievement.Progress)
				progressBar.Refresh()
			} else {
				progressBar.Hide()
			}

			// Points
			pointsLabel := box.Objects[2].(*widget.Label)
			if achievement.Unlocked {
				unlockedTime := achievement.UnlockedAt.Format("Jan 2, 2006")
				pointsLabel.SetText(fmt.Sprintf("✓ %d pts\n%s", achievement.Points, unlockedTime))
			} else {
				pointsLabel.SetText(fmt.Sprintf("%d pts", achievement.Points))
			}
		},
	)

	achievementList.OnSelected = func(id widget.ListItemID) {
		if id < len(allAchievements) {
			a.showAchievementDetails(allAchievements[id])
		}
		achievementList.UnselectAll()
	}

	// Build layout
	filters := container.NewHBox(
		widget.NewLabel("Category:"),
		categorySelect,
		unlockedOnlyCheck,
	)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Achievements"),
			widget.NewSeparator(),
			statsLabel,
			filters,
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		achievementList,
	)

	d := dialog.NewCustom("Achievements", "Close", content, a.window)
	d.Resize(fyne.NewSize(600, 500))
	d.Show()
}

// showAchievementDetails shows details for a specific achievement
func (a *AchievementsPanel) showAchievementDetails(achievement *achievements.Achievement) {
	var status string
	if achievement.Unlocked {
		status = fmt.Sprintf("✓ Unlocked on %s", achievement.UnlockedAt.Format("January 2, 2006 at 3:04 PM"))
	} else {
		status = "🔒 Locked"
		if achievement.MaxProgress > 1 {
			status += fmt.Sprintf("\nProgress: %d/%d (%.1f%%)",
				achievement.Progress, achievement.MaxProgress,
				float64(achievement.Progress)/float64(achievement.MaxProgress)*100.0)
		}
	}

	details := fmt.Sprintf(`
%s

%s

Category: %s
Points: %d

Status:
%s
`,
		achievement.Name,
		achievement.Description,
		achievement.Category,
		achievement.Points,
		status,
	)

	detailsEntry := widget.NewMultiLineEntry()
	detailsEntry.SetText(details)
	detailsEntry.Wrapping = fyne.TextWrapWord
	detailsEntry.Disable()

	d := dialog.NewCustom("Achievement Details", "Close", detailsEntry, a.window)
	d.Resize(fyne.NewSize(400, 300))
	d.Show()
}

// ShowUnlockNotification shows a notification when an achievement is unlocked
func ShowUnlockNotification(window fyne.Window, achievement *achievements.Achievement) {
	message := fmt.Sprintf(
		"🏆 Achievement Unlocked!\n\n%s\n\n%s\n\n+%d points",
		achievement.Name,
		achievement.Description,
		achievement.Points,
	)

	dialog.ShowInformation("Achievement Unlocked!", message, window)
}

// InitializeAchievementTracking sets up achievement tracking callbacks
func InitializeAchievementTracking(window fyne.Window, tracker *achievements.Tracker) {
	tracker.SetUnlockCallback(func(achievement *achievements.Achievement) {
		// Show notification on main thread
		ShowUnlockNotification(window, achievement)
	})
}
