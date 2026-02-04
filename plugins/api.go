package plugins

import (
	"context"
	ioPR "ffvi_editor/io/pr"
	"ffvi_editor/models"
	modelsPR "ffvi_editor/models/pr"
)

// PluginAPI provides safe access to save editor functionality for plugins
type PluginAPI interface {
	// Save Data Access
	GetCharacter(ctx context.Context, name string) (*models.Character, error)
	SetCharacter(ctx context.Context, name string, ch *models.Character) error
	GetInventory(ctx context.Context) (*modelsPR.Inventory, error)
	SetInventory(ctx context.Context, inv *modelsPR.Inventory) error
	GetParty(ctx context.Context) (*modelsPR.Party, error)
	SetParty(ctx context.Context, party *modelsPR.Party) error
	GetEquipment(ctx context.Context) (*models.Equipment, error)
	SetEquipment(ctx context.Context, eq *models.Equipment) error

	// Batch Operations
	ApplyBatchOperation(ctx context.Context, op string, params map[string]interface{}) (int, error)

	// Query Operations
	FindCharacter(ctx context.Context, predicate func(*models.Character) bool) *models.Character
	FindItems(ctx context.Context, predicate func(*modelsPR.Row) bool) []*modelsPR.Row

	// Events
	RegisterHook(event string, callback func(interface{}) error) error
	FireEvent(ctx context.Context, event string, data interface{}) error

	// UI
	ShowDialog(ctx context.Context, title, message string) error
	ShowConfirm(ctx context.Context, title, message string) (bool, error)
	ShowInput(ctx context.Context, prompt string) (string, error)

	// Logging
	Log(ctx context.Context, level string, message string) error

	// Settings
	GetSetting(key string) interface{}
	SetSetting(key string, value interface{}) error

	// Permissions
	HasPermission(permission string) bool
}

// APIImpl provides a default implementation of PluginAPI
type APIImpl struct {
	prData        *ioPR.PR
	hooks         map[string][]func(interface{}) error
	settings      map[string]interface{}
	permissions   map[string]bool
	logger        func(level, msg string)
	showDialogFn  func(title, message string) error
	showConfirmFn func(title, message string) bool
	showInputFn   func(prompt string) (string, error)
}
