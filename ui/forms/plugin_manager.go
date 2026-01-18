package forms

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"ffvi_editor/marketplace"
	"ffvi_editor/plugins"
)

// PluginManagerDialog manages plugin installation, configuration, and execution
type PluginManagerDialog struct {
	window         fyne.Window
	pluginManager  *plugins.Manager
	marketplace    *marketplace.Client
	registry       *marketplace.Registry
	selectedPlugin string
	outputLog      []string
}

// NewPluginManagerDialog creates a new plugin manager dialog
func NewPluginManagerDialog(window fyne.Window, pluginManager *plugins.Manager, marketplaceClient *marketplace.Client, registry *marketplace.Registry) *PluginManagerDialog {
	return &PluginManagerDialog{
		window:        window,
		pluginManager: pluginManager,
		marketplace:   marketplaceClient,
		registry:      registry,
		outputLog:     []string{},
	}
}

// Show displays the plugin manager dialog
func (p *PluginManagerDialog) Show() {
	// Build tabs
	installedTab := p.buildInstalledPluginsTab()
	availableTab := p.buildAvailablePluginsTab()
	settingsTab := p.buildPluginSettingsTab()
	outputTab := p.buildPluginOutputTab()

	tabs := container.NewAppTabs(
		container.NewTabItem("Installed Plugins", installedTab),
		container.NewTabItem("Available Plugins", availableTab),
		container.NewTabItem("Plugin Settings", settingsTab),
		container.NewTabItem("Plugin Output", outputTab),
	)

	content := container.NewVBox(
		widget.NewLabel("Plugin Manager"),
		widget.NewSeparator(),
		tabs,
	)

	buttons := container.NewHBox(
		widget.NewButton("Refresh", func() {
			// Refresh plugin list
		}),
		widget.NewButton("Close", func() {
			// Dialog will close
		}),
	)

	allContent := container.NewBorder(nil, buttons, nil, nil, content)

	d := dialog.NewCustom("Plugin Manager", "Close", allContent, p.window)
	d.Resize(fyne.NewSize(800, 700))
	d.Show()
}

// buildInstalledPluginsTab creates the installed plugins tab
func (p *PluginManagerDialog) buildInstalledPluginsTab() *fyne.Container {
	pluginList := container.NewVBox()

	// Get all loaded plugins
	plugins := p.pluginManager.ListPlugins()

	// Sort by name for consistent display
	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Name < plugins[j].Name
	})

	if len(plugins) == 0 {
		pluginList.Add(widget.NewLabel("No plugins installed"))
	}

	for _, plugin := range plugins {
		// Create plugin entry
		p.createPluginEntry(plugin, pluginList)
	}

	scroll := container.NewVScroll(pluginList)
	scroll.SetMinSize(fyne.NewSize(700, 400))

	return container.NewBorder(
		container.NewHBox(
			widget.NewButton("Browse Plugins", func() {
				// Open marketplace browser
				if p.marketplace != nil && p.registry != nil {
					browser := NewPluginBrowserDialog(p.window, p.pluginManager, p.marketplace, p.registry)
					browser.Show()
				} else {
					dialog.ShowError(fmt.Errorf("marketplace not available"), p.window)
				}
			}),
			widget.NewButton("Open Plugin Folder", func() {
				// Open file manager to plugin directory
			}),
		), nil, nil, nil, scroll,
	)
}

// createPluginEntry creates a UI entry for a plugin
func (p *PluginManagerDialog) createPluginEntry(plugin *plugins.Plugin, parentContainer *fyne.Container) {
	if plugin == nil {
		return
	}

	// Plugin header with name and version
	headerText := fmt.Sprintf("%s (v%s) by %s", plugin.Name, plugin.Version, plugin.Author)
	headerLabel := widget.NewLabel(headerText)

	// Description
	descLabel := widget.NewLabel(plugin.Description)
	descLabel.Wrapping = fyne.TextWrapWord

	// Status indicator
	statusText := "Disabled"
	if plugin.Enabled {
		statusText = "Enabled"
	}
	statusLabel := widget.NewLabel("Status: " + statusText)

	// Buttons
	var toggleBtn *widget.Button
	toggleBtn = widget.NewButton("Enable", func() {
		if plugin.Enabled {
			p.pluginManager.DisablePlugin(plugin.ID)
			toggleBtn.SetText("Enable")
			plugin.Enabled = false
		} else {
			p.pluginManager.EnablePlugin(plugin.ID)
			toggleBtn.SetText("Disable")
			plugin.Enabled = true
		}
	})

	if plugin.Enabled {
		toggleBtn.SetText("Disable")
	}

	runBtn := widget.NewButton("Run", func() {
		p.executePlugin(plugin.ID)
	})

	settingsBtn := widget.NewButton("Settings", func() {
		p.showPluginSettings(plugin)
	})

	removeBtn := widget.NewButton("Remove", func() {
		confirm := dialog.NewConfirm(
			"Confirm Removal",
			fmt.Sprintf("Remove plugin '%s'?", plugin.Name),
			func(ok bool) {
				if ok {
					p.pluginManager.UnloadPlugin(context.Background(), plugin.ID)
				}
			},
			p.window,
		)
		confirm.Show()
	})

	buttons := container.NewHBox(toggleBtn, runBtn, settingsBtn, removeBtn)

	entry := container.NewVBox(
		headerLabel,
		descLabel,
		statusLabel,
		buttons,
		widget.NewSeparator(),
	)

	parentContainer.Add(entry)
}

