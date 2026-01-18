package plugins

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"
)

// ExecutionMetrics tracks detailed metrics for a plugin execution
type ExecutionMetrics struct {
	PluginID        string        // Plugin identifier
	Timestamp       time.Time     // When execution started
	Duration        time.Duration // Total execution time
	CPUTime         time.Duration // CPU time (estimated)
	MemoryAllocated uint64        // Memory allocated by plugin (bytes)
	MemoryReleased  uint64        // Memory freed by plugin (bytes)
	GoroutinesStart int           // Goroutines before execution
	GoroutinesEnd   int           // Goroutines after execution
	AllocsBefore    uint64        // Memory allocations before
	AllocsAfter     uint64        // Memory allocations after
	CallCount       int           // Number of times this execution was called
	CPUPercent      float64       // CPU usage as percentage
	MemoryMB        float64       // Peak memory in MB
	Status          string        // "success", "error", "timeout"
	ErrorMessage    string        // If status is error
}

// PerformanceSample represents a single performance sample for flame graphs
type PerformanceSample struct {
	Timestamp  time.Time
	CPUTime    time.Duration
	Memory     uint64
	DurationNs int64
}

// PluginProfiler collects and analyzes plugin performance metrics
type PluginProfiler struct {
	metrics       map[string][]*ExecutionMetrics  // pluginID -> metrics history
	samples       map[string][]*PerformanceSample // pluginID -> performance samples
	aggregates    map[string]*AggregateMetrics    // pluginID -> aggregate stats
	maxHistoryLen int                             // Max metrics to keep per plugin
	samplingRate  float64                         // 0-1, fraction of executions to sample
	mu            sync.RWMutex
}

// AggregateMetrics tracks aggregate performance statistics
type AggregateMetrics struct {
	PluginID           string
	ExecutionCount     int
	SuccessCount       int
	FailureCount       int
	TimeoutCount       int
	AverageDuration    time.Duration
	MaxDuration        time.Duration
	MinDuration        time.Duration
	AverageCPUTime     time.Duration
	AverageMemoryMB    float64
	PeakMemoryMB       float64
	ErrorRate          float64
	LastExecution      time.Time
	FirstExecution     time.Time
	TotalExecutionTime time.Duration
	CPUUtilization     float64
}

// FlameGraphNode represents a node in a flame graph
type FlameGraphNode struct {
	Name     string        // Plugin ID or function name
	Duration time.Duration // Total time in this node
	Children []*FlameGraphNode
}

// NewPluginProfiler creates a new plugin profiler
func NewPluginProfiler(maxHistoryLen int, samplingRate float64) *PluginProfiler {
	if samplingRate < 0 {
		samplingRate = 0
	}
	if samplingRate > 1 {
		samplingRate = 1
	}

	return &PluginProfiler{
		metrics:       make(map[string][]*ExecutionMetrics),
		samples:       make(map[string][]*PerformanceSample),
		aggregates:    make(map[string]*AggregateMetrics),
		maxHistoryLen: maxHistoryLen,
		samplingRate:  samplingRate,
	}
}

// RecordExecution records metrics for a plugin execution
func (p *PluginProfiler) RecordExecution(pluginID string, duration time.Duration, success bool, err error) *ExecutionMetrics {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Capture memory metrics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	status := "success"
	errorMsg := ""
	if !success {
		status = "error"
		if err != nil {
			errorMsg = err.Error()
		}
	}

	metric := &ExecutionMetrics{
		PluginID:        pluginID,
		Timestamp:       time.Now(),
		Duration:        duration,
		CPUTime:         duration, // Approximate; actual CPU profiling would be more accurate
		MemoryAllocated: m.Alloc,
		GoroutinesStart: runtime.NumGoroutine(),
		GoroutinesEnd:   runtime.NumGoroutine(),
		Status:          status,
		ErrorMessage:    errorMsg,
		MemoryMB:        float64(m.Alloc) / 1024 / 1024,
	}

	// Add to metrics history
	if _, exists := p.metrics[pluginID]; !exists {
		p.metrics[pluginID] = make([]*ExecutionMetrics, 0)
		p.aggregates[pluginID] = &AggregateMetrics{
			PluginID:       pluginID,
			FirstExecution: time.Now(),
		}
	}

	p.metrics[pluginID] = append(p.metrics[pluginID], metric)

	// Limit history size
	if len(p.metrics[pluginID]) > p.maxHistoryLen {
		p.metrics[pluginID] = p.metrics[pluginID][1:]
	}

	// Sample for performance analysis
	if p.shouldSample() {
		if _, exists := p.samples[pluginID]; !exists {
			p.samples[pluginID] = make([]*PerformanceSample, 0)
		}

		sample := &PerformanceSample{
			Timestamp:  metric.Timestamp,
			CPUTime:    metric.CPUTime,
			Memory:     m.Alloc,
			DurationNs: int64(duration),
		}
		p.samples[pluginID] = append(p.samples[pluginID], sample)

		// Limit samples size
		if len(p.samples[pluginID]) > p.maxHistoryLen {
			p.samples[pluginID] = p.samples[pluginID][1:]
		}
	}

	// Update aggregates
	p.updateAggregates(pluginID, metric)

	return metric
}

