package pr

import (
	"encoding/json"
	"fmt"

	"ffvi_editor/global"
	"ffvi_editor/io/file"
	"ffvi_editor/models/consts"
	pri "ffvi_editor/models/pr"

	jo "gitlab.com/c0b/go-ordered-json"
)

// loadSection is a helper to load a specific section of the save file
type loadSection func() error

// loadCharactersFromData loads character data from the UserData section
func (p *PR) loadCharactersFromData() error {
	i, err := p.getFromTarget(p.UserData, OwnedCharacterList)
	if err != nil {
		return err
	}

	chars, err := decodeInterfaceSlice(i, "OwnedCharacterList")
	if err != nil {
		return err
	}

	for j, c := range chars {
		if p.Characters[j] == nil {
			p.Characters[j] = jo.NewOrderedMap()
		}
		s, err := decodeString(c, "OwnedCharacterList entry")
		if err != nil {
			return err
		}
		if err = p.Characters[j].UnmarshalJSON([]byte(s)); err != nil {
			return err
		}
	}
	return nil
}

// loadGameData loads all game data sections after character data is loaded
func (p *PR) loadGameData() error {
	sections := []struct {
		name string
		fn   loadSection
	}{
		{"characters", p.loadCharacters},
		{"party", p.loadParty},
		{"espers", p.loadEspers},
		{"misc stats", p.loadMiscStats},
		{"normal inventory", func() error { return p.loadInventory(NormalOwnedItemList, pri.GetInventory()) }},
		{"important inventory", func() error { return p.loadInventory(importantOwnedItemList, pri.GetImportantInventory()) }},
		{"veldt", p.loadVeldt},
		{"cheats", p.loadCheats},
		{"map data", p.loadMapData},
		{"transportation", p.loadTransportation},
	}

	for _, section := range sections {
		if err := section.fn(); err != nil {
			return fmt.Errorf("failed to load %s: %w", section.name, err)
		}
	}
	return nil
}

func (p *PR) Load(fromFile string, saveType global.SaveFileType) error {
	out, fileTrimmed, err := file.LoadFile(fromFile, saveType)
	if err != nil {
		return err
	}
	p.fileTrimmed = fileTrimmed

	s := string(out)

	// Load base structures
	if err := p.loadBase(s); err != nil {
		return fmt.Errorf("failed to load base: %w", err)
	}

	if err := p.unmarshalFrom(p.Base, UserData, p.UserData); err != nil {
		return fmt.Errorf("failed to load UserData: %w", err)
	}

	if err := p.unmarshalFrom(p.Base, MapData, p.MapData); err != nil {
		return fmt.Errorf("failed to load MapData: %w", err)
	}

	// Load character raw data
	if err := p.loadCharactersFromData(); err != nil {
		return fmt.Errorf("failed to load character data: %w", err)
	}

	// Clear party before loading game data
	pri.GetParty().Clear()

	// Load all game data sections
	if err := p.loadGameData(); err != nil {
		return err
	}

	return nil
}

func (p *PR) loadParty() (err error) {
	var (
		party  = pri.GetParty()
		i      interface{}
		member pri.Member
	)

	if i, err = p.getFromTarget(p.UserData, CorpsList); err != nil {
		return
	}
	var slots []interface{}
	slots, err = decodeInterfaceSlice(i, "CorpsList")
	if err != nil {
		return err
	}
	for slot, c := range slots {
		var str string
		str, err = decodeString(c, "CorpsList entry")
		if err != nil {
			return err
		}
		if err = json.Unmarshal([]byte(str), &member); err != nil {
			return
		}
		if err = party.SetMemberByID(slot, member.CharacterID); err != nil {
			return
		}
	}
	return
}

func (p *PR) loadBase(s string) (err error) {
	return p.Base.UnmarshalJSON([]byte(s))
}

func (p *PR) uncheckAll(rages []*consts.NameValueChecked) {
	for _, r := range rages {
		r.Checked = false
	}
}

func (p *PR) unmarshalFrom(from *jo.OrderedMap, key string, m *jo.OrderedMap) (err error) {
	i, ok := from.GetValue(key)
	if !ok || i == nil {
		return fmt.Errorf("unable to find %s", key)
	}
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("unable to unmarshal %s: expected string", key)
	}
	return m.UnmarshalJSON([]byte(s))
}

func (p *PR) unmarshal(i interface{}, m *map[string]interface{}) error {
	return json.Unmarshal([]byte(i.(string)), m)
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

type unicodeNameReplace struct {
	Original string
	Replaced string
}

type idCount struct {
	ContentID int `json:"contentId"`
	Count     int `json:"count"`
}
