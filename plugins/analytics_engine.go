package plugins

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// AnalyticsEvent represents a single analytics event
type AnalyticsEvent struct {
	EventType string                 // "plugin_load", "plugin_execute", "plugin_error", "plugin_unload"
	PluginID  string                 // Plugin identifier
	Timestamp time.Time              // When event occurred
	Duration  time.Duration          // Duration for execution events
	Metadata  map[string]interface{} // Additional data
	Status    string                 // "success", "error", "timeout"
	Error     string                 // Error message if applicable
}

// AnalyticsWindow represents a time window for analytics
type AnalyticsWindow struct {
	StartTime time.Time
	EndTime   time.Time
	Events    []*AnalyticsEvent
}

// PluginStats represents analytics statistics for a plugin
type PluginStats struct {
	PluginID           string
	LoadCount          int
	UnloadCount        int
	ExecutionCount     int
	ErrorCount         int
	SuccessRate        float64
	AverageLoadTime    time.Duration
	AverageExecTime    time.Duration
	TotalExecutionTime time.Duration
	PeakConcurrency    int
	LastActivityTime   time.Time
	FirstSeenTime      time.Time
	UptimePercent      float64
}

// AnalyticsEngine collects and analyzes plugin usage analytics
type AnalyticsEngine struct {
	events          []*AnalyticsEvent            // All events in order
	eventsByPlugin  map[string][]*AnalyticsEvent // Events per plugin
	eventsByType    map[string][]*AnalyticsEvent // Events by type
	pluginStats     map[string]*PluginStats      // Per-plugin statistics
	windows         []*AnalyticsWindow           // Time window snapshots
	maxEventsStored int                          // Maximum events to keep
	mu              sync.RWMutex
}

// TrendAnalysis represents trend analysis for metrics
type TrendAnalysis struct {
	MetricName      string
	Trend           string // "increasing", "decreasing", "stable"
	ChangePercent   float64
	RecentAverage   float64
	PreviousAverage float64
	SampleSize      int
}

// NewAnalyticsEngine creates a new analytics engine
func NewAnalyticsEngine(maxEventsStored int) *AnalyticsEngine {
	return &AnalyticsEngine{
		events:          make([]*AnalyticsEvent, 0),
		eventsByPlugin:  make(map[string][]*AnalyticsEvent),
		eventsByType:    make(map[string][]*AnalyticsEvent),
		pluginStats:     make(map[string]*PluginStats),
		windows:         make([]*AnalyticsWindow, 0),
		maxEventsStored: maxEventsStored,
	}
}

// RecordEvent records an analytics event
func (a *AnalyticsEngine) RecordEvent(event *AnalyticsEvent) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if event == nil {
		return
	}

	event.Timestamp = time.Now()

	// Add to main events list
	a.events = append(a.events, event)

	// Limit events storage
	if len(a.events) > a.maxEventsStored {
		a.events = a.events[1:]
	}

	// Index by plugin
	if _, exists := a.eventsByPlugin[event.PluginID]; !exists {
		a.eventsByPlugin[event.PluginID] = make([]*AnalyticsEvent, 0)
	}
	a.eventsByPlugin[event.PluginID] = append(a.eventsByPlugin[event.PluginID], event)

	// Index by event type
	if _, exists := a.eventsByType[event.EventType]; !exists {
		a.eventsByType[event.EventType] = make([]*AnalyticsEvent, 0)
	}
	a.eventsByType[event.EventType] = append(a.eventsByType[event.EventType], event)

	// Update plugin stats
	a.updatePluginStats(event)
}

