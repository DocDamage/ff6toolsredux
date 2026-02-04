package plugins

import "context"

// RegisterHook registers an event hook
func (a *APIImpl) RegisterHook(event string, callback func(interface{}) error) error {
	if a.hooks == nil {
		a.hooks = make(map[string][]func(interface{}) error)
	}
	a.hooks[event] = append(a.hooks[event], callback)
	return nil
}

// FireEvent triggers an event
func (a *APIImpl) FireEvent(ctx context.Context, event string, data interface{}) error {
	if a.hooks == nil {
		return nil
	}

	callbacks, ok := a.hooks[event]
	if !ok {
		return nil
	}

	for _, callback := range callbacks {
		if err := callback(data); err != nil {
			return err
		}
	}
	return nil
}
