package pr

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"ffvi_editor/global"
	"ffvi_editor/io/config"
	"ffvi_editor/io/file"
	"ffvi_editor/models"
	"ffvi_editor/models/consts"
	"ffvi_editor/models/consts/pr"
	pri "ffvi_editor/models/pr"

	jo "gitlab.com/c0b/go-ordered-json"
)

func (p *PR) Load(fromFile string, saveType global.SaveFileType) (err error) {
	var (
		out   []byte
		s     string
		names []unicodeNameReplace
	)

	if out, p.fileTrimmed, err = file.LoadFile(fromFile, saveType); err != nil {
		return
	}
	s = string(out)

	if err = p.loadBase(s); err != nil {
		return
	}

	if err = p.unmarshalFrom(p.Base, UserData, p.UserData); err != nil {
		return
	}

	if err = p.unmarshalFrom(p.Base, MapData, p.MapData); err != nil {
		return
	}

	// Safely extract character list with type validation
	characterStrings, err := SafeGetFromTarget(p.UserData, OwnedCharacterList)
	if err != nil {
		return fmt.Errorf("failed to extract character list: %w", err)
	}

	for j, charJSON := range characterStrings {
		if j >= len(p.Characters) {
			break
		}
		if p.Characters[j] == nil {
			p.Characters[j] = jo.NewOrderedMap()
		}
		if err = p.Characters[j].UnmarshalJSON([]byte(charJSON)); err != nil {
			return fmt.Errorf("failed to unmarshal character %d: %w", j, err)
		}
	}

	pri.GetParty().Clear()

	if err = p.loadCharacters(); err != nil {
		return
	}
	if err = p.loadParty(); err != nil {
		return
	}
	if err = p.loadEspers(); err != nil {
		return
	}
	if err = p.loadMiscStats(); err != nil {
		return
	}
	if err = p.loadInventory(NormalOwnedItemList, pri.GetInventory()); err != nil {
		return
	}
	if err = p.loadInventory(importantOwnedItemList, pri.GetImportantInventory()); err != nil {
		return
	}
	if err = p.loadVeldt(); err != nil {
		return
	}
	if err = p.loadCheats(); err != nil {
		return
	}
	if err = p.loadMapData(); err != nil {
		return
	}
	if err = p.loadTransportation(); err != nil {
		return
	}

	if len(names) > 0 {
		p.names = names
	}
	return
}

func (p *PR) loadParty() (err error) {
	party := pri.GetParty()

	// Safely extract corps list with type validation
	corpsList, err := SafeGetFromTarget(p.UserData, CorpsList)
	if err != nil {
		return fmt.Errorf("failed to extract corps list: %w", err)
	}

	for slot, corpJSON := range corpsList {
		member := pri.Member{}
		if err = json.Unmarshal([]byte(corpJSON), &member); err != nil {
			return fmt.Errorf("failed to unmarshal corps at slot %d: %w", slot, err)
		}
		if err = party.SetMemberByID(slot, member.CharacterID); err != nil {
			return fmt.Errorf("failed to set party member at slot %d: %w", slot, err)
		}
	}
	return
}

