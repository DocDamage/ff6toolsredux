package consts

import (
	"testing"
)

// TestStatusEffectsStructure tests StatusEffects structure
func TestStatusEffectsStructure(t *testing.T) {
	s := NewStatusEffects()

	if s == nil {
		t.Fatal("NewStatusEffects() returned nil")
	}

	// Result is a slice, not a struct with Effects field
	if len(s) == 0 {
		t.Error("StatusEffects slice is empty")
	}
}

// TestStatusEffectsCount tests that we have the expected number of status effects
func TestStatusEffectsCount(t *testing.T) {
	s := NewStatusEffects()

	// FF6 has 9 status effects defined
	expectedCount := len(StatusEffects)
	if len(s) != expectedCount {
		t.Errorf("StatusEffects count = %d, want %d", len(s), expectedCount)
	}
}

// TestStatusEffectsNames tests that status effects have names
func TestStatusEffectsNames(t *testing.T) {
	s := NewStatusEffects()

	for _, se := range s {
		if se.Name == "" {
			t.Error("StatusEffect has empty name")
		}
	}
}

// TestStatusEffectsSlots tests that status effects have valid slots
func TestStatusEffectsSlots(t *testing.T) {
	s := NewStatusEffects()

	for _, se := range s {
		if se.Slot < 0 {
			t.Errorf("StatusEffect %s has negative slot: %d", se.Name, se.Slot)
		}
	}
}

// TestStatusEffectsMasks tests that status effects have valid masks
func TestStatusEffectsMasks(t *testing.T) {
	s := NewStatusEffects()

	for _, se := range s {
		if se.Mask == 0 {
			t.Errorf("StatusEffect %s has zero mask", se.Name)
		}
	}
}

// TestNameSlotMask8Structure tests NameSlotMask8 structure
func TestNameSlotMask8Structure(t *testing.T) {
	mask := &NameSlotMask8{
		Name:    "Test",
		Slot:    1,
		Mask:    0xFF,
		Checked: true,
	}

	if mask.Name != "Test" {
		t.Errorf("NameSlotMask8.Name = %s, want Test", mask.Name)
	}
	if mask.Slot != 1 {
		t.Errorf("NameSlotMask8.Slot = %d, want 1", mask.Slot)
	}
	if mask.Mask != 0xFF {
		t.Errorf("NameSlotMask8.Mask = %d, want 255", mask.Mask)
	}
	if !mask.Checked {
		t.Error("NameSlotMask8.Checked = false, want true")
	}
}

// TestNewNameSlotMask8 tests NameSlotMask8 creation
func TestNewNameSlotMask8(t *testing.T) {
	mask := NewNameSlotMask8("TestName", 2, 0x04)

	if mask.Name != "TestName" {
		t.Errorf("NameSlotMask8.Name = %s, want TestName", mask.Name)
	}
	if mask.Slot != 2 {
		t.Errorf("NameSlotMask8.Slot = %d, want 2", mask.Slot)
	}
	if mask.Mask != 0x04 {
		t.Errorf("NameSlotMask8.Mask = %d, want 4", mask.Mask)
	}
	if mask.Checked {
		t.Error("NameSlotMask8.Checked should be false by default")
	}
}

// TestNewNameSlotMask8s tests creating multiple NameSlotMask8s
func TestNewNameSlotMask8s(t *testing.T) {
	names := []string{"A", "B", "C", "D"}
	masks := NewNameSlotMask8s(names...)

	if len(masks) != len(names) {
		t.Errorf("Created %d masks, want %d", len(masks), len(names))
	}

	for i, mask := range masks {
		if mask.Name != names[i] {
			t.Errorf("Mask[%d].Name = %s, want %s", i, mask.Name, names[i])
		}
	}
}

// TestNameSlotMask8SetChecked tests SetChecked method
func TestNameSlotMask8SetChecked(t *testing.T) {
	mask := &NameSlotMask8{
		Slot: 0,
		Mask: 0x01,
	}

	// Set checked with byte that includes the mask
	mask.SetChecked(0x01)
	if !mask.Checked {
		t.Error("SetChecked(0x01) should set Checked to true")
	}

	// Reset
	mask.Checked = false

	// Set checked with byte that doesn't include the mask
	mask.SetChecked(0x02)
	if mask.Checked {
		t.Error("SetChecked(0x02) should leave Checked as false")
	}
}

