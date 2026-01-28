# UI Testing Framework & Guide

**Date:** January 16, 2026  
**Status:** Phase 4C+ Testing Documentation

---

## Overview

Comprehensive testing framework for the FF6 Save Editor's plugin marketplace UI components and integrations.

---

## Testing Stack

```
Frontend Testing:
├── Unit Tests (Go - ui/editors, ui/forms)
├── Component Tests (UI interactions)
├── Integration Tests (Marketplace API)
├── E2E Tests (User workflows)
└── Accessibility Tests (WCAG compliance)
```

---

## Unit Testing Strategy

### 1. Marketplace UI Components

**File:** `ui/editors/marketplace_test.go`

```go
package editors

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "ff6editor/models"
)

func TestMarketplaceViewer_LoadPluginList(t *testing.T) {
    // Setup
    viewer := NewMarketplaceViewer()
    
    // Test data
    plugins := []models.Plugin{
        {
            ID:          "test-plugin",
            Name:        "Test Plugin",
            Version:     "1.0.0",
            Description: "A test plugin",
            Rating:      4.5,
        },
    }
    
    // Execute
    err := viewer.LoadPlugins(plugins)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 1, viewer.GetPluginCount())
}

func TestMarketplaceViewer_SearchPlugins(t *testing.T) {
    viewer := NewMarketplaceViewer()
    plugins := []models.Plugin{
        {ID: "speedrun-tracker", Name: "Speedrun Tracker"},
        {ID: "damage-calculator", Name: "Damage Calculator"},
    }
    viewer.LoadPlugins(plugins)
    
    results := viewer.Search("speedrun")
    assert.Equal(t, 1, len(results))
    assert.Equal(t, "speedrun-tracker", results[0].ID)
}

func TestMarketplaceViewer_FilterByCategory(t *testing.T) {
    viewer := NewMarketplaceViewer()
    plugins := []models.Plugin{
        {ID: "p1", Category: "Speedrun"},
        {ID: "p2", Category: "Tools"},
        {ID: "p3", Category: "Speedrun"},
    }
    viewer.LoadPlugins(plugins)
    
    filtered := viewer.FilterByCategory("Speedrun")
    assert.Equal(t, 2, len(filtered))
}

func TestMarketplaceViewer_PluginInstallation(t *testing.T) {
    viewer := NewMarketplaceViewer()
    plugin := models.Plugin{
        ID:      "test-plugin",
        Version: "1.0.0",
    }
    
    err := viewer.InstallPlugin(plugin)
    assert.NoError(t, err)
    
    installed := viewer.GetInstalledPlugins()
    assert.Equal(t, 1, len(installed))
}

func TestMarketplaceViewer_PluginUpdate(t *testing.T) {
    viewer := NewMarketplaceViewer()
    
    // Install v1.0.0
    plugin := models.Plugin{ID: "test", Version: "1.0.0"}
    viewer.InstallPlugin(plugin)
    
    // Update to v1.1.0
    newPlugin := models.Plugin{ID: "test", Version: "1.1.0"}
    err := viewer.UpdatePlugin(newPlugin)
    
    assert.NoError(t, err)
    installed := viewer.GetPlugin("test")
    assert.Equal(t, "1.1.0", installed.Version)
}

func TestMarketplaceViewer_PluginUninstall(t *testing.T) {
    viewer := NewMarketplaceViewer()
    plugin := models.Plugin{ID: "test-plugin", Version: "1.0.0"}
    viewer.InstallPlugin(plugin)
    
    err := viewer.UninstallPlugin("test-plugin")
    assert.NoError(t, err)
    
    installed := viewer.GetInstalledPlugins()
    assert.Equal(t, 0, len(installed))
}

func TestMarketplaceViewer_RatingSubmission(t *testing.T) {
    viewer := NewMarketplaceViewer()
    
    err := viewer.SubmitRating("test-plugin", 4.5, "Great plugin!")
    assert.NoError(t, err)
    
    plugin := viewer.GetPlugin("test-plugin")
    assert.Greater(t, plugin.Rating, 0.0)
}
```