func (p *PR) loadCharacters() (err error) {
	for _, d := range p.Characters {
		if d == nil {
			continue
		}

		var id int
		if id, err = p.getInt(d, ID); err != nil {
			return
		}

		var jobID int
		if jobID, err = p.getInt(d, JobID); err != nil {
			return
		}

		o, found := pri.GetCharacterBaseOffset(id, jobID)
		if !found {
			continue
		}

		c := pri.GetCharacter(o.Name)
		c.EnableCommandsSave = config.AutoEnableCmd()

		if c.Name, err = p.getString(d, Name); err != nil {
			return
		}

		// if pr.IsMainCharacter(c.Name) {
		pri.GetParty().AddPossibleMember(&pri.Member{
			CharacterID: id,
			Name:        c.Name,
		})
		// }

		if c.IsEnabled, err = p.getBool(d, IsEnableCorps); err != nil {
			return
		}

		params := jo.NewOrderedMap()
		if err = p.unmarshalFrom(d, Parameter, params); err != nil {
			return
		}

		if c.Level, err = p.getInt(params, AdditionalLevel); err != nil {
			return
		}

		if c.HP.Current, err = p.getInt(params, CurrentHP); err != nil {
			return
		}
		if c.HP.Max, err = p.getInt(params, AdditionalMaxHp); err != nil {
			return
		}
		c.HP.Max += o.HPBase

		if c.MP.Current, err = p.getInt(params, CurrentMP); err != nil {
			return
		}
		if c.MP.Max, err = p.getInt(params, AdditionalMaxMp); err != nil {
			return
		}
		c.MP.Max += o.MPBase

		if c.Exp, err = p.getInt(d, CurrentExp); err != nil {
			return
		}

		// TODO Status

		var values interface{}
		if values, err = p.getFromTarget(d, CommandList); err != nil {
			return
		}
		for i, v := range values.([]interface{}) {
			var j int64
			if j, err = v.(json.Number).Int64(); err != nil {
				return
			}
			if i >= len(c.Commands) {
				c.Commands = append(c.Commands, pr.CommandLookupByValue[int(j)])
			} else {
				c.Commands[i] = pr.CommandLookupByValue[int(j)]
			}
		}

		if c.Vigor, err = p.getInt(params, AdditionalPower); err != nil {
			return
		}

		if c.Stamina, err = p.getInt(params, AdditionalVitality); err != nil {
			return
		}

		if c.Speed, err = p.getInt(params, AdditionalAgility); err != nil {
			return
		}

		if c.Magic, err = p.getInt(params, AdditionMagic); err != nil {
			return
		}

		if err = p.loadEquipment(d, c); err != nil {
			return
		}

		if err = p.loadSpells(d, c); err != nil {
			return
		}

		// Cyan
		if jobID == 3 {
			p.uncheckAll(pr.Bushidos)
			if err = p.loadSkills(d, pr.BushidoFrom, pr.BushidoTo, pr.BushidoLookupByID); err != nil {
				return
			}
		}
		// Sabin
		if jobID == 6 {
			p.uncheckAll(pr.Blitzes)
			if err = p.loadSkills(d, pr.BlitzFrom, pr.BlitzTo, pr.BlitzLookupByID); err != nil {
				return
			}
		}
		// Mog
		if id == 16 {
			p.uncheckAll(pr.Dances)
			if err = p.loadSkills(d, pr.DanceFrom, pr.DanceTo, pr.DanceLookupByID); err != nil {
				return
			}
		}
		// Strago
		if jobID == 8 {
			p.uncheckAll(pr.Lores)
			if err = p.loadSkills(d, pr.LoreFrom, pr.LoreTo, pr.LoreLookupByID); err != nil {
				return
			}
		}
		// Gau
		if jobID == 12 {
			p.uncheckAll(pr.Rages)
			if err = p.loadSkills(d, pr.RageFrom, pr.RageTo, pr.RageLookupByID); err != nil {
				return
			}
		}
	}
	return
}

func (p *PR) loadEquipment(d *jo.OrderedMap, c *models.Character) (err error) {
	c.Equipment.WeaponID = 93
	c.Equipment.ShieldID = 93
	c.Equipment.ArmorID = 199
	c.Equipment.HelmetID = 198
	c.Equipment.Relic1ID = 200
	c.Equipment.Relic2ID = 200

	var eqIDCounts []idCount
	if eqIDCounts, err = p.unmarshalEquipment(d); err != nil {
		return
	}

	if len(eqIDCounts) > 0 {
		c.Equipment.WeaponID = eqIDCounts[0].ContentID
	}
	if len(eqIDCounts) > 1 {
		c.Equipment.ShieldID = eqIDCounts[1].ContentID
	}
	if len(eqIDCounts) > 2 {
		c.Equipment.ArmorID = eqIDCounts[2].ContentID
	}
	if len(eqIDCounts) > 3 {
		c.Equipment.HelmetID = eqIDCounts[3].ContentID
	}
	if len(eqIDCounts) > 4 {
		c.Equipment.Relic1ID = eqIDCounts[4].ContentID
	}
	if len(eqIDCounts) > 5 {
		c.Equipment.Relic2ID = eqIDCounts[5].ContentID
	}
	return
}

/*func (p *PR) getEquipmentKeys(m *jo.OrderedMap) (keys []int, err error) {
	i, ok := m.GetValue(EquipmentList)
	if !ok {
		return nil, fmt.Errorf("%s not found", EquipmentList)
	}

	eq := jo.NewOrderedMap()
	if err = eq.UnmarshalJSON([]byte(i.(string))); err != nil {
		return
	}

	if i, ok = eq.GetValue("keys"); ok && i != nil {
		keys = make([]int, len(i.([]interface{})))
		for j, v := range i.([]interface{}) {
			var k int64
			if k, err = v.(json.Number).Int64(); err != nil {
				return
			}
			keys[j] = int(k)
		}
	}
	return
}*/

