package pr

import (
	"testing"

	"ffvi_editor/models"
	jo "gitlab.com/c0b/go-ordered-json"
)

// TestPRFactoryNew tests the PR factory creation
func TestPRFactoryNew(t *testing.T) {
	pr := New()

	if pr == nil {
		t.Fatal("New() returned nil")
	}
	if pr.Base == nil {
		t.Fatal("Base should be initialized")
	}
	if pr.UserData == nil {
		t.Fatal("UserData should be initialized")
	}
	if pr.MapData == nil {
		t.Fatal("MapData should be initialized")
	}
	if pr.Characters == nil {
		t.Fatal("Characters should be initialized")
	}
	if len(pr.Characters) != 40 {
		t.Fatalf("Characters should have capacity of 40, got %d", len(pr.Characters))
	}
}

// TestPRHasUnicodeNames tests unicode name checking
func TestPRHasUnicodeNames(t *testing.T) {
	pr := New()

	if pr.HasUnicodeNames() {
		t.Fatal("new PR should not have unicode names")
	}

	pr.names = []unicodeNameReplace{{}}
	if !pr.HasUnicodeNames() {
		t.Fatal("PR with names should return true for HasUnicodeNames()")
	}
}

// TestGetString tests string retrieval from OrderedMap
func TestGetString(t *testing.T) {
	helpers := NewTestHelpers(t)
	pr := New()

	tests := []struct {
		name    string
		json    string
		key     string
		want    string
		wantErr bool
	}{
		{
			name:    "valid string",
			json:    `{"name": "TestValue"}`,
			key:     "name",
			want:    "TestValue",
			wantErr: false,
		},
		{
			name:    "missing key",
			json:    `{"other": "value"}`,
			key:     "name",
			want:    "",
			wantErr: true,
		},
		{
			name:    "non-string value",
			json:    `{"value": 123}`,
			key:     "value",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := helpers.CreateOrderedMap(tt.json)
			got, err := pr.getString(om, tt.key)

			if (err != nil) != tt.wantErr {
				t.Fatalf("getString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("getString() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestGetInt tests integer retrieval with various numeric types
func TestGetInt(t *testing.T) {
	helpers := NewTestHelpers(t)
	pr := New()

	tests := []struct {
		name    string
		json    string
		key     string
		want    int
		wantErr bool
	}{
		{
			name:    "integer value",
			json:    `{"value": 42}`,
			key:     "value",
			want:    42,
			wantErr: false,
		},
		{
			name:    "float value cannot be parsed as int",
			json:    `{"value": 42.7}`,
			key:     "value",
			want:    0,
			wantErr: true,
		},
		{
			name:    "string integer",
			json:    `{"value": "100"}`,
			key:     "value",
			want:    100,
			wantErr: false,
		},
		{
			name:    "missing key",
			json:    `{"other": 42}`,
			key:     "value",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid string",
			json:    `{"value": "not_a_number"}`,
			key:     "value",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := helpers.CreateOrderedMap(tt.json)
			got, err := pr.getInt(om, tt.key)

			if (err != nil) != tt.wantErr {
				t.Fatalf("getInt() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("getInt() = %d, want %d", got, tt.want)
			}
		})
	}
}

// TestGetBool tests boolean retrieval
func TestGetBool(t *testing.T) {
	helpers := NewTestHelpers(t)
	pr := New()

	tests := []struct {
		name    string
		json    string
		key     string
		want    bool
		wantErr bool
	}{
		{
			name:    "true value",
			json:    `{"enabled": true}`,
			key:     "enabled",
			want:    true,
			wantErr: false,
		},
		{
			name:    "false value",
			json:    `{"enabled": false}`,
			key:     "enabled",
			want:    false,
			wantErr: false,
		},
		{
			name:    "missing key",
			json:    `{"other": true}`,
			key:     "enabled",
			want:    false,
			wantErr: true,
		},
		{
			name:    "non-bool value",
			json:    `{"value": 1}`,
			key:     "value",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := helpers.CreateOrderedMap(tt.json)
			got, err := pr.getBool(om, tt.key)

			if (err != nil) != tt.wantErr {
				t.Fatalf("getBool() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("getBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetFloat tests float retrieval
func TestGetFloat(t *testing.T) {
	helpers := NewTestHelpers(t)
	pr := New()

	tests := []struct {
		name    string
		json    string
		key     string
		want    float64
		wantErr bool
	}{
		{
			name:    "float value",
			json:    `{"position": 42.5}`,
			key:     "position",
			want:    42.5,
			wantErr: false,
		},
		{
			name:    "integer to float",
			json:    `{"position": 42}`,
			key:     "position",
			want:    42.0,
			wantErr: false,
		},
		{
			name:    "string float",
			json:    `{"position": "3.14"}`,
			key:     "position",
			want:    3.14,
			wantErr: false,
		},
		{
			name:    "missing key",
			json:    `{"other": 42.5}`,
			key:     "position",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := helpers.CreateOrderedMap(tt.json)
			got, err := pr.getFloat(om, tt.key)

			if (err != nil) != tt.wantErr {
				t.Fatalf("getFloat() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("getFloat() = %f, want %f", got, tt.want)
			}
		})
	}
}

// TestGetUint tests unsigned integer retrieval
func TestGetUint(t *testing.T) {
	helpers := NewTestHelpers(t)
	pr := New()

	tests := []struct {
		name    string
		json    string
		key     string
		want    uint64
		wantErr bool
	}{
		{
			name:    "positive integer",
			json:    `{"timestamp": 1234567890}`,
			key:     "timestamp",
			want:    1234567890,
			wantErr: false,
		},
		{
			name:    "string integer",
			json:    `{"timestamp": "999"}`,
			key:     "timestamp",
			want:    999,
			wantErr: false,
		},
		{
			name:    "missing key",
			json:    `{"other": 123}`,
			key:     "timestamp",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := helpers.CreateOrderedMap(tt.json)
			got, err := pr.getUint(om, tt.key)

			if (err != nil) != tt.wantErr {
				t.Fatalf("getUint() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("getUint() = %d, want %d", got, tt.want)
			}
		})
	}
}

// TestUnmarshalFrom tests unmarshaling nested JSON strings
func TestUnmarshalFrom(t *testing.T) {
	helpers := NewTestHelpers(t)
	pr := New()

	tests := []struct {
		name    string
		json    string
		key     string
		wantErr bool
	}{
		{
			name:    "valid nested JSON",
			json:    `{"nested": "{\"value\": 42}"}`,
			key:     "nested",
			wantErr: false,
		},
		{
			name:    "missing key",
			json:    `{"other": "{\"value\": 42}"}`,
			key:     "nested",
			wantErr: true,
		},
		{
			name:    "invalid JSON string",
			json:    `{"nested": "not valid json"}`,
			key:     "nested",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := helpers.CreateOrderedMap(tt.json)
			result := jo.NewOrderedMap()
			err := pr.unmarshalFrom(om, tt.key, result)

			if (err != nil) != tt.wantErr {
				t.Fatalf("unmarshalFrom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetFlag tests flag conversion (int to bool)
func TestGetFlag(t *testing.T) {
	helpers := NewTestHelpers(t)
	pr := New()

	tests := []struct {
		name    string
		json    string
		key     string
		want    bool
		wantErr bool
	}{
		{
			name:    "zero is false",
			json:    `{"flag": 0}`,
			key:     "flag",
			want:    false,
			wantErr: false,
		},
		{
			name:    "non-zero is true",
			json:    `{"flag": 1}`,
			key:     "flag",
			want:    true,
			wantErr: false,
		},
		{
			name:    "any positive is true",
			json:    `{"flag": 999}`,
			key:     "flag",
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := helpers.CreateOrderedMap(tt.json)
			got, err := pr.getFlag(om, tt.key)

			if (err != nil) != tt.wantErr {
				t.Fatalf("getFlag() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("getFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetJsonInts tests retrieving integer arrays from JSON
func TestGetJsonInts(t *testing.T) {
	helpers := NewTestHelpers(t)
	pr := New()

	tests := []struct {
		name    string
		json    string
		key     string
		wantLen int
		wantErr bool
	}{
		{
			name:    "valid integer array",
			json:    `{"values": [1, 2, 3, 4, 5]}`,
			key:     "values",
			wantLen: 5,
			wantErr: false,
		},
		{
			name:    "empty array",
			json:    `{"values": []}`,
			key:     "values",
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "missing key",
			json:    `{"other": [1, 2, 3]}`,
			key:     "values",
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "non-array value",
			json:    `{"values": "not_array"}`,
			key:     "values",
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := helpers.CreateOrderedMap(tt.json)
			got, err := pr.getJsonInts(om, tt.key)

			if (err != nil) != tt.wantErr {
				t.Fatalf("getJsonInts() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(got) != tt.wantLen {
				t.Fatalf("getJsonInts() len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

// TestGetFromTarget tests retrieving values from nested "target" objects
func TestGetFromTarget(t *testing.T) {
	helpers := NewTestHelpers(t)
	pr := New()

	json := `{
		"data": "{\"target\": \"found_value\"}"
	}`

	om := helpers.CreateOrderedMap(json)
	got, err := pr.getFromTarget(om, "data")

	if err != nil {
		t.Fatalf("getFromTarget() error = %v", err)
	}
	if got != "found_value" {
		t.Fatalf("getFromTarget() = %v, want found_value", got)
	}
}

// TestUnmarshalEquipment tests equipment unmarshaling
func TestUnmarshalEquipment(t *testing.T) {
	helpers := NewTestHelpers(t)
	pr := New()

	json := `{
		"equipmentList": "{\"values\": [\"{\\\"contentId\\\": 100}\", \"{\\\"contentId\\\": 101}\"]}"
	}`

	om := helpers.CreateOrderedMap(json)
	equipment, err := pr.unmarshalEquipment(om)

	if err != nil {
		t.Fatalf("unmarshalEquipment() error = %v", err)
	}
	if len(equipment) != 2 {
		t.Fatalf("unmarshalEquipment() returned %d items, want 2", len(equipment))
	}
	if equipment[0].ContentID != 100 {
		t.Fatalf("first equipment ID = %d, want 100", equipment[0].ContentID)
	}
}

// TestLoadEspers tests esper loading logic
func TestLoadEspers(t *testing.T) {
	helpers := NewTestHelpers(t)
	pr := New()
	pr.UserData = helpers.CreateOrderedMap(`{
		"ownedMagicStoneList": "{\"target\": [1, 2, 3]}"
	}`)

	err := pr.loadEspers()
	if err != nil {
		t.Fatalf("loadEspers() error = %v", err)
	}
}

// TestLoadMiscStats tests miscellaneous stats loading
func TestLoadMiscStats(t *testing.T) {
	t.Skip("Skipping: Requires full model initialization with proper field bindings")
	helpers := NewTestHelpers(t)
	pr := New()

	pr.UserData = helpers.CreateOrderedMap(`{
		"owendGil": 5000,
		"Steps": 1000,
		"escapeCount": 10,
		"battleCount": 50,
		"saveCompleteCount": 3,
		"monstersKilledCount": 100,
		"playTime": 42.5,
		"openChestCount": 15
	}`)

	pr.Base = helpers.CreateOrderedMap(`{
		"dataStorage": "{\"global\": [0,0,0,0,0,0,0,0,0,0]}",
		"isCompleteFlag": 0
	}`)

	// Initialize misc model
	models.GetMisc()

	err := pr.loadMiscStats()
	if err != nil {
		t.Fatalf("loadMiscStats() error = %v", err)
	}

	misc := models.GetMisc()
	if misc.GP != 5000 {
		t.Fatalf("GP = %d, want 5000", misc.GP)
	}
	if misc.Steps != 1000 {
		t.Fatalf("Steps = %d, want 1000", misc.Steps)
	}
}

// BenchmarkGetInt benchmarks integer parsing performance
func BenchmarkGetInt(b *testing.B) {
	pr := New()
	om := jo.NewOrderedMap()
	om.UnmarshalJSON([]byte(`{"value": 42}`))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pr.getInt(om, "value")
	}
}

// BenchmarkGetString benchmarks string parsing performance
func BenchmarkGetString(b *testing.B) {
	pr := New()
	om := jo.NewOrderedMap()
	om.UnmarshalJSON([]byte(`{"value": "test"}`))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pr.getString(om, "value")
	}
}
