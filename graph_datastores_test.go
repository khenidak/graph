package graph

import (
	"testing"
)

func fillDataStore(store GraphDataStore) error {
	for vertix, _ := range graph_small_data {
		if err := store.AddVertix(stringVertix(vertix)); nil != err {
			return err
		}
	}
	// add adjacent (now that we have all vertices in
	// place)
	for vertix, adjacentlist := range graph_small_data {
		for _, adjacentVertix := range adjacentlist {
			if err := store.AddEdge(stringVertix(vertix), stringVertix(adjacentVertix)); nil != err {
				return err
			}
		}
	}
	return nil
}

// on the correct vertices
// are store by the data store
func validateDataStore_data(store GraphDataStore, t *testing.T) {
	seen := make([]bool, len(ordered))
	store.ForEachVertix(func(gdm GraphDataMap, current uint) {
		gv, err := gdm.VertixByIdx(current)
		if nil != err {
			t.Fatalf("failed to validate data store, idx can not be mapped to graph vertix with error:%v", err.Error())
		}

		for idx, value := range ordered {
			var sv stringVertix
			var ok bool
			if sv, ok = gv.(stringVertix); !ok {
				t.Fatalf("failed to validate data store, can not convert GraphVertix to stringVertix")
			}
			if value == string(sv) {
				if true == seen[idx] {
					t.Fatalf("failed to validate data store, store is returning duplicate vertices")
				}
				seen[idx] = true
				break
			}
		}
	})

	for idx, _ := range seen {
		if false == seen[idx] {
			t.Fatalf("failed to validate data store, store is not returning all the vertices")
		}
	}
}

// for each vertix, we enum the adjacent and we check that
// they all exist
func validateDataStore_adjacent(store GraphDataStore, t *testing.T) {
	store.ForEachVertix(func(gdm GraphDataMap, current uint) {
		var sv stringVertix
		var ok bool

		//convert it back to vertix
		gv, err := gdm.VertixByIdx(current)
		if nil != err {
			t.Fatalf("failed to validate data store, idx can not be mapped to graph vertix with error:%v", err.Error())
		}
		// convert it back to string vertix
		if sv, ok = gv.(stringVertix); !ok {
			t.Fatalf("failed to validate data store, can not convert GraphVertix to stringVertix")
		}
		if !ok {
			t.Fatalf("failed to validate data store, can not find vertix in testing data")
		}

		seen := make([]bool, len(graph_small_data[string(sv)])) // used to mark seen adjacent vertices
		err = store.ForEachAdjacentVertixByIdx(current, func(gdm GraphDataMap, src uint, weight EdgeWeight, adjacent uint) {
			dstGv, err := gdm.VertixByIdx(adjacent)
			if nil != err {
				t.Fatalf("failed to validate data store, adjacent idx can not be converted to GraphVertix")
			}
			dstSv, ok := dstGv.(stringVertix)
			if !ok {
				t.Fatalf("failed to validate data store, can not convert dst GraphVertix to stringVertix")
			}

			//mark it as seen
			for idx, value := range graph_small_data[string(sv)] {
				if value == string(dstSv) {
					if true == seen[idx] {
						t.Fatalf("failed to validate data store, there are duplicate adjacent")
					}
					seen[idx] = true
				}
			}
		})
		// if any marked as not seen then
		for idx, value := range seen {
			if false == value {
				t.Fatalf("failed to validate data store, an adjacent vertix was not found: %v -> %v", string(sv), graph_small_data[string(sv)][idx])
			}
		}

		if nil != err {
			t.Fatalf("failed to enumarate adjacent vertices with error:%v", err.Error())
		}
	})
}

func TestAdjacencyList(t *testing.T) {
	gdm := CreateDataGraphMap()
	al := CreateAdjacencyList(gdm, true)
	err := fillDataStore(al)
	if nil != err {
		t.Fatalf("failed to fill AdjacencyList with error:%v ", err)
	}
	validateDataStore_data(al, t)
	validateDataStore_adjacent(al, t)
}

func TestAdjacencyMatrix(t *testing.T) {
	gdm := CreateDataGraphMap()
	am := CreateAdjacencyMatrix(gdm, 6, false)
	err := fillDataStore(am)
	if nil != err {
		t.Fatalf("failed to fill AdjacencyMatrix with error:%v ", err)
	}
	validateDataStore_data(am, t)
	validateDataStore_adjacent(am, t)
}