func (p *PR) loadSpells(d *jo.OrderedMap, c *models.Character) (err error) {
	// Safely extract ability list with type validation
	abilityStrings, err := SafeGetFromTarget(d, AbilityList)
	if err != nil {
		return fmt.Errorf("failed to extract ability list: %w", err)
	}

	for j, abilityJSON := range abilityStrings {
		m := jo.NewOrderedMap()
		if err = m.UnmarshalJSON([]byte(abilityJSON)); err != nil {
			return fmt.Errorf("ability[%d] unmarshal failed: %w", j, err)
		}

		abilityIDValue, ok := m.GetValue("abilityId")
		if !ok {
			continue
		}

		abilityID, err := ExtractInt64(abilityIDValue)
		if err != nil {
			continue
		}

		// Check if this is a spell within the spell range
		if abilityID >= pr.SpellFrom && abilityID <= pr.SpellTo {
			spell, ok := c.SpellsByID[int(abilityID)]
			if !ok {
				continue
			}

			skillLevelValue := m.Get("skillLevel")
			if skillLevelValue != nil {
				if skillLevel, err := ExtractInt64(skillLevelValue); err == nil {
					spell.Value = int(skillLevel)
				}
			}
		}
	}
	return
}

func (p *PR) loadSkills(d *jo.OrderedMap, from int64, to int64, nvcLookup map[int]*consts.NameValueChecked) (err error) {
	// Safely extract ability list with type validation
	abilityStrings, err := SafeGetFromTarget(d, AbilityList)
	if err != nil {
		return fmt.Errorf("failed to extract ability list: %w", err)
	}

	for j, abilityJSON := range abilityStrings {
		m := jo.NewOrderedMap()
		if err = m.UnmarshalJSON([]byte(abilityJSON)); err != nil {
			return fmt.Errorf("ability[%d] unmarshal failed: %w", j, err)
		}

		abilityIDValue, ok := m.GetValue("abilityId")
		if !ok {
			continue
		}

		abilityID, err := ExtractInt64(abilityIDValue)
		if err != nil {
			continue
		}

		// Check if this ability is in the requested range
		if abilityID >= from && abilityID <= to {
			skillLevelValue := m.Get("skillLevel")
			if skillLevelValue != nil {
				if skillLevel, err := ExtractInt64(skillLevelValue); err == nil && skillLevel == 100 {
					if nvc, found := nvcLookup[int(abilityID)]; found {
						nvc.Checked = true
					}
				}
			}
		}
	}
	return
}

func (p *PR) loadEspers() (err error) {
	// Safely extract esper list with type validation
	esperList, err := SafeGetFromTargetRaw(p.UserData, OwnedMagicStoneList)
	if err != nil {
		return fmt.Errorf("failed to extract esper list: %w", err)
	}

	// Reset all espers to unchecked
	for _, e := range pr.Espers {
		e.Checked = false
	}

	// Extract esper IDs from the list
	esperIDs, err := ExtractInt64Array(esperList)
	if err != nil {
		return fmt.Errorf("failed to parse esper IDs: %w", err)
	}

	// Mark acquired espers as checked
	for _, esperID := range esperIDs {
		if esper, found := pr.EspersByValue[int(esperID)]; found {
			esper.Checked = true
		}
	}
	return nil
}

func (p *PR) loadMiscStats() (err error) {
	if models.GetMisc().GP, err = p.getInt(p.UserData, OwnedGil); err != nil {
		return
	}
	if models.GetMisc().Steps, err = p.getInt(p.UserData, Steps); err != nil {
		return
	}
	if models.GetMisc().EscapeCount, err = p.getInt(p.UserData, EscapeCount); err != nil {
		return
	}
	if models.GetMisc().BattleCount, err = p.getInt(p.UserData, BattleCount); err != nil {
		return
	}
	if models.GetMisc().NumberOfSaves, err = p.getInt(p.UserData, SaveCompleteCount); err != nil {
		return
	}
	if models.GetMisc().MonstersKilledCount, err = p.getInt(p.UserData, MonstersKilledCount); err != nil {
		return
	}
	if ds, ok := p.Base.GetValue(DataStorage); ok {
		m := jo.NewOrderedMap()
		if err = m.UnmarshalJSON([]byte(ds.(string))); err != nil {
			return
		}
		if models.GetMisc().CursedShieldFightCount, err = p.getIntFromSlice(m, "global"); err != nil {
			return
		}
	}
	return
}

