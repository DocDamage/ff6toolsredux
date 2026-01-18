package forms

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/plugins"
)

// PerformanceMetric represents cached metric data for efficient UI rendering
type PerformanceMetric struct {
	PluginID         string
	AvgExecutionTime float64 // milliseconds
	MaxExecutionTime float64
	MinExecutionTime float64
	ExecutionCount   int
	SuccessRate      float64 // percentage 0-100
	ErrorCount       int
	AvgMemoryMB      float64
	PeakMemoryMB     float64
	AvgCPUPercent    float64
	Bottlenecks      []string // List of detected bottlenecks
	LastUpdated      time.Time
}

// PluginPerformanceDashboard is the main performance metrics visualization
type PluginPerformanceDashboard struct {
	profiler    *plugins.PluginProfiler
	analytics   *plugins.AnalyticsEngine
	manager     *plugins.Manager
	window      fyne.Window
	refreshRate time.Duration
	metrics     map[string]*PerformanceMetric
	mu          sync.RWMutex
	stopChan    chan struct{}
	running     bool
}

// NewPluginPerformanceDashboard creates a new dashboard instance
func NewPluginPerformanceDashboard(manager *plugins.Manager, window fyne.Window) *PluginPerformanceDashboard {
	return &PluginPerformanceDashboard{
		profiler:    manager.GetProfiler(),
		analytics:   manager.GetAnalytics(),
		manager:     manager,
		window:      window,
		refreshRate: 1 * time.Second,
		metrics:     make(map[string]*PerformanceMetric),
		stopChan:    make(chan struct{}),
		running:     false,
	}
}

// Show displays the performance dashboard
func (d *PluginPerformanceDashboard) Show() fyne.CanvasObject {
	// Create main container with all sections
	content := container.NewVBox(
		d.createHeaderSection(),
		widget.NewSeparator(),
		d.createSystemHealthSection(),
		widget.NewSeparator(),
		d.createPluginPerformanceTable(),
		widget.NewSeparator(),
		d.createBottleneckAlertsSection(),
		widget.NewSeparator(),
		d.createExportSection(),
	)

	// Wrap in scrollable container
	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(1000, 600))

	// Start auto-refresh
	go d.autoRefresh()
	d.running = true

	return scroll
}

// createHeaderSection creates the dashboard title and refresh status
func (d *PluginPerformanceDashboard) createHeaderSection() *fyne.Container {
	title := canvas.NewText("Plugin Performance Dashboard", color.White)
	title.TextSize = 24
	title.TextStyle.Bold = true

	subtitle := canvas.NewText(fmt.Sprintf("Real-time monitoring | Refresh: %v", d.refreshRate), color.Gray{Y: 128})
	subtitle.TextSize = 12

	return container.NewVBox(title, subtitle)
}

// createSystemHealthSection displays overall system health
func (d *PluginPerformanceDashboard) createSystemHealthSection() *fyne.Container {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Calculate system-wide metrics
	totalPlugins := len(d.metrics)
	totalExecutions := 0
	totalErrors := 0
	totalSuccesses := 0
	totalMemory := 0.0
	averageSuccessRate := 0.0

	if totalPlugins > 0 {
		for _, m := range d.metrics {
			totalExecutions += m.ExecutionCount
			totalErrors += m.ErrorCount
			totalSuccesses += m.ExecutionCount - m.ErrorCount
			totalMemory += m.PeakMemoryMB
			averageSuccessRate += m.SuccessRate
		}
		averageSuccessRate /= float64(totalPlugins)
	}

	// Create health indicators
	healthLabel := canvas.NewText(fmt.Sprintf(
		"System Health: %d Plugins | %d Executions | Success Rate: %.1f%% | Memory: %.1f MB",
		totalPlugins, totalExecutions, averageSuccessRate, totalMemory,
	), d.getHealthColor(averageSuccessRate))
	healthLabel.TextSize = 14
	healthLabel.TextStyle.Bold = true

	// Error summary
	errorLabel := canvas.NewText(fmt.Sprintf("Total Errors: %d | Success: %d", totalErrors, totalSuccesses), color.Gray{Y: 200})
	errorLabel.TextSize = 12

	return container.NewVBox(healthLabel, errorLabel)
}

