package plugins

import (
	"context"

	"ffvi_editor/models"
	modelsPR "ffvi_editor/models/pr"
)

// MockAPI is a mock implementation of the PluginAPI for testing
type MockAPI struct{}

// NewMockAPI creates a new mock API
func NewMockAPI() *MockAPI {
	return &MockAPI{}
}

// OpenSaveFile mocks the OpenSaveFile function
func (m *MockAPI) OpenSaveFile(path string) error {
	return nil
}

// SaveFile mocks the SaveFile function
func (m *MockAPI) SaveFile(path string) error {
	return nil
}

// GetCharacter mocks the GetCharacter function
func (m *MockAPI) GetCharacter(ctx context.Context, name string) (*models.Character, error) {
	return nil, nil
}

// SetCharacter mocks the SetCharacter function
func (m *MockAPI) SetCharacter(ctx context.Context, name string, ch *models.Character) error {
	return nil
}

// GetInventory mocks the GetInventory function
func (m *MockAPI) GetInventory(ctx context.Context) (*modelsPR.Inventory, error) {
	return nil, nil
}

// SetInventory mocks the SetInventory function
func (m *MockAPI) SetInventory(ctx context.Context, inv *modelsPR.Inventory) error {
	return nil
}

// GetParty mocks the GetParty function
func (m *MockAPI) GetParty(ctx context.Context) (*modelsPR.Party, error) {
	return nil, nil
}

// SetParty mocks the SetParty function
func (m *MockAPI) SetParty(ctx context.Context, party *modelsPR.Party) error {
	return nil
}

// GetEquipment mocks the GetEquipment function
func (m *MockAPI) GetEquipment(ctx context.Context) (*models.Equipment, error) {
	return nil, nil
}

// SetEquipment mocks the SetEquipment function
func (m *MockAPI) SetEquipment(ctx context.Context, eq *models.Equipment) error {
	return nil
}

// ApplyBatchOperation mocks the ApplyBatchOperation function
func (m *MockAPI) ApplyBatchOperation(ctx context.Context, op string, params map[string]interface{}) (int, error) {
	return 0, nil
}

// FindCharacter mocks the FindCharacter function
func (m *MockAPI) FindCharacter(ctx context.Context, predicate func(*models.Character) bool) *models.Character {
	return nil
}

// FindItems mocks the FindItems function
func (m *MockAPI) FindItems(ctx context.Context, predicate func(*modelsPR.Row) bool) []*modelsPR.Row {
	return nil
}

// RegisterHook mocks the RegisterHook function
func (m *MockAPI) RegisterHook(event string, callback func(interface{}) error) error {
	return nil
}

// FireEvent mocks the FireEvent function
func (m *MockAPI) FireEvent(ctx context.Context, event string, data interface{}) error {
	return nil
}

// ShowDialog mocks the ShowDialog function
func (m *MockAPI) ShowDialog(ctx context.Context, title, message string) error {
	return nil
}

// ShowConfirm mocks the ShowConfirm function
func (m *MockAPI) ShowConfirm(ctx context.Context, title, message string) (bool, error) {
	return true, nil
}

// ShowInput mocks the ShowInput function
func (m *MockAPI) ShowInput(ctx context.Context, prompt string) (string, error) {
	return "", nil
}

// Log mocks the Log function
func (m *MockAPI) Log(ctx context.Context, level string, message string) error {
	return nil
}

// GetSetting mocks the GetSetting function
func (m *MockAPI) GetSetting(key string) interface{} {
	return nil
}

// SetSetting mocks the SetSetting function
func (m *MockAPI) SetSetting(key string, value interface{}) error {
	return nil
}

// HasPermission mocks the HasPermission function
func (m *MockAPI) HasPermission(permission string) bool {
	return true
}

// ShowMessage mocks the ShowMessage function (legacy)
func (m *MockAPI) ShowMessage(message string) {
	// No-op
}

// LogMessage mocks the LogMessage function (legacy)
func (m *MockAPI) LogMessage(message string) {
	// No-op
}
