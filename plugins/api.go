package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"ffvi_editor/models"
	constsPR "ffvi_editor/models/consts/pr"
	modelsPR "ffvi_editor/models/pr"
	ioPR "ffvi_editor/io/pr"
	jo "gitlab.com/c0b/go-ordered-json"
)

// PluginAPI provides safe access to save editor functionality for plugins
type PluginAPI interface {
	// Save Data Access
	GetCharacter(ctx context.Context, name string) (*models.Character, error)
	SetCharacter(ctx context.Context, name string, ch *models.Character) error
	GetInventory(ctx context.Context) (*modelsPR.Inventory, error)
	SetInventory(ctx context.Context, inv *modelsPR.Inventory) error
	GetParty(ctx context.Context) (*modelsPR.Party, error)
	SetParty(ctx context.Context, party *modelsPR.Party) error
	GetEquipment(ctx context.Context) (*models.Equipment, error)
	SetEquipment(ctx context.Context, eq *models.Equipment) error

	// Batch Operations
	ApplyBatchOperation(ctx context.Context, op string, params map[string]interface{}) (int, error)

	// Query Operations
	FindCharacter(ctx context.Context, predicate func(*models.Character) bool) *models.Character
	FindItems(ctx context.Context, predicate func(*modelsPR.Row) bool) []*modelsPR.Row

	// Events
	RegisterHook(event string, callback func(interface{}) error) error
	FireEvent(ctx context.Context, event string, data interface{}) error

	// UI
	ShowDialog(ctx context.Context, title, message string) error
	ShowConfirm(ctx context.Context, title, message string) (bool, error)
	ShowInput(ctx context.Context, prompt string) (string, error)

	// Logging
	Log(ctx context.Context, level string, message string) error

	// Settings
	GetSetting(key string) interface{}
	SetSetting(key string, value interface{}) error

	// Permissions
	HasPermission(permission string) bool
}

// APIImpl provides a default implementation of PluginAPI
type APIImpl struct {
	prData        *ioPR.PR
	hooks         map[string][]func(interface{}) error
	settings      map[string]interface{}
	permissions   map[string]bool
	logger        func(level, msg string)
	showDialogFn  func(title, message string) error
	showConfirmFn func(title, message string) bool
	showInputFn   func(prompt string) (string, error)
}

// NewAPIImpl creates a new API implementation
func NewAPIImpl(prData *ioPR.PR, permissions []string) *APIImpl {
	api := &APIImpl{
		prData:      prData,
		hooks:       make(map[string][]func(interface{}) error),
		settings:    make(map[string]interface{}),
		permissions: make(map[string]bool),
	}

	// Set permissions
	for _, perm := range permissions {
		api.permissions[perm] = true
	}

	// Set default logger
	api.logger = func(level, msg string) {
		// Default: silent
	}

	// Set default dialog functions (no-op)
	api.showDialogFn = func(title, message string) error {
		return nil
	}
	api.showConfirmFn = func(title, message string) bool {
		return false
	}
	api.showInputFn = func(prompt string) (string, error) {
		return "", nil
	}

	return api
}

// SetLogger sets the logging function
func (a *APIImpl) SetLogger(logger func(level, msg string)) {
	if logger != nil {
		a.logger = logger
	}
}

// SetDialogFunctions sets UI dialog functions
func (a *APIImpl) SetDialogFunctions(
	showDialog func(title, message string) error,
	showConfirm func(title, message string) bool,
	showInput func(prompt string) (string, error),
) {
	if showDialog != nil {
		a.showDialogFn = showDialog
	}
	if showConfirm != nil {
		a.showConfirmFn = showConfirm
	}
	if showInput != nil {
		a.showInputFn = showInput
	}
}

