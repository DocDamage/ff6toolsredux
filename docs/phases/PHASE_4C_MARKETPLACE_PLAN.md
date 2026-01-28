# Phase 4C: Marketplace System — Complete Implementation Plan

**Date:** January 16, 2026  
**Phase:** 4C (Marketplace System)  
**Status:** Planning & Implementation  
**Version:** 1.0  
**Build Status:** Go 1.25.6, Fyne v2.5.2 (ready for extension)

---

## Executive Summary

Phase 4C builds a **community-driven plugin marketplace** enabling users to discover, download, install, rate, and review community-created plugins. The marketplace consists of:

1. **Cloud Registry** - Central repository for plugin metadata (GitHub-hosted)
2. **Client API** - Local marketplace client for discovery and installation
3. **Browser UI** - User-friendly plugin browser with search, filters, ratings
4. **Version Management** - Automatic updates and upgrade detection
5. **Rating System** - Community feedback on plugin quality

**Scope:** ~2,000-2,500 lines of Go + UI code  
**Timeline:** 3-4 development sessions  
**Integration Points:** Plugin Manager, Marketplace Client, Registry Storage

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    FF6 Save Editor 4.1.0                    │
│                  (Plugin Manager + Marketplace)              │
└──────────────────────┬──────────────────────────────────────┘
                       │
        ┌──────────────┴──────────────┐
        │                             │
   ┌────▼────┐              ┌────────▼──────────┐
   │  Local  │              │   Plugin Browser   │
   │ Plugin  │              │  UI (Fyne)         │
   │ Manager │              │  - Search          │
   │         │              │  - Filter          │
   └────┬────┘              │  - Install         │
        │                   │  - View Ratings    │
        │                   └────────┬───────────┘
        │                            │
        └──────────────┬─────────────┘
                       │
          ┌────────────▼─────────────┐
          │  Marketplace Client API  │
          │  marketplace/client.go   │
          │  - Registry discovery    │
          │  - Plugin listing        │
          │  - Download/Install      │
          │  - Rating submission     │
          └────────────┬─────────────┘
                       │
        ┌──────────────▼──────────────┐
        │  Plugin Registry Management  │
        │  marketplace/registry.go     │
        │  - Cache metadata           │
        │  - Track installations      │
        │  - Version history          │
        │  - Update detection         │
        └────────────┬─────────────────┘
                     │
        ┌────────────▼──────────────┐
        │  Cloud Registry (GitHub)   │
        │  - Plugin packages         │
        │  - Metadata (JSON)         │
        │  - Download URLs           │
        │  - Rating aggregates       │
        └────────────────────────────┘
```

---

## Phase 4C Deliverables

### 1. Backend: Marketplace Client API
**File:** `marketplace/client.go` (~600 lines)

```go
// Core types
type Marketplace struct {
    baseURL      string              // Registry base URL
    client       *http.Client        // HTTP client for API calls
    registry     *Registry           // Local registry
    cache        map[string]interface{} // Response cache
    cacheTTL     time.Duration       // Cache lifetime
}

type RemotePlugin struct {
    ID            string
    Name          string
    Description   string
    Author        string
    Version       string
    MinVersion    string           // Min FF6Editor version
    MaxVersion    string           // Max FF6Editor version
    Category      string           // Tool, Utility, Game, Community
    Tags          []string
    Rating        float32
    RatingCount   int
    Downloads     int
    License       string
    RepositoryURL string
    DocumentationURL string
    UpdatedAt     time.Time
    CreatedAt     time.Time
}

type PluginRating struct {
    PluginID    string
    UserID      string            // Anonymous hash
    Rating      float32           // 1-5 stars
    Review      string
    Timestamp   time.Time
}

