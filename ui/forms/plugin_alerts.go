package forms

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// BottleneckAlert represents a performance or security alert
type BottleneckAlert struct {
	AlertID     string
	PluginID    string
	AlertType   string // "HIGH_ERROR_RATE", "SLOW_EXECUTION", "HIGH_MEMORY", "HIGH_CPU", "UNRELIABLE"
	Severity    string // "CRITICAL", "WARNING", "INFO"
	Message     string
	Timestamp   time.Time
	Metric      interface{}
	Dismissed   bool
	DismissedAt time.Time
}

// AlertManager manages and displays performance alerts
type AlertManager struct {
	alerts    []*BottleneckAlert
	mu        sync.RWMutex
	maxAlerts int // Keep last N alerts
	alertChan chan *BottleneckAlert
	window    fyne.Window
	stopChan  chan struct{}
	running   bool
}

// NewAlertManager creates a new alert manager
func NewAlertManager(maxAlerts int, window fyne.Window) *AlertManager {
	return &AlertManager{
		alerts:    make([]*BottleneckAlert, 0, maxAlerts),
		maxAlerts: maxAlerts,
		alertChan: make(chan *BottleneckAlert, 100),
		window:    window,
		stopChan:  make(chan struct{}),
		running:   false,
	}
}

// AddAlert adds a new alert
func (am *AlertManager) AddAlert(pluginID, alertType, severity, message string, metric interface{}) {
	alert := &BottleneckAlert{
		AlertID:   fmt.Sprintf("alert_%d", time.Now().UnixNano()),
		PluginID:  pluginID,
		AlertType: alertType,
		Severity:  severity,
		Message:   message,
		Timestamp: time.Now(),
		Metric:    metric,
		Dismissed: false,
	}

	am.mu.Lock()
	am.alerts = append(am.alerts, alert)

	// Keep only last N alerts
	if len(am.alerts) > am.maxAlerts {
		am.alerts = am.alerts[len(am.alerts)-am.maxAlerts:]
	}
	am.mu.Unlock()

	// Send to alert channel if running
	if am.running {
		select {
		case am.alertChan <- alert:
		default:
			// Channel full, skip
		}
	}
}

// DismissAlert marks an alert as dismissed
func (am *AlertManager) DismissAlert(alertID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	for _, a := range am.alerts {
		if a.AlertID == alertID {
			a.Dismissed = true
			a.DismissedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("alert not found: %s", alertID)
}

// GetActiveAlerts returns non-dismissed alerts
func (am *AlertManager) GetActiveAlerts() []*BottleneckAlert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	active := make([]*BottleneckAlert, 0)
	for _, a := range am.alerts {
		if !a.Dismissed {
			active = append(active, a)
		}
	}
	return active
}

// GetAlertsByPluginID returns alerts for a specific plugin
func (am *AlertManager) GetAlertsByPluginID(pluginID string) []*BottleneckAlert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	pluginAlerts := make([]*BottleneckAlert, 0)
	for _, a := range am.alerts {
		if a.PluginID == pluginID && !a.Dismissed {
			pluginAlerts = append(pluginAlerts, a)
		}
	}
	return pluginAlerts
}

// GetAlertsBySeverity returns alerts of a specific severity level
func (am *AlertManager) GetAlertsBySeverity(severity string) []*BottleneckAlert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	filtered := make([]*BottleneckAlert, 0)
	for _, a := range am.alerts {
		if a.Severity == severity && !a.Dismissed {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

// GetAllAlerts returns all alerts including dismissed ones
func (am *AlertManager) GetAllAlerts() []*BottleneckAlert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	alertsCopy := make([]*BottleneckAlert, len(am.alerts))
	copy(alertsCopy, am.alerts)
	return alertsCopy
}

// ClearAlerts removes all alerts
func (am *AlertManager) ClearAlerts() {
	am.mu.Lock()
	am.alerts = make([]*BottleneckAlert, 0, am.maxAlerts)
	am.mu.Unlock()
}

// ClearDismissedAlerts removes dismissed alerts
func (am *AlertManager) ClearDismissedAlerts() {
	am.mu.Lock()
	defer am.mu.Unlock()

	active := make([]*BottleneckAlert, 0)
	for _, a := range am.alerts {
		if !a.Dismissed {
			active = append(active, a)
		}
	}
	am.alerts = active
}

// Show displays the alert panel
func (am *AlertManager) Show() fyne.CanvasObject {
	content := container.NewVBox(
		am.createHeaderSection(),
		widget.NewSeparator(),
		am.createAlertListSection(),
		widget.NewSeparator(),
		am.createControlsSection(),
	)

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(900, 400))

	am.running = true
	go am.alertListener()

	return scroll
}

// createHeaderSection creates alert panel header
func (am *AlertManager) createHeaderSection() *fyne.Container {
	title := canvas.NewText("Performance & Security Alerts", color.White)
	title.TextSize = 20
	title.TextStyle.Bold = true

	activeCount := len(am.GetActiveAlerts())
	subtitle := canvas.NewText(fmt.Sprintf("Active Alerts: %d", activeCount), color.Gray{Y: 200})
	subtitle.TextSize = 12

	return container.NewVBox(title, subtitle)
}

