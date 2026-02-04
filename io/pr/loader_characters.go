package pr

import (
	"encoding/json"

	"ffvi_editor/models"
	"ffvi_editor/models/consts"
	"ffvi_editor/models/consts/pr"
	pri "ffvi_editor/models/pr"
	"ffvi_editor/settings"

	jo "gitlab.com/c0b/go-ordered-json"
)

// characterIdentity holds the basic identity info for a character
type characterIdentity struct {
	id      int
	jobID   int
	name    string
	enabled bool
}

// loadCharacterIdentity extracts basic character identity info
func (p *PR) loadCharacterIdentity(d *jo.OrderedMap) (characterIdentity, error) {
	var ci characterIdentity
	var err error

	if ci.id, err = p.getInt(d, ID); err != nil {
		return ci, err
	}
	if ci.jobID, err = p.getInt(d, JobID); err != nil {
		return ci, err
	}
	if ci.name, err = p.getString(d, Name); err != nil {
		return ci, err
	}
	if ci.enabled, err = p.getBool(d, IsEnableCorps); err != nil {
		return ci, err
	}
	return ci, nil
}

// loadCharacterStats loads HP, MP, and base stats from parameters
func (p *PR) loadCharacterStats(params *jo.OrderedMap, c *models.Character, baseOffset *pri.CharacterBase) error {
	var err error

	if c.Level, err = p.getInt(params, AdditionalLevel); err != nil {
		return err
	}

	if c.HP.Current, err = p.getInt(params, CurrentHP); err != nil {
		return err
	}
	if c.HP.Max, err = p.getInt(params, AdditionalMaxHp); err != nil {
		return err
	}
	c.HP.Max += baseOffset.HPBase

	if c.MP.Current, err = p.getInt(params, CurrentMP); err != nil {
		return err
	}
	if c.MP.Max, err = p.getInt(params, AdditionalMaxMp); err != nil {
		return err
	}
	c.MP.Max += baseOffset.MPBase

	if c.Vigor, err = p.getInt(params, AdditionalPower); err != nil {
		return err
	}
	if c.Stamina, err = p.getInt(params, AdditionalVitality); err != nil {
		return err
	}
	if c.Speed, err = p.getInt(params, AdditionalAgility); err != nil {
		return err
	}
	if c.Magic, err = p.getInt(params, AdditionMagic); err != nil {
		return err
	}
	return nil
}

// loadCharacterCommands loads the command list for a character
func (p *PR) loadCharacterCommands(d *jo.OrderedMap, c *models.Character) error {
	values, err := p.getFromTarget(d, CommandList)
	if err != nil {
		return err
	}
	vals, err := decodeInterfaceSlice(values, "CommandList")
	if err != nil {
		return err
	}
	for i, v := range vals {
		j, err := v.(json.Number).Int64()
		if err != nil {
			return err
		}
		if i >= len(c.Commands) {
			c.Commands = append(c.Commands, pr.CommandLookupByValue[int(j)])
		} else {
			c.Commands[i] = pr.CommandLookupByValue[int(j)]
		}
	}
	return nil
}

// skillLoader defines a job-specific skill loader
type skillLoader struct {
	check   func(int) bool
	from    int64
	to      int64
	lookup  map[int]*consts.NameValueChecked
	uncheck func()
	useID   bool // if true, use character ID instead of jobID
}

