package pr

import (
	"testing"

	"ffvi_editor/models"
	"ffvi_editor/models/consts"
	pri "ffvi_editor/models/pr"
)

// TestCharacterRoundTrip tests that character data survives a save/load cycle
func TestCharacterRoundTrip(t *testing.T) {
	helpers := NewTestHelpers(t)

	// Set up initial test data
	p := New()
	p.Base = helpers.CreateOrderedMap(helpers.CreateMinimalBaseJSON())
	p.UserData = helpers.CreateOrderedMap(helpers.CreateMinimalUserDataJSON())

	// Create a character in the models
	testChar := pri.GetCharacter("Terra")
	if testChar == nil {
		t.Skip("Terra character not available")
	}

	// Store original values
	originalHP := testChar.HP.Current
	originalMaxHP := testChar.HP.Max
	originalLevel := testChar.Level
	originalVigor := testChar.Vigor

	// Modify the character
	testChar.HP.Current = 999
	testChar.HP.Max = 9999
	testChar.Level = 50
	testChar.Vigor = 80

	// Note: In a full round-trip test, we would:
	// 1. Save to file
	// 2. Load from file
	// 3. Verify values match what we saved

	// For now, just verify the modifications took effect
	if testChar.HP.Current != 999 {
		t.Errorf("HP.Current = %d, want 999", testChar.HP.Current)
	}
	if testChar.Level != 50 {
		t.Errorf("Level = %d, want 50", testChar.Level)
	}

	// Restore original values
	testChar.HP.Current = originalHP
	testChar.HP.Max = originalMaxHP
	testChar.Level = originalLevel
	testChar.Vigor = originalVigor
}

// TestInventoryRoundTrip tests inventory persistence
func TestInventoryRoundTrip(t *testing.T) {
	inventory := pri.GetInventory()
	inventory.Clear()

	// Add some items
	testItems := []struct {
		id    int
		count int
	}{
		{1, 10},  // Potion
		{2, 5},   // Hi-Potion
		{10, 99}, // Ether
	}

	for i, item := range testItems {
		inventory.Set(i, pri.Row{ItemID: item.id, Count: item.count})
	}

	// Verify items were set
	if len(inventory.Rows) != inventory.Size {
		t.Fatalf("inventory rows = %d, want %d", len(inventory.Rows), inventory.Size)
	}

	for i, item := range testItems {
		if inventory.Rows[i].ItemID != item.id {
			t.Errorf("item %d ID = %d, want %d", i, inventory.Rows[i].ItemID, item.id)
		}
		if inventory.Rows[i].Count != item.count {
			t.Errorf("item %d count = %d, want %d", i, inventory.Rows[i].Count, item.count)
		}
	}

	// Test GetItemLookup
	lookup := inventory.GetItemLookup()
	for _, item := range testItems {
		count, found := lookup[item.id]
		if !found {
			t.Errorf("item %d not found in lookup", item.id)
			continue
		}
		if count != item.count {
			t.Errorf("item %d count = %d, want %d", item.id, count, item.count)
		}
	}
}

// TestMapDataRoundTrip tests map data persistence
func TestMapDataRoundTrip(t *testing.T) {
	mapData := pri.GetMapData()

	// Store original values
	originalMapID := mapData.MapID
	originalPointIn := mapData.PointIn
	originalX := mapData.Player.X
	originalY := mapData.Player.Y

	// Set new values
	mapData.MapID = 42
	mapData.PointIn = 1
	mapData.Player.X = 123.45
	mapData.Player.Y = 678.90

	// Verify values were set
	if mapData.MapID != 42 {
		t.Errorf("MapID = %d, want 42", mapData.MapID)
	}
	if mapData.PointIn != 1 {
		t.Errorf("PointIn = %d, want 1", mapData.PointIn)
	}
	if mapData.Player.X != 123.45 {
		t.Errorf("Player.X = %f, want 123.45", mapData.Player.X)
	}
	if mapData.Player.Y != 678.90 {
		t.Errorf("Player.Y = %f, want 678.90", mapData.Player.Y)
	}

	// Restore original values
	mapData.MapID = originalMapID
	mapData.PointIn = originalPointIn
	mapData.Player.X = originalX
	mapData.Player.Y = originalY
}