func (p *PR) loadInventory(key string, inventory *pri.Inventory) (err error) {
	// Safely extract inventory list with type validation
	itemStrings, err := SafeGetFromTarget(p.UserData, key)
	if err != nil {
		return fmt.Errorf("failed to extract %s: %w", key, err)
	}

	inventory.Reset()
	for i, itemJSON := range itemStrings {
		row := pri.Row{}
		if err = json.Unmarshal([]byte(itemJSON), &row); err != nil {
			return fmt.Errorf("inventory[%d] unmarshal failed: %w", i, err)
		}
		inventory.Set(i, row)
	}
	return nil
}

func (p *PR) loadMapData() (err error) {
	md := pri.GetMapData()
	if md.MapID, err = p.getInt(p.MapData, MapID); err != nil {
		return
	}
	if md.PointIn, err = p.getInt(p.MapData, PointIn); err != nil {
		return
	}
	if md.TransportationID, err = p.getInt(p.MapData, TransportationID); err != nil {
		return
	}
	if md.CarryingHoverShip, err = p.getBool(p.MapData, CarryingHoverShip); err != nil {
		return
	}
	if md.PlayableCharacterCorpsID, err = p.getInt(p.MapData, PlayableCharacterCorpsId); err != nil {
		return
	}

	pe := jo.NewOrderedMap()
	if err = p.unmarshalFrom(p.MapData, PlayerEntity, pe); err != nil {
		return
	}
	var pos *jo.OrderedMap
	if pos = pe.Get(PlayerPosition).(*jo.OrderedMap); pos == nil {
		err = errors.New("unable to get transportation position")
		return
	}
	if md.Player.X, err = p.getFloat(pos, "x"); err != nil {
		return
	}
	if md.Player.Y, err = p.getFloat(pos, "y"); err != nil {
		return
	}
	if md.Player.Z, err = p.getFloat(pos, "z"); err != nil {
		return
	}
	if md.PlayerDirection, err = p.getInt(pe, PlayerDirection); err != nil {
		return
	}

	gps := jo.NewOrderedMap()
	if err = p.unmarshalFrom(p.MapData, GpsData, gps); err != nil {
		return
	}
	if md.Gps.TransportationID, err = p.getInt(gps, GpsTransportationID); err != nil {
		err = nil
	}
	if md.Gps.MapID, err = p.getInt(gps, GpsDataMapID); err != nil {
		return
	}
	if md.Gps.AreaID, err = p.getInt(gps, GpsDataAreaID); err != nil {
		return
	}
	if md.Gps.GpsID, err = p.getInt(gps, GpsDataID); err != nil {
		return
	}
	if md.Gps.Width, err = p.getInt(gps, GpsDataWidth); err != nil {
		return
	}
	if md.Gps.Height, err = p.getInt(gps, GpsDataHeight); err != nil {
		return
	}
	return
}

func (p *PR) loadTransportation() (err error) {
	// Safely extract transportation list with type validation
	transStrings, err := SafeGetFromTarget(p.UserData, OwnedTransportationList)
	if err != nil {
		return fmt.Errorf("failed to extract transportation list: %w", err)
	}

	pri.Transportations = make([]*pri.Transportation, len(transStrings))
	for index, transJSON := range transStrings {
		om := jo.NewOrderedMap()
		if err = om.UnmarshalJSON([]byte(transJSON)); err != nil {
			return fmt.Errorf("transportation[%d] unmarshal failed: %w", index, err)
		}

		t := &pri.Transportation{}
		if t.ID, err = p.getInt(om, TransID); err != nil {
			return fmt.Errorf("transportation[%d] ID extraction failed: %w", index, err)
		}
		if t.MapID, err = p.getInt(om, TransMapID); err != nil {
			return fmt.Errorf("transportation[%d] MapID extraction failed: %w", index, err)
		}
		if t.Direction, err = p.getInt(om, TransDirection); err != nil {
			return fmt.Errorf("transportation[%d] Direction extraction failed: %w", index, err)
		}
		if t.TimeStampTicks, err = p.getUint(om, TransTimeStampTicks); err != nil {
			return fmt.Errorf("transportation[%d] TimeStampTicks extraction failed: %w", index, err)
		}

		pos, ok := om.Get(TransPosition).(*jo.OrderedMap)
		if !ok {
			return fmt.Errorf("transportation[%d]: position is not an OrderedMap", index)
		}
		if pos == nil {
			return fmt.Errorf("transportation[%d]: position is nil", index)
		}

		if t.Position.X, err = p.getFloat(pos, "x"); err != nil {
			return fmt.Errorf("transportation[%d] X position failed: %w", index, err)
		}
		if t.Position.Y, err = p.getFloat(pos, "y"); err != nil {
			return fmt.Errorf("transportation[%d] Y position failed: %w", index, err)
		}
		if t.Position.Z, err = p.getFloat(pos, "z"); err != nil {
			return fmt.Errorf("transportation[%d] Z position failed: %w", index, err)
		}

		t.Enabled = t.TimeStampTicks > 0 && t.MapID > 0 && t.Position.X > 0 && t.Position.Y > 0 && t.Position.Z > 0
		pri.Transportations[index] = t
	}
	return nil
}

