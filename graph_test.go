package graph

import (
	"fmt"
	"testing"
)

// GRAPH DATA MAPPER TESTS
func buildGrapMap(gm GraphDataMap, ordered []string, mapped map[string]uint) error {
	for idx, key := range ordered {
		inIdx, err := gm.AddMappedVertix(stringVertix(key))
		if inIdx != uint(idx) {
			return fmt.Errorf("When adding vertix found %v and expted :%v", inIdx, idx)
		}

		if nil != err {
			return err
		}
	}
	return nil
}

func validateCorrectGraphMap(gm GraphDataMap, mapped map[string]uint) error {
	for key, value := range mapped {
		inKey, err := gm.VertixByIdx(value)
		if nil != err {
			return err
		}
		if inKey.Id() != key {
			return fmt.Errorf("Found Id:%v for idx:%v expected:%v", inKey.Id(), value, key)
		}

		inValue, err := gm.IdxOfVertix(stringVertix(key))
		if nil != err {
			return err
		}

		if inValue != value {
			return fmt.Errorf("found idx:%v for Id:%v expected:%v", inValue, key, value)
		}
	}
	return nil
}
func TestCapGraghDataMapper(t *testing.T) {

	withCap := CreateDataGraphMapWithCapacity(3)
	if err := buildGrapMap(withCap, ordered, mapped); nil != err {
		t.Errorf("failed to create map with capaicty")
	}
	if err := validateCorrectGraphMap(withCap, mapped); nil != err {
		t.Errorf(err.Error())
	}
}

func TestGraphDataMapper(t *testing.T) {
	withoutCap := CreateDataGraphMap()
	if err := buildGrapMap(withoutCap, ordered, mapped); nil != err {
		t.Errorf("failed to create map without capaicty")
	}

	if err := validateCorrectGraphMap(withoutCap, mapped); nil != err {
		t.Errorf(err.Error())
	}
}
