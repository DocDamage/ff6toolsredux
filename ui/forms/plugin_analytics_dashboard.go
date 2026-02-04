package forms

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"sort"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/plugins"
)

// AnalyticsStat represents a single plugin's analytics statistics
type AnalyticsStat struct {
	PluginID       string
	LoadCount      int
	UnloadCount    int
	ExecutionCount int
	ErrorCount     int
	SuccessRate    float64
	UptimePercent  float64
	LastLoadTime   time.Time
	LastExecTime   time.Time
	Trend          string // "increasing", "decreasing", "stable"
	TrendPercent   float64
}

// PluginAnalyticsDashboard displays usage patterns and reliability metrics
type PluginAnalyticsDashboard struct {
	analytics   *plugins.AnalyticsEngine
	manager     *plugins.Manager
	window      fyne.Window
	refreshRate time.Duration
	stats       map[string]*AnalyticsStat
	mu          sync.RWMutex
	stopChan    chan struct{}
	running     bool
}

// NewPluginAnalyticsDashboard creates a new analytics dashboard
func NewPluginAnalyticsDashboard(manager *plugins.Manager, window fyne.Window) *PluginAnalyticsDashboard {
	return &PluginAnalyticsDashboard{
		analytics:   manager.GetAnalytics(),
		manager:     manager,
		window:      window,
		refreshRate: 2 * time.Second,
		stats:       make(map[string]*AnalyticsStat),
		stopChan:    make(chan struct{}),
		running:     false,
	}
}

// Show displays the analytics dashboard
func (d *PluginAnalyticsDashboard) Show() fyne.CanvasObject {
	content := container.NewVBox(
		d.createHeaderSection(),
		widget.NewSeparator(),
		d.createUsageStatsSection(),
		widget.NewSeparator(),
		d.createMostUsedPluginsSection(),
		widget.NewSeparator(),
		d.createReliabilityRankingSection(),
		widget.NewSeparator(),
		d.createTrendAnalysisSection(),
		widget.NewSeparator(),
		d.createExportSection(),
	)

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(1000, 600))

	go d.autoRefresh()
	d.running = true

	return scroll
}

// createHeaderSection creates dashboard title
func (d *PluginAnalyticsDashboard) createHeaderSection() *fyne.Container {
	title := canvas.NewText("Plugin Analytics Dashboard", color.White)
	title.TextSize = 24
	title.TextStyle.Bold = true

	subtitle := canvas.NewText(fmt.Sprintf("Usage patterns & reliability | Refresh: %v", d.refreshRate), color.Gray{Y: 128})
	subtitle.TextSize = 12

	return container.NewVBox(title, subtitle)
}

// createUsageStatsSection displays overall usage statistics
func (d *PluginAnalyticsDashboard) createUsageStatsSection() *fyne.Container {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Calculate totals
	totalLoads := 0
	totalUnloads := 0
	totalExecutions := 0
	totalErrors := 0
	activePlugins := 0

	for _, stat := range d.stats {
		totalLoads += stat.LoadCount
		totalUnloads += stat.UnloadCount
		totalExecutions += stat.ExecutionCount
		totalErrors += stat.ErrorCount
		if stat.LoadCount > stat.UnloadCount {
			activePlugins++
		}
	}

	totalSuccesses := totalExecutions - totalErrors
	overallSuccess := 0.0
	if totalExecutions > 0 {
		overallSuccess = (float64(totalSuccesses) / float64(totalExecutions)) * 100
	}

	stats := container.NewVBox(
		canvas.NewText("Overall Usage Statistics", color.White),
		canvas.NewText(
			fmt.Sprintf("Total Plugins: %d | Active: %d | Total Loads: %d | Total Unloads: %d",
				len(d.stats), activePlugins, totalLoads, totalUnloads),
			color.Gray{Y: 200}),
		canvas.NewText(
			fmt.Sprintf("Total Executions: %d | Successes: %d | Errors: %d | Overall Success Rate: %.1f%%",
				totalExecutions, totalSuccesses, totalErrors, overallSuccess),
			d.getSuccessColor(overallSuccess)),
	)

	return stats
}

// getSuccessColor returns color based on success rate
func (d *PluginAnalyticsDashboard) getSuccessColor(rate float64) color.Color {
	if rate >= 95 {
		return color.NRGBA{G: 200, A: 255} // Green
	} else if rate >= 80 {
		return color.NRGBA{R: 255, G: 165, A: 255} // Orange
	}
	return color.NRGBA{R: 255, A: 255} // Red
}