func (p *PR) loadVeldt() (err error) {
	veldt := pri.GetVeldt()

	// Safely extract veldt encounter flags with type validation
	encounterIDs, err := ExtractInt64Array(p.MapData.Get(BeastFieldEncountExchangeFlags))
	if err != nil {
		return fmt.Errorf("failed to extract veldt encounters: %w", err)
	}

	veldt.Encounters = make([]bool, len(encounterIDs))
	for i, encounterID := range encounterIDs {
		veldt.Encounters[i] = (encounterID == 1)
	}
	return nil
}

func (p *PR) loadCheats() (err error) {
	c := pri.GetCheats()
	if c.OpenedChestCount, err = p.getInt(p.UserData, OpenChestCount); err != nil {
		return
	}
	if c.IsCompleteFlag, err = p.getFlag(p.Base, IsCompleteFlag); err != nil {
		return
	}
	if c.PlayTime, err = p.getFloat(p.UserData, PlayTime); err != nil {
		return
	}
	return
}

func (p *PR) getString(c *jo.OrderedMap, key string) (s string, err error) {
	j, ok := c.GetValue(key)
	if !ok {
		err = fmt.Errorf("unable to find %s", key)
	}
	if s, ok = j.(string); !ok {
		err = fmt.Errorf("unable to parse field %s value %v ", key, j)
	}
	return
}

func (p *PR) getBool(c *jo.OrderedMap, key string) (b bool, err error) {
	j, ok := c.GetValue(key)
	if !ok {
		err = fmt.Errorf("unable to find %s", key)
	}
	if b, ok = j.(bool); !ok {
		err = fmt.Errorf("unable to parse field %s value %v", key, j)
	}
	return
}

func (p *PR) getInt(c *jo.OrderedMap, key string) (i int, err error) {
	j, ok := c.GetValue(key)
	if !ok {
		err = fmt.Errorf("unable to find %s", key)
	}

	k := reflect.ValueOf(j)
	switch k.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(k.Int()), nil
	case reflect.Float32, reflect.Float64:
		return int(k.Float()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(k.Uint()), nil
	case reflect.String:
		var l int64
		l, err = strconv.ParseInt(k.String(), 10, 32)
		if err == nil {
			i = int(l)
		}
	default:
		err = fmt.Errorf("unable to parse field %s value %v ", key, j)
	}
	return
}

func (p *PR) getUint(c *jo.OrderedMap, key string) (i uint64, err error) {
	j, ok := c.GetValue(key)
	if !ok {
		err = fmt.Errorf("unable to find %s", key)
	}

	k := reflect.ValueOf(j)
	switch k.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(k.Int()), nil
	case reflect.Float32, reflect.Float64:
		return uint64(k.Float()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint64(k.Uint()), nil
	case reflect.String:
		var l int64
		l, err = strconv.ParseInt(k.String(), 10, 64)
		if err == nil {
			i = uint64(l)
		}
	default:
		err = fmt.Errorf("unable to parse field %s value %v ", key, j)
	}
	return
}

func (p *PR) getFloat(c *jo.OrderedMap, key string) (f float64, err error) {
	j, ok := c.GetValue(key)
	if !ok {
		err = fmt.Errorf("unable to find %s", key)
	}

	k := reflect.ValueOf(j)
	switch k.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return k.Float(), nil
	case reflect.Float32, reflect.Float64:
		return k.Float(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return k.Float(), nil
	case reflect.String:
		f, err = strconv.ParseFloat(k.String(), 64)
	default:
		err = fmt.Errorf("unable to parse field %s value %v ", key, j)
	}
	return
}