// updatePluginStats updates statistics for a plugin
func (a *AnalyticsEngine) updatePluginStats(event *AnalyticsEvent) {
	stats, exists := a.pluginStats[event.PluginID]
	if !exists {
		stats = &PluginStats{
			PluginID:      event.PluginID,
			FirstSeenTime: event.Timestamp,
		}
		a.pluginStats[event.PluginID] = stats
	}

	stats.LastActivityTime = event.Timestamp

	switch event.EventType {
	case "plugin_load":
		stats.LoadCount++
		if event.Duration > 0 {
			currentAvg := stats.AverageLoadTime.Milliseconds() * int64(stats.LoadCount-1)
			stats.AverageLoadTime = time.Duration((currentAvg+event.Duration.Milliseconds())/int64(stats.LoadCount)) * time.Millisecond
		}

	case "plugin_unload":
		stats.UnloadCount++

	case "plugin_execute":
		stats.ExecutionCount++
		stats.TotalExecutionTime += event.Duration
		if event.Status == "error" || event.Status == "timeout" {
			stats.ErrorCount++
		}
		currentAvg := stats.AverageExecTime.Milliseconds() * int64(stats.ExecutionCount-1)
		stats.AverageExecTime = time.Duration((currentAvg+event.Duration.Milliseconds())/int64(stats.ExecutionCount)) * time.Millisecond

	case "plugin_error":
		stats.ErrorCount++
	}

	// Calculate success rate
	if stats.ExecutionCount > 0 {
		successCount := stats.ExecutionCount - stats.ErrorCount
		stats.SuccessRate = float64(successCount) / float64(stats.ExecutionCount) * 100
	}

	// Calculate uptime (days loaded / days tracked)
	if stats.LoadCount > 0 && stats.UnloadCount > 0 {
		totalTime := stats.LastActivityTime.Sub(stats.FirstSeenTime)
		daysTracked := totalTime.Hours() / 24
		if daysTracked > 0 {
			stats.UptimePercent = float64(stats.LoadCount) / daysTracked * 100
			if stats.UptimePercent > 100 {
				stats.UptimePercent = 100
			}
		}
	} else if stats.LoadCount > 0 {
		stats.UptimePercent = 100 // Still loaded
	}
}

// GetEvents returns all events in a time range
func (a *AnalyticsEngine) GetEvents(startTime, endTime time.Time) []*AnalyticsEvent {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var result []*AnalyticsEvent
	for _, event := range a.events {
		if event.Timestamp.After(startTime) && event.Timestamp.Before(endTime) {
			result = append(result, event)
		}
	}
	return result
}

// GetPluginEvents returns events for a specific plugin
func (a *AnalyticsEngine) GetPluginEvents(pluginID string) []*AnalyticsEvent {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if events, exists := a.eventsByPlugin[pluginID]; exists {
		result := make([]*AnalyticsEvent, len(events))
		copy(result, events)
		return result
	}
	return nil
}

// GetEventsByType returns events of a specific type
func (a *AnalyticsEngine) GetEventsByType(eventType string) []*AnalyticsEvent {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if events, exists := a.eventsByType[eventType]; exists {
		result := make([]*AnalyticsEvent, len(events))
		copy(result, events)
		return result
	}
	return nil
}

// GetPluginStats returns statistics for a specific plugin
func (a *AnalyticsEngine) GetPluginStats(pluginID string) *PluginStats {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if stats, exists := a.pluginStats[pluginID]; exists {
		// Return copy
		statsCopy := *stats
		return &statsCopy
	}
	return nil
}

// GetAllPluginStats returns statistics for all plugins
func (a *AnalyticsEngine) GetAllPluginStats() map[string]*PluginStats {
	a.mu.RLock()
	defer a.mu.RUnlock()

	result := make(map[string]*PluginStats)
	for id, stats := range a.pluginStats {
		statsCopy := *stats
		result[id] = &statsCopy
	}
	return result
}

// GetMostUsedPlugins returns plugins ranked by execution count
func (a *AnalyticsEngine) GetMostUsedPlugins(limit int) []*PluginStats {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var stats []*PluginStats
	for _, s := range a.pluginStats {
		stats = append(stats, s)
	}

	// Sort by execution count descending
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].ExecutionCount > stats[j].ExecutionCount
	})

	if limit > 0 && len(stats) > limit {
		stats = stats[:limit]
	}

	return stats
}