// shouldSample determines if this execution should be sampled
func (p *PluginProfiler) shouldSample() bool {
	if p.samplingRate >= 1.0 {
		return true
	}
	if p.samplingRate <= 0 {
		return false
	}

	// Simple probabilistic sampling
	return time.Now().UnixNano()%100 < int64(p.samplingRate*100)
}

// updateAggregates updates aggregate statistics
func (p *PluginProfiler) updateAggregates(pluginID string, metric *ExecutionMetrics) {
	agg := p.aggregates[pluginID]
	if agg == nil {
		agg = &AggregateMetrics{
			PluginID:       pluginID,
			FirstExecution: metric.Timestamp,
		}
		p.aggregates[pluginID] = agg
	}

	agg.ExecutionCount++
	agg.LastExecution = metric.Timestamp
	agg.TotalExecutionTime += metric.Duration

	if metric.Status == "success" {
		agg.SuccessCount++
	} else if metric.Status == "error" {
		agg.FailureCount++
	} else if metric.Status == "timeout" {
		agg.TimeoutCount++
	}

	// Update duration stats
	if agg.MaxDuration == 0 || metric.Duration > agg.MaxDuration {
		agg.MaxDuration = metric.Duration
	}
	if agg.MinDuration == 0 || metric.Duration < agg.MinDuration {
		agg.MinDuration = metric.Duration
	}

	// Recalculate averages
	if agg.ExecutionCount > 0 {
		agg.AverageDuration = agg.TotalExecutionTime / time.Duration(agg.ExecutionCount)
	}

	// Update memory stats
	if metric.MemoryMB > agg.PeakMemoryMB {
		agg.PeakMemoryMB = metric.MemoryMB
	}

	totalMemory := 0.0
	for _, m := range p.metrics[pluginID] {
		totalMemory += m.MemoryMB
	}
	agg.AverageMemoryMB = totalMemory / float64(len(p.metrics[pluginID]))

	// Calculate error rate
	totalAttempts := agg.SuccessCount + agg.FailureCount + agg.TimeoutCount
	if totalAttempts > 0 {
		agg.ErrorRate = float64(agg.FailureCount+agg.TimeoutCount) / float64(totalAttempts) * 100
	}

	// Calculate CPU utilization (simple approximation)
	if agg.ExecutionCount > 0 {
		agg.CPUUtilization = float64(agg.TotalExecutionTime.Milliseconds()) / float64(time.Since(agg.FirstExecution).Milliseconds()) * 100
	}
}

// GetMetrics returns metrics for a specific plugin
func (p *PluginProfiler) GetMetrics(pluginID string) []*ExecutionMetrics {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if metrics, exists := p.metrics[pluginID]; exists {
		// Return copy to prevent external modification
		result := make([]*ExecutionMetrics, len(metrics))
		copy(result, metrics)
		return result
	}
	return nil
}

// GetAggregates returns aggregate statistics for all plugins
func (p *PluginProfiler) GetAggregates() map[string]*AggregateMetrics {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make(map[string]*AggregateMetrics)
	for id, agg := range p.aggregates {
		result[id] = agg
	}
	return result
}

// GetAggregateForPlugin returns aggregate statistics for a specific plugin
func (p *PluginProfiler) GetAggregateForPlugin(pluginID string) *AggregateMetrics {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.aggregates[pluginID]
}

// GetSamples returns performance samples for a plugin
func (p *PluginProfiler) GetSamples(pluginID string) []*PerformanceSample {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if samples, exists := p.samples[pluginID]; exists {
		result := make([]*PerformanceSample, len(samples))
		copy(result, samples)
		return result
	}
	return nil
}

