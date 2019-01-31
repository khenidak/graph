package graph

// *********** Adjaceny List ****************
type edge struct {
	dstIdx uint
	weight EdgeWeight
}

type edgeList []*edge

type adjacencyList struct {
	gdm        GraphDataMap
	data       []edgeList
	isDirected bool
}

func CreateAdjacencyList(gdm GraphDataMap, isDirected bool) GraphDataStore {
	al := &adjacencyList{
		gdm:        gdm,
		isDirected: isDirected,
	}

	al.data = make([]edgeList, 0) // <- TODO: improve this
	return al
}

// AddVertix adds a vertix to store
func (al *adjacencyList) AddVertix(gv GraphVertix) error {
	idx, err := al.gdm.AddMappedVertix(gv)
	if nil != err {
		return err
	}

	if idx < uint(len(al.data)) {
		return ErrVertixAlreadyExist
	}

	// add it - we use nil slice to allow us to
	// use adjacencyList for lazy graphs
	al.data = append(al.data, nil)
	return nil
}

// AddEdge adds an edge to store
func (al *adjacencyList) AddEdge(src GraphVertix, dst GraphVertix) error {
	return al.AddEdgeWeighted(src, 0, dst)
}

// AddEdgeWeighted adds an edge to store with weight
func (al *adjacencyList) AddEdgeWeighted(src GraphVertix, weight EdgeWeight, dst GraphVertix) error {
	vertixCount := uint(len(al.data))
	srcIdx, err := al.gdm.IdxOfVertix(src)
	if nil != err {
		return err
	}

	if vertixCount <= srcIdx {
		return ErrVertixNotFound
	}
	dstIdx, err := al.gdm.IdxOfVertix(dst)
	if nil != err {
		return err
	}

	if vertixCount <= dstIdx {
		return ErrVertixNotFound
	}

	// add Edge
	if nil == al.data[srcIdx] {
		al.data[srcIdx] = make([]*edge, 0)
	}
	al.data[srcIdx] = append(al.data[srcIdx], &edge{dstIdx: dstIdx, weight: weight})
	if !al.isDirected {
		if nil == al.data[dstIdx] {
			al.data[dstIdx] = make([]*edge, 0)
		}
		al.data[dstIdx] = append(al.data[dstIdx], &edge{dstIdx: srcIdx, weight: weight})
	}
	return nil
}

// ForEachVertix visits each vertix in the store (no order is guranteed)
func (al *adjacencyList) ForEachVertix(visitor VertixVisitor) {
	for idx, _ := range al.data {
		visitor(al.gdm, uint(idx))
	}
}

// ForEachAdjacentVertix visits each of an adjacent vertix of src
func (al *adjacencyList) ForEachAdjacentVertix(src GraphVertix, adjacentVisitor VertixAdjacentVisitor) error {
	srcIdx, err := al.gdm.IdxOfVertix(src)
	if nil != err {
		return err
	}

	return al.ForEachAdjacentVertixByIdx(srcIdx, adjacentVisitor)

}

// ForEachAdjacentVertixByIdx visits each of an adjacent vertix of src
func (al *adjacencyList) ForEachAdjacentVertixByIdx(srcIdx uint, adjacentVisitor VertixAdjacentVisitor) error {
	if uint(len(al.data)) <= srcIdx {
		return ErrVertixNotFound
	}

	for _, val := range al.data[srcIdx] {
		adjacentVisitor(al.gdm, srcIdx, val.weight, val.dstIdx)
	}
	return nil
}

// ListVertices returns a slice of all vertices in store
func (al *adjacencyList) ListVertices() (GraphVertices, error) {
	list := make([]GraphVertix, len(al.data))
	for idx, _ := range al.data {
		gv, err := al.gdm.VertixByIdx(uint(idx))
		if nil != err {
			return nil, err
		}
		list[idx] = gv
	}

	return list, nil
}

func (al *adjacencyList) IsDirected() bool {
	return al.isDirected
}

// *********** Adjaceny Matrix  ****************
type adjacencyMatrix struct {
	gdm         GraphDataMap
	data        [][]*EdgeWeight
	vertixCount uint
	isDirected  bool
}

