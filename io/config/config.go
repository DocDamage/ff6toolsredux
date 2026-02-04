package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"ffvi_editor/global"
)

const (
	file = "ff6editor.config"
)

var (
	data d
)

type (
	d struct {
		WindowX             float32                     `json:"width"`
		WindowY             float32                     `json:"height"`
		SaveDir             string                      `json:"dir"`
		AutoEnableCmd       bool                        `json:"autoEnableCmd"`
		EnablePlayStation   bool                        `json:"ps"`
		FontSize            int                         `json:"fontSize,omitempty"`
		WorldMapPoints      map[int]map[string]MapPoint `json:"worldMapPoints,omitempty"`
		WorldMapLocations   map[int][]string            `json:"worldMapLocations,omitempty"`
		MarketplaceSettings MarketplaceConfig           `json:"marketplaceSettings,omitempty"`
	}
)

func init() {
	if b, err := os.ReadFile(filepath.Join(global.PWD, file)); err == nil {
		if err := json.Unmarshal(b, &data); err != nil {
			log.Printf("[config] failed to parse config file: %v", err)
		}
	}
	if data.WindowX == 0 {
		data.WindowX = global.WindowWidth
	}
	if data.WindowY == 0 {
		data.WindowY = global.WindowHeight
	}
}

func WindowSize() (x, y float32) {
	x = data.WindowX
	y = data.WindowY
	return
}

func SaveDir() string {
	return data.SaveDir
}

func AutoEnableCmd() bool {
	return data.AutoEnableCmd
}

func EnablePlayStation() bool { return data.EnablePlayStation }

func SetWindowSize(x, y float32) {
	data.WindowX = x
	data.WindowY = y
	save()
}

func SetSaveDir(dir string) {
	data.SaveDir = dir
	save()
}

func SetAutoEnableCmd(v bool) {
	data.AutoEnableCmd = v
	save()
}

func SetEnablePlayStation(v bool) {
	data.EnablePlayStation = v
	save()
}

// FontSize returns the configured font size
func FontSize() int {
	if data.FontSize == 0 {
		return 12 // Default font size
	}
	return data.FontSize
}

// SetFontSize sets the font size
func SetFontSize(size int) {
	data.FontSize = size
	save()
}

func save() {
	if data.WindowX == 0 {
		data.WindowX = global.WindowWidth
	}
	if data.WindowY == 0 {
		data.WindowY = global.WindowHeight
	}
	b, err := json.Marshal(&data)
	if err != nil {
		log.Printf("[config] failed to marshal config data: %v", err)
		return
	}
	if err := os.WriteFile(filepath.Join(global.PWD, file), b, 0755); err != nil {
		log.Printf("[config] failed to write config file: %v", err)
	}
}
