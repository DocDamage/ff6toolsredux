package plugins

import (
	"fmt"
	"sync"
)

// DependencyResolver resolves plugin dependencies and detects conflicts
type DependencyResolver struct {
	pluginGraph  *PluginDependencyGraph
	versionCache map[string][]*Version
	mu           sync.RWMutex
}

// PluginDependencyGraph represents the dependency relationships between plugins
type PluginDependencyGraph struct {
	nodes map[string]*PluginNode
	edges map[string][]string // pluginID -> list of dependent pluginIDs
	mu    sync.RWMutex
}

// PluginNode represents a plugin in the dependency graph
type PluginNode struct {
	ID           string
	Version      *Version
	Dependencies map[string]*VersionConstraint
}

// ConflictInfo describes a dependency conflict
type ConflictInfo struct {
	PluginID     string
	DependencyID string
	Required     *VersionConstraint
	Actual       *Version
	ConflictType string // "missing", "version_mismatch", "circular"
}

// NewDependencyResolver creates a new dependency resolver
func NewDependencyResolver() *DependencyResolver {
	return &DependencyResolver{
		pluginGraph: &PluginDependencyGraph{
			nodes: make(map[string]*PluginNode),
			edges: make(map[string][]string),
		},
		versionCache: make(map[string][]*Version),
	}
}

// AddPlugin adds a plugin to the dependency graph
func (d *DependencyResolver) AddPlugin(plugin *Plugin) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	version, err := ParseVersion(plugin.Version)
	if err != nil {
		return fmt.Errorf("invalid plugin version: %w", err)
	}

	node := &PluginNode{
		ID:           plugin.ID,
		Version:      version,
		Dependencies: make(map[string]*VersionConstraint),
	}

	// Get dependencies from metadata if available
	metadata := plugin.GetMetadata()
	if len(metadata.Dependencies) > 0 {
		for _, dep := range metadata.Dependencies {
			// Parse dependency string (format: "pluginID@constraint")
			// For now, assume Dependencies is a list of plugin IDs
			// We'll need to enhance this to support version constraints
			// For now, just mark as required with wildcard constraint
			constraint := &VersionConstraint{
				Operator: "*",
			}
			node.Dependencies[dep] = constraint
		}
	}

	d.pluginGraph.mu.Lock()
	d.pluginGraph.nodes[plugin.ID] = node
	d.pluginGraph.mu.Unlock()

	// Build edges
	for depID := range node.Dependencies {
		d.pluginGraph.mu.Lock()
		d.pluginGraph.edges[depID] = append(d.pluginGraph.edges[depID], plugin.ID)
		d.pluginGraph.mu.Unlock()
	}

	// Cache version
	d.versionCache[plugin.ID] = append(d.versionCache[plugin.ID], version)

	return nil
}

// RemovePlugin removes a plugin from the dependency graph
func (d *DependencyResolver) RemovePlugin(pluginID string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.pluginGraph.mu.Lock()
	defer d.pluginGraph.mu.Unlock()

	delete(d.pluginGraph.nodes, pluginID)
	delete(d.pluginGraph.edges, pluginID)

	// Remove from edges where this plugin is a dependent
	for depID, dependents := range d.pluginGraph.edges {
		newDependents := make([]string, 0)
		for _, dep := range dependents {
			if dep != pluginID {
				newDependents = append(newDependents, dep)
			}
		}
		d.pluginGraph.edges[depID] = newDependents
	}

	delete(d.versionCache, pluginID)
}

// ResolveDependencies resolves all dependencies for a plugin
// Returns a map of pluginID -> version string
func (d *DependencyResolver) ResolveDependencies(pluginID string, version string) (map[string]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	node, err := d.getNode(pluginID)
	if err != nil {
		return nil, err
	}

	resolved := make(map[string]string)
	visited := make(map[string]bool)

	if err := d.resolveDependenciesRecursive(node, resolved, visited); err != nil {
		return nil, err
	}

	return resolved, nil
}

// resolveDependenciesRecursive recursively resolves dependencies
func (d *DependencyResolver) resolveDependenciesRecursive(node *PluginNode, resolved map[string]string, visited map[string]bool) error {
	if visited[node.ID] {
		// Already processed
		return nil
	}

	visited[node.ID] = true

	for depID, constraint := range node.Dependencies {
		// Get the dependency node
		depNode, err := d.getNode(depID)
		if err != nil {
			return &ConflictInfo{
				PluginID:     node.ID,
				DependencyID: depID,
				Required:     constraint,
				ConflictType: "missing",
			}
		}

		// Check if version satisfies constraint
		if !constraint.Satisfies(depNode.Version) {
			return &ConflictInfo{
				PluginID:     node.ID,
				DependencyID: depID,
				Required:     constraint,
				Actual:       depNode.Version,
				ConflictType: "version_mismatch",
			}
		}

		// Add to resolved
		resolved[depID] = depNode.Version.String()

		// Recursively resolve transitive dependencies
		if err := d.resolveDependenciesRecursive(depNode, resolved, visited); err != nil {
			return err
		}
	}

	return nil
}

// GetTransitiveDeps gets all transitive dependencies for a plugin
func (d *DependencyResolver) GetTransitiveDeps(pluginID string) map[string]string {
	resolved, err := d.ResolveDependencies(pluginID, "")
	if err != nil {
		return make(map[string]string)
	}
	return resolved
}

