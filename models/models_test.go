package models

import (
	"testing"
	"time"
)

// TestCharacterCreation tests basic character creation
func TestCharacterCreation(t *testing.T) {
	c := &Character{
		ID:        1,
		RootName:  "Terra",
		Name:      "Terra",
		Level:     10,
		Exp:       1000,
		HP:        CurrentMax{Current: 100, Max: 200},
		MP:        CurrentMax{Current: 50, Max: 100},
		Vigor:     30,
		Stamina:   25,
		Speed:     20,
		Magic:     35,
		IsEnabled: true,
		IsNPC:     false,
	}

	if c.ID != 1 {
		t.Errorf("Character ID = %d, want 1", c.ID)
	}
	if c.Name != "Terra" {
		t.Errorf("Character Name = %s, want Terra", c.Name)
	}
	if c.Level != 10 {
		t.Errorf("Character Level = %d, want 10", c.Level)
	}
	if c.HP.Current != 100 {
		t.Errorf("Character HP.Current = %d, want 100", c.HP.Current)
	}
	if c.HP.Max != 200 {
		t.Errorf("Character HP.Max = %d, want 200", c.HP.Max)
	}
}

// TestSpellCreation tests spell creation
func TestSpellCreation(t *testing.T) {
	s := &Spell{
		Name:  "Fire",
		Index: 1,
		Value: 10,
	}

	if s.Name != "Fire" {
		t.Errorf("Spell Name = %s, want Fire", s.Name)
	}
	if s.Index != 1 {
		t.Errorf("Spell Index = %d, want 1", s.Index)
	}
	if s.Value != 10 {
		t.Errorf("Spell Value = %d, want 10", s.Value)
	}
}

// TestEquipmentCreation tests equipment creation
func TestEquipmentCreation(t *testing.T) {
	e := Equipment{
		WeaponID: 1,
		ShieldID: 2,
		HelmetID: 3,
		ArmorID:  4,
		Relic1ID: 5,
		Relic2ID: 6,
	}

	if e.WeaponID != 1 {
		t.Errorf("Equipment WeaponID = %d, want 1", e.WeaponID)
	}
	if e.ShieldID != 2 {
		t.Errorf("Equipment ShieldID = %d, want 2", e.ShieldID)
	}
	if e.HelmetID != 3 {
		t.Errorf("Equipment HelmetID = %d, want 3", e.HelmetID)
	}
	if e.ArmorID != 4 {
		t.Errorf("Equipment ArmorID = %d, want 4", e.ArmorID)
	}
	if e.Relic1ID != 5 {
		t.Errorf("Equipment Relic1ID = %d, want 5", e.Relic1ID)
	}
	if e.Relic2ID != 6 {
		t.Errorf("Equipment Relic2ID = %d, want 6", e.Relic2ID)
	}
}

// TestEquipmentDefaults tests equipment default values
func TestEquipmentDefaults(t *testing.T) {
	e := Equipment{}

	if e.WeaponID != 0 {
		t.Errorf("Equipment WeaponID default = %d, want 0", e.WeaponID)
	}
	if e.ShieldID != 0 {
		t.Errorf("Equipment ShieldID default = %d, want 0", e.ShieldID)
	}
	if e.HelmetID != 0 {
		t.Errorf("Equipment HelmetID default = %d, want 0", e.HelmetID)
	}
	if e.ArmorID != 0 {
		t.Errorf("Equipment ArmorID default = %d, want 0", e.ArmorID)
	}
	if e.Relic1ID != 0 {
		t.Errorf("Equipment Relic1ID default = %d, want 0", e.Relic1ID)
	}
	if e.Relic2ID != 0 {
		t.Errorf("Equipment Relic2ID default = %d, want 0", e.Relic2ID)
	}
}

// TestMiscCreation tests misc data creation
func TestMiscCreation(t *testing.T) {
	m := &Misc{
		GP:                     1000,
		Steps:                  500,
		NumberOfSaves:          10,
		SaveCountRollOver:      0,
		MapXAxis:               100,
		MapYAxis:               200,
		AirshipXAxis:           300,
		AirshipYAxis:           400,
		IsAirshipVisible:       true,
		CursedShieldFightCount: 50,
		EscapeCount:            20,
		BattleCount:            100,
		MonstersKilledCount:    75,
	}

	if m.GP != 1000 {
		t.Errorf("Misc GP = %d, want 1000", m.GP)
	}
	if m.Steps != 500 {
		t.Errorf("Misc Steps = %d, want 500", m.Steps)
	}
	if m.NumberOfSaves != 10 {
		t.Errorf("Misc NumberOfSaves = %d, want 10", m.NumberOfSaves)
	}
	if !m.IsAirshipVisible {
		t.Error("Misc IsAirshipVisible = false, want true")
	}
}