// TestGenerateBytes tests GenerateBytes function
func TestGenerateBytes(t *testing.T) {
	masks := []*NameSlotMask8{
		{Name: "A", Slot: 0, Mask: 0x01, Checked: true},
		{Name: "B", Slot: 0, Mask: 0x02, Checked: false},
		{Name: "C", Slot: 0, Mask: 0x04, Checked: true},
	}

	result := GenerateBytes(masks)

	// Should have 1 byte (slot 0)
	if len(result) != 1 {
		t.Errorf("GenerateBytes returned %d bytes, want 1", len(result))
	}

	// Should have bits 0 and 2 set: 0x01 | 0x04 = 0x05
	expected := byte(0x05)
	if result[0] != expected {
		t.Errorf("GenerateBytes result = 0x%02X, want 0x%02X", result[0], expected)
	}
}

// TestGenerateBytesMultiSlot tests GenerateBytes with multiple slots
func TestGenerateBytesMultiSlot(t *testing.T) {
	masks := []*NameSlotMask8{
		{Name: "A", Slot: 0, Mask: 0x80, Checked: true}, // Last bit of slot 0
		{Name: "B", Slot: 1, Mask: 0x01, Checked: true}, // First bit of slot 1
	}

	result := GenerateBytes(masks)

	// Should have 2 bytes (slot 0 and 1)
	if len(result) != 2 {
		t.Errorf("GenerateBytes returned %d bytes, want 2", len(result))
	}

	if result[0] != 0x80 {
		t.Errorf("Byte 0 = 0x%02X, want 0x80", result[0])
	}
	if result[1] != 0x01 {
		t.Errorf("Byte 1 = 0x%02X, want 0x01", result[1])
	}
}

// TestNameValueStructure tests NameValue structure
func TestNameValueStructure(t *testing.T) {
	nv := &NameValue{
		Name:  "Test",
		Value: 42,
	}

	if nv.Name != "Test" {
		t.Errorf("NameValue.Name = %s, want Test", nv.Name)
	}
	if nv.Value != 42 {
		t.Errorf("NameValue.Value = %d, want 42", nv.Value)
	}
}

// TestNewNameValue tests NewNameValue function
func TestNewNameValue(t *testing.T) {
	nv := NewNameValue("TestName", 100)

	if nv.Name != "TestName" {
		t.Errorf("NameValue.Name = %s, want TestName", nv.Name)
	}
	if nv.Value != 100 {
		t.Errorf("NameValue.Value = %d, want 100", nv.Value)
	}
}

// TestNewValueName tests NewValueName function
func TestNewValueName(t *testing.T) {
	nv := NewValueName(200, "TestName")

	if nv.Name != "TestName" {
		t.Errorf("NameValue.Name = %s, want TestName", nv.Name)
	}
	if nv.Value != 200 {
		t.Errorf("NameValue.Value = %d, want 200", nv.Value)
	}
}

// TestNameValueCheckedStructure tests NameValueChecked structure
func TestNameValueCheckedStructure(t *testing.T) {
	nvc := &NameValueChecked{
		NameValue: NameValue{Name: "Test", Value: 42},
		Checked:   true,
	}

	if nvc.Name != "Test" {
		t.Errorf("NameValueChecked.Name = %s, want Test", nvc.Name)
	}
	if nvc.Value != 42 {
		t.Errorf("NameValueChecked.Value = %d, want 42", nvc.Value)
	}
	if !nvc.Checked {
		t.Error("NameValueChecked.Checked = false, want true")
	}
}

// TestNewNameValueChecked tests NewNameValueChecked function
func TestNewNameValueChecked(t *testing.T) {
	nvc := NewNameValueChecked("TestName", 100)

	if nvc.Name != "TestName" {
		t.Errorf("NameValueChecked.Name = %s, want TestName", nvc.Name)
	}
	if nvc.Value != 100 {
		t.Errorf("NameValueChecked.Value = %d, want 100", nvc.Value)
	}
	if nvc.Checked {
		t.Error("NameValueChecked.Checked should be false by default")
	}
}

// TestNewNameValues tests NewNameValues function
func TestNewNameValues(t *testing.T) {
	names := []string{"A", "B", "C"}
	nvs := NewNameValues(names...)

	if len(nvs) != len(names) {
		t.Errorf("Created %d NameValues, want %d", len(nvs), len(names))
	}

	for i, nv := range nvs {
		if nv.Name != names[i] {
			t.Errorf("NameValue[%d].Name = %s, want %s", i, nv.Name, names[i])
		}
		if nv.Value != i {
			t.Errorf("NameValue[%d].Value = %d, want %d", i, nv.Value, i)
		}
	}
}

// TestSortByName tests SortByName function
func TestSortByName(t *testing.T) {
	items := []*NameSlotMask8{
		{Name: "Charlie"},
		{Name: "Alpha"},
		{Name: "Bravo"},
	}

	sorted := SortByName(items)

	if sorted[0].Name != "Alpha" {
		t.Errorf("After sort, first item = %s, want Alpha", sorted[0].Name)
	}
	if sorted[1].Name != "Bravo" {
		t.Errorf("After sort, second item = %s, want Bravo", sorted[1].Name)
	}
	if sorted[2].Name != "Charlie" {
		t.Errorf("After sort, third item = %s, want Charlie", sorted[2].Name)
	}
}