// API methods
func (m *Marketplace) ListPlugins(ctx context.Context, opts ListOptions) ([]RemotePlugin, error)
func (m *Marketplace) SearchPlugins(ctx context.Context, query string) ([]RemotePlugin, error)
func (m *Marketplace) GetPluginDetails(ctx context.Context, pluginID string) (*RemotePlugin, error)
func (m *Marketplace) DownloadPlugin(ctx context.Context, pluginID, version string) ([]byte, error)
func (m *Marketplace) InstallPlugin(ctx context.Context, pluginID, version string) error
func (m *Marketplace) SubmitRating(ctx context.Context, rating *PluginRating) error
func (m *Marketplace) GetPluginRatings(ctx context.Context, pluginID string) ([]*PluginRating, error)
func (m *Marketplace) CheckForUpdates(ctx context.Context) (map[string]*UpdateInfo, error)
```

**Features:**
- ✅ List/search plugins from cloud registry
- ✅ Filter by category, tags, compatibility
- ✅ Download and cache plugin packages
- ✅ Submit ratings and reviews
- ✅ Automatic update detection
- ✅ Version compatibility checking
- ✅ Fallback to local cache on network error
- ✅ Rate limiting and throttling

### 2. Backend: Plugin Registry Management
**File:** `marketplace/registry.go` (~400 lines)

```go
type Registry struct {
    db                 *json.Decoder          // Local registry database
    installedPlugins   map[string]*InstallRecord
    ratings            map[string][]*PluginRating
    updateChecks       map[string]time.Time   // Last check timestamps
}

type InstallRecord struct {
    PluginID      string
    Version       string
    InstalledAt   time.Time
    UpdatedAt     time.Time
    Enabled       bool
    AutoUpdate    bool
    Location      string
}

type UpdateInfo struct {
    PluginID      string
    CurrentVersion string
    LatestVersion  string
    ReleaseNotes  string
    UpdatedAt     time.Time
}

// Methods
func (r *Registry) GetInstalledPlugins() []InstallRecord
func (r *Registry) TrackInstallation(pluginID, version, location string) error
func (r *Registry) CheckUpdates(plugins []string) (map[string]*UpdateInfo, error)
func (r *Registry) MarkPluginUpdated(pluginID, newVersion string) error
func (r *Registry) SaveRatings(ratings []*PluginRating) error
func (r *Registry) LoadRatings(pluginID string) ([]*PluginRating, error)
func (r *Registry) GetCachedPlugins() ([]RemotePlugin, error)
func (r *Registry) CachePluginList(plugins []RemotePlugin) error
```

**Features:**
- ✅ Track locally installed plugins
- ✅ Store version history
- ✅ Cache rating/review data locally
- ✅ Detect available updates
- ✅ Support auto-update setting per plugin
- ✅ Version rollback support
- ✅ JSON-based persistence

### 3. UI: Plugin Browser Dialog
**File:** `ui/forms/plugin_browser.go` (~800 lines)

```go
type PluginBrowserDialog struct {
    window              fyne.Window
    searchEntry         *widget.Entry
    categorySelect      *widget.Select
    sortSelect          *widget.Select
    pluginList          *container.VBox
    pluginScroll        *container.Scroll
    detailsPanel        *fyne.Container
    installButton       *widget.Button
    statusLabel         *widget.Label
    downloadProgress    *widget.ProgressBar
    ratingStars         *widget.RichText
    reviewList          *container.VBox
    selectedPlugin      *RemotePlugin
    marketplace         *Marketplace
    pluginManager       *plugins.Manager
}

// Core methods
func NewPluginBrowserDialog(window fyne.Window, mgr *plugins.Manager) *PluginBrowserDialog
func (d *PluginBrowserDialog) Show() dialog.Dialog
func (d *PluginBrowserDialog) refreshPluginList() error
func (d *PluginBrowserDialog) searchPlugins(query string) error
func (d *PluginBrowserDialog) selectPlugin(plugin *RemotePlugin)
func (d *PluginBrowserDialog) installPlugin() error
func (d *PluginBrowserDialog) submitRating(rating float32, review string) error
func (d *PluginBrowserDialog) checkForUpdates() error
func (d *PluginBrowserDialog) displayPluginDetails(plugin *RemotePlugin)
```

**UI Components:**
- **Header:** Search bar, category filter, sort dropdown
- **Plugin List:** Scrollable list with name, author, rating, download count
- **Plugin Details Panel:**
  - Name, author, version
  - Description and documentation link
  - Rating (5-star display)
  - User reviews
  - Version history
  - Install button
- **Installation Progress:**
  - Download progress bar
  - Status messages
  - Cancel button
- **Update Check:** Display available updates with upgrade buttons

### 4. Cloud Registry Structure (GitHub-based)
**Location:** GitHub repository (separate or integrated)

```
gh:ff6-marketplace/registry/
├── plugins.json                    # Master registry metadata
├── ratings.json                    # Community ratings
├── categories.json                 # Plugin categories
└── plugins/
    ├── plugin-id-1/
    │   ├── 1.0.0/
    │   │   ├── plugin.lua          # Plugin code
    │   │   ├── metadata.json       # Version metadata
    │   │   └── changelog.md        # Version notes
    │   └── 1.0.1/
    │       └── ...
    └── plugin-id-2/
        └── ...