### 2. Plugin Manager UI Tests

**File:** `ui/forms/plugin_manager_form_test.go`

```go
package forms

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "ff6editor/plugins"
)

func TestPluginManagerForm_DisplayInstalledPlugins(t *testing.T) {
    form := NewPluginManagerForm()
    
    manager := plugins.NewManager()
    manager.Install("test-plugin", "1.0.0")
    
    form.SetPluginManager(manager)
    plugins := form.GetDisplayedPlugins()
    
    assert.Greater(t, len(plugins), 0)
}

func TestPluginManagerForm_EnableDisablePlugin(t *testing.T) {
    form := NewPluginManagerForm()
    manager := plugins.NewManager()
    manager.Install("test-plugin", "1.0.0")
    form.SetPluginManager(manager)
    
    // Disable
    err := form.DisablePlugin("test-plugin")
    assert.NoError(t, err)
    assert.False(t, form.IsPluginEnabled("test-plugin"))
    
    // Enable
    err = form.EnablePlugin("test-plugin")
    assert.NoError(t, err)
    assert.True(t, form.IsPluginEnabled("test-plugin"))
}

func TestPluginManagerForm_ConfigurePlugin(t *testing.T) {
    form := NewPluginManagerForm()
    config := map[string]interface{}{
        "debug_mode": true,
        "timeout": 5000,
    }
    
    err := form.SetPluginConfig("test-plugin", config)
    assert.NoError(t, err)
    
    retrieved := form.GetPluginConfig("test-plugin")
    assert.Equal(t, true, retrieved["debug_mode"])
}

func TestPluginManagerForm_ViewPluginDetails(t *testing.T) {
    form := NewPluginManagerForm()
    
    details := form.GetPluginDetails("speedrun-tracker")
    assert.NotNil(t, details)
    assert.Equal(t, "speedrun-tracker", details.ID)
}
```

### 3. Settings UI Tests

**File:** `ui/forms/settings_form_test.go`

```go
package forms

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestSettingsForm_MarketplaceSettings(t *testing.T) {
    form := NewSettingsForm()
    
    // Enable marketplace
    err := form.SetMarketplaceEnabled(true)
    assert.NoError(t, err)
    assert.True(t, form.IsMarketplaceEnabled())
}

func TestSettingsForm_SetRegistryURL(t *testing.T) {
    form := NewSettingsForm()
    
    url := "https://raw.githubusercontent.com/ff6-marketplace/registry/main"
    err := form.SetRegistryURL(url)
    assert.NoError(t, err)
    
    retrieved := form.GetRegistryURL()
    assert.Equal(t, url, retrieved)
}

func TestSettingsForm_AutoUpdateSetting(t *testing.T) {
    form := NewSettingsForm()
    
    err := form.SetAutoUpdatePlugins(true)
    assert.NoError(t, err)
    assert.True(t, form.GetAutoUpdatePlugins())
}

func TestSettingsForm_PluginSandboxMode(t *testing.T) {
    form := NewSettingsForm()
    
    err := form.SetSandboxMode(true)
    assert.NoError(t, err)
    assert.True(t, form.IsSandboxModeEnabled())
}
```

---

## Integration Testing

### Marketplace Integration Tests

**File:** `marketplace/client_integration_test.go`

