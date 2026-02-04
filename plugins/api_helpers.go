package plugins

import (
	"encoding/json"
	"ffvi_editor/models"

	jo "gitlab.com/c0b/go-ordered-json"
)

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
						// Look up command name by ID
						cmdName := getCommandNameByID(int(cmdID))
						cmd := &models.Command{}
						cmd.Value = int(cmdID)
						cmd.Name = cmdName
						char.Commands = append(char.Commands, cmd)
					}
				}
			}
		}
	}
}

// getCommandNameByID returns command name for a command ID
func getCommandNameByID(id int) string {
	commands := map[int]string{
		1:  "Fight",
		2:  "Item",
		3:  "Magic",
		4:  "Morph",
		5:  "Revert",
		6:  "Steal",
		7:  "Mug",
		8:  "Bushido",
		9:  "Throw",
		10: "Tools",
		11: "Blitz",
		12: "Runic",
		13: "Lore",
		14: "Sketch",
		15: "Control",
		16: "Slot",
		17: "Gil Toss",
		18: "Dance",
		19: "Rage",
		20: "Leap",
		21: "Mimic",
		22: "Row",
		23: "Defend",
		24: "Jump",
		25: "Dualcast",
		26: "Pray",
		27: "Shock",
		28: "Possess",
		29: "Magitek",
	}
	if name, ok := commands[id]; ok {
		return name
	}
	return "Unknown"
}