// GetCharacter retrieves a character by name
func (a *APIImpl) GetCharacter(ctx context.Context, name string) (*models.Character, error) {
	if !a.HasPermission(CommonPermissions.ReadSave) {
		return nil, ErrInsufficientPermissions
	}

	if a.prData == nil {
		return nil, ErrNilPRData
	}

	// Extract character from PR Characters OrderedMaps
	for _, charMap := range a.prData.Characters {
		if charMap == nil {
			continue
		}
		
		// Get character name from OrderedMap
		charNameVal := charMap.Get("name")
		if charNameVal == nil {
			continue
		}
		charName, ok := charNameVal.(string)
		if !ok || charName != name {
			continue
		}

		// Get base character model
		char := modelsPR.GetCharacter(name)
		if char == nil {
			continue
		}

		// Extract enabled status
		if enabledVal := charMap.Get("isEnableCorps"); enabledVal != nil {
			if enabled, ok := enabledVal.(bool); ok {
				char.IsEnabled = enabled
			}
		}

		// Extract experience
		if expVal := charMap.Get("currentExp"); expVal != nil {
			if expNum, ok := expVal.(json.Number); ok {
				if exp, err := expNum.Int64(); err == nil {
					char.Exp = int(exp)
				}
			}
		}

		// Extract character ID and job ID for base stats lookup
		var charID, jobID int
		if idVal := charMap.Get("id"); idVal != nil {
			if idNum, ok := idVal.(json.Number); ok {
				if id, err := idNum.Int64(); err == nil {
					charID = int(id)
					char.ID = charID
				}
			}
		}
		if jobIDVal := charMap.Get("jobId"); jobIDVal != nil {
			if jobIDNum, ok := jobIDVal.(json.Number); ok {
				if id, err := jobIDNum.Int64(); err == nil {
					jobID = int(id)
				}
			}
		}

		// Get base stats for HP/MP calculation
		baseOffset, _ := modelsPR.GetCharacterBaseOffset(charID, jobID)

		// Extract parameter nested object for stats
		if paramVal := charMap.Get("parameter"); paramVal != nil {
			if paramStr, ok := paramVal.(string); ok {
				// Parameter is stored as JSON string, need to unmarshal
				paramMap := jo.NewOrderedMap()
				if err := paramMap.UnmarshalJSON([]byte(paramStr)); err == nil {
					// Extract HP
					if hpVal := paramMap.Get("currentHP"); hpVal != nil {
						if hpNum, ok := hpVal.(json.Number); ok {
							if hp, err := hpNum.Int64(); err == nil {
								char.HP.Current = int(hp)
							}
						}
					}
					if maxHPVal := paramMap.Get("addtionalMaxHp"); maxHPVal != nil {
						if maxHPNum, ok := maxHPVal.(json.Number); ok {
							if maxHP, err := maxHPNum.Int64(); err == nil {
								char.HP.Max = int(maxHP)
								if baseOffset != nil {
									char.HP.Max += baseOffset.HPBase
								}
							}
						}
					}

					// Extract MP
					if mpVal := paramMap.Get("currentMP"); mpVal != nil {
						if mpNum, ok := mpVal.(json.Number); ok {
							if mp, err := mpNum.Int64(); err == nil {
								char.MP.Current = int(mp)
							}
						}
					}
					if maxMPVal := paramMap.Get("addtionalMaxMp"); maxMPVal != nil {
						if maxMPNum, ok := maxMPVal.(json.Number); ok {
							if maxMP, err := maxMPNum.Int64(); err == nil {
								char.MP.Max = int(maxMP)
								if baseOffset != nil {
									char.MP.Max += baseOffset.MPBase
								}
							}
						}
					}

					// Extract Level
					if levelVal := paramMap.Get("addtionalLevel"); levelVal != nil {
						if levelNum, ok := levelVal.(json.Number); ok {
							if level, err := levelNum.Int64(); err == nil {
								char.Level = int(level)
							}
						}
					}

					// Extract Vigor (Power)
					if vigorVal := paramMap.Get("addtionalPower"); vigorVal != nil {
						if vigorNum, ok := vigorVal.(json.Number); ok {
							if vigor, err := vigorNum.Int64(); err == nil {
								char.Vigor = int(vigor)
							}
						}
					}

					// Extract Stamina (Vitality)
					if staminaVal := paramMap.Get("addtionalVitality"); staminaVal != nil {
						if staminaNum, ok := staminaVal.(json.Number); ok {
							if stamina, err := staminaNum.Int64(); err == nil {
								char.Stamina = int(stamina)
							}
						}
					}

					// Extract Speed (Agility)
					if speedVal := paramMap.Get("addtionalAgility"); speedVal != nil {
						if speedNum, ok := speedVal.(json.Number); ok {
							if speed, err := speedNum.Int64(); err == nil {
								char.Speed = int(speed)
							}
						}
					}

					// Extract Magic
					if magicVal := paramMap.Get("addtionMagic"); magicVal != nil {
						if magicNum, ok := magicVal.(json.Number); ok {
							if magic, err := magicNum.Int64(); err == nil {
								char.Magic = int(magic)
							}
						}
					}
				}
			}
		}

		// Extract equipment
		a.extractEquipment(charMap, char)

		// Extract spells/magic
		a.extractSpells(charMap, char)

		// Extract commands
		a.extractCommands(charMap, char)

		return char, nil
	}

	return nil, ErrCharacterNotFound
}