// getHealthColor returns color based on success rate
func (d *PluginPerformanceDashboard) getHealthColor(successRate float64) color.Color {
	if successRate >= 95 {
		return color.NRGBA{G: 200, A: 255} // Green
	} else if successRate >= 80 {
		return color.NRGBA{R: 255, G: 165, A: 255} // Orange
	}
	return color.NRGBA{R: 255, A: 255} // Red
}

// createPluginPerformanceTable creates the main metrics table
func (d *PluginPerformanceDashboard) createPluginPerformanceTable() *fyne.Container {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Create table with headers
	headers := []string{"Plugin ID", "Executions", "Avg Time (ms)", "Max Time (ms)", "Success Rate", "Errors", "Avg Memory (MB)", "Status"}

	// Collect plugin data
	type pluginRow struct {
		id     string
		metric *PerformanceMetric
		status string
	}

	rows := []pluginRow{}
	for id, m := range d.metrics {
		status := "✓ OK"
		if len(m.Bottlenecks) > 0 {
			status = "⚠ Alert"
		}
		rows = append(rows, pluginRow{id, m, status})
	}

	// Sort by plugin ID for consistency
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].id < rows[j].id
	})

	// Build table data
	tableData := make([]string, 0, len(headers)+len(rows)*len(headers))
	tableData = append(tableData, headers...)

	for _, row := range rows {
		m := row.metric
		tableData = append(tableData, []string{
			row.id,
			fmt.Sprintf("%d", m.ExecutionCount),
			fmt.Sprintf("%.2f", m.AvgExecutionTime),
			fmt.Sprintf("%.2f", m.MaxExecutionTime),
			fmt.Sprintf("%.1f%%", m.SuccessRate),
			fmt.Sprintf("%d", m.ErrorCount),
			fmt.Sprintf("%.2f", m.AvgMemoryMB),
			row.status,
		}...)
	}

	// Create table widget
	table := widget.NewTable(
		func() (rows, cols int) {
			if len(tableData) == 0 {
				return 0, 0
			}
			rows = (len(tableData) / len(headers))
			cols = len(headers)
			return
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			idx := id.Row*len(headers) + id.Col
			if idx < len(tableData) {
				obj.(*widget.Label).SetText(tableData[idx])
			}
		},
	)

	table.SetColumnWidth(0, 150) // Plugin ID
	table.SetColumnWidth(1, 100) // Executions
	table.SetColumnWidth(2, 120) // Avg Time
	table.SetColumnWidth(3, 120) // Max Time
	table.SetColumnWidth(4, 120) // Success Rate
	table.SetColumnWidth(5, 80)  // Errors
	table.SetColumnWidth(6, 130) // Memory
	table.SetColumnWidth(7, 100) // Status

	return container.NewVBox(
		canvas.NewText("Plugin Performance Metrics", color.White),
		table,
	)
}