// GetSlowPlugins returns plugins with execution time above threshold
func (p *PluginProfiler) GetSlowPlugins(threshold time.Duration) []*AggregateMetrics {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var slow []*AggregateMetrics
	for _, agg := range p.aggregates {
		if agg.AverageDuration > threshold {
			slow = append(slow, agg)
		}
	}

	// Sort by average duration descending
	sort.Slice(slow, func(i, j int) bool {
		return slow[i].AverageDuration > slow[j].AverageDuration
	})

	return slow
}

// GetFailedPlugins returns plugins with highest error rates
func (p *PluginProfiler) GetFailedPlugins() []*AggregateMetrics {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var failed []*AggregateMetrics
	for _, agg := range p.aggregates {
		if agg.ErrorRate > 0 {
			failed = append(failed, agg)
		}
	}

	// Sort by error rate descending
	sort.Slice(failed, func(i, j int) bool {
		return failed[i].ErrorRate > failed[j].ErrorRate
	})

	return failed
}

// GetMemoryHogs returns plugins with highest memory usage
func (p *PluginProfiler) GetMemoryHogs() []*AggregateMetrics {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var hogs []*AggregateMetrics
	for _, agg := range p.aggregates {
		if agg.PeakMemoryMB > 0 {
			hogs = append(hogs, agg)
		}
	}

	// Sort by peak memory descending
	sort.Slice(hogs, func(i, j int) bool {
		return hogs[i].PeakMemoryMB > hogs[j].PeakMemoryMB
	})

	return hogs
}

// GetBottlenecks detects performance bottlenecks
func (p *PluginProfiler) GetBottlenecks() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var bottlenecks []string

	for id, agg := range p.aggregates {
		if agg.ErrorRate > 25 {
			bottlenecks = append(bottlenecks, fmt.Sprintf("HIGH_ERROR_RATE: %s has %.1f%% error rate", id, agg.ErrorRate))
		}

		if agg.AverageDuration > 5*time.Second {
			bottlenecks = append(bottlenecks, fmt.Sprintf("SLOW_EXECUTION: %s takes avg %.2fs", id, agg.AverageDuration.Seconds()))
		}

		if agg.PeakMemoryMB > 100 {
			bottlenecks = append(bottlenecks, fmt.Sprintf("HIGH_MEMORY: %s peaks at %.1fMB", id, agg.PeakMemoryMB))
		}

		if agg.CPUUtilization > 80 {
			bottlenecks = append(bottlenecks, fmt.Sprintf("HIGH_CPU: %s using %.1f%% CPU", id, agg.CPUUtilization))
		}

		if agg.ExecutionCount > 1000 && agg.FailureCount > 100 {
			bottlenecks = append(bottlenecks, fmt.Sprintf("UNRELIABLE: %s failed %d times out of %d executions", id, agg.FailureCount, agg.ExecutionCount))
		}
	}

	return bottlenecks
}

// GenerateFlameGraph generates a flame graph representation
func (p *PluginProfiler) GenerateFlameGraph(pluginID string) *FlameGraphNode {
	p.mu.RLock()
	defer p.mu.RUnlock()

	samples, exists := p.samples[pluginID]
	if !exists || len(samples) == 0 {
		return nil
	}

	root := &FlameGraphNode{
		Name: pluginID,
	}

	// Sum up all execution times
	var totalTime time.Duration
	for _, sample := range samples {
		totalTime += time.Duration(sample.DurationNs)
	}
	root.Duration = totalTime

	return root
}