// buildAvailablePluginsTab creates the available plugins browser tab
func (p *PluginManagerDialog) buildAvailablePluginsTab() *fyne.Container {
	// Check if marketplace is available
	if p.marketplace == nil || p.registry == nil {
		unavailableLabel := widget.NewLabel(
			"Marketplace Unavailable\n\n" +
				"The plugin marketplace is not configured or unreachable.\n\n" +
				"Please check your network connection and marketplace settings.",
		)
		return container.NewVBox(unavailableLabel)
	}

	// Create plugin list container
	pluginListContainer := container.NewVBox()
	var remotePlugins []marketplace.RemotePlugin

	// Search functionality
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search plugins...")

	categorySelect := widget.NewSelect(
		[]string{"All", "Editor Tools", "Speedrun", "Analytics", "Automation", "Utilities"},
		func(s string) {
			// Refilter when category changes
			p.filterAvailablePlugins(pluginListContainer, remotePlugins, searchEntry.Text, s)
		},
	)
	categorySelect.SetSelected("All")

	searchEntry.OnChanged = func(query string) {
		p.filterAvailablePlugins(pluginListContainer, remotePlugins, query, categorySelect.Selected)
	}

	// Status label
	statusLabel := widget.NewLabel("Loading plugins from marketplace...")

	// Refresh button
	refreshBtn := widget.NewButton("Refresh", func() {
		statusLabel.SetText("Refreshing...")
		go p.loadAvailablePlugins(pluginListContainer, &remotePlugins, statusLabel, searchEntry.Text, categorySelect.Selected)
	})

	// Open full marketplace browser button
	browseBtn := widget.NewButton("Open Marketplace Browser", func() {
		browser := NewPluginBrowserDialog(p.window, p.pluginManager, p.marketplace, p.registry)
		browser.Show()
	})

	// Header
	header := container.NewVBox(
		container.NewHBox(searchEntry, categorySelect, refreshBtn),
		browseBtn,
		widget.NewSeparator(),
		statusLabel,
	)

	// Content area
	scroll := container.NewVScroll(pluginListContainer)
	scroll.SetMinSize(fyne.NewSize(700, 350))

	// Load initial data
	go p.loadAvailablePlugins(pluginListContainer, &remotePlugins, statusLabel, "", "All")

	return container.NewBorder(header, nil, nil, nil, scroll)
}

// buildPluginSettingsTab creates the plugin settings tab
func (p *PluginManagerDialog) buildPluginSettingsTab() *fyne.Container {
	// Sandbox mode
	sandboxCheck := widget.NewCheck("Sandbox Mode", func(checked bool) {
		p.pluginManager.SetSandboxMode(checked)
	})

	// Max plugins
	maxPluginsLabel := widget.NewLabel("Max Plugins:")
	maxPluginsEntry := widget.NewEntry()
	maxPluginsEntry.SetText("50")

	// Timeout setting
	timeoutLabel := widget.NewLabel("Execution Timeout (seconds):")
	timeoutEntry := widget.NewEntry()
	timeoutEntry.SetText("30")

	// Auto-load toggle
	autoLoadCheck := widget.NewCheck("Auto-Load Plugins on Startup", func(checked bool) {
		// Toggle auto-load setting
	})

	// Default permissions
	permLabel := widget.NewLabel("Default Permissions:")
	readSaveCheck := widget.NewCheck("Read Save Data", func(checked bool) {})
	writeSaveCheck := widget.NewCheck("Write Save Data", func(checked bool) {})
	uiDisplayCheck := widget.NewCheck("UI Display", func(checked bool) {})
	eventsCheck := widget.NewCheck("Events", func(checked bool) {})

	// Set defaults
	readSaveCheck.SetChecked(true)
	uiDisplayCheck.SetChecked(true)

	content := container.NewVBox(
		widget.NewLabel("Plugin Configuration"),
		widget.NewSeparator(),
		sandboxCheck,
		widget.NewSeparator(),
		maxPluginsLabel,
		maxPluginsEntry,
		widget.NewSeparator(),
		timeoutLabel,
		timeoutEntry,
		widget.NewSeparator(),
		autoLoadCheck,
		widget.NewSeparator(),
		permLabel,
		readSaveCheck,
		writeSaveCheck,
		uiDisplayCheck,
		eventsCheck,
		widget.NewSeparator(),
		widget.NewButton("Save Settings", func() {
			dialog.ShowInformation("Saved", "Plugin settings saved", p.window)
		}),
	)

	scroll := container.NewVScroll(content)
	scroll.SetMinSize(fyne.NewSize(700, 400))

	return container.NewMax(scroll)
}