// createBottleneckAlertsSection displays detected bottlenecks
func (d *PluginPerformanceDashboard) createBottleneckAlertsSection() *fyne.Container {
	d.mu.RLock()
	defer d.mu.RUnlock()

	alertBox := container.NewVBox(
		canvas.NewText("Bottleneck Alerts", color.White),
	)

	// Collect all bottlenecks
	type alert struct {
		plugin   string
		issue    string
		severity color.Color
	}

	alerts := []alert{}

	for id, m := range d.metrics {
		for _, bottleneck := range m.Bottlenecks {
			severity := color.NRGBA{G: 150, A: 255} // Green/Yellow
			if strings.Contains(bottleneck, "CRITICAL") || strings.Contains(bottleneck, "HIGH") {
				severity = color.NRGBA{R: 255, A: 255} // Red
			}
			alerts = append(alerts, alert{id, bottleneck, severity})
		}
	}

	// Sort by severity (critical first)
	sort.Slice(alerts, func(i, j int) bool {
		isCriticalI := strings.Contains(alerts[i].issue, "CRITICAL")
		isCriticalJ := strings.Contains(alerts[j].issue, "CRITICAL")
		return isCriticalI && !isCriticalJ
	})

	// Display alerts (limit to 20 most recent)
	if len(alerts) == 0 {
		alertBox.Add(canvas.NewText("No bottlenecks detected", color.Gray{Y: 128}))
	} else {
		if len(alerts) > 20 {
			alerts = alerts[:20]
		}
		for i, a := range alerts {
			text := canvas.NewText(fmt.Sprintf("[%d] %s: %s", i+1, a.plugin, a.issue), a.severity)
			text.TextSize = 12
			alertBox.Add(text)
		}
	}

	return alertBox
}

// createExportSection creates export controls
func (d *PluginPerformanceDashboard) createExportSection() *fyne.Container {
	exportContainer := container.NewHBox(
		canvas.NewText("Export Metrics:", color.White),
		widget.NewButton("CSV", func() {
			d.exportAsCSV()
		}),
		widget.NewButton("JSON", func() {
			d.exportAsJSON()
		}),
		widget.NewButton("Refresh Now", func() {
			d.RefreshMetrics()
		}),
	)

	return exportContainer
}

// RefreshMetrics updates metrics from profiler and analytics
func (d *PluginPerformanceDashboard) RefreshMetrics() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Clear old metrics
	d.metrics = make(map[string]*PerformanceMetric)

	// Get current plugins
	if d.profiler == nil {
		return
	}

	// Get aggregates from profiler
	aggregates := d.profiler.GetAggregates()
	if aggregates == nil {
		return
	}

	// Build metrics from profiler data
	for pluginID, agg := range aggregates {
		duration := float64(agg.AverageDuration.Milliseconds())
		errorRate := 0.0
		if agg.ExecutionCount > 0 {
			errorRate = (float64(agg.FailureCount) / float64(agg.ExecutionCount)) * 100
		}

		metric := &PerformanceMetric{
			PluginID:         pluginID,
			AvgExecutionTime: duration,
			MaxExecutionTime: float64(agg.MaxDuration.Milliseconds()),
			MinExecutionTime: float64(agg.MinDuration.Milliseconds()),
			ExecutionCount:   agg.ExecutionCount,
			SuccessRate:      100.0 - errorRate,
			ErrorCount:       agg.FailureCount,
			AvgMemoryMB:      agg.AverageMemoryMB,
			PeakMemoryMB:     agg.PeakMemoryMB,
			AvgCPUPercent:    agg.CPUUtilization,
			LastUpdated:      time.Now(),
			Bottlenecks:      d.detectBottlenecks(pluginID, agg),
		}

		d.metrics[pluginID] = metric
	}
}

// detectBottlenecks identifies performance issues
func (d *PluginPerformanceDashboard) detectBottlenecks(pluginID string, agg *plugins.AggregateMetrics) []string {
	var bottlenecks []string

	// HIGH_ERROR_RATE: >25% errors
	if agg.ExecutionCount > 0 {
		errorRate := float64(agg.FailureCount) / float64(agg.ExecutionCount)
		if errorRate > 0.25 {
			bottlenecks = append(bottlenecks, fmt.Sprintf("HIGH_ERROR_RATE: %.1f%%", errorRate*100))
		}
	}

	// SLOW_EXECUTION: >5000ms average
	if agg.AverageDuration.Milliseconds() > 5000 {
		bottlenecks = append(bottlenecks, fmt.Sprintf("SLOW_EXECUTION: %.2f ms avg", float64(agg.AverageDuration.Milliseconds())))
	}

	// HIGH_MEMORY: >100MB peak
	if agg.PeakMemoryMB > 100 {
		bottlenecks = append(bottlenecks, fmt.Sprintf("HIGH_MEMORY: %.2f MB peak", agg.PeakMemoryMB))
	}

	// HIGH_CPU: >80% average
	if agg.CPUUtilization > 80 {
		bottlenecks = append(bottlenecks, fmt.Sprintf("HIGH_CPU: %.1f%%", agg.CPUUtilization))
	}

	// UNRELIABLE: >100 failures
	if agg.FailureCount > 100 {
		bottlenecks = append(bottlenecks, fmt.Sprintf("UNRELIABLE: %d failures", agg.FailureCount))
	}

	return bottlenecks
}

