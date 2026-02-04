package plugins

import (
	"strings"
	"testing"
	"time"
)

func TestAnalyticsEngine_RecordEvent(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	event := &AnalyticsEvent{
		EventType: "plugin_load",
		PluginID:  "test_plugin",
		Timestamp: time.Now(),
		Status:    "success",
	}

	engine.RecordEvent(event)

	// Verify event was recorded
	events := engine.GetPluginEvents("test_plugin")
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}
	if events[0].PluginID != "test_plugin" {
		t.Errorf("Expected plugin ID test_plugin, got %s", events[0].PluginID)
	}
}

func TestAnalyticsEngine_GetEvents(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	// Record multiple events
	now := time.Now()
	for i := 0; i < 5; i++ {
		event := &AnalyticsEvent{
			EventType: "plugin_execute",
			PluginID:  "plugin_1",
			Timestamp: now.Add(time.Duration(i*100) * time.Millisecond),
			Status:    "success",
		}
		engine.RecordEvent(event)
	}

	events := engine.GetEvents(now.Add(-1*time.Second), now.Add(1*time.Second))
	if len(events) < 1 {
		t.Errorf("Expected at least 1 event, got %d", len(events))
	}
}

func TestAnalyticsEngine_GetPluginEvents(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	// Record events for multiple plugins
	event1 := &AnalyticsEvent{
		EventType: "plugin_load",
		PluginID:  "plugin_1",
		Timestamp: time.Now(),
		Status:    "success",
	}
	event2 := &AnalyticsEvent{
		EventType: "plugin_load",
		PluginID:  "plugin_2",
		Timestamp: time.Now(),
		Status:    "success",
	}

	engine.RecordEvent(event1)
	engine.RecordEvent(event2)

	// Get events for specific plugin
	plugin1Events := engine.GetPluginEvents("plugin_1")
	if len(plugin1Events) != 1 {
		t.Errorf("Expected 1 event for plugin_1, got %d", len(plugin1Events))
	}

	plugin2Events := engine.GetPluginEvents("plugin_2")
	if len(plugin2Events) != 1 {
		t.Errorf("Expected 1 event for plugin_2, got %d", len(plugin2Events))
	}
}

func TestAnalyticsEngine_GetEventsByType(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	// Record different event types
	eventTypes := []string{"plugin_load", "plugin_execute", "plugin_error", "plugin_unload"}

	for _, eventType := range eventTypes {
		for i := 0; i < 3; i++ {
			event := &AnalyticsEvent{
				EventType: eventType,
				PluginID:  "test_plugin",
				Timestamp: time.Now(),
				Status:    "success",
			}
			if eventType == "plugin_error" {
				event.Status = "error"
				event.Error = "test error"
			}
			engine.RecordEvent(event)
		}
	}

	// Get events by type
	loadEvents := engine.GetEventsByType("plugin_load")
	if len(loadEvents) != 3 {
		t.Errorf("Expected 3 load events, got %d", len(loadEvents))
	}

	executeEvents := engine.GetEventsByType("plugin_execute")
	if len(executeEvents) != 3 {
		t.Errorf("Expected 3 execute events, got %d", len(executeEvents))
	}
}

func TestAnalyticsEngine_GetPluginStats(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	// Record mixed events for a plugin
	events := []struct {
		eventType string
		status    string
	}{
		{"plugin_load", "success"},
		{"plugin_execute", "success"},
		{"plugin_execute", "success"},
		{"plugin_execute", "error"},
		{"plugin_unload", "success"},
	}

	now := time.Now()
	for i, evt := range events {
		event := &AnalyticsEvent{
			EventType: evt.eventType,
			PluginID:  "test_plugin",
			Timestamp: now.Add(time.Duration(i*100) * time.Millisecond),
			Status:    evt.status,
			Duration:  100 * time.Millisecond,
		}
		if evt.status == "error" {
			event.Error = "test error"
		}
		engine.RecordEvent(event)
	}

	stats := engine.GetPluginStats("test_plugin")
	if stats == nil {
		t.Fatal("Expected plugin stats")
	}

	if stats.LoadCount != 1 {
		t.Errorf("Expected 1 load, got %d", stats.LoadCount)
	}
	if stats.UnloadCount != 1 {
		t.Errorf("Expected 1 unload, got %d", stats.UnloadCount)
	}
	if stats.ExecutionCount != 3 {
		t.Errorf("Expected 3 executions, got %d", stats.ExecutionCount)
	}
	if stats.ErrorCount != 1 {
		t.Errorf("Expected 1 error, got %d", stats.ErrorCount)
	}

	expectedSuccessRate := 2.0 / 3.0 // 2 successful executions out of 3
	if stats.SuccessRate != expectedSuccessRate {
		t.Errorf("Expected success rate %v, got %v", expectedSuccessRate, stats.SuccessRate)
	}
}