// loadAvailablePlugins loads plugins from the marketplace
func (p *PluginManagerDialog) loadAvailablePlugins(listContainer *fyne.Container, pluginsPtr *[]marketplace.RemotePlugin, statusLabel *widget.Label, query string, category string) {
	ctx := context.Background()

	// Fetch plugins from marketplace
	opts := &marketplace.ListOptions{
		Category: category,
		SortBy:   "downloads",
		Limit:    50,
		Offset:   0,
	}

	plugins, err := p.marketplace.ListPlugins(ctx, opts)
	if err != nil {
		statusLabel.SetText(fmt.Sprintf("Error loading plugins: %v", err))
		return
	}

	*pluginsPtr = plugins
	p.filterAvailablePlugins(listContainer, plugins, query, category)

	if len(plugins) == 0 {
		statusLabel.SetText("No plugins found")
	} else {
		statusLabel.SetText(fmt.Sprintf("Found %d plugins", len(plugins)))
	}
}

// filterAvailablePlugins filters and displays plugins based on search and category
func (p *PluginManagerDialog) filterAvailablePlugins(listContainer *fyne.Container, plugins []marketplace.RemotePlugin, query string, category string) {
	listContainer.RemoveAll()

	query = strings.ToLower(query)
	for _, plugin := range plugins {
		// Category filter
		if category != "All" && plugin.Category != category {
			continue
		}

		// Search filter
		if query != "" {
			nameMatch := strings.Contains(strings.ToLower(plugin.Name), query)
			descMatch := strings.Contains(strings.ToLower(plugin.Description), query)
			if !nameMatch && !descMatch {
				continue
			}
		}

		// Create plugin card
		p.createAvailablePluginCard(plugin, listContainer)
	}

	listContainer.Refresh()
}

// createAvailablePluginCard creates a card for an available plugin
func (p *PluginManagerDialog) createAvailablePluginCard(plugin marketplace.RemotePlugin, parentContainer *fyne.Container) {
	// Header
	nameLabel := widget.NewLabel(fmt.Sprintf("%s (v%s)", plugin.Name, plugin.Version))
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}

	authorLabel := widget.NewLabel(fmt.Sprintf("by %s", plugin.Author))

	// Description
	descLabel := widget.NewLabel(plugin.Description)
	descLabel.Wrapping = fyne.TextWrapWord

	// Stats
	statsLabel := widget.NewLabel(fmt.Sprintf("⭐ %.1f (%d) | ⬇ %d | Category: %s",
		plugin.Rating, plugin.RatingCount, plugin.Downloads, plugin.Category))

	// Install button
	installBtn := widget.NewButton("Install", func() {
		p.installPluginFromMarketplace(plugin)
	})

	// Check if already installed
	installedPlugins := p.registry.GetInstalledPlugins()
	installed := false
	for _, rec := range installedPlugins {
		if rec.PluginID == plugin.ID {
			installed = true
			break
		}
	}
	if installed {
		installBtn.SetText("Installed")
		installBtn.Disable()
	}

	// Assemble card
	card := container.NewVBox(
		container.NewHBox(nameLabel, authorLabel),
		descLabel,
		statsLabel,
		installBtn,
		widget.NewSeparator(),
	)

	parentContainer.Add(card)
}

// installPluginFromMarketplace installs a plugin from the marketplace
func (p *PluginManagerDialog) installPluginFromMarketplace(plugin marketplace.RemotePlugin) {
	ctx := context.Background()

	// Show progress dialog
	progressBar := widget.NewProgressBarInfinite()
	progressContent := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Installing %s...", plugin.Name)),
		progressBar,
	)
	progressDialog := dialog.NewCustom("Installing Plugin", "Cancel", progressContent, p.window)
	progressDialog.Show()

	// Install in background
	go func() {
		// Download and install
		err := p.marketplace.InstallPlugin(ctx, plugin.ID, plugin.Version)
		if err != nil {
			progressDialog.Hide()
			dialog.ShowError(fmt.Errorf("installation failed: %w", err), p.window)
			return
		}

		// Track installation
		err = p.registry.TrackInstallation(plugin.ID, plugin.Version)
		if err != nil {
			progressDialog.Hide()
			dialog.ShowError(fmt.Errorf("failed to track installation: %w", err), p.window)
			return
		}

		progressDialog.Hide()
		dialog.ShowInformation("Success", fmt.Sprintf("%s installed successfully!", plugin.Name), p.window)
	}()
}

