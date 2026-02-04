package pr

import (
	"encoding/json"
	"fmt"

	jo "gitlab.com/c0b/go-ordered-json"
)

// TypeSafeExtractors provides type-safe value extraction with validation
// This replaces unsafe interface{} patterns throughout the codebase

// ExtractStringArray safely extracts a string array from interface{} with validation
func ExtractStringArray(value interface{}) ([]string, error) {
	if value == nil {
		return nil, fmt.Errorf("value is nil")
	}

	rawArray, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", value)
	}

	result := make([]string, 0, len(rawArray))
	for i, item := range rawArray {
		str, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("array[%d]: expected string, got %T", i, item)
		}
		result = append(result, str)
	}
	return result, nil
}

// ExtractInt64Array safely extracts an int64 array from interface{} with validation
func ExtractInt64Array(value interface{}) ([]int64, error) {
	if value == nil {
		return nil, fmt.Errorf("value is nil")
	}

	rawArray, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", value)
	}

	result := make([]int64, 0, len(rawArray))
	for i, item := range rawArray {
		num, ok := item.(json.Number)
		if !ok {
			return nil, fmt.Errorf("array[%d]: expected json.Number, got %T", i, item)
		}
		val, err := num.Int64()
		if err != nil {
			return nil, fmt.Errorf("array[%d]: cannot parse as int64: %w", i, err)
		}
		result = append(result, val)
	}
	return result, nil
}

// ExtractFloat64Array safely extracts a float64 array from interface{} with validation
func ExtractFloat64Array(value interface{}) ([]float64, error) {
	if value == nil {
		return nil, fmt.Errorf("value is nil")
	}

	rawArray, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", value)
	}

	result := make([]float64, 0, len(rawArray))
	for i, item := range rawArray {
		num, ok := item.(json.Number)
		if !ok {
			return nil, fmt.Errorf("array[%d]: expected json.Number, got %T", i, item)
		}
		val, err := num.Float64()
		if err != nil {
			return nil, fmt.Errorf("array[%d]: cannot parse as float64: %w", i, err)
		}
		result = append(result, val)
	}
	return result, nil
}

// ExtractString safely extracts a string from interface{} with validation
func ExtractString(value interface{}) (string, error) {
	if value == nil {
		return "", fmt.Errorf("value is nil")
	}

	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %T", value)
	}
	return str, nil
}

// ExtractInt64 safely extracts an int64 from interface{} with validation
func ExtractInt64(value interface{}) (int64, error) {
	if value == nil {
		return 0, fmt.Errorf("value is nil")
	}

	num, ok := value.(json.Number)
	if !ok {
		return 0, fmt.Errorf("expected json.Number, got %T", value)
	}

	val, err := num.Int64()
	if err != nil {
		return 0, fmt.Errorf("cannot parse as int64: %w", err)
	}
	return val, nil
}

// UnmarshalOrderedMapFromString safely unmarshals a JSON string into an OrderedMap
func UnmarshalOrderedMapFromString(jsonStr string) (*jo.OrderedMap, error) {
	if jsonStr == "" {
		return nil, fmt.Errorf("empty JSON string")
	}

	om := jo.NewOrderedMap()
	if err := om.UnmarshalJSON([]byte(jsonStr)); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return om, nil
}

// SafeGetFromTarget retrieves values from the "target" key in nested structures
// with proper type validation. Returns typed slice of strings.
func SafeGetFromTarget(data *jo.OrderedMap, key string) ([]string, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	slTarget := jo.NewOrderedMap()
	value, ok := data.GetValue(key)
	if !ok {
		return nil, fmt.Errorf("key %q not found", key)
	}

	valueStr, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("expected string value for key %q, got %T", key, value)
	}

	if err := slTarget.UnmarshalJSON([]byte(valueStr)); err != nil {
		return nil, fmt.Errorf("failed to unmarshal value for key %q: %w", key, err)
	}

	targetValue, ok := slTarget.GetValue("target")
	if !ok {
		return nil, fmt.Errorf("target key not found in %q", key)
	}

	return ExtractStringArray(targetValue)
}

// SafeGetFromTargetRaw retrieves raw values from the "target" key
// This is used when the target contains mixed or unknown types
func SafeGetFromTargetRaw(data *jo.OrderedMap, key string) (interface{}, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	slTarget := jo.NewOrderedMap()
	value, ok := data.GetValue(key)
	if !ok {
		return nil, fmt.Errorf("key %q not found", key)
	}

	valueStr, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("expected string value for key %q, got %T", key, value)
	}

	if err := slTarget.UnmarshalJSON([]byte(valueStr)); err != nil {
		return nil, fmt.Errorf("failed to unmarshal value for key %q: %w", key, err)
	}

	targetValue, ok := slTarget.GetValue("target")
	if !ok {
		return nil, fmt.Errorf("target key not found in %q", key)
	}

	if targetValue == nil {
		return nil, fmt.Errorf("target value is nil for key %q", key)
	}

	return targetValue, nil
}

// ValidateArray ensures a value is an array and returns its length
func ValidateArray(value interface{}, description string) (int, error) {
	if value == nil {
		return 0, fmt.Errorf("%s: value is nil", description)
	}

	arr, ok := value.([]interface{})
	if !ok {
		return 0, fmt.Errorf("%s: expected array, got %T", description, value)
	}

	return len(arr), nil
}

// SafeUnmarshalJSON safely unmarshals JSON with error context
func SafeUnmarshalJSON(dest interface{}, jsonData string, context string) error {
	if jsonData == "" {
		return fmt.Errorf("%s: empty JSON data", context)
	}

	data := []byte(jsonData)
	if len(data) == 0 {
		return fmt.Errorf("%s: empty JSON data", context)
	}

	// Type assertion based on dest type
	switch v := dest.(type) {
	case *jo.OrderedMap:
		return v.UnmarshalJSON(data)
	case interface{} /* generic JSON */ :
		return json.Unmarshal(data, v)
	default:
		return json.Unmarshal(data, dest)
	}
}