func (p *PR) getJsonInts(data *jo.OrderedMap, key string) (ints []interface{}, err error) {
	// Deprecated: Use ExtractInt64Array() instead for type-safe extraction of arrays
	// This legacy method kept for backward compatibility with existing code paths
	j, ok := data.GetValue(key)
	if !ok {
		err = fmt.Errorf("unable to find %s", key)
		return
	}

	ints, ok = j.([]interface{})
	if !ok {
		err = fmt.Errorf("unable to load %s: expected array, got %T", key, j)
		return
	}

	// Validate that all elements are json.Number for safety
	for i, v := range ints {
		if _, ok := v.(json.Number); !ok {
			err = fmt.Errorf("invalid element at index %d: expected json.Number, got %T", i, v)
			return
		}
	}

	return
}

func (p *PR) getFlag(m *jo.OrderedMap, key string) (b bool, err error) {
	var i int
	if i, err = p.getInt(m, key); err != nil {
		return
	}
	if i == 0 {
		b = false
	} else {
		b = true
	}
	return
}

func (p *PR) getIntFromSlice(from *jo.OrderedMap, key string) (v int, err error) {
	sliceValue, ok := from.GetValue(key)
	if !ok {
		return 0, fmt.Errorf("unable to find %s", key)
	}

	// Safely extract the slice with type validation
	sliceArray, ok := sliceValue.([]interface{})
	if !ok {
		return 0, fmt.Errorf("expected array for %s, got %T", key, sliceValue)
	}

	if len(sliceArray) < 10 {
		return 0, fmt.Errorf("slice %s has insufficient elements: %d (need at least 10)", key, len(sliceArray))
	}

	// Get the element at index 9
	element := sliceArray[9]
	numValue, ok := element.(json.Number)
	if !ok {
		return 0, fmt.Errorf("element at index 9 is not json.Number, got %T", element)
	}

	i64, err := numValue.Int64()
	if err != nil {
		return 0, fmt.Errorf("failed to parse element at index 9 as int64: %w", err)
	}

	return int(i64), nil
}

func (p *PR) unmarshalFrom(from *jo.OrderedMap, key string, m *jo.OrderedMap) (err error) {
	i, ok := from.GetValue(key)
	if !ok {
		return fmt.Errorf("unable to find %s", key)
	}
	return m.UnmarshalJSON([]byte(i.(string)))
}

func (p *PR) unmarshal(i interface{}, m *map[string]interface{}) error {
	// s := strings.ReplaceAll(i.(string), `\\"`, `\"`)
	return json.Unmarshal([]byte(i.(string)), m)
}

func (p *PR) unmarshalEquipment(m *jo.OrderedMap) (idCounts []idCount, err error) {
	eqValue, ok := m.GetValue(EquipmentList)
	if !ok {
		return nil, fmt.Errorf("%s not found", EquipmentList)
	}

	// Type-safe extraction of equipment list string
	eqStr, err := ExtractString(eqValue)
	if err != nil {
		return nil, fmt.Errorf("equipment list value extraction failed: %w", err)
	}

	eq := jo.NewOrderedMap()
	if err = eq.UnmarshalJSON([]byte(eqStr)); err != nil {
		return nil, fmt.Errorf("equipment JSON unmarshal failed: %w", err)
	}

	valuesValue, ok := eq.GetValue("values")
	if !ok || valuesValue == nil {
		return idCounts, nil
	}

	// Validate and extract values array
	valuesArray, ok := valuesValue.([]interface{})
	if !ok {
		return nil, fmt.Errorf("equipment values: expected array, got %T", valuesValue)
	}

	idCounts = make([]idCount, len(valuesArray))
	for j, v := range valuesArray {
		valueStr, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("equipment[%d]: expected string, got %T", j, v)
		}
		if err = json.Unmarshal([]byte(valueStr), &idCounts[j]); err != nil {
			return nil, fmt.Errorf("equipment[%d] unmarshal failed: %w", j, err)
		}
	}
	return idCounts, nil
}

func (p *PR) getFromTarget(data *jo.OrderedMap, key string) (i interface{}, err error) {
	var (
		slTarget = jo.NewOrderedMap()
		ok       bool
	)
	if err = p.unmarshalFrom(data, key, slTarget); err != nil {
		return
	}
	if i, ok = slTarget.GetValue(targetKey); !ok {
		err = fmt.Errorf("unable to find %s", targetKey)
	}
	return
}

