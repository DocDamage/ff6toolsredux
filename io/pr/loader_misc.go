package pr

import (
	"encoding/json"
	"fmt"

	"ffvi_editor/models"
	"ffvi_editor/models/consts/pr"
	pri "ffvi_editor/models/pr"

	jo "gitlab.com/c0b/go-ordered-json"
)

func (p *PR) loadEspers() (err error) {
	var espers interface{}
	if espers, err = p.getFromTarget(p.UserData, OwnedMagicStoneList); err != nil {
		return
	}
	var id int64
	for _, e := range pr.Espers {
		e.Checked = false
	}
	if espers != nil {
		var esperSlice []interface{}
		esperSlice, err = decodeInterfaceSlice(espers, "espers")
		if err != nil {
			return err
		}
		for _, n := range esperSlice {
			if id, err = n.(json.Number).Int64(); err != nil {
				return
			}
			if e, found := pr.EspersByValue[int(id)]; found {
				e.Checked = true
			}
		}
	}
	return
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

func (p *PR) getIntFromSlice(from *jo.OrderedMap, key string) (v int, err error) {
	var (
		i, ok = from.GetValue(key)
		sl    []interface{}
		i64   int64
	)
	if !ok {
		err = fmt.Errorf("unable to find %s", key)
		return
	}

	if sl, ok = i.([]interface{}); !ok || len(sl) < 9 {
		err = fmt.Errorf("unable to load cursed shield battle count")
		return
	}

	if i64, err = sl[9].(json.Number).Int64(); err != nil {
		return
	}
	v = int(i64)
	return
}