// ExportMetricsJSON exports metrics as JSON-compatible map
func (p *PluginProfiler) ExportMetricsJSON() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make(map[string]interface{})

	// Export aggregates
	aggregates := make(map[string]map[string]interface{})
	for id, agg := range p.aggregates {
		aggregates[id] = map[string]interface{}{
			"execution_count":     agg.ExecutionCount,
			"success_count":       agg.SuccessCount,
			"failure_count":       agg.FailureCount,
			"timeout_count":       agg.TimeoutCount,
			"average_duration_ms": agg.AverageDuration.Milliseconds(),
			"max_duration_ms":     agg.MaxDuration.Milliseconds(),
			"min_duration_ms":     agg.MinDuration.Milliseconds(),
			"average_memory_mb":   fmt.Sprintf("%.2f", agg.AverageMemoryMB),
			"peak_memory_mb":      fmt.Sprintf("%.2f", agg.PeakMemoryMB),
			"error_rate":          fmt.Sprintf("%.2f%%", agg.ErrorRate),
			"cpu_utilization":     fmt.Sprintf("%.2f%%", agg.CPUUtilization),
			"last_execution":      agg.LastExecution.Format(time.RFC3339),
			"first_execution":     agg.FirstExecution.Format(time.RFC3339),
		}
	}
	result["aggregates"] = aggregates

	// Export recent metrics
	recentMetrics := make(map[string][]map[string]interface{})
	for id, metrics := range p.metrics {
		if len(metrics) > 0 {
			// Last 10 executions
			start := len(metrics) - 10
			if start < 0 {
				start = 0
			}

			for _, metric := range metrics[start:] {
				recentMetrics[id] = append(recentMetrics[id], map[string]interface{}{
					"timestamp":     metric.Timestamp.Format(time.RFC3339),
					"duration_ms":   metric.Duration.Milliseconds(),
					"memory_mb":     fmt.Sprintf("%.2f", metric.MemoryMB),
					"status":        metric.Status,
					"goroutines":    metric.GoroutinesEnd,
					"error_message": metric.ErrorMessage,
				})
			}
		}
	}
	result["recent_metrics"] = recentMetrics

	return result
}

// ClearMetrics clears all metrics for a plugin
func (p *PluginProfiler) ClearMetrics(pluginID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.metrics, pluginID)
	delete(p.samples, pluginID)
	delete(p.aggregates, pluginID)
}

// ClearAllMetrics clears all metrics
func (p *PluginProfiler) ClearAllMetrics() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.metrics = make(map[string][]*ExecutionMetrics)
	p.samples = make(map[string][]*PerformanceSample)
	p.aggregates = make(map[string]*AggregateMetrics)
}

// GetPerformanceReport generates a performance report
func (p *PluginProfiler) GetPerformanceReport() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	report := "=== Plugin Performance Report ===\n\n"

	if len(p.aggregates) == 0 {
		return report + "No metrics collected yet.\n"
	}

	// Overall stats
	report += "=== Overall Statistics ===\n"
	totalExecutions := 0
	totalFailures := 0
	totalTime := time.Duration(0)

	for _, agg := range p.aggregates {
		totalExecutions += agg.ExecutionCount
		totalFailures += agg.FailureCount + agg.TimeoutCount
		totalTime += agg.TotalExecutionTime
	}

	report += fmt.Sprintf("Total Executions: %d\n", totalExecutions)
	report += fmt.Sprintf("Total Failures: %d (%.1f%%)\n", totalFailures, float64(totalFailures)/float64(totalExecutions)*100)
	report += fmt.Sprintf("Total Execution Time: %.2fs\n\n", totalTime.Seconds())

	// Slow plugins
	slowPlugins := p.GetSlowPlugins(1 * time.Second)
	if len(slowPlugins) > 0 {
		report += "=== Slow Plugins (avg > 1s) ===\n"
		for _, agg := range slowPlugins {
			report += fmt.Sprintf("- %s: avg %.2fs (min %.2fs, max %.2fs)\n",
				agg.PluginID, agg.AverageDuration.Seconds(), agg.MinDuration.Seconds(), agg.MaxDuration.Seconds())
		}
		report += "\n"
	}

	// Failed plugins
	failedPlugins := p.GetFailedPlugins()
	if len(failedPlugins) > 0 {
		report += "=== Plugins with Failures ===\n"
		for _, agg := range failedPlugins {
			report += fmt.Sprintf("- %s: %.1f%% error rate (%d/%d)\n",
				agg.PluginID, agg.ErrorRate, agg.FailureCount, agg.ExecutionCount)
		}
		report += "\n"
	}

	// Memory hogs
	memoryHogs := p.GetMemoryHogs()
	if len(memoryHogs) > 0 {
		report += "=== Memory Usage ===\n"
		for _, agg := range memoryHogs {
			report += fmt.Sprintf("- %s: peak %.1fMB, avg %.1fMB\n",
				agg.PluginID, agg.PeakMemoryMB, agg.AverageMemoryMB)
		}
		report += "\n"
	}

	// Bottlenecks
	bottlenecks := p.GetBottlenecks()
	if len(bottlenecks) > 0 {
		report += "=== Detected Bottlenecks ===\n"
		for _, bottleneck := range bottlenecks {
			report += fmt.Sprintf("- %s\n", bottleneck)
		}
		report += "\n"
	}

	return report
}