// SetCharacter updates a character with full data modification
func (a *APIImpl) SetCharacter(ctx context.Context, name string, ch *models.Character) error {
	if !a.HasPermission(CommonPermissions.WriteSave) {
		return ErrInsufficientPermissions
	}

	if a.prData == nil {
		return ErrNilPRData
	}

	// Find and update character in PR Characters OrderedMaps
	for _, charMap := range a.prData.Characters {
		if charMap == nil {
			continue
		}

		// Check if this is the target character
		charNameVal := charMap.Get("name")
		if charNameVal == nil {
			continue
		}
		charName, ok := charNameVal.(string)
		if !ok || charName != name {
			continue
		}

		// Update basic character fields
		charMap.Set("name", ch.Name)
		charMap.Set("isEnableCorps", ch.IsEnabled)
		
		// Update experience
		charMap.Set("currentExp", json.Number(fmt.Sprintf("%d", ch.Exp)))

		// Get base stats for HP/MP calculation
		var charID, jobID int
		if idVal := charMap.Get("id"); idVal != nil {
			if idNum, ok := idVal.(json.Number); ok {
				if id, err := idNum.Int64(); err == nil {
					charID = int(id)
				}
			}
		}
		if jobIDVal := charMap.Get("jobId"); jobIDVal != nil {
			if jobIDNum, ok := jobIDVal.(json.Number); ok {
				if id, err := jobIDNum.Int64(); err == nil {
					jobID = int(id)
				}
			}
		}
		baseOffset, _ := modelsPR.GetCharacterBaseOffset(charID, jobID)

		// Update parameter nested object for stats
		if err := a.updateParameter(charMap, ch, baseOffset); err != nil {
			return fmt.Errorf("failed to update parameter: %w", err)
		}

		// Update equipment
		if err := a.updateEquipment(charMap, ch); err != nil {
			return fmt.Errorf("failed to update equipment: %w", err)
		}

		// Update spells
		if err := a.updateSpells(charMap, ch); err != nil {
			return fmt.Errorf("failed to update spells: %w", err)
		}

		// Update commands
		if err := a.updateCommands(charMap, ch); err != nil {
			return fmt.Errorf("failed to update commands: %w", err)
		}

		return nil
	}

	return ErrCharacterNotFound
}

// GetInventory retrieves the inventory
func (a *APIImpl) GetInventory(ctx context.Context) (*modelsPR.Inventory, error) {
	if !a.HasPermission(CommonPermissions.ReadSave) {
		return nil, ErrInsufficientPermissions
	}

	if a.prData == nil {
		return nil, ErrNilPRData
	}

	// Return the cached inventory instance
	// In production, this would extract from UserData OrderedMap
	inventory := modelsPR.GetInventory()
	return inventory, nil
}

// SetInventory updates the inventory
func (a *APIImpl) SetInventory(ctx context.Context, inv *modelsPR.Inventory) error {
	if !a.HasPermission(CommonPermissions.WriteSave) {
		return ErrInsufficientPermissions
	}

	if a.prData == nil {
		return ErrNilPRData
	}

	// Serialize inventory rows to JSON strings
	rows := inv.GetRows()
	itemStrings := make([]interface{}, 0, len(rows))
	
	for _, r := range rows {
		// Skip empty rows
		if r.ItemID == 0 || r.Count == 0 {
			continue
		}
		
		// Create item object
		itemJSON, err := json.Marshal(r)
		if err != nil {
			return fmt.Errorf("failed to marshal inventory item: %w", err)
		}
		itemStrings = append(itemStrings, string(itemJSON))
	}
	
	// Update normalOwnedItemList in UserData
	itemListObj := jo.NewOrderedMap()
	itemListObj.Set("target", itemStrings)
	
	itemListJSON, err := itemListObj.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal item list: %w", err)
	}
	
	a.prData.UserData.Set("normalOwnedItemList", string(itemListJSON))
	return nil
}

