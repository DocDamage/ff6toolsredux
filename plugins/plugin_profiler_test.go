package plugins

import (
	"fmt"
	"testing"
	"time"
)

func TestProfiler_RecordExecution(t *testing.T) {
	profiler := NewPluginProfiler(100, 1.0)
	profiler.RecordExecution("test_plugin", 100*time.Millisecond, true, nil)

	results := profiler.GetMetrics("test_plugin")
	if results == nil || len(results) == 0 {
		t.Fatal("Expected metrics to be recorded")
	}
	if results[0].PluginID != "test_plugin" {
		t.Errorf("Expected plugin ID test_plugin, got %s", results[0].PluginID)
	}
}

func TestProfiler_GetAggregates(t *testing.T) {
	profiler := NewPluginProfiler(100, 1.0)

	for i := 0; i < 10; i++ {
		profiler.RecordExecution("plugin_1", time.Duration(50+i*10)*time.Millisecond, true, nil)
	}

	aggs := profiler.GetAggregates()
	if aggs == nil {
		t.Fatal("Expected aggregates")
	}
	agg, exists := aggs["plugin_1"]
	if !exists {
		t.Fatal("Expected aggregates for plugin_1")
	}
	if agg.ExecutionCount != 10 {
		t.Errorf("Expected 10 executions, got %d", agg.ExecutionCount)
	}
	if agg.SuccessCount != 10 {
		t.Errorf("Expected 10 successes, got %d", agg.SuccessCount)
	}
}

func TestProfiler_Failures(t *testing.T) {
	profiler := NewPluginProfiler(100, 1.0)

	for i := 0; i < 10; i++ {
		success := i%2 != 0
		var err error
		if !success {
			err = fmt.Errorf("test error")
		}
		profiler.RecordExecution("plugin_1", 100*time.Millisecond, success, err)
	}

	aggs := profiler.GetAggregates()
	agg, exists := aggs["plugin_1"]
	if !exists {
		t.Fatal("Expected aggregates for plugin_1")
	}
	if agg.ExecutionCount != 10 {
		t.Errorf("Expected 10 executions, got %d", agg.ExecutionCount)
	}
	if agg.FailureCount != 5 {
		t.Errorf("Expected 5 failures, got %d", agg.FailureCount)
	}
}

func TestProfiler_ConcurrentAccess(t *testing.T) {
	profiler := NewPluginProfiler(1000, 1.0)

	done := make(chan bool, 10)
	for g := 0; g < 10; g++ {
		go func(id int) {
			for i := 0; i < 100; i++ {
				profiler.RecordExecution("plugin_1", time.Duration(i*10)*time.Millisecond, true, nil)
			}
			done <- true
		}(g)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	aggs := profiler.GetAggregates()
	agg, exists := aggs["plugin_1"]
	if !exists {
		t.Fatal("Expected aggregates for plugin_1")
	}
	if agg.ExecutionCount != 1000 {
		t.Errorf("Expected 1000 executions, got %d", agg.ExecutionCount)
	}
}

func TestProfiler_ExportJSON(t *testing.T) {
	profiler := NewPluginProfiler(100, 1.0)
	profiler.RecordExecution("test_plugin", 100*time.Millisecond, true, nil)

	jsonData := profiler.ExportMetricsJSON()
	if jsonData == nil {
		t.Fatal("Expected JSON export")
	}
}

func TestProfiler_GetReport(t *testing.T) {
	profiler := NewPluginProfiler(100, 1.0)
	profiler.RecordExecution("test_plugin", 100*time.Millisecond, true, nil)

	report := profiler.GetPerformanceReport()
	if report == "" {
		t.Error("Expected non-empty report")
	}
}

func TestProfiler_ClearMetrics(t *testing.T) {
	profiler := NewPluginProfiler(100, 1.0)
	profiler.RecordExecution("test_plugin", 100*time.Millisecond, true, nil)

	profiler.ClearMetrics("test_plugin")

	result := profiler.GetMetrics("test_plugin")
	if result != nil {
		t.Error("Expected metrics cleared")
	}
}
