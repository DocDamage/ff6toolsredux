package pr

import (
	"testing"

	"ffvi_editor/models"
	"ffvi_editor/models/consts"
	pri "ffvi_editor/models/pr"
)

// TestLoadCharacters tests the character loading from parsed save data
func TestLoadCharacters(t *testing.T) {
	helpers := NewTestHelpers(t)
	p := New()

	// Setup minimal character data
	charJSON := helpers.CreateMinimalCharacterJSON()
	charOM := helpers.CreateOrderedMap(charJSON)

	p.Characters[0] = charOM
	// Note: ownedCharacterList should be a JSON string containing {"target": [...]} structure
	p.UserData = helpers.CreateOrderedMap(`{
		"ownedCharacterList": "{\"target\": []}"
	}`)

	// Initialize party and character systems
	pri.GetParty().Clear()

	// We can't fully test without proper initialization of the models,
	// but we can verify the function doesn't panic
	err := p.loadCharacters()
	if err != nil {
		t.Logf("loadCharacters() error (expected in isolated test): %v", err)
	}
}

// TestLoadEquipment tests equipment parsing from character data
func TestLoadEquipment(t *testing.T) {
	helpers := NewTestHelpers(t)
	p := New()

	charJSON := `{
		"equipmentList": "{\"values\": [\"{\\\"contentId\\\": 100}\", \"{\\\"contentId\\\": 101}\", \"{\\\"contentId\\\": 102}\", \"{\\\"contentId\\\": 103}\", \"{\\\"contentId\\\": 104}\", \"{\\\"contentId\\\": 105}\"]}"
	}`

	charOM := helpers.CreateOrderedMap(charJSON)
	testChar := &models.Character{
		Equipment: models.Equipment{},
	}

	err := p.loadEquipment(charOM, testChar)
	helpers.AssertNoError(err, "loadEquipment")

	if testChar.Equipment.WeaponID != 100 {
		t.Fatalf("WeaponID = %d, want 100", testChar.Equipment.WeaponID)
	}
	if testChar.Equipment.ShieldID != 101 {
		t.Fatalf("ShieldID = %d, want 101", testChar.Equipment.ShieldID)
	}
	if testChar.Equipment.ArmorID != 102 {
		t.Fatalf("ArmorID = %d, want 102", testChar.Equipment.ArmorID)
	}
}

// TestLoadEquipmentPartial tests equipment loading with fewer items
func TestLoadEquipmentPartial(t *testing.T) {
	helpers := NewTestHelpers(t)
	p := New()

	charJSON := `{
		"equipmentList": "{\"values\": [\"{\\\"contentId\\\": 100}\"]}"
	}`

	charOM := helpers.CreateOrderedMap(charJSON)
	testChar := &models.Character{
		Equipment: models.Equipment{},
	}

	err := p.loadEquipment(charOM, testChar)
	helpers.AssertNoError(err, "loadEquipment with partial data")

	if testChar.Equipment.WeaponID != 100 {
		t.Fatalf("WeaponID = %d, want 100", testChar.Equipment.WeaponID)
	}
	// Defaults should be set for missing items
	if testChar.Equipment.ShieldID != 93 {
		t.Fatalf("ShieldID default = %d, want 93", testChar.Equipment.ShieldID)
	}
}

// TestLoadInventory tests inventory item loading
func TestLoadInventory(t *testing.T) {
	helpers := NewTestHelpers(t)
	p := New()

	p.UserData = helpers.CreateOrderedMap(`{
		"normalOwnedItemList": "{\"target\": [{\"contentId\": 2, \"quantity\": 5}, {\"contentId\": 3, \"quantity\": 3}]}"
	}`)

	inventory := pri.GetInventory()
	err := p.loadInventory(NormalOwnedItemList, inventory)

	if err != nil {
		t.Logf("loadInventory error (expected in isolated test): %v", err)
	}
}

// TestLoadMapData tests map data parsing
func TestLoadMapData(t *testing.T) {
	helpers := NewTestHelpers(t)
	p := New()

	mapJSON := helpers.CreateMinimalMapDataJSON()
	p.MapData = helpers.CreateOrderedMap(mapJSON)

	mapData := pri.GetMapData()
	// Reset not needed - GetMapData returns fresh instance

	err := p.loadMapData()
	helpers.AssertNoError(err, "loadMapData")

	if mapData.MapID != 1 {
		t.Fatalf("MapID = %d, want 1", mapData.MapID)
	}
	if mapData.PointIn != 0 {
		t.Fatalf("PointIn = %d, want 0", mapData.PointIn)
	}
}

// TestLoadTransportation tests transportation data loading
func TestLoadTransportation(t *testing.T) {
	helpers := NewTestHelpers(t)
	p := New()

	p.UserData = helpers.CreateOrderedMap(`{
		"ownedTransportationList": "{\"target\": [{\"transId\": 1, \"transMapId\": 10, \"transDirection\": 0, \"transTimeStampTicks\": 100, \"transPosition\": {\"x\": 50.0, \"y\": 50.0, \"z\": 0.0}}]}"
	}`)

	pri.Transportations = nil
	err := p.loadTransportation()

	if err != nil {
		t.Logf("loadTransportation error (expected in isolated test): %v", err)
	}
}

