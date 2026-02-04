package pr

import (
	"encoding/json"
	"errors"
	"fmt"

	pri "ffvi_editor/models/pr"

	jo "gitlab.com/c0b/go-ordered-json"
)

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

// convertToOrderedMap converts various types to *jo.OrderedMap
func convertToOrderedMap(i interface{}) (*jo.OrderedMap, error) {
	om := jo.NewOrderedMap()

	switch v := i.(type) {
	case string:
		if err := om.UnmarshalJSON([]byte(v)); err != nil {
			return nil, err
		}
	case *jo.OrderedMap:
		return v, nil
	case map[string]interface{}:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		if err := om.UnmarshalJSON(b); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unexpected transportation row type %T", i)
	}

	return om, nil
}

// parseTransportation extracts transportation data from an OrderedMap
func (p *PR) parseTransportation(om *jo.OrderedMap) (*pri.Transportation, error) {
	t := &pri.Transportation{}
	var err error

	if t.ID, err = p.getInt(om, TransID); err != nil {
		return nil, fmt.Errorf("failed to get transportation ID: %w", err)
	}
	if t.MapID, err = p.getInt(om, TransMapID); err != nil {
		return nil, fmt.Errorf("failed to get transportation MapID: %w", err)
	}
	if t.Direction, err = p.getInt(om, TransDirection); err != nil {
		return nil, fmt.Errorf("failed to get transportation Direction: %w", err)
	}
	if t.TimeStampTicks, err = p.getUint(om, TransTimeStampTicks); err != nil {
		return nil, fmt.Errorf("failed to get transportation TimeStampTicks: %w", err)
	}

	pos := om.Get(TransPosition).(*jo.OrderedMap)
	if pos == nil {
		return nil, errors.New("unable to get transportation position")
	}

	if t.Position.X, err = p.getFloat(pos, "x"); err != nil {
		return nil, fmt.Errorf("failed to get transportation X: %w", err)
	}
	if t.Position.Y, err = p.getFloat(pos, "y"); err != nil {
		return nil, fmt.Errorf("failed to get transportation Y: %w", err)
	}
	if t.Position.Z, err = p.getFloat(pos, "z"); err != nil {
		return nil, fmt.Errorf("failed to get transportation Z: %w", err)
	}

	t.Enabled = t.TimeStampTicks > 0 && t.MapID > 0 &&
		t.Position.X > 0 && t.Position.Y > 0 && t.Position.Z > 0

	return t, nil
}

func (p *PR) loadTransportation() error {
	sl, err := p.getFromTarget(p.UserData, OwnedTransportationList)
	if err != nil {
		return err
	}

	slSlice, err := decodeInterfaceSlice(sl, "transportation")
	if err != nil {
		return err
	}

	pri.Transportations = make([]*pri.Transportation, len(slSlice))

	for index, i := range slSlice {
		om, err := convertToOrderedMap(i)
		if err != nil {
			return fmt.Errorf("transportation[%d]: %w", index, err)
		}

		t, err := p.parseTransportation(om)
		if err != nil {
			return fmt.Errorf("transportation[%d]: %w", index, err)
		}

		pri.Transportations[index] = t
	}

	return nil
}

func (p *PR) loadVeldt() (err error) {
	var (
		veldt = pri.GetVeldt()
		sl    []interface{}
	)
	if sl, err = p.getJsonInts(p.MapData, BeastFieldEncountExchangeFlags); err != nil {
		return
	}
	veldt.Encounters = make([]bool, len(sl))
	for i, v := range sl {
		if v.(json.Number).String() == "1" {
			veldt.Encounters[i] = true
		}
	}
	return
}