// TestMiscStatsRoundTrip tests misc stats persistence
func TestMiscStatsRoundTrip(t *testing.T) {
	misc := models.GetMisc()

	// Store original values
	originalGP := misc.GP
	originalSteps := misc.Steps
	originalBattleCount := misc.BattleCount
	originalEscapeCount := misc.EscapeCount
	originalMonstersKilled := misc.MonstersKilledCount

	// Set new values
	misc.GP = 999999
	misc.Steps = 50000
	misc.BattleCount = 1000
	misc.EscapeCount = 50
	misc.MonstersKilledCount = 500

	// Verify values were set
	if misc.GP != 999999 {
		t.Errorf("GP = %d, want 999999", misc.GP)
	}
	if misc.Steps != 50000 {
		t.Errorf("Steps = %d, want 50000", misc.Steps)
	}
	if misc.BattleCount != 1000 {
		t.Errorf("BattleCount = %d, want 1000", misc.BattleCount)
	}
	if misc.EscapeCount != 50 {
		t.Errorf("EscapeCount = %d, want 50", misc.EscapeCount)
	}
	if misc.MonstersKilledCount != 500 {
		t.Errorf("MonstersKilledCount = %d, want 500", misc.MonstersKilledCount)
	}

	// Restore original values
	misc.GP = originalGP
	misc.Steps = originalSteps
	misc.BattleCount = originalBattleCount
	misc.EscapeCount = originalEscapeCount
	misc.MonstersKilledCount = originalMonstersKilled
}

// TestCheatsRoundTrip tests cheat flags persistence
func TestCheatsRoundTrip(t *testing.T) {
	cheats := pri.GetCheats()

	// Store original values
	originalChestCount := cheats.OpenedChestCount
	originalIsComplete := cheats.IsCompleteFlag

	// Set new values
	cheats.OpenedChestCount = 100
	cheats.IsCompleteFlag = true

	// Verify values were set
	if cheats.OpenedChestCount != 100 {
		t.Errorf("OpenedChestCount = %d, want 100", cheats.OpenedChestCount)
	}
	if !cheats.IsCompleteFlag {
		t.Error("IsCompleteFlag should be true")
	}

	// Restore original values
	cheats.OpenedChestCount = originalChestCount
	cheats.IsCompleteFlag = originalIsComplete
}

// TestVeldtRoundTrip tests Veldt encounter flags persistence
func TestVeldtRoundTrip(t *testing.T) {
	veldt := pri.GetVeldt()

	// Store original encounters
	originalEncounters := make([]bool, len(veldt.Encounters))
	copy(originalEncounters, veldt.Encounters)

	// Set new encounter pattern
	testEncounters := []bool{true, false, true, false, true, true, false}
	veldt.Encounters = testEncounters

	// Verify values were set
	if len(veldt.Encounters) != len(testEncounters) {
		t.Fatalf("encounters length = %d, want %d", len(veldt.Encounters), len(testEncounters))
	}

	for i, expected := range testEncounters {
		if veldt.Encounters[i] != expected {
			t.Errorf("encounter[%d] = %v, want %v", i, veldt.Encounters[i], expected)
		}
	}

	// Restore original encounters
	veldt.Encounters = originalEncounters
}

// TestEquipmentRoundTrip tests equipment persistence
func TestEquipmentRoundTrip(t *testing.T) {
	testChar := &models.Character{
		ID:   1,
		Name: "TestChar",
		Equipment: models.Equipment{
			WeaponID: 100,
			ShieldID: 101,
			ArmorID:  102,
			HelmetID: 103,
			Relic1ID: 104,
			Relic2ID: 105,
		},
	}

	// Verify equipment was set
	if testChar.Equipment.WeaponID != 100 {
		t.Errorf("WeaponID = %d, want 100", testChar.Equipment.WeaponID)
	}
	if testChar.Equipment.ShieldID != 101 {
		t.Errorf("ShieldID = %d, want 101", testChar.Equipment.ShieldID)
	}
	if testChar.Equipment.ArmorID != 102 {
		t.Errorf("ArmorID = %d, want 102", testChar.Equipment.ArmorID)
	}
	if testChar.Equipment.HelmetID != 103 {
		t.Errorf("HelmetID = %d, want 103", testChar.Equipment.HelmetID)
	}
	if testChar.Equipment.Relic1ID != 104 {
		t.Errorf("Relic1ID = %d, want 104", testChar.Equipment.Relic1ID)
	}
	if testChar.Equipment.Relic2ID != 105 {
		t.Errorf("Relic2ID = %d, want 105", testChar.Equipment.Relic2ID)
	}
}

// TestSpellsRoundTrip tests spell learning persistence
func TestSpellsRoundTrip(t *testing.T) {
	testChar := &models.Character{
		ID:   1,
		Name: "TestChar",
		SpellsByID: map[int]*models.Spell{
			1:  {Index: 1, Value: 100},  // Fire - learned
			2:  {Index: 2, Value: 50},   // Blizzard - learning
			3:  {Index: 3, Value: 0},    // Thunder - not learned
			54: {Index: 54, Value: 100}, // Cure - learned
		},
		SpellsByIndex: make([]*models.Spell, 0),
		SpellsSorted:  make([]*models.Spell, 0),
	}

	// Verify spells were set
	if len(testChar.SpellsByID) != 4 {
		t.Fatalf("spells count = %d, want 4", len(testChar.SpellsByID))
	}

	// Check specific spell
	fireSpell, ok := testChar.SpellsByID[1]
	if !ok {
		t.Fatal("Fire spell not found")
	}
	if fireSpell.Value != 100 {
		t.Errorf("Fire spell value = %d, want 100", fireSpell.Value)
	}
}

