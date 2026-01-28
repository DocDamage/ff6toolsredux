package plugins

import "fmt"

// Plugin-related errors
var (
	ErrNilAPI                    = fmt.Errorf("plugin API is nil")
	ErrPluginNotInitialized      = fmt.Errorf("plugin is not initialized")
	ErrInvalidPluginName         = fmt.Errorf("plugin name is invalid")
	ErrInvalidPluginVersion      = fmt.Errorf("plugin version is invalid")
	ErrInvalidPluginAuthor       = fmt.Errorf("plugin author is invalid")
	ErrNilPRData                 = fmt.Errorf("PR data is nil")
	ErrCharacterNotFound         = fmt.Errorf("character not found")
	ErrInsufficientPermissions   = fmt.Errorf("insufficient permissions for this operation")
	ErrNilCallback               = fmt.Errorf("callback function is nil")
	ErrBatchOpNotSupported       = fmt.Errorf("batch operation is not supported")
	ErrPluginNotFound            = fmt.Errorf("plugin not found")
	ErrPluginAlreadyLoaded       = fmt.Errorf("plugin is already loaded")
	ErrPluginLoadFailed          = fmt.Errorf("failed to load plugin")
	ErrPluginExecutionTimeout    = fmt.Errorf("plugin execution timeout")
	ErrPluginExecutionError      = fmt.Errorf("plugin execution error")
	ErrInvalidPluginCode         = fmt.Errorf("invalid plugin code")
	ErrMaxPluginsExceeded        = fmt.Errorf("maximum number of plugins exceeded")
	ErrPluginDependencyMissing   = fmt.Errorf("required plugin dependency is missing")
	ErrPluginVersionMismatch     = fmt.Errorf("plugin version is incompatible with editor")
)