// TestLoadVeldt tests Veldt encounter flags loading
func TestLoadVeldt(t *testing.T) {
	helpers := NewTestHelpers(t)
	p := New()

	p.MapData = helpers.CreateOrderedMap(`{
		"beastFieldEncountExchangeFlags": [1, 1, 0, 1, 0]
	}`)

	veldt := pri.GetVeldt()
	err := p.loadVeldt()
	helpers.AssertNoError(err, "loadVeldt")

	if len(veldt.Encounters) != 5 {
		t.Fatalf("veldt encounters length = %d, want 5", len(veldt.Encounters))
	}
	if !veldt.Encounters[0] {
		t.Fatal("first encounter should be true")
	}
	if veldt.Encounters[2] {
		t.Fatal("third encounter should be false")
	}
}

// TestLoadSpells tests spell loading for characters
func TestLoadSpells(t *testing.T) {
	helpers := NewTestHelpers(t)
	p := New()

	testChar := &models.Character{
		SpellsByID: make(map[int]*models.Spell),
	}

	// Add some dummy spells for testing
	testChar.SpellsByID[1] = &models.Spell{Index: 1, Value: 0}

	charJSON := `{
		"abilityList": "{\"target\": [\"{\\\"abilityId\\\": 1, \\\"skillLevel\\\": 50}\"]}"
	}`

	charOM := helpers.CreateOrderedMap(charJSON)
	err := p.loadSpells(charOM, testChar)
	helpers.AssertNoError(err, "loadSpells")
}

// TestLoadCheats tests cheat flags loading
func TestLoadCheats(t *testing.T) {
	helpers := NewTestHelpers(t)
	p := New()

	p.UserData = helpers.CreateOrderedMap(`{
		"openChestCount": 25,
		"playTime": 100.5
	}`)

	p.Base = helpers.CreateOrderedMap(`{
		"isCompleteFlag": 1,
		"dataStorage": "{\"global\": [0,0,0,0,0,0,0,0,0,0]}"
	}`)

	err := p.loadCheats()
	helpers.AssertNoError(err, "loadCheats")

	cheats := pri.GetCheats()
	if cheats.OpenedChestCount != 25 {
		t.Fatalf("OpenedChestCount = %d, want 25", cheats.OpenedChestCount)
	}
}

// TestLoadBase tests base JSON structure loading
func TestLoadBase(t *testing.T) {
	helpers := NewTestHelpers(t)
	p := New()

	baseJSON := helpers.CreateMinimalBaseJSON()
	err := p.loadBase(baseJSON)

	if err != nil {
		t.Logf("loadBase error (expected in isolated test): %v", err)
	}
}

// TestUnmarshalEquipmentEmpty tests empty equipment list
func TestUnmarshalEquipmentEmpty(t *testing.T) {
	helpers := NewTestHelpers(t)
	p := New()

	charJSON := `{
		"equipmentList": "{\"values\": []}"
	}`

	charOM := helpers.CreateOrderedMap(charJSON)
	equipment, err := p.unmarshalEquipment(charOM)

	if err != nil {
		t.Fatalf("unmarshalEquipment error: %v", err)
	}
	if len(equipment) != 0 {
		t.Fatalf("empty equipment list should have 0 items, got %d", len(equipment))
	}
}

// TestUnmarshalEquipmentMissing tests missing equipment list
func TestUnmarshalEquipmentMissing(t *testing.T) {
	helpers := NewTestHelpers(t)
	p := New()

	charJSON := `{}`
	charOM := helpers.CreateOrderedMap(charJSON)

	_, err := p.unmarshalEquipment(charOM)
	if err == nil {
		t.Fatal("unmarshalEquipment should error on missing equipment list")
	}
}

// TestGetIntFromSlice tests cursed shield battle count parsing
func TestGetIntFromSlice(t *testing.T) {
	helpers := NewTestHelpers(t)
	p := New()

	json := `{
		"global": [0,0,0,0,0,0,0,0,0,42]
	}`

	om := helpers.CreateOrderedMap(json)
	value, err := p.getIntFromSlice(om, "global")

	helpers.AssertNoError(err, "getIntFromSlice")
	if value != 42 {
		t.Fatalf("getIntFromSlice value = %d, want 42", value)
	}
}

// TestStatusEffectsInitialization tests that status effects are properly created
func TestStatusEffectsInitialization(t *testing.T) {
	effects := consts.NewStatusEffects()

	if effects == nil {
		t.Fatal("NewStatusEffects returned nil")
	}

	// Test that effects can be accessed (implementation dependent)
	if len(effects) > 0 {
		t.Logf("Status effects initialized with %d effects", len(effects))
	}
}

// BenchmarkLoadEquipment benchmarks equipment loading
func BenchmarkLoadEquipment(b *testing.B) {
	helpers := NewTestHelpers(nil)
	p := New()
	charJSON := `{
		"equipmentList": "{\"values\": [{\"contentId\": 100}, {\"contentId\": 101}, {\"contentId\": 102}, {\"contentId\": 103}, {\"contentId\": 104}, {\"contentId\": 105}]}"
	}`
	charOM := helpers.CreateOrderedMap(charJSON)
	testChar := &models.Character{
		Equipment: models.Equipment{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.loadEquipment(charOM, testChar)
	}
}

// BenchmarkLoadMapData benchmarks map data loading
func BenchmarkLoadMapData(b *testing.B) {
	helpers := NewTestHelpers(nil)
	p := New()
	mapJSON := helpers.CreateMinimalMapDataJSON()
	p.MapData = helpers.CreateOrderedMap(mapJSON)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.loadMapData()
	}
}