// GetParty retrieves the party
func (a *APIImpl) GetParty(ctx context.Context) (*modelsPR.Party, error) {
	if !a.HasPermission(CommonPermissions.ReadSave) {
		return nil, ErrInsufficientPermissions
	}

	if a.prData == nil {
		return nil, ErrNilPRData
	}

	// Return the cached party instance
	// In production, this would extract from UserData OrderedMap
	party := modelsPR.GetParty()
	return party, nil
}

// SetParty updates the party
func (a *APIImpl) SetParty(ctx context.Context, party *modelsPR.Party) error {
	if !a.HasPermission(CommonPermissions.WriteSave) {
		return ErrInsufficientPermissions
	}

	if a.prData == nil {
		return ErrNilPRData
	}

	// Get existing party ID (or default to 1)
	partyID := 1
	corpsListValue := a.prData.UserData.Get("corpsList")
	if corpsListValue != nil {
		corpsListStr, ok := corpsListValue.(string)
		if ok {
			corpsListObj := jo.NewOrderedMap()
			if err := corpsListObj.UnmarshalJSON([]byte(corpsListStr)); err == nil {
				targetValue := corpsListObj.Get("target")
				if targetArray, ok := targetValue.([]interface{}); ok && len(targetArray) > 0 {
					if firstMemberStr, ok := targetArray[0].(string); ok {
						var pm struct {
							ID          int `json:"id"`
							CharacterID int `json:"characterId"`
						}
						if err := json.Unmarshal([]byte(firstMemberStr), &pm); err == nil {
							partyID = pm.ID
						}
					}
				}
			}
		}
	}
	
	// Serialize party members to JSON strings
	memberStrings := make([]interface{}, 4)
	for i, member := range party.Members {
		pm := struct {
			ID          int `json:"id"`
			CharacterID int `json:"characterId"`
		}{
			ID:          partyID,
			CharacterID: member.CharacterID,
		}
		
		memberJSON, err := json.Marshal(pm)
		if err != nil {
			return fmt.Errorf("failed to marshal party member %d: %w", i, err)
		}
		memberStrings[i] = string(memberJSON)
	}
	
	// Update corpsList in UserData
	corpsListObj := jo.NewOrderedMap()
	corpsListObj.Set("target", memberStrings)
	
	corpsListJSON, err := corpsListObj.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal corps list: %w", err)
	}
	
	a.prData.UserData.Set("corpsList", string(corpsListJSON))
	return nil
}

// GetEquipment retrieves equipment
func (a *APIImpl) GetEquipment(ctx context.Context) (*models.Equipment, error) {
	if !a.HasPermission(CommonPermissions.ReadSave) {
		return nil, ErrInsufficientPermissions
	}

	if a.prData == nil {
		return nil, ErrNilPRData
	}

	// Extract equipment from UserData OrderedMap
	equipment := &models.Equipment{}
	// In production, this would extract equipment data from UserData
	// For now, return empty equipment structure
	return equipment, nil
}

// SetEquipment updates equipment for the first character (standalone method)
// Note: For character-specific equipment updates, use SetCharacter() instead
func (a *APIImpl) SetEquipment(ctx context.Context, eq *models.Equipment) error {
	if !a.HasPermission(CommonPermissions.WriteSave) {
		return ErrInsufficientPermissions
	}

	if a.prData == nil {
		return ErrNilPRData
	}

	// Find first character in the party to update equipment
	// This is a convenience method - use SetCharacter for specific character updates
	if len(a.prData.Characters) == 0 {
		return fmt.Errorf("no characters found")
	}
	
	firstChar := a.prData.Characters[0]
	if firstChar == nil {
		return fmt.Errorf("first character is nil")
	}
	
	// Create equipment OrderedMap
	equipmentMap := jo.NewOrderedMap()
	values := make([]interface{}, 6)
	
	// Build equipment array in order: Weapon, Shield, Armor, Helmet, Relic1, Relic2
	equipIDs := []int{
		eq.WeaponID,
		eq.ShieldID,
		eq.ArmorID,
		eq.HelmetID,
		eq.Relic1ID,
		eq.Relic2ID,
	}
	
	for i, equipID := range equipIDs {
		itemMap := make(map[string]interface{})
		itemMap["contentId"] = json.Number(fmt.Sprintf("%d", equipID))
		itemMap["count"] = json.Number("1")
		
		itemJSON, err := json.Marshal(itemMap)
		if err != nil {
			return fmt.Errorf("failed to marshal equipment slot %d: %w", i, err)
		}
		values[i] = string(itemJSON)
	}
	
	equipmentMap.Set("values", values)
	
	// Marshal to JSON string
	equipmentJSON, err := equipmentMap.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal equipment: %w", err)
	}
	
	// Update equipmentList field in character
	firstChar.Set("equipmentList", string(equipmentJSON))
	return nil
}

