package pr

import (
	"encoding/json"
	"testing"
)

func TestDecodeInterfaceSlice(t *testing.T) {
	good := []interface{}{1, 2, 3}
	bad := map[string]interface{}{"a": 1}

	_, err := decodeInterfaceSlice(good, "test")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, err = decodeInterfaceSlice(bad, "test")
	if err == nil {
		t.Errorf("expected error for bad input, got nil")
	}
}

func TestDecodeString(t *testing.T) {
	good := "hello"
	bad := 123

	_, err := decodeString(good, "test")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, err = decodeString(bad, "test")
	if err == nil {
		t.Errorf("expected error for bad input, got nil")
	}
}

func TestDecodeRawMessage(t *testing.T) {
	good := `{"foo": 1}`
	bad := 123

	msg, err := decodeRawMessage(good, "test")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(msg, &m); err != nil {
		t.Errorf("expected valid json, got %v", err)
	}

	_, err = decodeRawMessage(bad, "test")
	if err == nil {
		t.Errorf("expected error for bad input, got nil")
	}
}
