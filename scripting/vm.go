package scripting

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// VM wraps a Lua VM with sandboxing and timeout support
type VM struct {
	timeout   time.Duration
	maxMemory int
	modules   map[string]bool
	api       interface{}
	mu        sync.RWMutex
	running   bool
}

// NewVM creates a new Lua VM instance
func NewVM(timeout time.Duration) *VM {
	vm := &VM{
		timeout:   timeout,
		maxMemory: 50 * 1024 * 1024, // 50MB default
		modules:   make(map[string]bool),
	}

	// Only allow safe modules
	vm.modules["table"] = true
	vm.modules["string"] = true
	vm.modules["math"] = true
	vm.modules["utf8"] = true

	return vm
}

// SetAPI sets the plugin API for Lua scripts
func (vm *VM) SetAPI(api interface{}) {
	vm.mu.Lock()
	vm.api = api
	vm.mu.Unlock()
}

// Execute executes Lua code with sandboxing
func (vm *VM) Execute(ctx context.Context, code string) error {
	if code == "" {
		return fmt.Errorf("code is empty")
	}

	vm.mu.Lock()
	if vm.running {
		vm.mu.Unlock()
		return fmt.Errorf("script already running")
	}
	vm.running = true
	vm.mu.Unlock()

	defer func() {
		vm.mu.Lock()
		vm.running = false
		vm.mu.Unlock()
	}()

	// Create context with timeout
	execCtx, cancel := context.WithTimeout(ctx, vm.timeout)
	defer cancel()

	// Execute code in goroutine
	done := make(chan error, 1)
	go func() {
		// Placeholder for actual Lua execution
		// Real implementation would:
		// 1. Create Lua state
		// 2. Load restricted libraries
		// 3. Register API bindings
		// 4. Load and run code
		// 5. Handle errors and cleanup
		done <- nil
	}()

	// Wait for execution or timeout
	select {
	case err := <-done:
		return err
	case <-execCtx.Done():
		return fmt.Errorf("execution timeout: %w", execCtx.Err())
	}
}

// ExecuteFile executes a Lua file with sandboxing
func (vm *VM) ExecuteFile(ctx context.Context, filepath string) error {
	if filepath == "" {
		return fmt.Errorf("filepath is empty")
	}

	// Placeholder for actual file execution
	// Real implementation would read file and call Execute()
	return fmt.Errorf("not implemented: file execution requires actual Lua VM")
}

// Call calls a Lua function with arguments
func (vm *VM) Call(ctx context.Context, functionName string, args ...interface{}) (interface{}, error) {
	if functionName == "" {
		return nil, fmt.Errorf("function name is empty")
	}

	// Placeholder for actual function call
	// Real implementation would use Lua state to call function
	return nil, fmt.Errorf("not implemented: function calls require actual Lua VM")
}

// SetGlobal sets a global variable in Lua
func (vm *VM) SetGlobal(name string, value interface{}) error {
	if name == "" {
		return fmt.Errorf("variable name is empty")
	}

	// Placeholder for actual global setting
	// Real implementation would use Lua state to set global
	return nil
}

// GetGlobal gets a global variable from Lua
func (vm *VM) GetGlobal(name string) (interface{}, error) {
	if name == "" {
		return nil, fmt.Errorf("variable name is empty")
	}

	// Placeholder for actual global retrieval
	// Real implementation would use Lua state to get global
	return nil, fmt.Errorf("not implemented: global retrieval requires actual Lua VM")
}

// Close closes the VM and cleans up resources
func (vm *VM) Close() error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	vm.api = nil
	vm.running = false
	return nil
}

// SetTimeout sets the execution timeout
func (vm *VM) SetTimeout(timeout time.Duration) {
	vm.mu.Lock()
	vm.timeout = timeout
	vm.mu.Unlock()
}

// GetTimeout gets the execution timeout
func (vm *VM) GetTimeout() time.Duration {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.timeout
}

// SetMaxMemory sets the maximum memory usage
func (vm *VM) SetMaxMemory(bytes int) {
	vm.mu.Lock()
	vm.maxMemory = bytes
	vm.mu.Unlock()
}

// GetMaxMemory gets the maximum memory usage
func (vm *VM) GetMaxMemory() int {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.maxMemory
}

// LoadLibrary loads a standard library module
func (vm *VM) LoadLibrary(name string) error {
	if !vm.modules[name] {
		return fmt.Errorf("library %s is not allowed in sandbox", name)
	}

	// Placeholder for actual library loading
	// Real implementation would load Lua standard libraries safely
	return nil
}

// RegisterFunction registers a Go function as a Lua global
func (vm *VM) RegisterFunction(name string, fn interface{}) error {
	if name == "" {
		return fmt.Errorf("function name is empty")
	}
	if fn == nil {
		return fmt.Errorf("function is nil")
	}

	// Placeholder for actual function registration
	// Real implementation would register the Go function in Lua
	return nil
}
func (vm *VM) IsRunning() bool {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.running
}