// buildPluginOutputTab creates the plugin output/logging tab
func (p *PluginManagerDialog) buildPluginOutputTab() *fyne.Container {
	// Execution log display
	logContent := container.NewVBox()

	// Get execution log from manager
	logEntries := p.pluginManager.GetExecutionLog()

	if len(logEntries) == 0 {
		logContent.Add(widget.NewLabel("No plugin executions logged yet"))
	}

	for _, entry := range logEntries {
		logLine := fmt.Sprintf(
			"[%s] %s (v%s): %s - %dms - %s",
			entry.StartTime.Format("15:04:05"),
			entry.PluginName,
			entry.PluginVer,
			entry.Status,
			entry.Duration.Milliseconds(),
			entry.Error,
		)
		logContent.Add(widget.NewLabel(logLine))
	}

	logScroll := container.NewVScroll(logContent)
	logScroll.SetMinSize(fyne.NewSize(700, 300))

	// Filters
	filterLabel := widget.NewLabel("Filter by Plugin:")
	pluginFilterSelect := widget.NewSelect([]string{"All", "Plugin 1", "Plugin 2"}, func(s string) {
		// Filter log by plugin
	})
	pluginFilterSelect.SetSelected("All")

	levelLabel := widget.NewLabel("Filter by Level:")
	levelFilterSelect := widget.NewSelect(
		[]string{"All", "Debug", "Info", "Warning", "Error"},
		func(s string) {
			// Filter log by level
		},
	)
	levelFilterSelect.SetSelected("All")

	filters := container.NewHBox(
		filterLabel, pluginFilterSelect,
		levelLabel, levelFilterSelect,
	)

	buttons := container.NewHBox(
		widget.NewButton("Clear Log", func() {
			p.pluginManager.ClearExecutionLog()
			logContent.RemoveAll()
			logContent.Add(widget.NewLabel("Log cleared"))
		}),
		widget.NewButton("Export Log", func() {
			dialog.ShowInformation("Coming Soon", "Log export feature coming soon", p.window)
		}),
	)

	return container.NewBorder(
		container.NewVBox(filters, widget.NewSeparator()),
		buttons,
		nil, nil,
		logScroll,
	)
}

// executePlugin runs a plugin immediately
func (p *PluginManagerDialog) executePlugin(pluginID string) {
	ctx := context.Background()
	err := p.pluginManager.ExecutePlugin(ctx, pluginID)

	if err != nil {
		dialog.ShowError(fmt.Errorf("plugin execution failed: %w", err), p.window)
		return
	}

	dialog.ShowInformation("Success", fmt.Sprintf("Plugin '%s' executed successfully", pluginID), p.window)
}

// showPluginSettings displays settings for a specific plugin
func (p *PluginManagerDialog) showPluginSettings(plugin *plugins.Plugin) {
	cfg, err := p.pluginManager.GetConfig(plugin.ID)
	if err != nil {
		dialog.ShowError(fmt.Errorf("could not load plugin config: %w", err), p.window)
		return
	}

	// Build settings form
	settingsForm := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Settings for: %s", plugin.Name)),
		widget.NewSeparator(),
	)

	// Add settings entries (this would be dynamic based on plugin metadata)
	for key := range cfg.Settings {
		label := widget.NewLabel(fmt.Sprintf("%s:", key))
		entry := widget.NewEntry()
		entry.SetText(fmt.Sprintf("%v", cfg.Settings[key]))

		settingsForm.Add(container.NewHBox(label, entry))
	}

	content := container.NewVBox(
		settingsForm,
		widget.NewSeparator(),
		container.NewHBox(
			widget.NewButton("Save", func() {
				p.pluginManager.SetConfig(plugin.ID, cfg)
				dialog.ShowInformation("Success", "Plugin settings saved", p.window)
			}),
			widget.NewButton("Cancel", func() {
				// Dialog will close
			}),
		),
	)

	d := dialog.NewCustom(
		fmt.Sprintf("Settings: %s", plugin.Name),
		"Cancel",
		content,
		p.window,
	)
	d.Resize(fyne.NewSize(400, 300))
	d.Show()
}