// createAlertListSection displays active alerts
func (am *AlertManager) createAlertListSection() *fyne.Container {
	container := container.NewVBox()

	activeAlerts := am.GetActiveAlerts()
	if len(activeAlerts) == 0 {
		container.Add(canvas.NewText("No active alerts ✓", color.NRGBA{G: 200, A: 255}))
		return container
	}

	// Sort by severity (CRITICAL first)
	severityOrder := map[string]int{"CRITICAL": 0, "WARNING": 1, "INFO": 2}
	for i := 0; i < len(activeAlerts)-1; i++ {
		for j := i + 1; j < len(activeAlerts); j++ {
			if severityOrder[activeAlerts[i].Severity] > severityOrder[activeAlerts[j].Severity] {
				activeAlerts[i], activeAlerts[j] = activeAlerts[j], activeAlerts[i]
			}
		}
	}

	// Display limited number (top 20)
	displayCount := len(activeAlerts)
	if displayCount > 20 {
		displayCount = 20
	}

	for i := 0; i < displayCount; i++ {
		alert := activeAlerts[i]
		container.Add(am.createAlertRow(alert))
	}

	if len(activeAlerts) > 20 {
		remaining := len(activeAlerts) - 20
		text := canvas.NewText(fmt.Sprintf("... and %d more alerts", remaining), color.Gray{Y: 128})
		text.TextSize = 11
		container.Add(text)
	}

	return container
}

// createAlertRow creates a single alert display row
func (am *AlertManager) createAlertRow(alert *BottleneckAlert) *fyne.Container {
	// Determine color based on severity
	alertColor := am.getSeverityColor(alert.Severity)

	// Build alert text
	headerText := fmt.Sprintf("[%s] %s - %s", alert.Severity, alert.PluginID, alert.AlertType)
	header := canvas.NewText(headerText, alertColor)
	header.TextSize = 12
	header.TextStyle.Bold = true

	messageText := fmt.Sprintf("  %s", alert.Message)
	message := canvas.NewText(messageText, color.Gray{Y: 200})
	message.TextSize = 11

	timeText := fmt.Sprintf("  Time: %s", alert.Timestamp.Format("15:04:05"))
	timeLabel := canvas.NewText(timeText, color.Gray{Y: 100})
	timeLabel.TextSize = 10

	// Create dismiss button
	dismissBtn := widget.NewButton("Dismiss", func() {
		am.DismissAlert(alert.AlertID)
	})

	// Create row container
	row := container.NewHBox(
		container.NewVBox(header, message, timeLabel),
		dismissBtn,
	)

	return row
}

// getSeverityColor returns color for alert severity
func (am *AlertManager) getSeverityColor(severity string) color.Color {
	switch severity {
	case "CRITICAL":
		return color.NRGBA{R: 255, A: 255} // Red
	case "WARNING":
		return color.NRGBA{R: 255, G: 165, A: 255} // Orange
	case "INFO":
		return color.NRGBA{G: 200, A: 255} // Green
	default:
		return color.White
	}
}

// createControlsSection creates control buttons
func (am *AlertManager) createControlsSection() *fyne.Container {
	return container.NewHBox(
		canvas.NewText("Actions:", color.White),
		widget.NewButton("Clear Dismissed", func() {
			am.ClearDismissedAlerts()
		}),
		widget.NewButton("Clear All", func() {
			am.ClearAlerts()
		}),
		widget.NewButton("Export Log", func() {
			am.exportAlertLog()
		}),
	)
}

// exportAlertLog exports alerts to CSV file
func (am *AlertManager) exportAlertLog() {
	alerts := am.GetAllAlerts()
	if len(alerts) == 0 {
		fmt.Println("No alerts to export")
		return
	}

	filename := fmt.Sprintf("plugin_alerts_%d.csv", time.Now().Unix())
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
		"AlertID", "Timestamp", "Severity", "Type", "PluginID", "Message", "Dismissed",
	}
	writer.Write(headers)

	// Write alert data
	for _, alert := range alerts {
		dismissed := "No"
		if alert.Dismissed {
			dismissed = "Yes"
		}
		row := []string{
			alert.AlertID,
			alert.Timestamp.Format(time.RFC3339),
			alert.Severity,
			alert.AlertType,
			alert.PluginID,
			alert.Message,
			dismissed,
		}
		writer.Write(row)
	}

	fmt.Printf("✓ Exported %d alerts to %s\n", len(alerts), filename)
}

// alertListener listens for new alerts on the channel
func (am *AlertManager) alertListener() {
	for {
		select {
		case <-am.stopChan:
			am.running = false
			return
		case alert := <-am.alertChan:
			// Process new alert (could trigger notifications, logging, etc.)
			_ = alert // Use alert parameter to avoid unused error
		}
	}
}

// Stop halts the alert manager
func (am *AlertManager) Stop() {
	if am.running {
		close(am.stopChan)
	}
}

// GetAlertStats returns statistics about alerts
func (am *AlertManager) GetAlertStats() map[string]interface{} {
	am.mu.RLock()
	defer am.mu.RUnlock()

	stats := map[string]interface{}{
		"total_alerts":  len(am.alerts),
		"active_alerts": 0,
		"dismissed":     0,
		"critical":      0,
		"warning":       0,
		"info":          0,
		"by_type":       make(map[string]int),
	}

	criticalCount := 0
	warningCount := 0
	infoCount := 0
	dismissedCount := 0
	activeCount := 0
	byType := make(map[string]int)

	for _, a := range am.alerts {
		switch a.Severity {
		case "CRITICAL":
			criticalCount++
		case "WARNING":
			warningCount++
		case "INFO":
			infoCount++
		}

		if a.Dismissed {
			dismissedCount++
		} else {
			activeCount++
		}

		byType[a.AlertType]++
	}

	stats["active_alerts"] = activeCount
	stats["dismissed"] = dismissedCount
	stats["critical"] = criticalCount
	stats["warning"] = warningCount
	stats["info"] = infoCount
	stats["by_type"] = byType

	return stats
}

// HasCriticalAlerts returns true if there are critical alerts
func (am *AlertManager) HasCriticalAlerts() bool {
	critical := am.GetAlertsBySeverity("CRITICAL")
	return len(critical) > 0
}