// TestGetMisc tests the singleton GetMisc function
func TestGetMisc(t *testing.T) {
	// Reset the singleton for testing
	misc = nil

	m1 := GetMisc()
	if m1 == nil {
		t.Fatal("GetMisc() returned nil")
	}

	m2 := GetMisc()
	if m1 != m2 {
		t.Error("GetMisc() returned different instances, expected singleton")
	}
}

// TestMiscDefaults tests misc default values
func TestMiscDefaults(t *testing.T) {
	m := &Misc{}

	if m.GP != 0 {
		t.Errorf("Misc GP default = %d, want 0", m.GP)
	}
	if m.Steps != 0 {
		t.Errorf("Misc Steps default = %d, want 0", m.Steps)
	}
	if m.IsAirshipVisible != false {
		t.Error("Misc IsAirshipVisible default = true, want false")
	}
}

// TestCurrentMax tests the CurrentMax struct
func TestCurrentMax(t *testing.T) {
	cm := CurrentMax{
		Current: 100,
		Max:     200,
	}

	if cm.Current != 100 {
		t.Errorf("CurrentMax Current = %d, want 100", cm.Current)
	}
	if cm.Max != 200 {
		t.Errorf("CurrentMax Max = %d, want 200", cm.Max)
	}
}

// TestCurrentMaxFix tests the Fix method
func TestCurrentMaxFix(t *testing.T) {
	cm := CurrentMax{
		Current: 250,
		Max:     200,
	}

	cm.Fix()

	if cm.Current != 200 {
		t.Errorf("CurrentMax.Current after Fix = %d, want 200", cm.Current)
	}
}

// TestChangeCreation tests change creation
func TestChangeCreation(t *testing.T) {
	c := NewChange("Character", "Level", 10, 20)

	if c.Target != "Character" {
		t.Errorf("Change Target = %s, want Character", c.Target)
	}
	if c.FieldName != "Level" {
		t.Errorf("Change FieldName = %s, want Level", c.FieldName)
	}
	if c.OldValue != 10 {
		t.Errorf("Change OldValue = %v, want 10", c.OldValue)
	}
	if c.NewValue != 20 {
		t.Errorf("Change NewValue = %v, want 20", c.NewValue)
	}
	if c.ID == "" {
		t.Error("Change ID should not be empty")
	}
	if c.Timestamp.IsZero() {
		t.Error("Change Timestamp should not be zero")
	}
}

// TestNewBatchChange tests batch change creation
func TestNewBatchChange(t *testing.T) {
	c := NewBatchChange("batch1", "Level Up", "Character", "Level", 10, 20)

	if !c.Batch {
		t.Error("Batch change should have Batch = true")
	}
	if c.BatchID != "batch1" {
		t.Errorf("Change BatchID = %s, want batch1", c.BatchID)
	}
	if c.BatchName != "Level Up" {
		t.Errorf("Change BatchName = %s, want Level Up", c.BatchName)
	}
}

// TestValidateChange tests change validation
func TestValidateChange(t *testing.T) {
	// Valid change
	c1 := NewChange("Character", "Level", 10, 20)
	if !ValidateChange(c1) {
		t.Error("ValidateChange should return true for valid change")
	}

	// Invalid - empty target
	c2 := Change{
		Target:    "",
		FieldName: "Level",
		OldValue:  10,
		NewValue:  20,
	}
	if ValidateChange(c2) {
		t.Error("ValidateChange should return false for empty target")
	}

	// Invalid - empty field name
	c3 := Change{
		Target:    "Character",
		FieldName: "",
		OldValue:  10,
		NewValue:  20,
	}
	if ValidateChange(c3) {
		t.Error("ValidateChange should return false for empty field name")
	}

	// Invalid - nil old value
	c4 := Change{
		Target:    "Character",
		FieldName: "Level",
		OldValue:  nil,
		NewValue:  20,
	}
	if ValidateChange(c4) {
		t.Error("ValidateChange should return false for nil old value")
	}

	// Invalid - nil new value
	c5 := Change{
		Target:    "Character",
		FieldName: "Level",
		OldValue:  10,
		NewValue:  nil,
	}
	if ValidateChange(c5) {
		t.Error("ValidateChange should return false for nil new value")
	}
}