// createMostUsedPluginsSection displays top plugins by execution count
func (d *PluginAnalyticsDashboard) createMostUsedPluginsSection() *fyne.Container {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Sort by execution count
	type pluginRank struct {
		id   string
		stat *AnalyticsStat
		rank int
	}

	var ranked []pluginRank
	for _, stat := range d.stats {
		ranked = append(ranked, pluginRank{id: stat.PluginID, stat: stat})
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].stat.ExecutionCount > ranked[j].stat.ExecutionCount
	})

	// Keep top 10
	if len(ranked) > 10 {
		ranked = ranked[:10]
	}

	// Add ranks
	for i := range ranked {
		ranked[i].rank = i + 1
	}

	container := container.NewVBox(
		canvas.NewText("Most Used Plugins (Top 10)", color.White),
	)

	for _, r := range ranked {
		stat := r.stat
		barLength := stat.ExecutionCount / 10 // Simple bar visualization
		if barLength > 50 {
			barLength = 50
		}

		bar := ""
		for i := 0; i < barLength; i++ {
			bar += "█"
		}

		text := canvas.NewText(
			fmt.Sprintf("[%d] %s: %d exec | %d loads | Success: %.1f%% %s",
				r.rank, stat.PluginID, stat.ExecutionCount, stat.LoadCount,
				stat.SuccessRate, bar),
			color.White)
		text.TextSize = 11
		container.Add(text)
	}

	return container
}

// createReliabilityRankingSection displays plugins by reliability
func (d *PluginAnalyticsDashboard) createReliabilityRankingSection() *fyne.Container {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Sort by success rate
	type reliabilityRank struct {
		id   string
		stat *AnalyticsStat
		rank int
	}

	var ranked []reliabilityRank
	for _, stat := range d.stats {
		ranked = append(ranked, reliabilityRank{id: stat.PluginID, stat: stat})
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].stat.SuccessRate > ranked[j].stat.SuccessRate
	})

	// Keep top 10
	if len(ranked) > 10 {
		ranked = ranked[:10]
	}

	// Add ranks
	for i := range ranked {
		ranked[i].rank = i + 1
	}

	container := container.NewVBox(
		canvas.NewText("Most Reliable Plugins (Top 10)", color.White),
	)

	for _, r := range ranked {
		stat := r.stat
		successColor := d.getSuccessColor(stat.SuccessRate)

		text := canvas.NewText(
			fmt.Sprintf("[%d] %s: %.1f%% success | %d errors out of %d executions",
				r.rank, stat.PluginID, stat.SuccessRate, stat.ErrorCount, stat.ExecutionCount),
			successColor)
		text.TextSize = 11
		container.Add(text)
	}

	return container
}

// createTrendAnalysisSection displays trend data
func (d *PluginAnalyticsDashboard) createTrendAnalysisSection() *fyne.Container {
	d.mu.RLock()
	defer d.mu.RUnlock()

	container := container.NewVBox(
		canvas.NewText("Plugin Trends (Improving/Degrading)", color.White),
	)

	// Collect trends
	type trendData struct {
		id        string
		stat      *AnalyticsStat
		indicator string
	}

	var trends []trendData
	for _, stat := range d.stats {
		indicator := "→" // stable

		if stat.Trend == "increasing" {
			indicator = "↑" // improving
		} else if stat.Trend == "decreasing" {
			indicator = "↓" // degrading
		}

		trends = append(trends, trendData{stat.PluginID, stat, indicator})
	}

	// Sort by trend type
	sort.Slice(trends, func(i, j int) bool {
		if trends[i].stat.Trend != trends[j].stat.Trend {
			// Put improving first, then stable, then degrading
			trendOrder := map[string]int{"increasing": 0, "stable": 1, "decreasing": 2}
			return trendOrder[trends[i].stat.Trend] < trendOrder[trends[j].stat.Trend]
		}
		return trends[i].id < trends[j].id
	})

	// Limit display to 15 most relevant
	if len(trends) > 15 {
		trends = trends[:15]
	}

	for _, t := range trends {
		stat := t.stat
		color := d.getTrendColor(stat.Trend)

		text := canvas.NewText(
			fmt.Sprintf("%s %s: %s (%.1f%% change)",
				t.indicator, stat.PluginID, stat.Trend, stat.TrendPercent),
			color)
		text.TextSize = 11
		container.Add(text)
	}

	return container
}

// getTrendColor returns color for trend
func (d *PluginAnalyticsDashboard) getTrendColor(trend string) color.Color {
	switch trend {
	case "increasing":
		return color.NRGBA{G: 200, A: 255} // Green
	case "decreasing":
		return color.NRGBA{R: 255, A: 255} // Red
	default:
		return color.Gray{Y: 200} // Gray for stable
	}
}