// GetMostReliablePlugins returns plugins with highest success rates
func (a *AnalyticsEngine) GetMostReliablePlugins(minExecutions int) []*PluginStats {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var stats []*PluginStats
	for _, s := range a.pluginStats {
		if s.ExecutionCount >= minExecutions {
			stats = append(stats, s)
		}
	}

	// Sort by success rate descending
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].SuccessRate > stats[j].SuccessRate
	})

	return stats
}

// AnalyzeTrends analyzes trends for a specific metric
func (a *AnalyticsEngine) AnalyzeTrends(pluginID string, sampleSize int) []*TrendAnalysis {
	a.mu.RLock()
	defer a.mu.RUnlock()

	events, exists := a.eventsByPlugin[pluginID]
	if !exists || len(events) < sampleSize*2 {
		return nil
	}

	var trends []*TrendAnalysis

	// Split events into two halves
	midpoint := len(events) / 2
	recentEvents := events[midpoint:]
	previousEvents := events[:midpoint]

	if sampleSize > 0 {
		if len(recentEvents) > sampleSize {
			recentEvents = recentEvents[len(recentEvents)-sampleSize:]
		}
		if len(previousEvents) > sampleSize {
			previousEvents = previousEvents[:sampleSize]
		}
	}

	// Analyze execution duration trend
	recentDurationAvg := calculateAverageDuration(recentEvents)
	previousDurationAvg := calculateAverageDuration(previousEvents)

	if recentDurationAvg > 0 && previousDurationAvg > 0 {
		changePercent := ((recentDurationAvg - previousDurationAvg) / previousDurationAvg) * 100
		trend := "stable"
		if changePercent > 5 {
			trend = "increasing"
		} else if changePercent < -5 {
			trend = "decreasing"
		}

		trends = append(trends, &TrendAnalysis{
			MetricName:      "execution_duration",
			Trend:           trend,
			ChangePercent:   changePercent,
			RecentAverage:   recentDurationAvg,
			PreviousAverage: previousDurationAvg,
			SampleSize:      len(recentEvents),
		})
	}

	// Analyze error rate trend
	recentErrorRate := calculateErrorRate(recentEvents)
	previousErrorRate := calculateErrorRate(previousEvents)

	if previousErrorRate >= 0 {
		changePercent := recentErrorRate - previousErrorRate
		trend := "stable"
		if changePercent > 5 {
			trend = "increasing"
		} else if changePercent < -5 {
			trend = "decreasing"
		}

		trends = append(trends, &TrendAnalysis{
			MetricName:      "error_rate",
			Trend:           trend,
			ChangePercent:   changePercent,
			RecentAverage:   recentErrorRate,
			PreviousAverage: previousErrorRate,
			SampleSize:      len(recentEvents),
		})
	}

	return trends
}