// DetectConflicts detects all dependency conflicts in the current graph
func (d *DependencyResolver) DetectConflicts(plugins map[string]string) []ConflictInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()

	conflicts := make([]ConflictInfo, 0)

	// Check each plugin's dependencies
	d.pluginGraph.mu.RLock()
	for pluginID, node := range d.pluginGraph.nodes {
		for depID, constraint := range node.Dependencies {
			depNode, err := d.getNode(depID)
			if err != nil {
				// Dependency missing
				conflicts = append(conflicts, ConflictInfo{
					PluginID:     pluginID,
					DependencyID: depID,
					Required:     constraint,
					ConflictType: "missing",
				})
				continue
			}

			// Check version compatibility
			if !constraint.Satisfies(depNode.Version) {
				conflicts = append(conflicts, ConflictInfo{
					PluginID:     pluginID,
					DependencyID: depID,
					Required:     constraint,
					Actual:       depNode.Version,
					ConflictType: "version_mismatch",
				})
			}
		}
	}
	d.pluginGraph.mu.RUnlock()

	// Detect circular dependencies
	circular := d.detectCircularDependencies()
	for _, cycle := range circular {
		for i := 0; i < len(cycle)-1; i++ {
			conflicts = append(conflicts, ConflictInfo{
				PluginID:     cycle[i],
				DependencyID: cycle[i+1],
				ConflictType: "circular",
			})
		}
	}

	return conflicts
}

// detectCircularDependencies detects circular dependency chains
func (d *DependencyResolver) detectCircularDependencies() [][]string {
	cycles := make([][]string, 0)
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	d.pluginGraph.mu.RLock()
	defer d.pluginGraph.mu.RUnlock()

	for nodeID := range d.pluginGraph.nodes {
		if !visited[nodeID] {
			path := make([]string, 0)
			if cycle := d.detectCycleUtil(nodeID, visited, recStack, path); cycle != nil {
				cycles = append(cycles, cycle)
			}
		}
	}

	return cycles
}

// detectCycleUtil is a helper for DFS-based cycle detection
func (d *DependencyResolver) detectCycleUtil(nodeID string, visited, recStack map[string]bool, path []string) []string {
	visited[nodeID] = true
	recStack[nodeID] = true
	path = append(path, nodeID)

	node, exists := d.pluginGraph.nodes[nodeID]
	if !exists {
		return nil
	}

	for depID := range node.Dependencies {
		if !visited[depID] {
			if cycle := d.detectCycleUtil(depID, visited, recStack, path); cycle != nil {
				return cycle
			}
		} else if recStack[depID] {
			// Found a cycle
			cycleStart := -1
			for i, id := range path {
				if id == depID {
					cycleStart = i
					break
				}
			}
			if cycleStart != -1 {
				return append(path[cycleStart:], depID)
			}
		}
	}

	recStack[nodeID] = false
	return nil
}

// ValidateVersionConstraint validates that a version satisfies a constraint
func (d *DependencyResolver) ValidateVersionConstraint(constraint *VersionConstraint, version string) bool {
	v, err := ParseVersion(version)
	if err != nil {
		return false
	}
	return constraint.Satisfies(v)
}

// GetDependents returns all plugins that depend on the specified plugin
func (d *DependencyResolver) GetDependents(pluginID string) []string {
	d.pluginGraph.mu.RLock()
	defer d.pluginGraph.mu.RUnlock()

	if dependents, ok := d.pluginGraph.edges[pluginID]; ok {
		result := make([]string, len(dependents))
		copy(result, dependents)
		return result
	}

	return make([]string, 0)
}

// GetDependencies returns all direct dependencies of a plugin
func (d *DependencyResolver) GetDependencies(pluginID string) map[string]*VersionConstraint {
	node, err := d.getNode(pluginID)
	if err != nil {
		return make(map[string]*VersionConstraint)
	}

	deps := make(map[string]*VersionConstraint)
	for k, v := range node.Dependencies {
		deps[k] = v
	}
	return deps
}

// getNode safely retrieves a node from the graph
func (d *DependencyResolver) getNode(pluginID string) (*PluginNode, error) {
	d.pluginGraph.mu.RLock()
	defer d.pluginGraph.mu.RUnlock()

	node, ok := d.pluginGraph.nodes[pluginID]
	if !ok {
		return nil, fmt.Errorf("plugin not found in graph: %s", pluginID)
	}
	return node, nil
}

// GenerateDotGraph generates a DOT graph representation for visualization
func (d *DependencyResolver) GenerateDotGraph() string {
	d.pluginGraph.mu.RLock()
	defer d.pluginGraph.mu.RUnlock()

	dot := "digraph PluginDependencies {\n"
	dot += "  rankdir=LR;\n"
	dot += "  node [shape=box];\n\n"

	// Add nodes
	for id, node := range d.pluginGraph.nodes {
		dot += fmt.Sprintf("  \"%s\" [label=\"%s\\nv%s\"];\n", id, id, node.Version.String())
	}

	dot += "\n"

	// Add edges
	for _, node := range d.pluginGraph.nodes {
		for depID, constraint := range node.Dependencies {
			dot += fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\"];\n",
				node.ID, depID, constraint.String())
		}
	}

	dot += "}\n"
	return dot
}

// Error implements error interface for ConflictInfo
func (c ConflictInfo) Error() string {
	switch c.ConflictType {
	case "missing":
		return fmt.Sprintf("plugin %s requires dependency %s (%s) but it is not installed",
			c.PluginID, c.DependencyID, c.Required.String())
	case "version_mismatch":
		return fmt.Sprintf("plugin %s requires %s %s but version %s is installed",
			c.PluginID, c.DependencyID, c.Required.String(), c.Actual.String())
	case "circular":
		return fmt.Sprintf("circular dependency detected: %s -> %s",
			c.PluginID, c.DependencyID)
	default:
		return fmt.Sprintf("dependency conflict in plugin %s", c.PluginID)
	}
}
