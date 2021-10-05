package main

import (
	"fmt"
	"time"
)

type Vertex int
type Edge struct {
	V1 Vertex
	V2 Vertex
}

type LWWSet struct {
	vertexAdded   map[Vertex]time.Time
	vertexRemoved map[Vertex]time.Time
	edgeAdded     map[Edge]time.Time
	edgeRemoved   map[Edge]time.Time
}

// New constructs a new LWW set
func New() *LWWSet {
	return &LWWSet{
		vertexAdded:   make(map[Vertex]time.Time),
		vertexRemoved: make(map[Vertex]time.Time),
		edgeAdded:     make(map[Edge]time.Time),
		edgeRemoved:   make(map[Edge]time.Time),
	}
}

// AddVertex adds a new vertex, recording the timestamp with it
func (l *LWWSet) AddVertex(v Vertex) {
	l.vertexAdded[v] = Now()
}

// RemoveVertex removes a vertex if it's in the graph and not used in an edge
func (l *LWWSet) RemoveVertex(v Vertex) error {
	if l.ContainsVertex(v) && !l.VertexInEdge(v) {
		l.vertexRemoved[v] = Now()
		return nil
	}
	return fmt.Errorf("cannot remove vertex %#v", v)
}

// ContainsVertex returns true if the graph contains the vertex
// (i.e. if in vertexAdded and not in vertexRemoved, or in vertexRemoved with an earlier timestamp)
func (l *LWWSet) ContainsVertex(v Vertex) bool {
	tsAdd, ok := l.vertexAdded[v]
	if !ok {
		return false
	}

	tsRem, ok := l.vertexRemoved[v]
	if !ok {
		return true
	}

	return tsAdd.After(tsRem)
}

// AddEdge adds and edge in the graph if both vertices are in it
func (l *LWWSet) AddEdge(v1 Vertex, v2 Vertex) error {
	if l.ContainsVertex(v1) && l.ContainsVertex(v2) {
		l.edgeAdded[Edge{
			V1: v1,
			V2: v2,
		}] = Now()
		return nil
	}
	return fmt.Errorf("cannot add edge {%#v, %#v}", v1, v2)
}

// RemoveEdge adds an edge in edgeRemoved with timestamp
func (l *LWWSet) RemoveEdge(e Edge) error {
	if l.ContainsEdge(e) {
		l.edgeRemoved[e] = Now()
		return nil
	}
	return fmt.Errorf("cannot remove edge %#v", e)
}

// ContainsEdge returns true if bothg vertices are in graph and edge is in graph (and not in edgeRemoved / in there with a higher timestamp)
func (l *LWWSet) ContainsEdge(e Edge) bool {
	ok := l.ContainsVertex(e.V1) && l.ContainsVertex(e.V2)
	if !ok {
		return false
	}

	tsAdd, ok := l.edgeAdded[e]
	if !ok {
		return false
	}

	tsRem, ok := l.edgeRemoved[e]
	if !ok {
		return true
	}

	return tsAdd.After(tsRem)
}

// VertexInEdge returns true if the vertex is contained within any edge in the graph
func (l *LWWSet) VertexInEdge(v Vertex) bool {
	for e := range l.edgeAdded {
		if l.ContainsEdge(e) {
			return e.V1 == v || e.V2 == v
		}
	}
	return false
}

// ConnectedVertices returns all vertices that have a connection with a single vertex
func (l *LWWSet) ConnectedVertices(v Vertex) []Vertex {
	var res []Vertex
	for e := range l.edgeAdded {
		if l.ContainsEdge(e) {
			if e.V1 == v {
				res = append(res, e.V2)
			} else if e.V2 == v {
				res = append(res, e.V1)
			}
		}
	}
	return res
}

// Path finds one path between two vertices
func (l *LWWSet) Path(source Vertex, target Vertex) []Vertex {
	var res []Vertex
	if !l.ContainsVertex(source) || !l.ContainsVertex(target) {
		return res
	}

	if source == target {
		return []Vertex{source}
	}

	visited := make(map[Vertex]bool)
	var conn [][]Vertex

	pop := func(a *[]Vertex) *Vertex {
		if len(*a) == 0 {
			return nil
		}
		rm := (*a)[len(*a)-1]
		*a = (*a)[:len(*a)-1]
		return &rm
	}

	popcon := func(a *[][]Vertex) *[]Vertex {
		if len(*a) == 0 {
			return nil
		}
		rm := (*a)[len(*a)-1]
		*a = (*a)[:len(*a)-1]
		return &rm
	}

	shift := func(a *[]Vertex) *Vertex {
		if len(*a) == 0 {
			return nil
		}
		rm := (*a)[0]
		*a = (*a)[1:]
		return &rm
	}

	fwd := func(v Vertex) {
		res = append(res, v)
		visited[v] = true

		var tmp []Vertex
		for _, c := range l.ConnectedVertices(v) {
			if _, ok := visited[c]; !ok {
				tmp = append(tmp, c)
			}
		}
		conn = append(conn, tmp)
	}

	bck := func() {
		rm := pop(&res)
		if rm != nil {
			delete(visited, *rm)
		}
		popcon(&conn)
	}

	fwd(source)
	for len(res) > 0 {
		c := popcon(&conn)
		if c != nil && len(*c) > 0 {
			next := shift(c)
			conn = append(conn, *c)
			if l.ContainsVertex(*next) {
				fwd(*next)
			}
		} else {
			conn = append(conn, *c)
			bck()
			continue
		}
		if res[len(res)-1] == target {
			return res
		}
	}

	return res
}

// Merge merges twqo LWW sets
func (l *LWWSet) Merge(r *LWWSet) {
	latest := func(local time.Time, remote time.Time) time.Time {
		if local.After(remote) {
			return local
		}
		return remote
	}

	mergeVertex := func(local map[Vertex]time.Time, remote map[Vertex]time.Time) {
		for key, rv := range remote {
			local[key] = latest(local[key], rv)
		}
	}

	mergeEdge := func(local map[Edge]time.Time, remote map[Edge]time.Time) {
		for key, rv := range remote {
			local[key] = latest(local[key], rv)
		}
	}

	mergeVertex(l.vertexAdded, r.vertexAdded)
	mergeVertex(l.vertexRemoved, r.vertexRemoved)
	mergeEdge(l.edgeAdded, r.edgeAdded)
	mergeEdge(l.edgeRemoved, r.edgeRemoved)
}