func CreateAdjacencyMatrix(gdm GraphDataMap, vertixCount uint, isDirected bool) GraphDataStore {
	am := &adjacencyMatrix{
		gdm:         gdm,
		isDirected:  isDirected,
		vertixCount: vertixCount,
	}

	am.data = make([][]*EdgeWeight, vertixCount)
	// TODO: there has to be a better way than this
	for idx := uint(0); idx < am.vertixCount; idx++ {
		am.data[idx] = make([]*EdgeWeight, vertixCount)
	}
	return am
}

// AddVertix adds a vertix to store
func (am *adjacencyMatrix) AddVertix(gv GraphVertix) error {
	// we only add the mapping off
	idx, err := am.gdm.AddMappedVertix(gv)
	if nil != err {
		return err
	}

	if idx >= am.vertixCount {
		return ErrVertixOutOfRange
	}

	return nil
}

// AddEdge adds an edge to store
func (am *adjacencyMatrix) AddEdge(src GraphVertix, dst GraphVertix) error {
	return am.AddEdgeWeighted(src, 0, dst)
}

// AddEdgeWeighted adds an edge to store with weight
func (am *adjacencyMatrix) AddEdgeWeighted(src GraphVertix, weight EdgeWeight, dst GraphVertix) error {
	srcIdx, err := am.gdm.IdxOfVertix(src)
	if nil != err {
		return err
	}

	if am.vertixCount <= srcIdx {
		return ErrVertixNotFound
	}
	dstIdx, err := am.gdm.IdxOfVertix(dst)
	if nil != err {
		return err
	}

	if am.vertixCount <= dstIdx {
		return ErrVertixNotFound
	}

	// add Edge
	am.data[srcIdx][dstIdx] = &weight
	if !am.isDirected {
		am.data[dstIdx][srcIdx] = &weight
	}
	return nil
}

// ForEachVertix visits each vertix in the store (no order is guranteed)
func (am *adjacencyMatrix) ForEachVertix(visitor VertixVisitor) {
	for idx := uint(0); idx < am.vertixCount; idx++ {
		visitor(am.gdm, idx)
	}
}

// ForEachAdjacentVertix visits each of an adjacent vertix of src
func (am *adjacencyMatrix) ForEachAdjacentVertix(src GraphVertix, adjacentVisitor VertixAdjacentVisitor) error {
	srcIdx, err := am.gdm.IdxOfVertix(src)
	if nil != err {
		return err
	}

	return am.ForEachAdjacentVertixByIdx(srcIdx, adjacentVisitor)
}

// ForEachAdjacentVertixByIdx visits each of an adjacent vertix of src
func (am *adjacencyMatrix) ForEachAdjacentVertixByIdx(srcIdx uint, adjacentVisitor VertixAdjacentVisitor) error {
	if am.vertixCount <= srcIdx {
		return ErrVertixNotFound
	}

	for dstIdx := uint(0); dstIdx < am.vertixCount; dstIdx++ {
		if nil != am.data[srcIdx][dstIdx] {
			adjacentVisitor(am.gdm, srcIdx, *(am.data[srcIdx][dstIdx]), dstIdx)
		}
	}
	return nil
}

// ListVertices returns a slice of all vertices in store
func (am *adjacencyMatrix) ListVertices() (GraphVertices, error) {
	list := make([]GraphVertix, am.vertixCount)
	for idx := uint(0); idx < am.vertixCount; idx++ {
		gv, err := am.gdm.VertixByIdx(idx)
		if nil != err {
			return nil, err
		}
		list[idx] = gv
	}

	return list, nil
}

func (am *adjacencyMatrix) IsDirected() bool {
	return am.isDirected
}

//TODO:
// *********** Lazy Graph  ****************
// Lazy graph is where the data is semi represent
// as a graph but the cost of querying edges is high
// example a data base of people, with some sort of
// relation (manager, employe etc).
// so you need a struct that caches edges as they are
// requested and you are not sure how many edges your
// graph will have may contain all edges on set or some of
// the edges in set
