package pr

import (
	"testing"

	jo "gitlab.com/c0b/go-ordered-json"
)

// TestHelpers provides utilities for testing PR save parsing functionality
type TestHelpers struct {
	t *testing.T
}

// NewTestHelpers creates a new TestHelpers instance
func NewTestHelpers(t *testing.T) *TestHelpers {
	return &TestHelpers{t: t}
}

// CreateOrderedMap creates an OrderedMap from a JSON string, failing the test if it can't be unmarshaled
func (th *TestHelpers) CreateOrderedMap(jsonStr string) *jo.OrderedMap {
	om := jo.NewOrderedMap()
	if err := om.UnmarshalJSON([]byte(jsonStr)); err != nil {
		th.t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	return om
}

// CreateMinimalCharacterJSON creates a minimal character JSON object for testing
func (th *TestHelpers) CreateMinimalCharacterJSON() string {
	return `{"corpseId":1,"jobId":1,"isEnableCorps":true,"name":"TestCharacter","currentExp":0,"currentHp":0,"parameter":{"additionalLevel":1,"currentHp":50,"additionalMaxHp":0,"currentMp":20,"additionalMaxMp":0,"additionalPower":10,"additionalVitality":10,"additionalAgility":10,"additionMagic":10},"commandList":[10,11,12,13,14,15,16,17,18],"equipmentList":"{\"values\": []}","abilityList":"{}","learnedSkillList":[]}`
}

// CreateMinimalMapDataJSON creates a minimal map data JSON object for testing
func (th *TestHelpers) CreateMinimalMapDataJSON() string {
	return `{"mapId":1,"pointIn":0,"transportationId":-1,"carryingHoverShip":false,"playableCharacterCorpsId":1,"playerEntity":"{\"position\":{\"x\":100.0,\"y\":100.0,\"z\":0.0},\"direction\":0}","gpsData":"{\"transportationId\":-1,\"mapId\":1,\"areaId\":0,\"gpsId\":0,\"width\":10,\"height\":10}","beastFieldEncountExchangeFlags":[]}`
}

// CreateMinimalUserDataJSON creates a minimal user data JSON object for testing
func (th *TestHelpers) CreateMinimalUserDataJSON() string {
	return `{"owendGil":0,"Steps":0,"escapeCount":0,"battleCount":0,"saveCompleteCount":0,"monstersKilledCount":0,"playTime":0.0,"openChestCount":0,"ownedCharacterList":[],"ownedMagicStoneList":[],"normalOwnedItemList":[],"importantOwendItemList":[],"ownedTransportationList":[]}`
}

// CreateMinimalBaseJSON creates a minimal base JSON structure
func (th *TestHelpers) CreateMinimalBaseJSON() string {
	return `{"userData":"{\"owendGil\":0}","mapData":"{\"mapId\":1}","isCompleteFlag":0,"dataStorage":"{\"global\":[0,0,0,0,0,0,0,0,0,0]}"}`
}

// AssertNoError is a helper to check errors during testing
func (th *TestHelpers) AssertNoError(err error, msg string) {
	if err != nil {
		th.t.Fatalf("%s: %v", msg, err)
	}
}

// AssertError checks that an error occurred as expected
func (th *TestHelpers) AssertError(err error, msg string) {
	if err == nil {
		th.t.Fatalf("%s: expected error but got none", msg)
	}
}

// AssertEquals compares two values
func (th *TestHelpers) AssertEquals(expected, actual interface{}, msg string) {
	if expected != actual {
		th.t.Fatalf("%s: expected %v but got %v", msg, expected, actual)
	}
}
