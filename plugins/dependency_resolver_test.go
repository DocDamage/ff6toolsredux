package plugins

import (
	"testing"
)

func TestDependencyResolver_AddPlugin(t *testing.T) {
	resolver := NewDependencyResolver()

	metadata := PluginMetadata{
		ID:      "plugin1",
		Version: "1.0.0",
	}

	plugin := &Plugin{
		ID:      "plugin1",
		Name:    "Plugin 1",
		Version: "1.0.0",
	}
	plugin.SetMetadata(metadata)

	err := resolver.AddPlugin(plugin)
	if err != nil {
		t.Fatalf("Failed to add plugin: %v", err)
	}

	resolver.mu.RLock()
	_, exists := resolver.pluginGraph.nodes["plugin1"]
	resolver.mu.RUnlock()

	if !exists {
		t.Error("Plugin should be added to graph")
	}
}

func TestDependencyResolver_RemovePlugin(t *testing.T) {
	resolver := NewDependencyResolver()

	metadata := PluginMetadata{
		ID:      "plugin1",
		Version: "1.0.0",
	}

	plugin := &Plugin{
		ID:      "plugin1",
		Name:    "Plugin 1",
		Version: "1.0.0",
	}
	plugin.SetMetadata(metadata)

	resolver.AddPlugin(plugin)
	resolver.RemovePlugin("plugin1")

	resolver.mu.RLock()
	_, exists := resolver.pluginGraph.nodes["plugin1"]
	resolver.mu.RUnlock()

	if exists {
		t.Error("Plugin should be removed from graph")
	}
}

func TestDependencyResolver_ResolveDependencies(t *testing.T) {
	resolver := NewDependencyResolver()

	// Add plugins with dependencies
	metadata1 := PluginMetadata{
		ID:           "plugin1",
		Version:      "1.0.0",
		Dependencies: []string{"plugin2"},
	}

	plugin1 := &Plugin{
		ID:      "plugin1",
		Name:    "Plugin 1",
		Version: "1.0.0",
	}
	plugin1.SetMetadata(metadata1)

	metadata2 := PluginMetadata{
		ID:      "plugin2",
		Version: "1.5.0",
	}

	plugin2 := &Plugin{
		ID:      "plugin2",
		Name:    "Plugin 2",
		Version: "1.5.0",
	}
	plugin2.SetMetadata(metadata2)

	resolver.AddPlugin(plugin1)
	resolver.AddPlugin(plugin2)

	// Note: ResolveDependencies signature changed to (pluginID, version)
	deps, err := resolver.ResolveDependencies("plugin1", "1.0.0")
	if err != nil {
		t.Fatalf("Failed to resolve dependencies: %v", err)
	}

	// Should include plugin2
	if _, exists := deps["plugin2"]; !exists {
		t.Error("Expected plugin2 in dependencies")
	}
}

func TestDependencyResolver_DetectConflicts(t *testing.T) {
	resolver := NewDependencyResolver()

	metadata1 := PluginMetadata{
		ID:           "plugin1",
		Version:      "1.0.0",
		Dependencies: []string{"missing-plugin"},
	}

	plugin1 := &Plugin{
		ID:      "plugin1",
		Name:    "Plugin 1",
		Version: "1.0.0",
	}
	plugin1.SetMetadata(metadata1)

	resolver.AddPlugin(plugin1)

	versionMap := map[string]string{
		"plugin1": "1.0.0",
	}

	conflicts := resolver.DetectConflicts(versionMap)

	if len(conflicts) == 0 {
		t.Error("Should detect conflicts for missing dependencies")
	}
}

func TestDependencyResolver_GenerateDotGraph(t *testing.T) {
	resolver := NewDependencyResolver()

	metadata1 := PluginMetadata{
		ID:           "plugin1",
		Version:      "1.0.0",
		Dependencies: []string{"plugin2"},
	}

	plugin1 := &Plugin{
		ID:      "plugin1",
		Name:    "Plugin 1",
		Version: "1.0.0",
	}
	plugin1.SetMetadata(metadata1)

	metadata2 := PluginMetadata{
		ID:      "plugin2",
		Version: "1.0.0",
	}

	plugin2 := &Plugin{
		ID:      "plugin2",
		Name:    "Plugin 2",
		Version: "1.0.0",
	}
	plugin2.SetMetadata(metadata2)

	resolver.AddPlugin(plugin1)
	resolver.AddPlugin(plugin2)

	dot := resolver.GenerateDotGraph()

	if len(dot) == 0 {
		t.Error("DOT graph should not be empty")
	}

	// Check for expected content
	if !containsString(dot, "digraph") {
		t.Error("DOT graph should contain 'digraph'")
	}

	if !containsString(dot, "plugin1") {
		t.Error("DOT graph should contain plugin1")
	}

	if !containsString(dot, "plugin2") {
		t.Error("DOT graph should contain plugin2")
	}

	if !containsString(dot, "->") {
		t.Error("DOT graph should contain dependency arrow")
	}
}

func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != "" && substr != "" &&
		(s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr || findSubstr(s, substr))))
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