```

**Master Registry (plugins.json):**
```json
{
  "plugins": [
    {
      "id": "speedrun-tracker",
      "name": "Speedrun Tracker Pro",
      "description": "Advanced speedrun metrics tracking",
      "author": "CommunityDev",
      "latestVersion": "2.1.0",
      "minVersion": "4.1.0",
      "maxVersion": "*",
      "category": "Tools",
      "tags": ["speedrun", "metrics", "timing"],
      "rating": 4.8,
      "ratingCount": 45,
      "downloads": 1250,
      "license": "MIT",
      "repositoryURL": "https://github.com/user/speedrun-tracker",
      "documentationURL": "https://docs.example.com/speedrun-tracker",
      "versions": [
        {
          "version": "2.1.0",
          "releaseDate": "2026-01-15",
          "downloadURL": "https://github.com/ff6-marketplace/registry/raw/main/plugins/speedrun-tracker/2.1.0/plugin.lua",
          "size": 45300,
          "checksum": "sha256:abc123..."
        },
        {
          "version": "2.0.5",
          "releaseDate": "2026-01-10",
          "downloadURL": "..."
        }
      ],
      "createdAt": "2025-12-01",
      "updatedAt": "2026-01-15"
    }
  ]
}
```

---

## Implementation Plan

### Phase 4C-1: Backend Foundation (Session 1)
1. ✅ Create `marketplace/client.go` with HTTP API methods
2. ✅ Implement plugin discovery and listing
3. ✅ Add download functionality with caching
4. ✅ Create error handling and retry logic

### Phase 4C-2: Registry Management (Session 1-2)
1. ✅ Create `marketplace/registry.go` for local state
2. ✅ Implement version tracking
3. ✅ Add update detection logic
4. ✅ Create JSON persistence layer

### Phase 4C-3: UI Browser (Session 2-3)
1. ✅ Create `ui/forms/plugin_browser.go` dialog
2. ✅ Build plugin list display
3. ✅ Implement search and filtering
4. ✅ Add installation workflow
5. ✅ Display ratings and reviews

### Phase 4C-4: Integration & Testing (Session 3-4)
1. ✅ Wire marketplace into main window menu
2. ✅ Connect with existing plugin manager
3. ✅ Write comprehensive test suite
4. ✅ Create documentation

### Phase 4C-5: Cloud Registry Setup
1. ✅ Create GitHub repository for plugin registry
2. ✅ Set up initial plugins.json master file
3. ✅ Document plugin submission process

---

## Data Models

### RemotePlugin
```go
type RemotePlugin struct {
    ID               string        // Unique plugin ID
    Name             string        // Display name
    Description      string        // Short description
    FullDescription  string        // Long description
    Author           string        // Author name
    Version          string        // Latest version (semver)
    MinVersion       string        // Min editor version
    MaxVersion       string        // Max editor version
    Category         string        // Tool|Utility|Game|Community
    Tags             []string      // Search tags
    Rating           float32       // Average rating (1-5)
    RatingCount      int           // Number of ratings
    Downloads        int           // Download count
    License          string        // License type (MIT, Apache, GPL, etc)
    RepositoryURL    string        // Link to source repo
    DocumentationURL string        // Link to docs
    IconURL          string        // Plugin icon/avatar
    VersionHistory   []VersionInfo // Previous versions
    UpdatedAt        time.Time     // Last update time
    CreatedAt        time.Time     // Creation time
}

type VersionInfo struct {
    Version      string    // Semantic version
    ReleaseDate  time.Time // Release date
    DownloadURL  string    // Direct download link
    Size         int64     // File size in bytes
    Checksum     string    // SHA256 checksum
    ReleaseNotes string    // Changelog/notes
}

