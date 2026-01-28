package plugins

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
func (m *MockAPI) GetCharacter(id int) (interface{}, error) {
	return nil, nil
}

// SetCharacter mocks the SetCharacter function
func (m *MockAPI) SetCharacter(id int, data interface{}) error {
	return nil
}

// ShowMessage mocks the ShowMessage function
func (m *MockAPI) ShowMessage(message string) {
	// No-op
}

// LogMessage mocks the LogMessage function
func (m *MockAPI) LogMessage(message string) {
	// No-op
}