// ApplyBatchOperation applies a batch operation
func (a *APIImpl) ApplyBatchOperation(ctx context.Context, op string, params map[string]interface{}) (int, error) {
	if !a.HasPermission(CommonPermissions.WriteSave) {
		return 0, ErrInsufficientPermissions
	}

	// Placeholder for batch operation support
	// Real implementation would call batch operation handlers
	return 0, ErrBatchOpNotSupported
}

// FindCharacter finds a character using a predicate
func (a *APIImpl) FindCharacter(ctx context.Context, predicate func(*models.Character) bool) *models.Character {
	if !a.HasPermission(CommonPermissions.ReadSave) || a.prData == nil {
		return nil
	}

	// Extract and search characters from PR data structure
	for _, charMap := range a.prData.Characters {
		if charMap == nil {
			continue
		}

		// Get character name
		charNameVal := charMap.Get("name")
		if charNameVal == nil {
			continue
		}
		charName, ok := charNameVal.(string)
		if !ok {
			continue
		}

		// Get base character model
		char := modelsPR.GetCharacter(charName)
		if char == nil {
			continue
		}

		// Extract full character data (reuse GetCharacter logic)
		// This ensures consistent data extraction
		if fullChar, err := a.GetCharacter(ctx, charName); err == nil {
			char = fullChar
		}

		// Test predicate
		if predicate(char) {
			return char
		}
	}

	return nil
}

// FindItems finds items using a predicate
func (a *APIImpl) FindItems(ctx context.Context, predicate func(*modelsPR.Row) bool) []*modelsPR.Row {
	if !a.HasPermission(CommonPermissions.ReadSave) || a.prData == nil {
		return nil
	}

	// Extract inventory rows from PR data structure
	inventory := modelsPR.GetInventory()
	if inventory == nil {
		return nil
	}

	var items []*modelsPR.Row
	for _, row := range inventory.Rows {
		if row != nil && predicate(row) {
			items = append(items, row)
		}
	}
	return items
}

// RegisterHook registers a hook callback
func (a *APIImpl) RegisterHook(event string, callback func(interface{}) error) error {
	if callback == nil {
		return ErrNilCallback
	}

	a.hooks[event] = append(a.hooks[event], callback)
	return nil
}

// FireEvent fires an event to all registered callbacks
func (a *APIImpl) FireEvent(ctx context.Context, event string, data interface{}) error {
	callbacks, ok := a.hooks[event]
	if !ok {
		return nil
	}

	for _, callback := range callbacks {
		if err := callback(data); err != nil {
			a.logger("error", "hook callback error: "+err.Error())
		}
	}

	return nil
}

// ShowDialog shows a dialog
func (a *APIImpl) ShowDialog(ctx context.Context, title, message string) error {
	if !a.HasPermission(CommonPermissions.UIDisplay) {
		return ErrInsufficientPermissions
	}

	return a.showDialogFn(title, message)
}

// ShowConfirm shows a confirmation dialog
func (a *APIImpl) ShowConfirm(ctx context.Context, title, message string) (bool, error) {
	if !a.HasPermission(CommonPermissions.UIDisplay) {
		return false, ErrInsufficientPermissions
	}

	return a.showConfirmFn(title, message), nil
}

// ShowInput shows an input dialog
func (a *APIImpl) ShowInput(ctx context.Context, prompt string) (string, error) {
	if !a.HasPermission(CommonPermissions.UIDisplay) {
		return "", ErrInsufficientPermissions
	}

	return a.showInputFn(prompt)
}

// Log logs a message
func (a *APIImpl) Log(ctx context.Context, level string, message string) error {
	a.logger(level, message)
	return nil
}

// GetSetting gets a setting
func (a *APIImpl) GetSetting(key string) interface{} {
	return a.settings[key]
}

// SetSetting sets a setting
func (a *APIImpl) SetSetting(key string, value interface{}) error {
	a.settings[key] = value
	return nil
}