// TestChangeGroupCreation tests change group creation
func TestChangeGroupCreation(t *testing.T) {
	cg := &ChangeGroup{
		ID:      "group1",
		Name:    "Level Up Group",
		Changes: make([]Change, 0),
		Time:    time.Now(),
	}

	if cg.ID != "group1" {
		t.Errorf("ChangeGroup ID = %s, want group1", cg.ID)
	}
	if cg.Name != "Level Up Group" {
		t.Errorf("ChangeGroup Name = %s, want Level Up Group", cg.Name)
	}
	if len(cg.Changes) != 0 {
		t.Errorf("ChangeGroup Changes length = %d, want 0", len(cg.Changes))
	}
}

// TestCommandCreation tests command creation
func TestCommandCreation(t *testing.T) {
	c := &Command{
		NameValue: struct {
			Name  string
			Value int
		}{Name: "Fight", Value: 4},
		SortedIndex: 1,
	}

	if c.Name != "Fight" {
		t.Errorf("Command Name = %s, want Fight", c.Name)
	}
	if c.Value != 4 {
		t.Errorf("Command Value = %d, want 4", c.Value)
	}
	if c.SortedIndex != 1 {
		t.Errorf("Command SortedIndex = %d, want 1", c.SortedIndex)
	}
}

// TestPartyPresetCreation tests party preset creation
func TestPartyPresetCreation(t *testing.T) {
	p := &PartyPreset{
		ID:          "preset1",
		Name:        "My Preset",
		Description: "A test preset",
		Members:     [4]uint8{1, 2, 3, 4},
		Favorite:    true,
		Version:     1,
	}

	if p.ID != "preset1" {
		t.Errorf("PartyPreset ID = %s, want preset1", p.ID)
	}
	if p.Name != "My Preset" {
		t.Errorf("PartyPreset Name = %s, want My Preset", p.Name)
	}
	if p.Description != "A test preset" {
		t.Errorf("PartyPreset Description = %s, want A test preset", p.Description)
	}
	if p.Members[0] != 1 {
		t.Errorf("PartyPreset Members[0] = %d, want 1", p.Members[0])
	}
	if !p.Favorite {
		t.Error("PartyPreset Favorite = false, want true")
	}
	if p.Version != 1 {
		t.Errorf("PartyPreset Version = %d, want 1", p.Version)
	}
}

// TestPartyPresetValidate tests preset validation
func TestPartyPresetValidate(t *testing.T) {
	// Valid preset
	p1 := &PartyPreset{
		Name:    "Valid",
		Members: [4]uint8{1, 2, 3, 4},
	}
	if err := p1.Validate(); err != nil {
		t.Errorf("Validate() returned error for valid preset: %v", err)
	}

	// Invalid - empty name
	p2 := &PartyPreset{
		Name:    "",
		Members: [4]uint8{1, 2, 3, 4},
	}
	if err := p2.Validate(); err == nil {
		t.Error("Validate() should return error for empty name")
	}

	// Invalid - member too high
	p3 := &PartyPreset{
		Name:    "Invalid",
		Members: [4]uint8{1, 2, 3, 20},
	}
	if err := p3.Validate(); err == nil {
		t.Error("Validate() should return error for invalid member index")
	}
}

// TestPartyPresetIsComplete tests IsComplete method
func TestPartyPresetIsComplete(t *testing.T) {
	// Complete party (all slots filled with non-zero)
	p1 := &PartyPreset{
		Members: [4]uint8{1, 2, 3, 4},
	}
	if !p1.IsComplete() {
		t.Error("IsComplete() should return true for full party")
	}

	// Incomplete party (some zeros)
	p2 := &PartyPreset{
		Members: [4]uint8{1, 0, 0, 0},
	}
	if p2.IsComplete() {
		t.Error("IsComplete() should return false for partial party")
	}
}

