package pr

import (
	"encoding/json"
	"strings"
	"testing"

	jo "gitlab.com/c0b/go-ordered-json"
)

// TestExtractStringArray tests type-safe string array extraction
func TestExtractStringArray(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		want    []string
		wantErr bool
	}{
		{
			name:    "valid string array",
			value:   []interface{}{"a", "b", "c"},
			want:    []string{"a", "b", "c"},
			wantErr: false,
		},
		{
			name:    "empty array",
			value:   []interface{}{},
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "nil value",
			value:   nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "not an array",
			value:   "not_array",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "mixed types in array",
			value:   []interface{}{"a", 1, "c"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractStringArray(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Fatalf("length mismatch: got %d, want %d", len(got), len(tt.want))
				}
				for i, v := range got {
					if v != tt.want[i] {
						t.Fatalf("element %d: got %q, want %q", i, v, tt.want[i])
					}
				}
			}
		})
	}
}

// TestExtractInt64 tests type-safe int64 extraction
func TestExtractInt64(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		want    int64
		wantErr bool
	}{
		{
			name:    "json.Number integer",
			value:   newJsonNumber("42"),
			want:    42,
			wantErr: false,
		},
		{
			name:    "large number",
			value:   newJsonNumber("9223372036854775807"), // max int64
			want:    9223372036854775807,
			wantErr: false,
		},
		{
			name:    "negative number",
			value:   newJsonNumber("-100"),
			want:    -100,
			wantErr: false,
		},
		{
			name:    "nil value",
			value:   nil,
			want:    0,
			wantErr: true,
		},
		{
			name:    "not a json.Number",
			value:   "not_number",
			want:    0,
			wantErr: true,
		},
		{
			name:    "float number",
			value:   newJsonNumber("3.14"),
			want:    3,
			wantErr: true, // Expect error when parsing float as int
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractInt64(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

// TestExtractString tests type-safe string extraction
func TestExtractString(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "valid string",
			value:   "hello",
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "empty string",
			value:   "",
			want:    "",
			wantErr: false,
		},
		{
			name:    "nil value",
			value:   nil,
			want:    "",
			wantErr: true,
		},
		{
			name:    "not a string",
			value:   42,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractString(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// TestUnmarshalOrderedMapFromString tests safe OrderedMap unmarshaling
func TestUnmarshalOrderedMapFromString(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name:    "valid JSON",
			json:    `{"key": "value"}`,
			wantErr: false,
		},
		{
			name:    "empty JSON object",
			json:    `{}`,
			wantErr: false,
		},
		{
			name:    "empty string",
			json:    "",
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			json:    `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalOrderedMapFromString(tt.json)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Fatal("got nil OrderedMap")
			}
		})
	}
}

// TestSafeGetFromTarget tests safe extraction from target key
func TestSafeGetFromTarget(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		key     string
		wantLen int
		wantErr bool
	}{
		{
			name:    "valid target with string array",
			json:    `{"data": "{\"target\": [\"a\", \"b\", \"c\"]}"}`,
			key:     "data",
			wantLen: 3,
			wantErr: false,
		},
		{
			name:    "missing key",
			json:    `{"other": "{\"target\": []}"}`,
			key:     "data",
			wantErr: true,
		},
		{
			name:    "missing target",
			json:    `{"data": "{\"values\": []}"}`,
			key:     "data",
			wantErr: true,
		},
		{
			name:    "empty target array",
			json:    `{"data": "{\"target\": []}"}`,
			key:     "data",
			wantLen: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := jo.NewOrderedMap()
			om.UnmarshalJSON([]byte(tt.json))

			got, err := SafeGetFromTarget(om, tt.key)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(got) != tt.wantLen {
				t.Fatalf("length: got %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

// TestValidateArray tests array validation
func TestValidateArray(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		want    int
		wantErr bool
	}{
		{
			name:    "valid array",
			value:   []interface{}{"a", "b", "c"},
			want:    3,
			wantErr: false,
		},
		{
			name:    "empty array",
			value:   []interface{}{},
			want:    0,
			wantErr: false,
		},
		{
			name:    "nil value",
			value:   nil,
			want:    0,
			wantErr: true,
		},
		{
			name:    "not an array",
			value:   "string",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateArray(tt.value, "test array")
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

// Helper function to create json.Number for testing
func newJsonNumber(s string) interface{} {
	// json.Number comes from JSON unmarshaling with UseNumber()
	var result map[string]interface{}
	dec := json.NewDecoder(strings.NewReader(`{"num":` + s + `}`))
	dec.UseNumber()
	_ = dec.Decode(&result)
	if num, ok := result["num"]; ok {
		return num
	}
	return nil
}

// BenchmarkExtractStringArray benchmarks string array extraction
func BenchmarkExtractStringArray(b *testing.B) {
	value := []interface{}{"a", "b", "c", "d", "e"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractStringArray(value)
	}
}

// BenchmarkSafeGetFromTarget benchmarks safe target extraction
func BenchmarkSafeGetFromTarget(b *testing.B) {
	om := jo.NewOrderedMap()
	om.UnmarshalJSON([]byte(`{"data": "{\"target\": [\"a\", \"b\", \"c\"]}"}`))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SafeGetFromTarget(om, "data")
	}
}
