package batch

import (
	"testing"

	"ffvi_editor/models"
	pri "ffvi_editor/models/pr"
)

// TestRegistryNotEmpty tests that the operation registry is populated
func TestRegistryNotEmpty(t *testing.T) {
	if len(Registry) == 0 {
		t.Error("Registry should not be empty")
	}

	t.Logf("Registry contains %d operations", len(Registry))
}

// TestGetOperationByID tests retrieving operations by ID
func TestGetOperationByID(t *testing.T) {
	tests := []struct {
		id      string
		wantNil bool
	}{
		{"max_all_stats", false},
		{"max_hp_mp", false},
		{"set_level_99", false},
		{"learn_all_magic", false},
		{"add_items_99", false},
		{"clear_inventory", false},
		{"unlock_all_espers", false},
		{"non_existent", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			op := GetOperationByID(tt.id)
			if tt.wantNil && op != nil {
				t.Errorf("GetOperationByID(%q) should return nil", tt.id)
			}
			if !tt.wantNil && op == nil {
				t.Errorf("GetOperationByID(%q) should not return nil", tt.id)
			}
		})
	}
}

// TestGetOperationsByCategory tests filtering by category
func TestGetOperationsByCategory(t *testing.T) {
	tests := []struct {
		category Category
		minCount int
	}{
		{CategoryCharacter, 1},
		{CategoryInventory, 1},
		{CategoryMagic, 1},
		{CategoryEsper, 1},
		{"NonExistent", 0},
	}

	for _, tt := range tests {
		t.Run(string(tt.category), func(t *testing.T) {
			ops := GetOperationsByCategory(tt.category)
			if len(ops) < tt.minCount {
				t.Errorf("Expected at least %d operations, got %d", tt.minCount, len(ops))
			}
		})
	}
}

// TestOperationMaxAllStats tests the max_all_stats operation
func TestOperationMaxAllStats(t *testing.T) {
	op := GetOperationByID("max_all_stats")
	if op == nil {
		t.Fatal("max_all_stats operation not found")
	}

	// Verify operation properties
	if op.Name != "Max All Stats" {
		t.Errorf("Name = %q, want Max All Stats", op.Name)
	}
	if op.Category != CategoryCharacter {
		t.Errorf("Category = %q, want Character", op.Category)
	}

	// Create test context
	ctx := &BatchContext{
		Characters: []*models.Character{
			{
				Name:    "TestChar",
				HP:      models.CurrentMax{Current: 100, Max: 100},
				MP:      models.CurrentMax{Current: 50, Max: 50},
				Vigor:   10,
				Stamina: 10,
				Speed:   10,
				Magic:   10,
			},
		},
		Changes: make(map[string]string),
	}

	// Test Apply
	err := op.Apply(ctx)
	if err != nil {
		t.Errorf("Apply error: %v", err)
	}

	// Verify character was modified
	char := ctx.Characters[0]
	if char.HP.Max != 9999 {
		t.Errorf("HP.Max = %d, want 9999", char.HP.Max)
	}
	if char.MP.Max != 9999 {
		t.Errorf("MP.Max = %d, want 9999", char.MP.Max)
	}
	if char.Vigor != 99 {
		t.Errorf("Vigor = %d, want 99", char.Vigor)
	}

	// Test Preview
	preview := op.Preview(ctx)
	if preview == "" {
		t.Error("Preview should not be empty")
	}
}

// TestOperationMaxHPMP tests the max_hp_mp operation
func TestOperationMaxHPMP(t *testing.T) {
	op := GetOperationByID("max_hp_mp")
	if op == nil {
		t.Fatal("max_hp_mp operation not found")
	}

	ctx := &BatchContext{
		Characters: []*models.Character{
			{
				Name:  "TestChar",
				HP:    models.CurrentMax{Current: 100, Max: 100},
				MP:    models.CurrentMax{Current: 50, Max: 50},
				Vigor: 10,
			},
		},
		Changes: make(map[string]string),
	}

	err := op.Apply(ctx)
	if err != nil {
		t.Errorf("Apply error: %v", err)
	}

	char := ctx.Characters[0]
	if char.HP.Max != 9999 {
		t.Errorf("HP.Max = %d, want 9999", char.HP.Max)
	}
	if char.MP.Max != 9999 {
		t.Errorf("MP.Max = %d, want 9999", char.MP.Max)
	}
	// Vigor should not be changed
	if char.Vigor != 10 {
		t.Errorf("Vigor should not be changed, got %d", char.Vigor)
	}
}

// TestOperationSetLevel99 tests the set_level_99 operation
func TestOperationSetLevel99(t *testing.T) {
	op := GetOperationByID("set_level_99")
	if op == nil {
		t.Fatal("set_level_99 operation not found")
	}

	ctx := &BatchContext{
		Characters: []*models.Character{
			{
				Name:  "TestChar",
				Level: 10,
			},
		},
		Changes: make(map[string]string),
	}

	err := op.Apply(ctx)
	if err != nil {
		t.Errorf("Apply error: %v", err)
	}

	if ctx.Characters[0].Level != 99 {
		t.Errorf("Level = %d, want 99", ctx.Characters[0].Level)
	}
}