// TestFF6SpriteCreation tests FF6Sprite creation
func TestFF6SpriteCreation(t *testing.T) {
	s := NewSprite("sprite1", "Test Sprite", SpriteTypeCharacter)

	if s.ID != "sprite1" {
		t.Errorf("FF6Sprite ID = %s, want sprite1", s.ID)
	}
	if s.Name != "Test Sprite" {
		t.Errorf("FF6Sprite Name = %s, want Test Sprite", s.Name)
	}
	if s.Type != SpriteTypeCharacter {
		t.Errorf("FF6Sprite Type = %v, want SpriteTypeCharacter", s.Type)
	}

	// Check dimensions were set
	w, h := SpriteTypeCharacter.GetDimensions()
	if s.Width != w {
		t.Errorf("FF6Sprite Width = %d, want %d", s.Width, w)
	}
	if s.Height != h {
		t.Errorf("FF6Sprite Height = %d, want %d", s.Height, h)
	}

	// Check defaults
	if s.Frames != 1 {
		t.Errorf("FF6Sprite Frames = %d, want 1", s.Frames)
	}
	if s.FrameRate != 12 {
		t.Errorf("FF6Sprite FrameRate = %d, want 12", s.FrameRate)
	}
	if s.Palette == nil {
		t.Error("FF6Sprite Palette should not be nil")
	}
}

// TestSpriteTypeGetDimensions tests GetDimensions for different sprite types
func TestSpriteTypeGetDimensions(t *testing.T) {
	tests := []struct {
		st             SpriteType
		expectedWidth  int
		expectedHeight int
	}{
		{SpriteTypeCharacter, 16, 24},
		{SpriteTypeNPC, 16, 24},
		{SpriteTypeBattle, 32, 32},
		{SpriteTypeEnemy, 32, 32},
		{SpriteTypePortrait, 48, 48},
		{SpriteTypeOverworld, 16, 16},
	}

	for _, tt := range tests {
		w, h := tt.st.GetDimensions()
		if w != tt.expectedWidth || h != tt.expectedHeight {
			t.Errorf("%s.GetDimensions() = (%d, %d), want (%d, %d)",
				tt.st.String(), w, h, tt.expectedWidth, tt.expectedHeight)
		}
	}
}

// TestSpriteTypeString tests String method for sprite types
func TestSpriteTypeString(t *testing.T) {
	tests := []struct {
		st       SpriteType
		expected string
	}{
		{SpriteTypeCharacter, "Character"},
		{SpriteTypeBattle, "Battle"},
		{SpriteTypePortrait, "Portrait"},
		{SpriteTypeNPC, "NPC"},
		{SpriteTypeEnemy, "Enemy"},
		{SpriteTypeOverworld, "Overworld"},
		{SpriteType(999), "Unknown"},
	}

	for _, tt := range tests {
		result := tt.st.String()
		if result != tt.expected {
			t.Errorf("SpriteType(%d).String() = %s, want %s", tt.st, result, tt.expected)
		}
	}
}

// TestRGB555 tests RGB555 color structure
func TestRGB555(t *testing.T) {
	c := RGB555{R: 31, G: 31, B: 31}

	if c.R != 31 {
		t.Errorf("RGB555.R = %d, want 31", c.R)
	}
	if c.G != 31 {
		t.Errorf("RGB555.G = %d, want 31", c.G)
	}
	if c.B != 31 {
		t.Errorf("RGB555.B = %d, want 31", c.B)
	}
}

// TestRGB555ToRGB888 tests conversion to RGB888
func TestRGB555ToRGB888(t *testing.T) {
	c := RGB555{R: 31, G: 0, B: 0}
	r, g, b := c.ToRGB888()

	if r != 255 {
		t.Errorf("Red channel = %d, want 255", r)
	}
	if g != 0 {
		t.Errorf("Green channel = %d, want 0", g)
	}
	if b != 0 {
		t.Errorf("Blue channel = %d, want 0", b)
	}
}

// TestFromRGB888 tests conversion from RGB888
func TestFromRGB888(t *testing.T) {
	c := FromRGB888(255, 0, 0)

	if c.R != 31 {
		t.Errorf("RGB555.R = %d, want 31", c.R)
	}
	if c.G != 0 {
		t.Errorf("RGB555.G = %d, want 0", c.G)
	}
	if c.B != 0 {
		t.Errorf("RGB555.B = %d, want 0", c.B)
	}
}

// TestPaletteCreation tests palette creation
func TestPaletteCreation(t *testing.T) {
	p := NewPalette("Test Palette")

	if p.Name != "Test Palette" {
		t.Errorf("Palette.Name = %s, want Test Palette", p.Name)
	}
	if p.Created.IsZero() {
		t.Error("Palette.Created should not be zero")
	}
}