// TestSortByNameEmpty tests sorting empty slice
func TestSortByNameEmpty(t *testing.T) {
	items := []*NameSlotMask8{}
	sorted := SortByName(items)

	if len(sorted) != 0 {
		t.Error("Empty slice should remain empty")
	}
}

// TestSortByNameSingleItem tests sorting single item
func TestSortByNameSingleItem(t *testing.T) {
	items := []*NameSlotMask8{
		{Name: "Only"},
	}
	sorted := SortByName(items)

	if sorted[0].Name != "Only" {
		t.Error("Single item should remain unchanged")
	}
}

// TestSortByNameDuplicateNames tests sorting with duplicate names
func TestSortByNameDuplicateNames(t *testing.T) {
	items := []*NameSlotMask8{
		{Name: "Alpha", Mask: 1},
		{Name: "Alpha", Mask: 2},
		{Name: "Beta", Mask: 3},
	}

	sorted := SortByName(items)

	// Items with same name - one will be overwritten in lookup
	// This is expected behavior
	if len(sorted) != 3 {
		t.Errorf("Sorted length = %d, want 3", len(sorted))
	}
}

// TestSortByNameChecked tests SortByNameChecked function
func TestSortByNameChecked(t *testing.T) {
	items := []*NameValueChecked{
		{NameValue: NameValue{Name: "Charlie"}},
		{NameValue: NameValue{Name: "Alpha"}},
		{NameValue: NameValue{Name: "Bravo"}},
	}

	sorted := SortByNameChecked(items)

	if sorted[0].Name != "Alpha" {
		t.Errorf("After sort, first item = %s, want Alpha", sorted[0].Name)
	}
	if sorted[1].Name != "Bravo" {
		t.Errorf("After sort, second item = %s, want Bravo", sorted[1].Name)
	}
	if sorted[2].Name != "Charlie" {
		t.Errorf("After sort, third item = %s, want Charlie", sorted[2].Name)
	}
}

// TestSortByNameValue tests SortByNameValue function
func TestSortByNameValue(t *testing.T) {
	items := []*NameValue{
		{Name: "Charlie"},
		{Name: "Alpha"},
		{Name: "Bravo"},
	}

	sorted := SortByNameValue(items)

	if sorted[0].Name != "Alpha" {
		t.Errorf("After sort, first item = %s, want Alpha", sorted[0].Name)
	}
	if sorted[1].Name != "Bravo" {
		t.Errorf("After sort, second item = %s, want Bravo", sorted[1].Name)
	}
	if sorted[2].Name != "Charlie" {
		t.Errorf("After sort, third item = %s, want Charlie", sorted[2].Name)
	}
}

// TestStatusEffectsList tests that StatusEffects list is defined
func TestStatusEffectsList(t *testing.T) {
	expectedEffects := []string{
		"Darkness",
		"Zombie",
		"Poison",
		"Magiteck",
		"Invisible",
		"Imp",
		"Stone",
		"Wounded",
		"Float",
	}

	if len(StatusEffects) != len(expectedEffects) {
		t.Errorf("StatusEffects has %d items, want %d", len(StatusEffects), len(expectedEffects))
	}

	for i, effect := range expectedEffects {
		if StatusEffects[i] != effect {
			t.Errorf("StatusEffects[%d] = %s, want %s", i, StatusEffects[i], effect)
		}
	}
}

// BenchmarkNewStatusEffects benchmarks status effects creation
func BenchmarkNewStatusEffects(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewStatusEffects()
	}
}

// BenchmarkSortByName benchmarks name sorting
func BenchmarkSortByName(b *testing.B) {
	items := []*NameSlotMask8{
		{Name: "Charlie"},
		{Name: "Alpha"},
		{Name: "Bravo"},
		{Name: "Delta"},
		{Name: "Echo"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Make a copy for each iteration
		testItems := make([]*NameSlotMask8, len(items))
		copy(testItems, items)
		SortByName(testItems)
	}
}

// BenchmarkGenerateBytes benchmarks GenerateBytes
func BenchmarkGenerateBytes(b *testing.B) {
	masks := []*NameSlotMask8{
		{Name: "A", Slot: 0, Mask: 0x01, Checked: true},
		{Name: "B", Slot: 0, Mask: 0x02, Checked: false},
		{Name: "C", Slot: 0, Mask: 0x04, Checked: true},
		{Name: "D", Slot: 1, Mask: 0x01, Checked: true},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateBytes(masks)
	}
}