// HasPermission checks if a permission is granted
func (a *APIImpl) HasPermission(permission string) bool {
	return a.permissions[permission]
}

// SetCharacterStat sets a character stat
func (a *APIImpl) SetCharacterStat(charID int, stat string, value int) error {
	// TODO: Set stat in PR save
	return nil
}

// extractEquipment extracts equipment slots from character data
func (a *APIImpl) extractEquipment(charMap *jo.OrderedMap, char *models.Character) {
	// Equipment is stored as JSON string in equipmentList field
	if eqListVal := charMap.Get("equipmentList"); eqListVal != nil {
		if eqListStr, ok := eqListVal.(string); ok {
			eqMap := jo.NewOrderedMap()
			if err := eqMap.UnmarshalJSON([]byte(eqListStr)); err == nil {
				// Equipment is stored as array of {contentId, count} objects
				if valuesVal := eqMap.Get("values"); valuesVal != nil {
					if valuesArr, ok := valuesVal.([]interface{}); ok {
						// Slots: 0=Weapon, 1=Shield, 2=Armor, 3=Helmet, 4=Relic1, 5=Relic2
						for i, itemVal := range valuesArr {
							if itemMap, ok := itemVal.(map[string]interface{}); ok {
								if contentIDVal, ok := itemMap["contentId"]; ok {
									if contentIDNum, ok := contentIDVal.(json.Number); ok {
										if contentID, err := contentIDNum.Int64(); err == nil {
											switch i {
											case 0:
												char.Equipment.WeaponID = int(contentID)
											case 1:
												char.Equipment.ShieldID = int(contentID)
											case 2:
												char.Equipment.ArmorID = int(contentID)
											case 3:
												char.Equipment.HelmetID = int(contentID)
											case 4:
												char.Equipment.Relic1ID = int(contentID)
											case 5:
												char.Equipment.Relic2ID = int(contentID)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

// extractSpells extracts learned spells/magic from character data
func (a *APIImpl) extractSpells(charMap *jo.OrderedMap, char *models.Character) {
	// Spells are stored in abilityList as JSON string
	if abilityListVal := charMap.Get("abilityList"); abilityListVal != nil {
		if abilityListStr, ok := abilityListVal.(string); ok {
			abilityMap := jo.NewOrderedMap()
			if err := abilityMap.UnmarshalJSON([]byte(abilityListStr)); err == nil {
				if valuesVal := abilityMap.Get("values"); valuesVal != nil {
					if valuesArr, ok := valuesVal.([]interface{}); ok {
						for _, abilityVal := range valuesArr {
							if abilityStr, ok := abilityVal.(string); ok {
								abilityObj := jo.NewOrderedMap()
								if err := abilityObj.UnmarshalJSON([]byte(abilityStr)); err == nil {
									// Get ability ID
									if abilityIDVal := abilityObj.Get("abilityId"); abilityIDVal != nil {
										if abilityIDNum, ok := abilityIDVal.(json.Number); ok {
											if abilityID, err := abilityIDNum.Int64(); err == nil {
												// Check if it's a spell (spells are in specific ID range)
												if abilityID >= 1 && abilityID <= 54 {
													if spell, ok := char.SpellsByID[int(abilityID)]; ok {
														// Get skill level (spell proficiency)
														if skillLevelVal := abilityObj.Get("skillLevel"); skillLevelVal != nil {
															if skillLevelNum, ok := skillLevelVal.(json.Number); ok {
																if skillLevel, err := skillLevelNum.Int64(); err == nil {
																	spell.Value = int(skillLevel)
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

// extractCommands extracts command list from character data
func (a *APIImpl) extractCommands(charMap *jo.OrderedMap, char *models.Character) {
	// Commands are stored as array of command IDs
	if commandListVal := charMap.Get("commandList"); commandListVal != nil {
		if commandArr, ok := commandListVal.([]interface{}); ok {
			char.Commands = make([]*models.Command, 0, len(commandArr))
			for _, cmdVal := range commandArr {
				if cmdNum, ok := cmdVal.(json.Number); ok {
					if cmdID, err := cmdNum.Int64(); err == nil {
						// Look up command by ID from constants
						if cmd := constsPR.CommandLookupByValue[int(cmdID)]; cmd != nil {
							char.Commands = append(char.Commands, cmd)
						}
					}
				}
			}
		}
	}
}