type PluginRating struct {
    ID          string    // Unique rating ID
    PluginID    string    // Plugin being rated
    UserID      string    // Anonymous user hash (for privacy)
    Rating      float32   // 1.0-5.0 stars
    Review      string    // Optional text review
    Helpful     int       // Helpful count
    Timestamp   time.Time // When submitted
}
```

### ListOptions
```go
type ListOptions struct {
    Category     string   // Filter by category
    Tags         []string // Filter by tags
    SortBy       string   // rating|downloads|recent|name
    SortOrder    string   // asc|desc
    Limit        int      // Results per page
    Offset       int      // Pagination offset
    MinRating    float32  // Minimum rating filter
    OnlyUpdates  bool     // Only show available updates
}
```

---

## API Endpoints (Cloud Registry)

```
GET  /api/v1/plugins                    → List all plugins
GET  /api/v1/plugins/search?q=...       → Search plugins
GET  /api/v1/plugins/:id                → Plugin details
GET  /api/v1/plugins/:id/versions       → Version history
GET  /api/v1/plugins/:id/ratings        → Plugin ratings/reviews
GET  /api/v1/plugins/:id/:version/download → Download plugin
POST /api/v1/plugins/:id/ratings        → Submit rating
GET  /api/v1/categories                 → Available categories
GET  /api/v1/stats                      → Marketplace statistics
```

---

## Marketplace Browser Features

### 1. Plugin Discovery
- **Search:** Full-text search across name, description, author
- **Categories:** Filter by plugin type (Tools, Utilities, Game, Community)
- **Tags:** Multi-select tag filtering
- **Sorting:** By rating, downloads, recently updated, alphabetical
- **Pagination:** Load plugins in batches of 20

### 2. Plugin Details
- **Metadata:** Name, author, version, category, license
- **Rating:** 5-star display with review count
- **Description:** Full description with markdown support
- **Version History:** Display past versions with changelogs
- **Installation:** One-click install with progress tracking
- **Links:** Documentation, repository, author profile

### 3. Rating System
- **Display:** Current average rating and vote count
- **Submission:** Rate and review plugins
- **Review List:** Display top/recent reviews
- **Helpful:** Mark reviews as helpful/unhelpful

### 4. Installation Management
- **One-Click Install:** Download and enable plugin
- **Auto-Update:** Optional automatic update checking
- **Update Available:** Show notification when updates available
- **Version Selection:** Allow downgrading to previous versions
- **Uninstall:** Remove plugin with confirmation

### 5. Status & Notifications
- **Download Progress:** Show download percentage
- **Installation Status:** Message updates during install
- **Update Notifications:** Alert when updates available
- **Error Handling:** Clear error messages with troubleshooting

---

## Integration with Existing Components

### Plugin Manager (`plugins/manager.go`)
- New method: `InstallFromMarketplace(pluginID, version string) error`
- New method: `UninstallPlugin(pluginID string) error`
- New method: `UpdatePlugin(pluginID string) error`
- Existing methods remain unchanged

### Main Window (`ui/window.go`)
- Add menu item: **Tools → Marketplace Browser**
- Callback: Open marketplace browser dialog
- Status bar update: Show marketplace sync status

### Startup Initialization (`main.go`)
- Initialize marketplace client on app startup
- Background check for plugin updates (optional)
- Display update notifications in status bar

---

## Testing Strategy

### Unit Tests (`marketplace/marketplace_test.go`)
```go
// Client tests
func TestMarketplaceListPlugins(t *testing.T)
func TestMarketplaceSearchPlugins(t *testing.T)
func TestMarketplaceDownloadPlugin(t *testing.T)
func TestMarketplaceInstallPlugin(t *testing.T)
func TestMarketplaceSubmitRating(t *testing.T)
func TestMarketplaceCheckUpdates(t *testing.T)
func TestMarketplaceErrorHandling(t *testing.T)
func TestMarketplaceCaching(t *testing.T)

// Registry tests
func TestRegistryTrackInstallation(t *testing.T)
func TestRegistryCheckUpdates(t *testing.T)
func TestRegistryVersionComparison(t *testing.T)
func TestRegistryPersistence(t *testing.T)