// calculateAverageDuration calculates average duration from events
func calculateAverageDuration(events []*AnalyticsEvent) float64 {
	if len(events) == 0 {
		return 0
	}

	totalDuration := time.Duration(0)
	count := 0

	for _, event := range events {
		if event.Duration > 0 {
			totalDuration += event.Duration
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return totalDuration.Seconds() / float64(count)
}

// calculateErrorRate calculates error rate from events
func calculateErrorRate(events []*AnalyticsEvent) float64 {
	if len(events) == 0 {
		return 0
	}

	errorCount := 0
	for _, event := range events {
		if event.Status == "error" || event.Status == "timeout" {
			errorCount++
		}
	}

	return float64(errorCount) / float64(len(events)) * 100
}

// CreateTimeWindow creates a snapshot of analytics for a time period
func (a *AnalyticsEngine) CreateTimeWindow(duration time.Duration) *AnalyticsWindow {
	endTime := time.Now()
	startTime := endTime.Add(-duration)

	a.mu.Lock()
	defer a.mu.Unlock()

	window := &AnalyticsWindow{
		StartTime: startTime,
		EndTime:   endTime,
		Events:    make([]*AnalyticsEvent, 0),
	}

	for _, event := range a.events {
		if event.Timestamp.After(startTime) && event.Timestamp.Before(endTime) {
			window.Events = append(window.Events, event)
		}
	}

	a.windows = append(a.windows, window)

	// Limit windows stored
	if len(a.windows) > 100 {
		a.windows = a.windows[1:]
	}

	return window
}

// ExportAnalyticsJSON exports analytics as JSON-compatible map
func (a *AnalyticsEngine) ExportAnalyticsJSON() map[string]interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	result := make(map[string]interface{})

	// Export plugin statistics
	statsMap := make(map[string]map[string]interface{})
	for id, stats := range a.pluginStats {
		statsMap[id] = map[string]interface{}{
			"load_count":           stats.LoadCount,
			"unload_count":         stats.UnloadCount,
			"execution_count":      stats.ExecutionCount,
			"error_count":          stats.ErrorCount,
			"success_rate":         fmt.Sprintf("%.2f%%", stats.SuccessRate),
			"average_load_time_ms": stats.AverageLoadTime.Milliseconds(),
			"average_exec_time_ms": stats.AverageExecTime.Milliseconds(),
			"total_exec_time_s":    stats.TotalExecutionTime.Seconds(),
			"last_activity":        stats.LastActivityTime.Format(time.RFC3339),
			"first_seen":           stats.FirstSeenTime.Format(time.RFC3339),
			"uptime_percent":       fmt.Sprintf("%.2f%%", stats.UptimePercent),
		}
	}
	result["plugin_stats"] = statsMap

	// Export event counts by type
	eventCounts := make(map[string]int)
	for eventType, events := range a.eventsByType {
		eventCounts[eventType] = len(events)
	}
	result["event_counts"] = eventCounts

	// Export recent events (last 50)
	recentEvents := make([]map[string]interface{}, 0)
	start := len(a.events) - 50
	if start < 0 {
		start = 0
	}

	for _, event := range a.events[start:] {
		recentEvents = append(recentEvents, map[string]interface{}{
			"event_type":  event.EventType,
			"plugin_id":   event.PluginID,
			"timestamp":   event.Timestamp.Format(time.RFC3339),
			"duration_ms": event.Duration.Milliseconds(),
			"status":      event.Status,
			"error":       event.Error,
		})
	}
	result["recent_events"] = recentEvents

	return result
}

// GetAnalyticsSummary generates an analytics summary
func (a *AnalyticsEngine) GetAnalyticsSummary() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	summary := "=== Analytics Summary ===\n\n"

	if len(a.pluginStats) == 0 {
		return summary + "No analytics data available.\n"
	}

	summary += fmt.Sprintf("Tracked Plugins: %d\n", len(a.pluginStats))
	summary += fmt.Sprintf("Total Events: %d\n\n", len(a.events))

	// Event type breakdown
	summary += "=== Event Breakdown ===\n"
	for eventType, events := range a.eventsByType {
		summary += fmt.Sprintf("- %s: %d\n", eventType, len(events))
	}
	summary += "\n"

	// Top plugins by usage
	mostUsed := a.GetMostUsedPlugins(5)
	if len(mostUsed) > 0 {
		summary += "=== Top 5 Most Used Plugins ===\n"
		for i, stats := range mostUsed {
			summary += fmt.Sprintf("%d. %s: %d executions (%.1f%% success rate)\n",
				i+1, stats.PluginID, stats.ExecutionCount, stats.SuccessRate)
		}
		summary += "\n"
	}

	// Plugin reliability
	reliable := a.GetMostReliablePlugins(1)
	if len(reliable) > 0 {
		summary += "=== Most Reliable Plugins (100% success) ===\n"
		for _, stats := range reliable {
			if stats.SuccessRate == 100 {
				summary += fmt.Sprintf("- %s: %d executions, 0 errors\n",
					stats.PluginID, stats.ExecutionCount)
			}
		}
		summary += "\n"
	}

	return summary
}

// ClearAnalytics clears all analytics data
func (a *AnalyticsEngine) ClearAnalytics() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.events = make([]*AnalyticsEvent, 0)
	a.eventsByPlugin = make(map[string][]*AnalyticsEvent)
	a.eventsByType = make(map[string][]*AnalyticsEvent)
	a.pluginStats = make(map[string]*PluginStats)
	a.windows = make([]*AnalyticsWindow, 0)
}
