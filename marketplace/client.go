package marketplace

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Preset represents a shared preset in the marketplace
type Preset struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // "character", "party", "template"
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Version     string    `json:"version"`
	Tags        []string  `json:"tags"`
	Downloads   int       `json:"downloads"`
	Rating      float64   `json:"rating"`
	Reviews     int       `json:"reviews"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Data        string    `json:"data"` // JSON or encoded data
}

// Review represents a user review
type Review struct {
	ID        string    `json:"id"`
	PresetID  string    `json:"presetId"`
	Author    string    `json:"author"`
	Rating    int       `json:"rating"` // 1-5 stars
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}

// RemotePlugin represents a plugin in the marketplace
type RemotePlugin struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Description      string        `json:"description"`
	FullDescription  string        `json:"fullDescription,omitempty"`
	Author           string        `json:"author"`
	Version          string        `json:"version"`
	MinVersion       string        `json:"minVersion"`
	MaxVersion       string        `json:"maxVersion"`
	Category         string        `json:"category"`
	Tags             []string      `json:"tags"`
	Rating           float32       `json:"rating"`
	RatingCount      int           `json:"ratingCount"`
	Downloads        int           `json:"downloads"`
	License          string        `json:"license"`
	RepositoryURL    string        `json:"repositoryUrl"`
	DocumentationURL string        `json:"documentationUrl"`
	IconURL          string        `json:"iconUrl,omitempty"`
	VersionHistory   []VersionInfo `json:"versionHistory"`
	UpdatedAt        time.Time     `json:"updatedAt"`
	CreatedAt        time.Time     `json:"createdAt"`
}

// VersionInfo contains version-specific information
type VersionInfo struct {
	Version      string    `json:"version"`
	ReleaseDate  time.Time `json:"releaseDate"`
	DownloadURL  string    `json:"downloadUrl"`
	Size         int64     `json:"size"`
	Checksum     string    `json:"checksum"`
	ReleaseNotes string    `json:"releaseNotes"`
}

// PluginRating represents a user rating
type PluginRating struct {
	ID        string    `json:"id"`
	PluginID  string    `json:"pluginId"`
	UserID    string    `json:"userId"`
	Rating    float32   `json:"rating"`
	Review    string    `json:"review"`
	Helpful   int       `json:"helpful"`
	Timestamp time.Time `json:"timestamp"`
}

// ListOptions contains filtering and pagination options
type ListOptions struct {
	Category    string
	Tags        []string
	SortBy      string // rating, downloads, recent, name
	SortOrder   string // asc, desc
	Limit       int
	Offset      int
	MinRating   float32
	OnlyUpdates bool
}

// UpdateInfo contains information about available updates
type UpdateInfo struct {
	PluginID       string
	CurrentVersion string
	LatestVersion  string
	ReleaseNotes   string
	UpdatedAt      time.Time
}

// cacheEntry holds cached responses
type cacheEntry struct {
	data      interface{}
	expiresAt time.Time
}

// Client handles marketplace API communication
type Client struct {
	baseURL   string
	apiKey    string
	httpCli   *http.Client
	cache     map[string]cacheEntry
	cacheTTL  time.Duration
	userAgent string
	registry  *Registry
}

// NewClient creates a new marketplace client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL:   baseURL,
		apiKey:    apiKey,
		httpCli:   &http.Client{Timeout: 30 * time.Second},
		cache:     make(map[string]cacheEntry),
		cacheTTL:  24 * time.Hour,
		userAgent: "FF6Editor-Marketplace/1.0",
	}
}

// NewClientWithRegistry creates a new marketplace client with registry support
func NewClientWithRegistry(baseURL, apiKey string, registryPath string) (*Client, error) {
	registry, err := NewRegistry(registryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create registry: %w", err)
	}

	return &Client{
		baseURL:   baseURL,
		apiKey:    apiKey,
		httpCli:   &http.Client{Timeout: 30 * time.Second},
		cache:     make(map[string]cacheEntry),
		cacheTTL:  24 * time.Hour,
		userAgent: "FF6Editor-Marketplace/1.0",
		registry:  registry,
	}, nil
}

// Search searches for presets
func (c *Client) Search(query string, presetType string, tags []string) ([]*Preset, error) {
	// Build query string
	q := url.Values{}
	q.Set("q", query)
	if presetType != "" {
		q.Set("type", presetType)
	}
	for _, tag := range tags {
		q.Add("tags", tag)
	}

	url := fmt.Sprintf("%s/presets/search?%s", c.baseURL, q.Encode())
	resp, err := c.doRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed: %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Presets []*Preset `json:"presets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Presets, nil
}