// Integration tests
func TestEndToEndInstallation(t *testing.T)
func TestEndToEndRating(t *testing.T)
func TestEndToEndUpdate(t *testing.T)
```

### Mock/Stub Server
- HTTP mock server for registry endpoints
- Test data fixtures (sample plugins, ratings)
- Error simulation (404, 500, timeout)

---

## Configuration

### Editor Config (`ff6editor.config`)
```ini
[marketplace]
enabled=true
registryURL=https://ff6-marketplace.github.io/registry/api/v1
cachePath=~/.ff6editor/marketplace/cache
cacheExpiry=24h
autoCheckUpdates=true
updateCheckInterval=12h
maxConcurrentDownloads=3
```

### Defaults
- Registry URL: GitHub-hosted registry
- Cache expiry: 24 hours
- Update check: Every 12 hours
- Max downloads: 3 concurrent

---

## Security Considerations

1. **HTTPS Only:** All marketplace communication over HTTPS
2. **Checksum Verification:** Verify plugin downloads with SHA256
3. **Version Compatibility:** Prevent installing incompatible versions
4. **Rating Privacy:** Hash user IDs to protect anonymity
5. **Plugin Isolation:** Marketplace plugins run in same sandbox as local plugins
6. **Rate Limiting:** Implement client-side throttling

---

## Error Handling

### Network Errors
```
- Connection timeout → Retry with backoff, show offline message
- DNS failure → Use cached data if available
- HTTP 404 → Plugin not found
- HTTP 500 → Server error, retry later
- Rate limited → Implement exponential backoff
```

### Installation Errors
```
- Checksum mismatch → Retry or report corruption
- Incompatible version → Show version requirement
- Disk space → Check available space before download
- Permission denied → Show permission error
- Plugin exists → Offer replace option
```

---

## Documentation Files

### PHASE_4C_MARKETPLACE_GUIDE.md
- User guide for discovering and installing plugins
- Step-by-step marketplace browser walkthrough
- How to rate and review plugins
- Update notification handling
- Troubleshooting common issues

### PHASE_4C_API_REFERENCE.md
- Cloud registry API documentation
- Client library reference
- Registry storage format
- Installation workflow diagram
- Version compatibility matrix

### PHASE_4C_PLUGIN_SUBMISSION_GUIDE.md
- How to submit plugins to marketplace
- Plugin metadata requirements
- Version tagging conventions
- Category selection guidelines
- Documentation requirements

### PHASE_4C_QUICK_REFERENCE.md
- Marketplace commands
- Configuration options
- Keyboard shortcuts
- Common workflows
- Performance metrics

---

## Success Criteria

- ✅ Marketplace client successfully fetches plugin list
- ✅ Plugin browser displays 20+ plugins with full details
- ✅ Search returns relevant results in <200ms
- ✅ One-click installation works seamlessly
- ✅ Rating/review submission functional
- ✅ Update detection identifies available upgrades
- ✅ All unit tests passing (20+ tests)
- ✅ Zero compilation errors
- ✅ Documentation complete (3,000+ lines)
- ✅ Ready for community plugin submissions

---

## Timeline Estimate

| Task | Duration | Status |
|------|----------|--------|
| Backend API (client.go) | 3-4 hours | Planned |
| Registry Management (registry.go) | 2-3 hours | Planned |
| Plugin Browser UI | 4-5 hours | Planned |
| Integration & Testing | 3-4 hours | Planned |
| Documentation | 3-4 hours | Planned |
| **Total** | **15-20 hours** | **Planning** |

---

## Next Steps

1. ✅ Review and approve Phase 4C plan
2. ⏳ Implement marketplace backend API
3. ⏳ Build plugin browser UI
4. ⏳ Create test suite
5. ⏳ Write comprehensive documentation
6. ⏳ Set up GitHub plugin registry
7. ⏳ Submit example plugins to marketplace

---

**Phase 4C Status: READY FOR IMPLEMENTATION** ✅

All specifications documented. Ready to begin backend development.

**Project Lead:** Development Team  
**Date Created:** January 16, 2026  
**Last Updated:** January 16, 2026 (Evening)

---