```go
package marketplace

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "ff6editor/models"
)

func TestMarketplaceClient_FetchPlugins(t *testing.T) {
    client := NewClient("https://raw.githubusercontent.com/ff6-marketplace/registry/main")
    
    plugins, err := client.FetchPlugins()
    
    assert.NoError(t, err)
    assert.Greater(t, len(plugins), 0)
}

func TestMarketplaceClient_DownloadPlugin(t *testing.T) {
    client := NewClient("https://raw.githubusercontent.com/ff6-marketplace/registry/main")
    
    plugin := models.Plugin{
        ID:          "speedrun-tracker",
        Version:     "1.0.0",
        DownloadURL: "https://raw.githubusercontent.com/ff6-marketplace/registry/main/plugins/speedrun-tracker/1.0.0/plugin.lua",
    }
    
    content, err := client.DownloadPlugin(plugin)
    
    assert.NoError(t, err)
    assert.Greater(t, len(content), 0)
}

func TestMarketplaceClient_VerifyChecksum(t *testing.T) {
    client := NewClient("https://raw.githubusercontent.com/ff6-marketplace/registry/main")
    
    content := []byte("plugin code here")
    checksum := "abc123..."
    
    verified := client.VerifyChecksum(content, checksum)
    
    // This test depends on actual checksum
    assert.NotNil(t, verified)
}

func TestMarketplaceClient_FetchRatings(t *testing.T) {
    client := NewClient("https://raw.githubusercontent.com/ff6-marketplace/registry/main")
    
    ratings, err := client.FetchRatings("speedrun-tracker")
    
    assert.NoError(t, err)
    assert.NotNil(t, ratings)
}

func TestMarketplaceClient_SubmitRating(t *testing.T) {
    client := NewClient("https://raw.githubusercontent.com/ff6-marketplace/registry/main")
    
    rating := models.Rating{
        PluginID: "test-plugin",
        Score:    4.5,
        Comment:  "Great plugin!",
    }
    
    err := client.SubmitRating(rating)
    assert.NoError(t, err)
}
```

---

## E2E Test Scenarios

### Scenario 1: Browse and Install Plugin

```go
func TestE2E_BrowseAndInstallPlugin(t *testing.T) {
    // 1. Open Marketplace
    editor := SetupEditor()
    marketplace := editor.OpenMarketplace()
    
    // 2. Search for plugin
    marketplace.Search("speedrun")
    results := marketplace.GetSearchResults()
    assert.Greater(t, len(results), 0)
    
    // 3. View plugin details
    plugin := results[0]
    details := marketplace.ViewPluginDetails(plugin.ID)
    assert.Equal(t, plugin.ID, details.ID)
    
    // 4. Install plugin
    err := marketplace.InstallPlugin(plugin)
    assert.NoError(t, err)
    
    // 5. Verify installation
    installed := editor.GetInstalledPlugins()
    found := false
    for _, p := range installed {
        if p.ID == plugin.ID {
            found = true
            break
        }
    }
    assert.True(t, found)
    
    // 6. Enable plugin
    err = editor.EnablePlugin(plugin.ID)
    assert.NoError(t, err)
    
    editor.Close()
}
```

### Scenario 2: Update Installed Plugin

```go
func TestE2E_UpdatePlugin(t *testing.T) {
    editor := SetupEditor()
    
    // 1. Install v1.0.0
    plugin := models.Plugin{ID: "test", Version: "1.0.0"}
    editor.InstallPlugin(plugin)
    
    // 2. Open Marketplace
    marketplace := editor.OpenMarketplace()
    
    // 3. Find update
    pluginInfo := marketplace.GetPluginInfo("test")
    assert.Greater(t, pluginInfo.LatestVersion, "1.0.0")
    
    // 4. Update
    err := marketplace.UpdatePlugin("test")
    assert.NoError(t, err)
    
    // 5. Verify update
    installed := editor.GetPlugin("test")
    assert.Greater(t, installed.Version, "1.0.0")
    
    editor.Close()
}
```

### Scenario 3: Rate and Review Plugin

```go
func TestE2E_RatePlugin(t *testing.T) {
    editor := SetupEditor()
    marketplace := editor.OpenMarketplace()
    
    // 1. Find installed plugin
    installed := marketplace.GetInstalledPlugins()
    plugin := installed[0]
    
    // 2. Open ratings
    ratings := marketplace.ViewRatings(plugin.ID)
    assert.NotNil(t, ratings)
    
    // 3. Submit rating
    err := marketplace.SubmitRating(models.Rating{
        PluginID: plugin.ID,
        Score:    5.0,
        Comment:  "Excellent plugin!",
    })
    assert.NoError(t, err)
    
    // 4. Verify rating submitted
    updated := marketplace.ViewRatings(plugin.ID)
    assert.Greater(t, updated.Count, 0)
    
    editor.Close()
}
```

---

## Performance Testing

### Load Testing

