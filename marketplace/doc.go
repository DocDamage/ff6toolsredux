// Package marketplace provides local and remote plugin/preset management.
//
// The marketplace supports:
//   - Plugin discovery, download, and installation (server-based)
//   - Local preset management (local-only)
//   - Rating and review storage (local-only)
//   - Update checking for installed plugins
//
// Content Types:
//
//	Plugins  - Extend editor functionality (requires server for discovery)
//	Presets  - Character/party templates (local storage)
//
// Implementation Status:
//
// Plugin Operations (Server-based):
//   - SearchPlugins, ListPlugins, GetPluginDetails - FULLY IMPLEMENTED
//   - DownloadPlugin, InstallPlugin                  - FULLY IMPLEMENTED
//   - SubmitRating, GetPluginRatings                 - FULLY IMPLEMENTED
//   - CheckForUpdates                                - FULLY IMPLEMENTED
//
// Preset Operations (Local-only):
//   - UploadPreset, UpdatePreset, DeletePreset       - LOCAL STORAGE ONLY
//   - GetPreset, DownloadPreset                      - LOCAL STORAGE ONLY
//   - RatePreset, AddReview, GetReviews              - LOCAL STORAGE ONLY
//   - GetPopular, GetRecent                          - LOCAL STORAGE ONLY
//
// Server-based preset sharing requires a marketplace backend that is not
// yet available. Presets can be exported/imported via files for sharing.
//
// Usage:
//
//	client, _ := marketplace.NewClientWithRegistry(baseURL, apiKey, registryPath)
//	plugins, _ := client.SearchPlugins(ctx, "combat")
//	preset, _ := client.GetPreset("my-preset")  // Local only
package marketplace
