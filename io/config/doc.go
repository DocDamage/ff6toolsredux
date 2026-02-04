// Package config provides configuration file management.
//
// The config package handles:
//   - Loading and saving application configuration
//   - Map point and location management
//   - Marketplace configuration
//   - Legacy config migration
//
// This package works alongside the settings package, providing
// lower-level configuration I/O operations.
//
// Example usage:
//
//	// Load configuration
//	cfg, err := config.Load(configPath)
//	if err != nil {
//	    return err
//	}
//
//	// Access map points
//	points := cfg.GetMapPoints(worldID)
package config