func (p *PR) fixEscapeCharsForLoad(s string) string {
	var (
		sb    strings.Builder
		going = false
		count int
		found = make(map[int]string)
	)
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			if !going {
				going = true
			}
			sb.WriteByte('\\')
			count++
		} else {
			if going {
				going = false
				sb.WriteByte('"')
				found[count] = sb.String()
				count = 0
				sb.Reset()
			}
		}
	}

	sorted := make([]int, 0, len(found))
	for k, _ := range found {
		sorted = append(sorted, k)
	}
	sort.Ints(sorted)

	for i := len(sorted) - 1; i >= 0; i-- {
		count = sorted[i]
		v, _ := found[count]
		sb.Reset()
		for j := 0; j < count; j++ {
			sb.WriteByte('~')
		}
		sb.WriteByte('"')
		found[count] = sb.String()
		s = strings.ReplaceAll(s, v, sb.String())
	}

	for i := len(sorted) - 1; i >= 0; i-- {
		count = sorted[i]
		v, _ := found[count]
		sb.Reset()
		half := (len(v) - 1) / 2
		for j := 0; j < half; j++ {
			sb.WriteByte('\\')
		}
		if sb.Len() == 0 {
			sb.WriteByte('\\')
		}
		sb.WriteByte('"')
		s = strings.ReplaceAll(s, v, sb.String())
	}
	return s
}

// {"keys":[1,2,3,4,5,6],"values":["{\\"contentId\\":113,\\"count\\":2}","{\\"contentId\\":214,\\"count\\":1}","{\\"contentId\\":268,\\"count\\":1}","{\\"contentId\\":233,\\"count\\":8}","{\\"contentId\\":200,\\"count\\":67}","{\\"contentId\\":315,\\"count\\":1}"]}
type equipment struct {
	Keys   []int    `json:"keys"`
	Values []string `json:"values"`
}

type idCount struct {
	ContentID int `json:"contentId"`
	Count     int `json:"count"`
}

func (p *PR) loadBase(s string) (err error) {
	return p.Base.UnmarshalJSON([]byte(s))
}

func (p *PR) uncheckAll(rages []*consts.NameValueChecked) {
	for _, r := range rages {
		r.Checked = false
	}
}

func (p *PR) replaceUnicodeNames(b []byte, names *[]unicodeNameReplace) ([]byte, error) {
	for i := 0; i < len(b)-10; i++ {
		if b[i] == '"' && b[i+1] == 'n' && b[i+2] == 'a' && b[i+3] == 'm' && b[i+4] == 'e' && b[i+5] == '\\' {
			i += 5
			for i < len(b)-1 && (b[i] == '\\' || b[i] == ':' || b[i] == ' ' || b[i] == '"') {
				i++
			}
			if b[i-1] != '"' {
				i--
			}

			// Start of name is i

			isUnicode := false
			for j := i; j < len(b)-50 && b[j] != '"' && !(b[j] == '\\' && b[j+1] == '\\'); j++ {
				if b[j] == '\\' && b[j+1] == 'x' {
					isUnicode = true
				}
			}
			if !isUnicode {
				continue
			}

			var original, replaced []byte
			for i < len(b)-1 && b[i] != '"' && !(b[i] == '\\' && b[i+1] == '\\') {
				original = append(original, b[i])
				replaced = append(replaced, '~')
				b[i] = '~'
				i++
			}
			r := unicodeNameReplace{
				Replaced: string(replaced),
			}
			var err error
			sb := strings.Builder{}
			for j := 0; j < len(original)-1; j++ {
				var char []byte
				if original[j] == '\\' && original[j+1] == 'x' {
					char = append(char, original[j])
					for j++; j < len(original); j++ {
						if original[j] == '\\' {
							break
						}
						char = append(char, original[j])
					}
					j--
					var s string
					if s, err = strconv.Unquote(strings.Replace(strconv.Quote(string(char)), `\\x`, `\x`, -1)); err != nil {
						return nil, err
					}
					sb.WriteString(s)
				} else {
					sb.WriteString(string(original[j]))
				}
			}
			r.Original = sb.String()
			*names = append(*names, r)
		}
	}
	return b, nil
}

type unicodeNameReplace struct {
	Original string
	Replaced string
}

func extractArchiveFile(dest string, f *zip.File) (err error) {
	var (
		rc   io.ReadCloser
		file *os.File
		path string
	)
	if rc, err = f.Open(); err != nil {
		return
	}
	defer func() { _ = rc.Close() }()

	path = filepath.Join(dest, f.Name)
	// Check for ZipSlip (Directory traversal)
	path = strings.ReplaceAll(path, "..", "")

	if f.FileInfo().IsDir() {
		if err = os.MkdirAll(path, f.Mode()); err != nil {
			return
		}
	} else {
		if err = os.MkdirAll(filepath.Dir(path), f.Mode()); err != nil {
			return
		}
		if file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode()); err != nil {
			return
		}
		defer func() { _ = file.Close() }()
		if _, err = io.Copy(file, rc); err != nil {
			return
		}
	}
	return
}