// exportAsCSV exports metrics to CSV file
func (d *PluginPerformanceDashboard) exportAsCSV() {
	d.mu.RLock()
	defer d.mu.RUnlock()

	filename := fmt.Sprintf("performance_metrics_%d.csv", time.Now().Unix())
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating CSV file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := []string{
		"PluginID", "ExecutionCount", "AvgExecutionTimeMS", "MaxExecutionTimeMS",
		"MinExecutionTimeMS", "SuccessRate", "ErrorCount", "AvgMemoryMB",
		"PeakMemoryMB", "AvgCPUPercent", "Bottlenecks", "LastUpdated",
	}
	writer.Write(headers)

	// Write data
	for _, m := range d.metrics {
		bottlenecksStr := strings.Join(m.Bottlenecks, "; ")
		writer.Write([]string{
			m.PluginID,
			fmt.Sprintf("%d", m.ExecutionCount),
			fmt.Sprintf("%.2f", m.AvgExecutionTime),
			fmt.Sprintf("%.2f", m.MaxExecutionTime),
			fmt.Sprintf("%.2f", m.MinExecutionTime),
			fmt.Sprintf("%.1f", m.SuccessRate),
			fmt.Sprintf("%d", m.ErrorCount),
			fmt.Sprintf("%.2f", m.AvgMemoryMB),
			fmt.Sprintf("%.2f", m.PeakMemoryMB),
			fmt.Sprintf("%.2f", m.AvgCPUPercent),
			bottlenecksStr,
			m.LastUpdated.Format(time.RFC3339),
		})
	}

	fmt.Printf("✓ Exported metrics to %s\n", filename)
}

// exportAsJSON exports metrics to JSON file
func (d *PluginPerformanceDashboard) exportAsJSON() {
	d.mu.RLock()
	defer d.mu.RUnlock()

	filename := fmt.Sprintf("performance_metrics_%d.json", time.Now().Unix())

	type exportData struct {
		Timestamp time.Time
		Metrics   map[string]*PerformanceMetric
	}

	data := exportData{
		Timestamp: time.Now(),
		Metrics:   d.metrics,
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

	fmt.Printf("✓ Exported metrics to %s\n", filename)
}

// autoRefresh periodically updates metrics
func (d *PluginPerformanceDashboard) autoRefresh() {
	ticker := time.NewTicker(d.refreshRate)
	defer ticker.Stop()

	for {
		select {
		case <-d.stopChan:
			d.running = false
			return
		case <-ticker.C:
			d.RefreshMetrics()
		}
	}
}

// SetRefreshRate updates the refresh interval
func (d *PluginPerformanceDashboard) SetRefreshRate(rate time.Duration) {
	d.refreshRate = rate
}

// Stop halts the dashboard refresh loop
func (d *PluginPerformanceDashboard) Stop() {
	if d.running {
		close(d.stopChan)
	}
}

// GetMetrics returns current cached metrics
func (d *PluginPerformanceDashboard) GetMetrics() map[string]*PerformanceMetric {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Return copy to prevent external mutation
	metricsCopy := make(map[string]*PerformanceMetric)
	for k, v := range d.metrics {
		metricsCopy[k] = v
	}
	return metricsCopy
}

// GetPluginMetric returns metrics for a specific plugin
func (d *PluginPerformanceDashboard) GetPluginMetric(pluginID string) *PerformanceMetric {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.metrics[pluginID]
}