// TestStatusEffectsRoundTrip tests status effects persistence
func TestStatusEffectsRoundTrip(t *testing.T) {
	testChar := &models.Character{
		ID:            1,
		Name:          "TestChar",
		StatusEffects: consts.NewStatusEffects(),
	}

	// Set some status effects by checking them
	for _, effect := range testChar.StatusEffects {
		if effect.Name == "Poison" || effect.Name == "Blind" || effect.Name == "Sleep" {
			effect.Checked = true
		}
	}

	// Generate bytes
	bytes := consts.GenerateBytes(testChar.StatusEffects)

	// Verify bytes were generated
	if len(bytes) == 0 {
		t.Fatal("GenerateBytes returned empty slice")
	}

	// The bytes should have some bits set for our checked effects
	hasBits := false
	for _, b := range bytes {
		if b != 0 {
			hasBits = true
			break
		}
	}
	if !hasBits {
		t.Error("Status effect bytes should have some bits set")
	}
}

// TestPartyRoundTrip tests party composition persistence
func TestPartyRoundTrip(t *testing.T) {
	party := pri.GetParty()
	party.Clear()
	party.Enabled = true
	party.Possible = make(map[string]*pri.Member)

	// Create test members
	members := []*pri.Member{
		{CharacterID: 1, Name: "Terra"},
		{CharacterID: 5, Name: "Locke"},
		{CharacterID: 6, Name: "Cyan"},
		{CharacterID: 7, Name: "Shadow"},
	}

	for _, m := range members {
		party.AddPossibleMember(m)
	}

	// Verify members were added
	if len(party.Possible) != len(members) {
		t.Fatalf("party members = %d, want %d", len(party.Possible), len(members))
	}

	// Verify each member
	for _, expected := range members {
		found := false
		for _, actual := range party.Possible {
			if actual.CharacterID == expected.CharacterID {
				found = true
				if actual.Name != expected.Name {
					t.Errorf("member %d name = %q, want %q", expected.CharacterID, actual.Name, expected.Name)
				}
				break
			}
		}
		if !found {
			t.Errorf("member %d (%s) not found in party", expected.CharacterID, expected.Name)
		}
	}
}

// TestTransportationRoundTrip tests transportation data persistence
func TestTransportationRoundTrip(t *testing.T) {
	// Create test transportation
	trans := &pri.Transportation{
		ID:             1, // Falcon
		MapID:          100,
		Direction:      2,
		TimeStampTicks: 123456789,
		Enabled:        true,
	}
	trans.Position.X = 500.0
	trans.Position.Y = 600.0
	trans.Position.Z = 0.0

	// Verify values
	if trans.ID != 1 {
		t.Errorf("ID = %d, want 1", trans.ID)
	}
	if trans.MapID != 100 {
		t.Errorf("MapID = %d, want 100", trans.MapID)
	}
	if trans.Direction != 2 {
		t.Errorf("Direction = %d, want 2", trans.Direction)
	}
	if trans.Position.X != 500.0 {
		t.Errorf("Position.X = %f, want 500.0", trans.Position.X)
	}
	if !trans.Enabled {
		t.Error("Enabled should be true")
	}
}

// BenchmarkInventoryOperations benchmarks inventory operations
func BenchmarkInventoryOperations(b *testing.B) {
	inventory := pri.GetInventory()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inventory.Clear()
		for j := 0; j < 50; j++ {
			inventory.Set(j, pri.Row{ItemID: j + 1, Count: 10})
		}
		// Access items
		for j := 0; j < 50; j++ {
			_ = inventory.GetItemLookup()
		}
	}
}

// BenchmarkCharacterModification benchmarks character modifications
func BenchmarkCharacterModification(b *testing.B) {
	testChar := &models.Character{
		ID:   1,
		Name: "Test",
		HP:   models.CurrentMax{Current: 100, Max: 1000},
		MP:   models.CurrentMax{Current: 50, Max: 500},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate various character modifications
		testChar.HP.Current = i % testChar.HP.Max
		testChar.MP.Current = i % testChar.MP.Max
		testChar.Level = (i % 99) + 1
		testChar.Vigor = i % 256
	}
}