// TestSpriteHistory tests sprite history
func TestSpriteHistory(t *testing.T) {
	h := NewSpriteHistory(10)

	if h.MaxSnapshots != 10 {
		t.Errorf("SpriteHistory.MaxSnapshots = %d, want 10", h.MaxSnapshots)
	}
	if h.CurrentIdx != -1 {
		t.Errorf("SpriteHistory.CurrentIdx = %d, want -1", h.CurrentIdx)
	}
}

// TestSpriteHistoryPush tests pushing to history
func TestSpriteHistoryPush(t *testing.T) {
	h := NewSpriteHistory(10)
	s := NewSprite("s1", "Sprite 1", SpriteTypeCharacter)

	h.Push(s)

	if h.CurrentIdx != 0 {
		t.Errorf("SpriteHistory.CurrentIdx = %d, want 0", h.CurrentIdx)
	}
	if len(h.Snapshots) != 1 {
		t.Errorf("SpriteHistory.Snapshots length = %d, want 1", len(h.Snapshots))
	}
}

// TestSpriteHistoryUndoRedo tests undo/redo functionality
func TestSpriteHistoryUndoRedo(t *testing.T) {
	h := NewSpriteHistory(10)

	s1 := NewSprite("s1", "Sprite 1", SpriteTypeCharacter)
	s2 := NewSprite("s2", "Sprite 2", SpriteTypeCharacter)

	h.Push(s1)
	h.Push(s2)

	// Can undo
	if !h.CanUndo() {
		t.Error("CanUndo() should return true")
	}

	// Undo
	undone := h.Undo()
	if undone == nil {
		t.Error("Undo() returned nil")
	}

	// Can redo
	if !h.CanRedo() {
		t.Error("CanRedo() should return true")
	}

	// Redo
	redone := h.Redo()
	if redone == nil {
		t.Error("Redo() returned nil")
	}
}

// TestSpriteImportOptionsCreation tests SpriteImportOptions creation
func TestSpriteImportOptionsCreation(t *testing.T) {
	opts := NewSpriteImportOptions()

	if opts.MaxColors != 16 {
		t.Errorf("SpriteImportOptions.MaxColors = %d, want 16", opts.MaxColors)
	}
	if opts.DitherMethod != "floyd-steinberg" {
		t.Errorf("SpriteImportOptions.DitherMethod = %s, want floyd-steinberg", opts.DitherMethod)
	}
	if opts.PreservePalette {
		t.Error("SpriteImportOptions.PreservePalette = true, want false")
	}
	if !opts.AutoDetectType {
		t.Error("SpriteImportOptions.AutoDetectType = false, want true")
	}
	if !opts.AutoPadding {
		t.Error("SpriteImportOptions.AutoPadding = false, want true")
	}
}

// TestValidationSeverity tests validation severity constants
func TestValidationSeverity(t *testing.T) {
	// Test that severity values are strings
	if SeverityError != "error" {
		t.Errorf("SeverityError = %s, want error", SeverityError)
	}
	if SeverityWarning != "warning" {
		t.Errorf("SeverityWarning = %s, want warning", SeverityWarning)
	}
	if SeverityInfo != "info" {
		t.Errorf("SeverityInfo = %s, want info", SeverityInfo)
	}
}

// TestValidationMode tests validation mode constants
func TestValidationMode(t *testing.T) {
	if StrictMode != "strict" {
		t.Errorf("StrictMode = %s, want strict", StrictMode)
	}
	if NormalMode != "normal" {
		t.Errorf("NormalMode = %s, want normal", NormalMode)
	}
	if LenientMode != "lenient" {
		t.Errorf("LenientMode = %s, want lenient", LenientMode)
	}
}