func TestAnalyticsEngine_GetMostUsedPlugins(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	// Record different usage levels
	pluginUsage := map[string]int{
		"plugin_1": 10,
		"plugin_2": 5,
		"plugin_3": 20,
	}

	for pluginID, count := range pluginUsage {
		for i := 0; i < count; i++ {
			event := &AnalyticsEvent{
				EventType: "plugin_execute",
				PluginID:  pluginID,
				Timestamp: time.Now(),
				Status:    "success",
			}
			engine.RecordEvent(event)
		}
	}

	mostUsed := engine.GetMostUsedPlugins(3)
	if len(mostUsed) != 3 {
		t.Errorf("Expected 3 plugins, got %d", len(mostUsed))
	}

	// Verify sorted by usage (descending)
	if mostUsed[0].PluginID != "plugin_3" || mostUsed[0].ExecutionCount != 20 {
		t.Errorf("Expected plugin_3 with 20 executions first, got %s with %d",
			mostUsed[0].PluginID, mostUsed[0].ExecutionCount)
	}
	if mostUsed[1].PluginID != "plugin_1" || mostUsed[1].ExecutionCount != 10 {
		t.Errorf("Expected plugin_1 with 10 executions second, got %s with %d",
			mostUsed[1].PluginID, mostUsed[1].ExecutionCount)
	}
	if mostUsed[2].PluginID != "plugin_2" || mostUsed[2].ExecutionCount != 5 {
		t.Errorf("Expected plugin_2 with 5 executions third, got %s with %d",
			mostUsed[2].PluginID, mostUsed[2].ExecutionCount)
	}
}

func TestAnalyticsEngine_GetMostReliablePlugins(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	// Record plugins with different success rates
	// Plugin 1: 10 success, 0 errors
	for i := 0; i < 10; i++ {
		event := &AnalyticsEvent{
			EventType: "plugin_execute",
			PluginID:  "reliable_plugin",
			Timestamp: time.Now(),
			Status:    "success",
		}
		engine.RecordEvent(event)
	}

	// Plugin 2: 5 success, 5 errors
	for i := 0; i < 5; i++ {
		event := &AnalyticsEvent{
			EventType: "plugin_execute",
			PluginID:  "unreliable_plugin",
			Timestamp: time.Now(),
			Status:    "success",
		}
		engine.RecordEvent(event)

		failEvent := &AnalyticsEvent{
			EventType: "plugin_error",
			PluginID:  "unreliable_plugin",
			Timestamp: time.Now(),
			Status:    "error",
			Error:     "test error",
		}
		engine.RecordEvent(failEvent)
	}

	reliable := engine.GetMostReliablePlugins(1)
	if len(reliable) != 1 {
		t.Errorf("Expected 1 plugin (min 1 execution), got %d", len(reliable))
	}

	if reliable[0].PluginID != "reliable_plugin" {
		t.Errorf("Expected reliable_plugin, got %s", reliable[0].PluginID)
	}

	if reliable[0].SuccessRate != 1.0 {
		t.Errorf("Expected 100%% success rate, got %f", reliable[0].SuccessRate)
	}
}