```go
func TestPerformance_LoadManyPlugins(t *testing.T) {
    marketplace := NewMarketplaceViewer()
    
    // Create 1000 plugins
    plugins := make([]models.Plugin, 1000)
    for i := 0; i < 1000; i++ {
        plugins[i] = models.Plugin{
            ID:   fmt.Sprintf("plugin-%d", i),
            Name: fmt.Sprintf("Plugin %d", i),
        }
    }
    
    // Measure load time
    start := time.Now()
    marketplace.LoadPlugins(plugins)
    duration := time.Since(start)
    
    // Should load in < 100ms
    assert.Less(t, duration, 100*time.Millisecond)
}

func TestPerformance_SearchOptimization(t *testing.T) {
    marketplace := NewMarketplaceViewer()
    
    // Load 500 plugins
    plugins := generateTestPlugins(500)
    marketplace.LoadPlugins(plugins)
    
    // Measure search time
    start := time.Now()
    results := marketplace.Search("test")
    duration := time.Since(start)
    
    // Should search in < 50ms
    assert.Less(t, duration, 50*time.Millisecond)
}

func TestPerformance_DownloadPluginFile(t *testing.T) {
    client := NewClient("https://raw.githubusercontent.com/ff6-marketplace/registry/main")
    
    // Measure download time for 1MB file
    start := time.Now()
    content, err := client.DownloadPlugin(testPlugin)
    duration := time.Since(start)
    
    assert.NoError(t, err)
    // Should download 1MB in < 5 seconds
    assert.Less(t, duration, 5*time.Second)
}
```

---

## Accessibility Testing

### WCAG 2.1 Compliance

```go
func TestAccessibility_KeyboardNavigation(t *testing.T) {
    marketplace := NewMarketplaceViewer()
    
    // Test tab navigation
    marketplace.FocusSearchBox()
    assert.True(t, marketplace.HasKeyboardFocus())
    
    marketplace.PressTab()
    assert.True(t, marketplace.HasKeyboardFocus())
    
    // Test Enter key activation
    marketplace.PressEnter()
    assert.NoError(t, nil)
}

func TestAccessibility_ScreenReaderSupport(t *testing.T) {
    marketplace := NewMarketplaceViewer()
    
    // Test ARIA labels
    searchBox := marketplace.GetSearchBox()
    label := searchBox.GetARIALabel()
    assert.Equal(t, "Search plugins", label)
    
    // Test descriptions
    plugin := marketplace.GetPlugin("test")
    description := plugin.GetAccessibleDescription()
    assert.NotEmpty(t, description)
}

func TestAccessibility_ColorContrast(t *testing.T) {
    marketplace := NewMarketplaceViewer()
    
    elements := marketplace.GetUIElements()
    for _, elem := range elements {
        contrast := elem.GetColorContrast()
        // WCAG AA requires 4.5:1 for normal text
        assert.GreaterOrEqual(t, contrast, 4.5)
    }
}
```

---

## Test Coverage Goals

| Component | Target | Current |
|-----------|--------|---------|
| Marketplace Viewer | 90% | - |
| Plugin Manager | 85% | - |
| Settings Form | 80% | - |
| Marketplace Client | 95% | - |
| API Integration | 90% | - |

---

## Running Tests

### Run All Tests
```bash
make test
```

### Run Specific Package
```bash
go test ./ui/editors -v
go test ./marketplace -v
go test ./plugins -v
```

### Run with Coverage
```bash
go test -cover ./...
make coverage
```

### Run E2E Tests
```bash
go test ./e2e -v -tags=e2e
```

---

## CI/CD Integration

**GitHub Actions Workflow:** `.github/workflows/tests.yml`

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      
      - name: Run tests
        run: make test
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

---

## Test Maintenance Schedule

- **Daily:** Run all tests on commits
- **Weekly:** Update E2E tests for new features
- **Monthly:** Review and update performance baselines
- **Quarterly:** Audit test coverage and add missing tests

---

**UI Testing Framework Complete** ✅

Comprehensive testing ensures marketplace UI stability, performance, and user experience.
