package graph

// general purpose error
type GraphError string

func (ge GraphError) Error() string { return string(ge) }

const ErrVertixNotFound = GraphError("Vertix not found")
const ErrVertixAlreadyExist = GraphError("Vertix already exist")
const ErrVertixOutOfRange = GraphError("Vertix idx out of range")

// VertixId is whatever the caller wants to use
// as identfier for a vertix. used in mapping vertixIdx(uint) to vertix
type VertixId interface{}
type EdgeWeight float64

// GraphVertix defines a vertix
type GraphVertix interface {
	Id() VertixId
}

// GraphVertices is a slice of GraphVertix
type GraphVertices []GraphVertix

// VertixVisitor general purpose func
// pointer for visiting a vertix
type VertixVisitor func(gdm GraphDataMap, current uint)
type VertixAdjacentVisitor func(gdm GraphDataMap, current uint, weight EdgeWeight, adjacent uint)

// GraphStore defines a typical store
// of vertices and edges. implemented
// as adjacency matrix or adjacency list
// or if the data is represented somewhere
// caller can just implement the interface
// to avoid memory duplication
type GraphDataStore interface {
	// AddVertix adds a vertix to store
	AddVertix(gv GraphVertix) error
	// AddEdge adds an edge to store
	AddEdge(src GraphVertix, dst GraphVertix) error
	// AddEdgeWeighted adds an edge to store with weight
	AddEdgeWeighted(src GraphVertix, weight EdgeWeight, dst GraphVertix) error
	/*
		// HasEdge retrun true if an edge that connects src to dst
		HasEdge(src GraphVertix, dst GraphVertix) (int64, error)
		// HasEdgeByIdx same as HasEdge but expected indexs
		HasEdgeByIdx(srcIdx uint, dstIdx uint) (int64, error)
	*/
	// ForEachVertix visits each vertix in the store (no order is guranteed)
	ForEachVertix(visitor VertixVisitor)
	// ForEachAdjacentVertix visits each of an adjacent vertix of src
	ForEachAdjacentVertix(src GraphVertix, adjacentVisitor VertixAdjacentVisitor) error
	// ForEachAdjacentVertixByIdx visits each of an adjacent vertix of src
	ForEachAdjacentVertixByIdx(srcIdx uint, adjacentVisitor VertixAdjacentVisitor) error
	// ListVertices returns a slice of all vertices in store
	ListVertices() (GraphVertices, error)
	// IsDirected return true if the graph is directed
	IsDirected() bool
}

// GraphDataMap defines the relation
// between a vertix and its idx
type GraphDataMap interface {
	VertixByIdx(idx uint) (GraphVertix, error)
	IdxOfVertix(gv GraphVertix) (uint, error)
	AddMappedVertix(gv GraphVertix) (uint, error)
}

type Graph interface {
	GraphDataStore
}

// typical implementation GraphDataMap
// user can replace this, specially if
// such mapped data exist some where
type graph_map struct {
	count    uint
	capacity uint
	_data    []GraphVertix
	_map     map[VertixId]uint
}

func CreateDataGraphMapWithCapacity(capacity uint) GraphDataMap {
	gm := &graph_map{}
	gm._data = make([]GraphVertix, capacity)
	gm._map = make(map[VertixId]uint)
	gm.capacity = capacity

	return gm
}

func CreateDataGraphMap() GraphDataMap {
	return CreateDataGraphMapWithCapacity(0)
}

func (gm *graph_map) VertixByIdx(idx uint) (GraphVertix, error) {
	if idx >= gm.count {
		return nil, ErrVertixNotFound
	}

	return gm._data[idx], nil
}
func (gm *graph_map) IdxOfVertix(gv GraphVertix) (uint, error) {
	var idx uint
	var ok bool
	if idx, ok = gm._map[gv.Id()]; !ok {
		return 0, ErrVertixNotFound
	}

	return idx, nil
}
func (gm *graph_map) AddMappedVertix(gv GraphVertix) (uint, error) {
	if _, ok := gm._map[gv.Id()]; ok {
		return 0, ErrVertixAlreadyExist
	}

	newCount := gm.count + 1
	if newCount > gm.capacity {
		gm._data = append(gm._data, gv)
	} else {
		gm._data[gm.count] = gv
	}
	gm._map[gv.Id()] = gm.count
	gm.count = newCount
	return gm.count - 1, nil
}

type graphObject struct {
	GraphDataMap
	GraphDataStore
}

// Create a graph using a store and map. Use that
// if your data is already repsented as the graph in question
func CreateGraph(gdm GraphDataMap, store GraphDataStore) Graph {
	return &graphObject{gdm, store}
}

//Creates graph based on AdjacencyList
func CreateAdjacencyListGraph(gdm GraphDataMap, isDirected bool) Graph {
	al := CreateAdjacencyList(gdm, isDirected)
	return CreateGraph(gdm, al)
}

// Creates graph based on AdjacencyMatrix
func CreateAdjacencyMatrixGraph(gdm GraphDataMap, isDirected bool, vertixCount uint) Graph {
	am := CreateAdjacencyMatrix(gdm, vertixCount, isDirected)
	return CreateGraph(gdm, am)
}