// loadCharacterSkills loads job-specific skills for a character
func (p *PR) loadCharacterSkills(d *jo.OrderedMap, id, jobID int) error {
	loaders := []skillLoader{
		{pr.IsJobWithBushido, pr.BushidoFrom, pr.BushidoTo, pr.BushidoLookupByID, func() { p.uncheckAll(pr.Bushidos) }, false},
		{pr.IsJobWithBlitz, pr.BlitzFrom, pr.BlitzTo, pr.BlitzLookupByID, func() { p.uncheckAll(pr.Blitzes) }, false},
		{pr.IsCharacterWithDance, pr.DanceFrom, pr.DanceTo, pr.DanceLookupByID, func() { p.uncheckAll(pr.Dances) }, true},
		{pr.IsJobWithLore, pr.LoreFrom, pr.LoreTo, pr.LoreLookupByID, func() { p.uncheckAll(pr.Lores) }, false},
		{pr.IsJobWithRage, pr.RageFrom, pr.RageTo, pr.RageLookupByID, func() { p.uncheckAll(pr.Rages) }, false},
	}

	for _, loader := range loaders {
		checkID := jobID
		if loader.useID {
			checkID = id
		}
		if loader.check(checkID) {
			loader.uncheck()
			if err := p.loadSkills(d, loader.from, loader.to, loader.lookup); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *PR) loadCharacters() error {
	// Load settings once outside the loop
	s := settings.NewManager("")
	_ = s.Load()
	autoEnableCmd := s.Get().AutoEnableCmd

	for _, d := range p.Characters {
		if d == nil {
			continue
		}

		// Load character identity
		identity, err := p.loadCharacterIdentity(d)
		if err != nil {
			return err
		}

		baseOffset, found := pri.GetCharacterBaseOffset(identity.id, identity.jobID)
		if !found {
			continue
		}

		c := pri.GetCharacter(baseOffset.Name)
		c.EnableCommandsSave = autoEnableCmd
		c.Name = identity.name
		c.IsEnabled = identity.enabled

		pri.GetParty().AddPossibleMember(&pri.Member{
			CharacterID: identity.id,
			Name:        c.Name,
		})

		// Load parameters
		params := jo.NewOrderedMap()
		if err := p.unmarshalFrom(d, Parameter, params); err != nil {
			return err
		}

		// Load stats
		if err := p.loadCharacterStats(params, c, baseOffset); err != nil {
			return err
		}

		if c.Exp, err = p.getInt(d, CurrentExp); err != nil {
			return err
		}

		// Load commands
		if err := p.loadCharacterCommands(d, c); err != nil {
			return err
		}

		// Load equipment
		if err := p.loadEquipment(d, c); err != nil {
			return err
		}

		// Load spells
		if err := p.loadSpells(d, c); err != nil {
			return err
		}

		// Load job-specific skills
		if err := p.loadCharacterSkills(d, identity.id, identity.jobID); err != nil {
			return err
		}
	}
	return nil
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

func (p *PR) loadSpells(d *jo.OrderedMap, c *models.Character) (err error) {
	var i interface{}
	if i, err = p.getFromTarget(d, AbilityList); err != nil {
		return
	}
	var sli []interface{}
	sli, err = decodeInterfaceSlice(i, "AbilityList")
	if err != nil {
		return err
	}
	for j := 0; j < len(sli); j++ {
		m := jo.NewOrderedMap()
		var s string
		s, err = decodeString(sli[j], "AbilityList entry")
		if err != nil {
			return err
		}
		if err = m.UnmarshalJSON([]byte(s)); err != nil {
			return
		}
		if v, ok := m.GetValue("abilityId"); ok {
			if iv, _ := v.(json.Number).Int64(); iv >= pr.SpellFrom && iv <= pr.SpellTo {
				if err = m.UnmarshalJSON([]byte(s)); err != nil {
					return
				}
				var spell *models.Spell
				if spell, ok = c.SpellsByID[int(iv)]; ok {
					k := m.Get("skillLevel")
					if kn, ok := k.(json.Number); ok {
						if iv, err = kn.Int64(); err == nil {
							spell.Value = int(iv)
						}
					}
				}
			}
		}
	}
	return
}

func (p *PR) loadSkills(d *jo.OrderedMap, from int64, to int64, nvcLookup map[int]*consts.NameValueChecked) (err error) {
	var (
		i     interface{}
		nvc   *consts.NameValueChecked
		found bool
		sli   []interface{}
	)
	if i, err = p.getFromTarget(d, AbilityList); err != nil {
		return
	}
	sli, err = decodeInterfaceSlice(i, "AbilityList")
	if err != nil {
		return err
	}
	for j := 0; j < len(sli); j++ {
		m := jo.NewOrderedMap()
		var s string
		s, err = decodeString(sli[j], "AbilityList entry")
		if err != nil {
			return err
		}
		if err = m.UnmarshalJSON([]byte(s)); err != nil {
			return
		}
		if v, ok := m.GetValue("abilityId"); ok {
			if iv, _ := v.(json.Number).Int64(); iv >= from && iv <= to {
				if err = m.UnmarshalJSON([]byte(s)); err != nil {
					return
				}
				k := m.Get("skillLevel")
				l, _ := k.(json.Number).Int64()
				if l == 100 {
					if nvc, found = nvcLookup[int(iv)]; found {
						nvc.Checked = true
					}
				}
			}
		}
	}
	return
}