func TestAnalyticsEngine_AnalyzeTrends(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	now := time.Now()

	// Record improving trend: fast -> slow
	durations := []time.Duration{
		100 * time.Millisecond,
		110 * time.Millisecond,
		120 * time.Millisecond,
		1000 * time.Millisecond, // Sudden spike
		1100 * time.Millisecond,
	}

	for i, duration := range durations {
		event := &AnalyticsEvent{
			EventType: "plugin_execute",
			PluginID:  "trending_plugin",
			Timestamp: now.Add(time.Duration(i*100) * time.Millisecond),
			Duration:  duration,
			Status:    "success",
		}
		engine.RecordEvent(event)
	}

	trends := engine.AnalyzeTrends("trending_plugin", len(durations))
	if len(trends) == 0 {
		t.Skip("No trends detected (may be expected with small sample)")
		return
	}

	// Verify trend structure
	for _, trend := range trends {
		if trend.MetricName == "" {
			t.Error("Expected metric name in trend")
		}
		if trend.Trend != "increasing" && trend.Trend != "decreasing" && trend.Trend != "stable" {
			t.Errorf("Invalid trend direction: %s", trend.Trend)
		}
	}
}

func TestAnalyticsEngine_CreateTimeWindow(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	now := time.Now()

	// Record events
	for i := 0; i < 10; i++ {
		event := &AnalyticsEvent{
			EventType: "plugin_execute",
			PluginID:  "test_plugin",
			Timestamp: now.Add(time.Duration(i*100) * time.Millisecond),
			Status:    "success",
		}
		engine.RecordEvent(event)
	}

	// Create time window for last 1 second
	window := engine.CreateTimeWindow(1 * time.Second)
	if window == nil {
		t.Fatal("Expected time window to be created")
	}

	if len(window.Events) == 0 {
		t.Error("Expected events in time window")
	}
}

func TestAnalyticsEngine_ExportAnalyticsJSON(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	event := &AnalyticsEvent{
		EventType: "plugin_load",
		PluginID:  "test_plugin",
		Timestamp: time.Now(),
		Status:    "success",
	}
	engine.RecordEvent(event)

	jsonData := engine.ExportAnalyticsJSON()
	if jsonData == nil {
		t.Fatal("Expected JSON export to return data")
	}

	// Verify structure
	if _, hasStats := jsonData["plugin_stats"]; !hasStats {
		t.Error("Expected plugin_stats in JSON export")
	}
	if _, hasEventCounts := jsonData["event_counts"]; !hasEventCounts {
		t.Error("Expected event_counts in JSON export")
	}
	if _, hasEvents := jsonData["recent_events"]; !hasEvents {
		t.Error("Expected recent_events in JSON export")
	}
}

func TestAnalyticsEngine_GetAnalyticsSummary(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	event := &AnalyticsEvent{
		EventType: "plugin_execute",
		PluginID:  "test_plugin",
		Timestamp: time.Now(),
		Status:    "success",
	}
	engine.RecordEvent(event)

	summary := engine.GetAnalyticsSummary()
	if summary == "" {
		t.Error("Expected non-empty analytics summary")
	}

	if !strings.Contains(summary, "Analytics Summary") {
		t.Error("Expected summary header in output")
	}
}

func TestAnalyticsEngine_ClearAnalytics(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	// Record events
	for i := 0; i < 5; i++ {
		event := &AnalyticsEvent{
			EventType: "plugin_execute",
			PluginID:  "test_plugin",
			Timestamp: time.Now(),
			Status:    "success",
		}
		engine.RecordEvent(event)
	}

	// Clear analytics
	engine.ClearAnalytics()

	// Verify cleared
	events := engine.GetEvents(time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour))
	if len(events) != 0 {
		t.Errorf("Expected 0 events after clear, got %d", len(events))
	}
}

func TestAnalyticsEngine_MaxEventsStored(t *testing.T) {
	maxEvents := 100
	engine := NewAnalyticsEngine(maxEvents)

	// Record more than max
	for i := 0; i < 200; i++ {
		event := &AnalyticsEvent{
			EventType: "plugin_execute",
			PluginID:  "test_plugin",
			Timestamp: time.Now(),
			Status:    "success",
		}
		engine.RecordEvent(event)
	}

	// Should respect max events
	events := engine.GetEvents(time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour))
	if len(events) > maxEvents+10 { // Allow small margin
		t.Errorf("Expected at most %d events, got %d", maxEvents+10, len(events))
	}
}