// TestDefaultValidationConfig tests default validation config
func TestDefaultValidationConfig(t *testing.T) {
	cfg := DefaultValidationConfig()

	if cfg.Mode != NormalMode {
		t.Errorf("DefaultValidationConfig Mode = %s, want %s", cfg.Mode, NormalMode)
	}
	if !cfg.RealTimeValidation {
		t.Error("DefaultValidationConfig RealTimeValidation = false, want true")
	}
	if !cfg.PreSaveValidation {
		t.Error("DefaultValidationConfig PreSaveValidation = false, want true")
	}
	if cfg.MaxCharacterLevel != 99 {
		t.Errorf("DefaultValidationConfig MaxCharacterLevel = %d, want 99", cfg.MaxCharacterLevel)
	}
	if cfg.MaxCharacterHP != 9999 {
		t.Errorf("DefaultValidationConfig MaxCharacterHP = %d, want 9999", cfg.MaxCharacterHP)
	}
	if cfg.MaxCharacterMP != 9999 {
		t.Errorf("DefaultValidationConfig MaxCharacterMP = %d, want 9999", cfg.MaxCharacterMP)
	}
	if cfg.MaxStatValue != 255 {
		t.Errorf("DefaultValidationConfig MaxStatValue = %d, want 255", cfg.MaxStatValue)
	}
	if cfg.AutoFixSimpleIssues {
		t.Error("DefaultValidationConfig AutoFixSimpleIssues = true, want false")
	}
}

// TestValidationResultCreation tests validation result creation
func TestValidationResultCreation(t *testing.T) {
	vr := ValidationResult{
		Valid:    true,
		Errors:   make([]ValidationIssue, 0),
		Warnings: make([]ValidationIssue, 0),
		Infomsgs: make([]ValidationIssue, 0),
	}

	if !vr.Valid {
		t.Error("ValidationResult Valid = false, want true")
	}
	if len(vr.Errors) != 0 {
		t.Errorf("ValidationResult Errors length = %d, want 0", len(vr.Errors))
	}
}

// TestValidationResultHasErrors tests HasErrors method
func TestValidationResultHasErrors(t *testing.T) {
	vr := ValidationResult{
		Errors: []ValidationIssue{},
	}
	if vr.HasErrors() {
		t.Error("HasErrors() should return false for empty errors")
	}

	vr.Errors = append(vr.Errors, ValidationIssue{})
	if !vr.HasErrors() {
		t.Error("HasErrors() should return true when errors exist")
	}
}

// TestValidationResultHasWarnings tests HasWarnings method
func TestValidationResultHasWarnings(t *testing.T) {
	vr := ValidationResult{
		Warnings: []ValidationIssue{},
	}
	if vr.HasWarnings() {
		t.Error("HasWarnings() should return false for empty warnings")
	}

	vr.Warnings = append(vr.Warnings, ValidationIssue{})
	if !vr.HasWarnings() {
		t.Error("HasWarnings() should return true when warnings exist")
	}
}

// TestValidationResultAllIssues tests AllIssues method
func TestValidationResultAllIssues(t *testing.T) {
	vr := ValidationResult{
		Errors:   []ValidationIssue{{Rule: "error1"}},
		Warnings: []ValidationIssue{{Rule: "warning1"}},
		Infomsgs: []ValidationIssue{{Rule: "info1"}},
	}

	all := vr.AllIssues()
	if len(all) != 3 {
		t.Errorf("AllIssues() returned %d issues, want 3", len(all))
	}
}

// TestValidationResultFixableIssues tests FixableIssues method
func TestValidationResultFixableIssues(t *testing.T) {
	vr := ValidationResult{
		Errors: []ValidationIssue{
			{Rule: "fixable", Fixable: true},
			{Rule: "not_fixable", Fixable: false},
		},
	}

	fixable := vr.FixableIssues()
	if len(fixable) != 1 {
		t.Errorf("FixableIssues() returned %d issues, want 1", len(fixable))
	}
	if fixable[0].Rule != "fixable" {
		t.Errorf("Fixable issue Rule = %s, want fixable", fixable[0].Rule)
	}
}

// BenchmarkCharacterCreation benchmarks character creation
func BenchmarkCharacterCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = &Character{
			ID:       i,
			RootName: "Test",
			Name:     "Test",
			Level:    50,
			HP:       CurrentMax{Current: 1000, Max: 2000},
			MP:       CurrentMax{Current: 500, Max: 1000},
		}
	}
}

// BenchmarkEquipmentCreation benchmarks equipment creation
func BenchmarkEquipmentCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Equipment{
			WeaponID: i,
			ShieldID: i + 1,
			HelmetID: i + 2,
			ArmorID:  i + 3,
		}
	}
}

// BenchmarkNewChange benchmarks change creation
func BenchmarkNewChange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewChange("Character", "Level", 10, 20)
	}
}

// BenchmarkNewSprite benchmarks sprite creation
func BenchmarkNewSprite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewSprite("sprite", "Test", SpriteTypeCharacter)
	}
}