// debug_WriteOut writes JSON output to file for debugging
// Development utility: Used to examine loaded JSON structure during development
// Disabled in production builds - consider moving to debug build tags
func (p *PR) debug_WriteOut(out []byte) {
	var m map[string]any
	if err := json.Unmarshal(out, &m); err == nil {
		p.debug_WriteSection("configData", &m)
		p.debug_WriteSection("dataStorage", &m)
		if md, loaded := p.debug_LoadMap("mapData", &m); loaded {
			p.debug_WriteSection("gpsData", &md)
			p.debug_WriteSection("playerEntity", &md)
			m["mapData"] = md
		}
		if ud, loaded := p.debug_LoadMap("userData", &m, true); loaded {
			p.debug_WriteSection("corpsList", &ud, true)
			if cp, l := p.debug_LoadMap("corpsSlots", &ud, true); l {
				p.debug_LoadValue("values", &cp)
				ud["corpsSlots"] = cp
			}
			// p.debug_WriteSection("learnedAbilityList", &ud)
			// p.debug_WriteSection("importantOwendItemList", &ud)
			// p.debug_WriteSection("normalOwnedItemList", &ud)
			// p.debug_WriteSection("normalOwnedItemSortIdList", &ud)
			if ow, l := p.debug_LoadMap("ownedCharacterList", &ud, true); l {
				p.debug_LoadValue("target", &ow)
				sl := ow["target"].([]map[string]any)
				for _, v := range sl {
					p.debug_WriteSection("abilitySlotDataList", &v, true)
				}
				ud["ownedCharacterList"] = ow
			}
			// p.debug_WriteSection("ownedKeyWaordList", &ud)
			// p.debug_WriteSection("ownedMagicList", &ud)
			// p.debug_WriteSection("ownedMagicStoneList", &ud)
			// p.debug_WriteSection("releasedJobs", &ud)
			// p.debug_WriteSection("warehouseItemList", &ud)
			m["userData"] = ud
		}
		out, _ = json.MarshalIndent(m, "", "  ")
	}
	_ = os.WriteFile("loaded.json", out, 0755)
}

// debug_WriteSection extracts and writes a JSON section to file
// Development utility: Helper for debug_WriteOut
func (p *PR) debug_WriteSection(key string, parent *map[string]any, skip ...bool) {
	if m, loaded := p.debug_LoadMap(key, parent, skip...); loaded {
		(*parent)[key] = m
	}
}

// debug_LoadMap unmarshals a JSON string from parent map
// Development utility: Helper for debug_WriteOut
func (p *PR) debug_LoadMap(key string, parent *map[string]any, skip ...bool) (m map[string]any, loaded bool) {
	var v any
	if v, loaded = (*parent)[key]; loaded {
		s := v.(string)
		// for strings.Contains(s, `\\"`) {
		if len(skip) == 0 {
			s = strings.ReplaceAll(s, `\"`, `"`)
		}
		// }
		loaded = json.Unmarshal([]byte(s), &m) == nil
	}
	return
}

// debug_LoadValue converts JSON string array to typed array
// Development utility: Helper for debug_WriteOut
func (p *PR) debug_LoadValue(key string, parent *map[string]any) {
	var v []any
	if j, ok := (*parent)[key]; ok {
		if v, ok = j.([]any); ok {
			a := make([]map[string]any, len(v))
			for i, s := range v {
				if loaded := json.Unmarshal([]byte(s.(string)), &a[i]) == nil; !loaded {
					return
				}
			}
			(*parent)[key] = a
		}
	}
}

/*
d := []byte(s)
		for i, c := range d {
			if c == 'x' && d[i-1] == '\\' {
				d[i] = '^'
				d[i-1] = '^'
			}
		}

func (p *PR) fixFile(s string) (bool, string) {
	if i := strings.Index(s, "clearFlag"); i != -1 {
		c := s[i+9]
		if c != ':' {
			cc := s[i+10]
			if cc >= 48 && c <= 57 {
				cc -= 48
			} else {
				cc = 0
			}
			s = s[:i+9] + fmt.Sprintf(`":%d}`, cc)
		}
		return true, s
	} else if i = strings.Index(s, `"playTime`); i != -1 && s[i+4] >= 48 && s[i+4] <= 57 {
		s = s[0:i] + `"playTime":0.0,"clearFlag":0}`
		return true, s
	}
	return false, s
}*/