// createExportSection creates export controls
func (d *PluginAnalyticsDashboard) createExportSection() *fyne.Container {
	exportContainer := container.NewHBox(
		canvas.NewText("Export Analytics:", color.White),
		widget.NewButton("JSON", func() {
			d.exportAsJSON()
		}),
		widget.NewButton("Refresh Now", func() {
			d.RefreshAnalytics()
		}),
	)

	return exportContainer
}

// RefreshAnalytics updates analytics from engine
func (d *PluginAnalyticsDashboard) RefreshAnalytics() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.analytics == nil {
		return
	}

	// Clear old stats
	d.stats = make(map[string]*AnalyticsStat)

	// Get all plugin stats
	allStats := d.analytics.GetAllPluginStats()

	for pluginID, engineStat := range allStats {
		stat := &AnalyticsStat{
			PluginID:       pluginID,
			LoadCount:      engineStat.LoadCount,
			UnloadCount:    engineStat.UnloadCount,
			ExecutionCount: engineStat.ExecutionCount,
			ErrorCount:     engineStat.ErrorCount,
			SuccessRate:    engineStat.SuccessRate * 100,
			UptimePercent:  engineStat.UptimePercent * 100,
			LastLoadTime:   time.Now(),
			LastExecTime:   time.Now(),
		}

		// Get trend analysis
		trends := d.analytics.AnalyzeTrends(pluginID, 10) // Use 10 as sample size
		if trends != nil && len(trends) > 0 {
			trend := trends[0]
			stat.Trend = trend.Trend
			stat.TrendPercent = trend.ChangePercent
		} else {
			stat.Trend = "stable"
			stat.TrendPercent = 0
		}

		d.stats[pluginID] = stat
	}
}

// exportAsJSON exports analytics to JSON file
func (d *PluginAnalyticsDashboard) exportAsJSON() {
	d.mu.RLock()
	defer d.mu.RUnlock()

	filename := fmt.Sprintf("analytics_%d.json", time.Now().Unix())

	type exportData struct {
		Timestamp time.Time
		Stats     map[string]*AnalyticsStat
	}

	data := exportData{
		Timestamp: time.Now(),
		Stats:     d.stats,
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	err = os.WriteFile(filename, bytes, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON file: %v\n", err)
		return
	}

	fmt.Printf("✓ Exported analytics to %s\n", filename)
}

// autoRefresh periodically updates analytics
func (d *PluginAnalyticsDashboard) autoRefresh() {
	ticker := time.NewTicker(d.refreshRate)
	defer ticker.Stop()

	for {
		select {
		case <-d.stopChan:
			d.running = false
			return
		case <-ticker.C:
			d.RefreshAnalytics()
		}
	}
}

// SetRefreshRate updates refresh interval
func (d *PluginAnalyticsDashboard) SetRefreshRate(rate time.Duration) {
	d.refreshRate = rate
}

// Stop halts the refresh loop
func (d *PluginAnalyticsDashboard) Stop() {
	if d.running {
		close(d.stopChan)
	}
}

// GetStats returns current cached stats
func (d *PluginAnalyticsDashboard) GetStats() map[string]*AnalyticsStat {
	d.mu.RLock()
	defer d.mu.RUnlock()

	statsCopy := make(map[string]*AnalyticsStat)
	for k, v := range d.stats {
		statsCopy[k] = v
	}
	return statsCopy
}

// GetPluginStat returns stats for a specific plugin
func (d *PluginAnalyticsDashboard) GetPluginStat(pluginID string) *AnalyticsStat {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.stats[pluginID]
}

// GetTopPlugins returns N most-used plugins
func (d *PluginAnalyticsDashboard) GetTopPlugins(count int) []*AnalyticsStat {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var stats []*AnalyticsStat
	for _, s := range d.stats {
		stats = append(stats, s)
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].ExecutionCount > stats[j].ExecutionCount
	})

	if len(stats) > count {
		stats = stats[:count]
	}

	return stats
}

// GetMostReliablePlugins returns N most-reliable plugins
func (d *PluginAnalyticsDashboard) GetMostReliablePlugins(count int) []*AnalyticsStat {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var stats []*AnalyticsStat
	for _, s := range d.stats {
		stats = append(stats, s)
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].SuccessRate > stats[j].SuccessRate
	})

	if len(stats) > count {
		stats = stats[:count]
	}

	return stats
}
