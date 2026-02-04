package pr

import (
	"encoding/json"
	"fmt"

	pri "ffvi_editor/models/pr"

	jo "gitlab.com/c0b/go-ordered-json"
)

func (p *PR) loadInventory(key string, inventory *pri.Inventory) (err error) {
	var (
		sl  interface{}
		row pri.Row
	)
	if sl, err = p.getFromTarget(p.UserData, key); err != nil {
		return
	}
	inventory.Reset()
	var slSlice []interface{}
	slSlice, err = decodeInterfaceSlice(sl, "inventory")
	if err != nil {
		return err
	}
	for i, r := range slSlice {
		switch v := r.(type) {
		case string:
			if err = json.Unmarshal([]byte(v), &row); err != nil {
				return
			}
		case *jo.OrderedMap:
			var b []byte
			if b, err = v.MarshalJSON(); err != nil {
				return
			}
			if err = json.Unmarshal(b, &row); err != nil {
				return
			}
		case map[string]interface{}:
			var b []byte
			if b, err = json.Marshal(v); err != nil {
				return
			}
			if err = json.Unmarshal(b, &row); err != nil {
				return
			}
		default:
			return fmt.Errorf("unexpected inventory row type %T", r)
		}
		inventory.Set(i, row)
	}
	return nil
}

func (p *PR) unmarshalEquipment(m *jo.OrderedMap) (idCounts []idCount, err error) {
	i, ok := m.GetValue(EquipmentList)
	if !ok {
		return nil, fmt.Errorf("%s not found", EquipmentList)
	}

	eq := jo.NewOrderedMap()
	if err = eq.UnmarshalJSON([]byte(i.(string))); err != nil {
		return
	}

	if i, ok = eq.GetValue("values"); ok && i != nil {
		valSlice, ok := i.([]interface{})
		if !ok {
			return nil, fmt.Errorf("expected []interface{} for equipment values, got %T", i)
		}
		idCounts = make([]idCount, len(valSlice))
		for j, v := range valSlice {
			str, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("expected string in equipment values, got %T", v)
			}
			if err = json.Unmarshal([]byte(str), &idCounts[j]); err != nil {
				return
			}
		}
	}
	return
}