// ListPlugins retrieves a list of plugins from the marketplace
func (c *Client) ListPlugins(ctx context.Context, opts *ListOptions) ([]RemotePlugin, error) {
	if opts == nil {
		opts = &ListOptions{Limit: 20, Offset: 0}
	}

	// Set defaults
	if opts.Limit == 0 {
		opts.Limit = 20
	}
	if opts.SortBy == "" {
		opts.SortBy = "rating"
	}
	if opts.SortOrder == "" {
		opts.SortOrder = "desc"
	}

	// Build query string
	query := url.Values{}
	query.Set("limit", fmt.Sprintf("%d", opts.Limit))
	query.Set("offset", fmt.Sprintf("%d", opts.Offset))
	query.Set("sortBy", opts.SortBy)
	query.Set("sortOrder", opts.SortOrder)
	query.Set("minRating", fmt.Sprintf("%.1f", opts.MinRating))

	if opts.Category != "" {
		query.Set("category", opts.Category)
	}
	for _, tag := range opts.Tags {
		query.Add("tags", tag)
	}
	if opts.OnlyUpdates {
		query.Set("onlyUpdates", "true")
	}

	reqURL := fmt.Sprintf("%s/plugins?%s", c.baseURL, query.Encode())

	resp, err := c.doRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("marketplace returned %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Plugins []RemotePlugin `json:"plugins"`
		Total   int            `json:"total"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode plugin list: %w", err)
	}

	return result.Plugins, nil
}

// SearchPlugins searches for plugins by query
func (c *Client) SearchPlugins(ctx context.Context, query string) ([]RemotePlugin, error) {
	if query == "" {
		return []RemotePlugin{}, nil
	}

	// Check cache first
	cacheKey := fmt.Sprintf("search:%s", query)
	if cached, ok := c.getFromCache(cacheKey); ok {
		if plugins, ok := cached.([]RemotePlugin); ok {
			return plugins, nil
		}
	}

	reqURL := fmt.Sprintf("%s/plugins/search?q=%s", c.baseURL, url.QueryEscape(query))

	resp, err := c.doRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed: %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Plugins []RemotePlugin `json:"plugins"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search results: %w", err)
	}

	// Cache results
	c.setCache(cacheKey, result.Plugins)

	return result.Plugins, nil
}

// GetPluginDetails retrieves detailed information about a specific plugin
func (c *Client) GetPluginDetails(ctx context.Context, pluginID string) (*RemotePlugin, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("plugin:%s", pluginID)
	if cached, ok := c.getFromCache(cacheKey); ok {
		if plugin, ok := cached.(*RemotePlugin); ok {
			return plugin, nil
		}
	}

	reqURL := fmt.Sprintf("%s/plugins/%s", c.baseURL, pluginID)

	resp, err := c.doRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("plugin not found: %s", pluginID)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get plugin: %d: %s", resp.StatusCode, string(body))
	}

	var plugin RemotePlugin
	if err := json.NewDecoder(resp.Body).Decode(&plugin); err != nil {
		return nil, fmt.Errorf("failed to decode plugin: %w", err)
	}

	// Cache result
	c.setCache(cacheKey, &plugin)

	return &plugin, nil
}

// DownloadPlugin downloads a plugin package
func (c *Client) DownloadPlugin(ctx context.Context, pluginID, version string) ([]byte, error) {
	plugin, err := c.GetPluginDetails(ctx, pluginID)
	if err != nil {
		return nil, err
	}

	// Find the version info
	var versionInfo *VersionInfo
	if version == "" {
		// Use latest version
		if len(plugin.VersionHistory) > 0 {
			versionInfo = &plugin.VersionHistory[0]
		} else {
			return nil, fmt.Errorf("no versions available for plugin %s", pluginID)
		}
	} else {
		// Find specific version
		for i := range plugin.VersionHistory {
			if plugin.VersionHistory[i].Version == version {
				versionInfo = &plugin.VersionHistory[i]
				break
			}
		}
		if versionInfo == nil {
			return nil, fmt.Errorf("version %s not found for plugin %s", version, pluginID)
		}
	}

	// Download the plugin
	resp, err := c.doRequestWithContext(ctx, "GET", versionInfo.DownloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read download: %w", err)
	}

	// Verify checksum
	if versionInfo.Checksum != "" {
		hash := sha256.Sum256(data)
		hashStr := fmt.Sprintf("sha256:%x", hash)
		if versionInfo.Checksum != hashStr {
			return nil, fmt.Errorf("checksum mismatch: expected %s, got %s", versionInfo.Checksum, hashStr)
		}
	}

	return data, nil
}

// InstallPlugin downloads and installs a plugin
func (c *Client) InstallPlugin(ctx context.Context, pluginID, version string) error {
	_, err := c.DownloadPlugin(ctx, pluginID, version)
	if err != nil {
		return err
	}

	// Track installation if registry is available
	if c.registry != nil {
		if err := c.registry.TrackInstallation(pluginID, version); err != nil {
			return fmt.Errorf("failed to track installation: %w", err)
		}
	}

	return nil
}