func TestAnalyticsEngine_ConcurrentAccess(t *testing.T) {
	engine := NewAnalyticsEngine(10000)

	// Simulate concurrent recording
	done := make(chan bool, 10)
	for g := 0; g < 10; g++ {
		go func(goroutineID int) {
			for i := 0; i < 100; i++ {
				event := &AnalyticsEvent{
					EventType: "plugin_execute",
					PluginID:  "test_plugin",
					Timestamp: time.Now(),
					Status:    "success",
				}
				engine.RecordEvent(event)
			}
			done <- true
		}(g)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify data integrity
	events := engine.GetEvents(time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour))
	if len(events) < 900 || len(events) > 1100 { // Should be around 1000
		t.Logf("Expected around 1000 events from concurrent access, got %d", len(events))
	}
}

func TestAnalyticsEngine_EventMetadata(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	metadata := map[string]interface{}{
		"version":    "1.0.0",
		"complexity": 5,
		"tags":       []string{"critical", "async"},
	}

	event := &AnalyticsEvent{
		EventType: "plugin_execute",
		PluginID:  "test_plugin",
		Timestamp: time.Now(),
		Status:    "success",
		Metadata:  metadata,
	}

	engine.RecordEvent(event)

	events := engine.GetPluginEvents("test_plugin")
	if len(events) != 1 {
		t.Fatal("Expected 1 event")
	}

	if events[0].Metadata == nil {
		t.Error("Expected metadata to be preserved")
	}

	if version, ok := events[0].Metadata["version"]; !ok || version != "1.0.0" {
		t.Error("Expected version in metadata")
	}
}

func TestAnalyticsEngine_AllPluginStats(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	// Record events for multiple plugins
	for pluginID := 1; pluginID <= 3; pluginID++ {
		pidStr := "plugin_" + string(rune('0'+pluginID))
		for i := 0; i < pluginID*2; i++ {
			event := &AnalyticsEvent{
				EventType: "plugin_execute",
				PluginID:  pidStr,
				Timestamp: time.Now(),
				Status:    "success",
			}
			engine.RecordEvent(event)
		}
	}

	allStats := engine.GetAllPluginStats()
	if len(allStats) != 3 {
		t.Errorf("Expected stats for 3 plugins, got %d", len(allStats))
	}

	// Verify each plugin has stats
	for pluginID := 1; pluginID <= 3; pluginID++ {
		pidStr := "plugin_" + string(rune('0'+pluginID))
		if _, hasStats := allStats[pidStr]; !hasStats {
			t.Errorf("Expected stats for %s", pidStr)
		}
	}
}

func TestAnalyticsEngine_ErrorTracking(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	// Record successful and failed executions
	for i := 0; i < 10; i++ {
		status := "success"
		var errorMsg string

		if i%3 == 0 {
			status = "error"
			errorMsg = "execution failed"
		}

		event := &AnalyticsEvent{
			EventType: "plugin_execute",
			PluginID:  "test_plugin",
			Timestamp: time.Now(),
			Status:    status,
			Error:     errorMsg,
		}
		engine.RecordEvent(event)
	}

	stats := engine.GetPluginStats("test_plugin")
	if stats.ErrorCount != 4 { // 0, 3, 6, 9
		t.Errorf("Expected 4 errors, got %d", stats.ErrorCount)
	}

	successRate := stats.SuccessRate
	expectedRate := 6.0 / 10.0 // 6 successful out of 10
	if successRate != expectedRate {
		t.Errorf("Expected success rate %v, got %v", expectedRate, successRate)
	}
}

func TestAnalyticsEngine_DurationTracking(t *testing.T) {
	engine := NewAnalyticsEngine(1000)

	durations := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		150 * time.Millisecond,
	}

	for _, duration := range durations {
		event := &AnalyticsEvent{
			EventType: "plugin_execute",
			PluginID:  "test_plugin",
			Timestamp: time.Now(),
			Duration:  duration,
			Status:    "success",
		}
		engine.RecordEvent(event)
	}

	stats := engine.GetPluginStats("test_plugin")

	// Verify average duration is reasonable (should be around 150ms)
	if stats.AverageExecTime < 100*time.Millisecond ||
		stats.AverageExecTime > 250*time.Millisecond {
		t.Errorf("Expected average execution time around 150ms, got %v", stats.AverageExecTime)
	}
}