// TestOperationLearnAllMagic tests the learn_all_magic operation
func TestOperationLearnAllMagic(t *testing.T) {
	op := GetOperationByID("learn_all_magic")
	if op == nil {
		t.Fatal("learn_all_magic operation not found")
	}

	ctx := &BatchContext{
		Characters: []*models.Character{
			{
				Name:       "TestChar",
				SpellsByID: make(map[int]*models.Spell),
			},
		},
		Changes: make(map[string]string),
	}

	err := op.Apply(ctx)
	if err != nil {
		t.Errorf("Apply error: %v", err)
	}
}

// TestOperationClearInventory tests the clear_inventory operation
func TestOperationClearInventory(t *testing.T) {
	op := GetOperationByID("clear_inventory")
	if op == nil {
		t.Fatal("clear_inventory operation not found")
	}

	ctx := &BatchContext{
		Inventory: &pri.Inventory{},
		Changes:   make(map[string]string),
	}

	err := op.Apply(ctx)
	if err != nil {
		t.Errorf("Apply error: %v", err)
	}
}

// TestOperationAddItems99 tests the add_items_99 operation
func TestOperationAddItems99(t *testing.T) {
	op := GetOperationByID("add_items_99")
	if op == nil {
		t.Fatal("add_items_99 operation not found")
	}

	ctx := &BatchContext{
		Inventory: &pri.Inventory{},
		Changes:   make(map[string]string),
	}

	err := op.Apply(ctx)
	if err != nil {
		t.Errorf("Apply error: %v", err)
	}
}

// TestOperationUnlockAllEspers tests the unlock_all_espers operation
func TestOperationUnlockAllEspers(t *testing.T) {
	op := GetOperationByID("unlock_all_espers")
	if op == nil {
		t.Fatal("unlock_all_espers operation not found")
	}

	ctx := &BatchContext{
		Characters: []*models.Character{
			{Name: "TestChar"},
		},
		Changes: make(map[string]string),
	}

	err := op.Apply(ctx)
	if err != nil {
		t.Errorf("Apply error: %v", err)
	}
}

// TestBatchContextNilCharacters tests operations with nil characters
func TestBatchContextNilCharacters(t *testing.T) {
	op := GetOperationByID("max_all_stats")
	if op == nil {
		t.Fatal("max_all_stats operation not found")
	}

	ctx := &BatchContext{
		Characters: []*models.Character{nil, nil},
		Changes:    make(map[string]string),
	}

	// Should not panic with nil characters
	err := op.Apply(ctx)
	if err != nil {
		t.Errorf("Apply error with nil characters: %v", err)
	}
}

// TestBatchContextEmptyCharacters tests operations with empty character list
func TestBatchContextEmptyCharacters(t *testing.T) {
	op := GetOperationByID("max_all_stats")
	if op == nil {
		t.Fatal("max_all_stats operation not found")
	}

	ctx := &BatchContext{
		Characters: []*models.Character{},
		Changes:    make(map[string]string),
	}

	err := op.Apply(ctx)
	if err != nil {
		t.Errorf("Apply error with empty characters: %v", err)
	}

	preview := op.Preview(ctx)
	if preview == "" {
		t.Error("Preview should not be empty even with no characters")
	}
}

// TestCategoryConstants tests category constant values
func TestCategoryConstants(t *testing.T) {
	if CategoryCharacter != "Character" {
		t.Errorf("CategoryCharacter = %q, want Character", CategoryCharacter)
	}
	if CategoryInventory != "Inventory" {
		t.Errorf("CategoryInventory = %q, want Inventory", CategoryInventory)
	}
	if CategoryMagic != "Magic" {
		t.Errorf("CategoryMagic = %q, want Magic", CategoryMagic)
	}
	if CategoryEsper != "Esper" {
		t.Errorf("CategoryEsper = %q, want Esper", CategoryEsper)
	}
}

// BenchmarkGetOperationByID benchmarks operation lookup
func BenchmarkGetOperationByID(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetOperationByID("max_all_stats")
	}
}

// BenchmarkGetOperationsByCategory benchmarks category filtering
func BenchmarkGetOperationsByCategory(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetOperationsByCategory(CategoryCharacter)
	}
}

// BenchmarkApplyMaxAllStats benchmarks the max_all_stats operation
func BenchmarkApplyMaxAllStats(b *testing.B) {
	ctx := &BatchContext{
		Characters: []*models.Character{
			{
				Name:    "TestChar",
				HP:      models.CurrentMax{Current: 100, Max: 100},
				MP:      models.CurrentMax{Current: 50, Max: 50},
				Vigor:   10,
				Stamina: 10,
				Speed:   10,
				Magic:   10,
			},
		},
		Changes: make(map[string]string),
	}

	op := GetOperationByID("max_all_stats")
	if op == nil {
		b.Fatal("max_all_stats operation not found")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = op.Apply(ctx)
	}
}
