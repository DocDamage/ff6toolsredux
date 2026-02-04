package pr

import (
	"encoding/json"
	"fmt"
)

// decodeInterfaceSlice safely decodes an interface{} to a []interface{} or returns an error.
func decodeInterfaceSlice(i interface{}, context string) ([]interface{}, error) {
	slice, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected []interface{} for %s, got %T", context, i)
	}
	return slice, nil
}

// decodeString safely decodes an interface{} to a string or returns an error.
func decodeString(i interface{}, context string) (string, error) {
	s, ok := i.(string)
	if !ok {
		return "", fmt.Errorf("expected string for %s, got %T", context, i)
	}
	return s, nil
}

// decodeRawMessage decodes an interface{} to json.RawMessage if possible.
func decodeRawMessage(i interface{}, context string) (json.RawMessage, error) {
	s, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("expected string for %s, got %T", context, i)
	}
	return json.RawMessage(s), nil
}