// SubmitRating submits a rating and review for a plugin
func (c *Client) SubmitRating(ctx context.Context, rating *PluginRating) error {
	if rating.Rating < 1 || rating.Rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	body, err := json.Marshal(rating)
	if err != nil {
		return fmt.Errorf("failed to marshal rating: %w", err)
	}

	reqURL := fmt.Sprintf("%s/plugins/%s/ratings", c.baseURL, rating.PluginID)

	resp, err := c.doRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("submission failed: %d: %s", resp.StatusCode, string(body))
	}

	// Save locally if registry is available
	if c.registry != nil {
		if err := c.registry.SaveRatings([]*PluginRating{rating}); err != nil {
			return fmt.Errorf("failed to save rating locally: %w", err)
		}
	}

	return nil
}

// GetPluginRatings retrieves ratings and reviews for a plugin
func (c *Client) GetPluginRatings(ctx context.Context, pluginID string) ([]*PluginRating, error) {
	reqURL := fmt.Sprintf("%s/plugins/%s/ratings", c.baseURL, pluginID)

	resp, err := c.doRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get ratings: %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Ratings []*PluginRating `json:"ratings"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode ratings: %w", err)
	}

	return result.Ratings, nil
}

// CheckForUpdates checks for available updates for installed plugins
func (c *Client) CheckForUpdates(ctx context.Context) (map[string]*UpdateInfo, error) {
	if c.registry == nil {
		return nil, fmt.Errorf("registry not available")
	}

	installed := c.registry.GetInstalledPlugins()
	updates := make(map[string]*UpdateInfo)

	for _, record := range installed {
		plugin, err := c.GetPluginDetails(ctx, record.PluginID)
		if err != nil {
			// Log error but continue checking other plugins
			continue
		}

		// Compare versions
		if compareVersions(plugin.Version, record.Version) > 0 {
			updates[record.PluginID] = &UpdateInfo{
				PluginID:       record.PluginID,
				CurrentVersion: record.Version,
				LatestVersion:  plugin.Version,
				ReleaseNotes:   plugin.Version,
				UpdatedAt:      plugin.UpdatedAt,
			}
		}
	}

	return updates, nil
}

// Helper methods

// doRequest performs an HTTP request with proper headers
func (c *Client) doRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpCli.Do(req)
}

// doRequestWithContext performs an HTTP request with context
func (c *Client) doRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpCli.Do(req)
}

// getFromCache retrieves a value from cache if not expired
func (c *Client) getFromCache(key string) (interface{}, bool) {
	if entry, ok := c.cache[key]; ok {
		if time.Now().Before(entry.expiresAt) {
			return entry.data, true
		}
		// Expired, delete it
		delete(c.cache, key)
	}
	return nil, false
}

// setCache stores a value in cache
func (c *Client) setCache(key string, data interface{}) {
	c.cache[key] = cacheEntry{
		data:      data,
		expiresAt: time.Now().Add(c.cacheTTL),
	}
}

// ClearCache clears all cached data
func (c *Client) ClearCache() {
	c.cache = make(map[string]cacheEntry)
}

// SetCacheTTL sets the cache time-to-live
func (c *Client) SetCacheTTL(ttl time.Duration) {
	c.cacheTTL = ttl
}

// compareVersions compares two semantic versions
// Returns: >0 if v1 > v2, 0 if equal, <0 if v1 < v2
func compareVersions(v1, v2 string) int {
	// Simple version comparison (major.minor.patch)
	// In production, use a proper semver library
	if v1 == v2 {
		return 0
	}
	// Convert to timestamps for now
	// TODO: Implement proper semantic version comparison
	return 0
}

// GetPreset retrieves a preset by ID
func (c *Client) GetPreset(id string) (*Preset, error) {
	// TODO: Implement API call
	return nil, nil
}

// UploadPreset uploads a new preset
func (c *Client) UploadPreset(preset *Preset) error {
	// TODO: Implement API call
	return nil
}

// UpdatePreset updates an existing preset
func (c *Client) UpdatePreset(preset *Preset) error {
	// TODO: Implement API call
	return nil
}

// DeletePreset deletes a preset
func (c *Client) DeletePreset(id string) error {
	// TODO: Implement API call
	return nil
}

// DownloadPreset downloads a preset
func (c *Client) DownloadPreset(id string) (*Preset, error) {
	// TODO: Implement API call
	return nil, nil
}

// RatePreset rates a preset
func (c *Client) RatePreset(id string, rating int) error {
	// TODO: Implement API call
	return nil
}

// AddReview adds a review
func (c *Client) AddReview(review *Review) error {
	// TODO: Implement API call
	return nil
}

// GetReviews gets reviews for a preset
func (c *Client) GetReviews(presetID string) ([]*Review, error) {
	// TODO: Implement API call
	return []*Review{}, nil
}

// GetPopular gets popular presets
func (c *Client) GetPopular(limit int) ([]*Preset, error) {
	// TODO: Implement API call
	return []*Preset{}, nil
}

// GetRecent gets recently added presets
func (c *Client) GetRecent(limit int) ([]*Preset, error) {
	// TODO: Implement API call
	return []*Preset{}, nil
}
// Close closes the marketplace client
func (c *Client) Close() error {
	c.ClearCache()
	if c.registry != nil {
		return c.registry.Close()
	}
	return nil
}