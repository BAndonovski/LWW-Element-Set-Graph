package main

import (
	"testing"
	"time"
)

func TestAddsVertexWithCurrentTimestamp(t *testing.T) {
	lww := New()
	MockTime(time.Now())
	lww.AddVertex(1)
	if !(len(lww.vertexAdded) == 1) ||
		!(lww.vertexAdded[1].Equal(GetMockTime())) {
		t.Fatal()
	}
}

func TestVertexIsNotFoundWhenNotExists(t *testing.T) {
	lww := New()
	if lww.ContainsVertex(1) {
		t.Fatal()
	}
}

func TestVertexIsFoundAfterAdded(t *testing.T) {
	lww := New()
	lww.AddVertex(1)
	if !lww.ContainsVertex(1) {
		t.Fatal()
	}
}

func TestCannotRemoveIfNotExist(t *testing.T) {
	lww := New()
	if lww.RemoveVertex(1) == nil {
		t.Fatal()
	}
}

func TestRemoveExistingAddsToVertexRemoved(t *testing.T) {
	lww := New()
	lww.AddVertex(1)

	MockTime(time.Now())
	if lww.RemoveVertex(1) != nil {
		t.Fail()
	}
	if !lww.vertexRemoved[1].Equal(GetMockTime()) {
		t.Fatal()
	}
}

func TestContainsFalseAfterRemove(t *testing.T) {
	lww := New()
	lww.AddVertex(1)

	MockTime(time.Now())
	if lww.RemoveVertex(1) != nil {
		t.Fail()
	}
	if lww.ContainsVertex(1) {
		t.Fatal()
	}
}

func TestCantAddEdgeIfVerticesDontExist(t *testing.T) {
	lww := New()
	if lww.AddEdge(1, 2) == nil {
		t.Fatal()
	}
}

func TestAddEdgeSuccess(t *testing.T) {
	lww := New()
	lww.AddVertex(1)
	lww.AddVertex(2)
	MockTime(time.Now())
	if lww.AddEdge(1, 2) != nil {
		t.Fail()
	}
	if !lww.edgeAdded[Edge{V1: 1, V2: 2}].Equal(GetMockTime()) {
		t.Fatal()
	}
}

func TestEdgeContainsFailsWhenVerticesExistEdgeNotExists(t *testing.T) {
	lww := New()
	lww.AddVertex(1)
	lww.AddVertex(2)
	if lww.ContainsEdge(Edge{V1: 1, V2: 2}) {
		t.Fatal()
	}
}

func TestEdgeContainsFailsWhenNotExists(t *testing.T) {
	lww := New()
	if lww.ContainsEdge(Edge{V1: 1, V2: 2}) {
		t.Fatal()
	}
}

func TestContainsTrueAfterAdd(t *testing.T) {
	lww := New()
	lww.AddVertex(1)
	lww.AddVertex(2)
	MockTime(time.Now())
	if lww.AddEdge(1, 2) != nil {
		t.Fail()
	}
	if !lww.ContainsEdge(Edge{V1: 1, V2: 2}) {
		t.Fatal()
	}
}

func TestCantRemoveEdgeIfNotExist(t *testing.T) {
	lww := New()
	if lww.RemoveEdge(Edge{V1: 1, V2: 2}) == nil {
		t.Fatal()
	}
}

func TestRemoveExistingEdgeAddsToEdgeRemoved(t *testing.T) {
	lww := New()
	lww.AddVertex(1)
	lww.AddVertex(2)
	MockTime(time.Now())
	if lww.AddEdge(1, 2) != nil {
		t.Fail()
	}
	MockTime(time.Now())
	if lww.RemoveEdge(Edge{V1: 1, V2: 2}) != nil {
		t.Fail()
	}
	if !lww.edgeRemoved[Edge{V1: 1, V2: 2}].Equal(GetMockTime()) {
		t.Fatal()
	}
}

func TestContainsEdgeFailsAfterRemove(t *testing.T) {
	lww := New()
	lww.AddVertex(1)
	lww.AddVertex(2)
	MockTime(time.Now())
	if lww.AddEdge(1, 2) != nil {
		t.Fail()
	}
	MockTime(time.Now())
	if lww.RemoveEdge(Edge{V1: 1, V2: 2}) != nil {
		t.Fail()
	}
	if lww.ContainsEdge(Edge{V1: 1, V2: 2}) {
		t.Fail()
	}
}

func TestCantRemoveVertexIfInEdge(t *testing.T) {
	lww := New()
	lww.AddVertex(1)
	lww.AddVertex(2)
	MockTime(time.Now())
	if lww.AddEdge(1, 2) != nil {
		t.Fail()
	}
	if lww.RemoveVertex(1) == nil {
		t.Fail()
	}
	if _, ok := lww.vertexRemoved[1]; ok {
		t.Fatal()
	}
}

func TestConnectedVertices(t *testing.T) {
	lww := New()
	lww.AddVertex(1)
	lww.AddVertex(2)
	lww.AddVertex(3)
	lww.AddEdge(1, 2)
	lww.AddEdge(2, 3)
	con := lww.ConnectedVertices(1)
	if len(con) != 1 || con[0] != 2 {
		t.Fatal()
	}
	con = lww.ConnectedVertices(2)
	has1 := false
	has3 := false
	for _, v := range con {
		if v == 1 {
			has1 = true
		} else if v == 3 {
			has3 = true
		}
	}
	if len(con) != 2 || !has1 || !has3 {
		t.Fatal()
	}
}

func TestVertexNotInEdge(t *testing.T) {
	lww := New()
	lww.AddVertex(1)
	if lww.VertexInEdge(1) {
		t.Fatal()
	}
}

func TestVertexInEdge(t *testing.T) {
	lww := New()
	lww.AddVertex(1)
	lww.AddVertex(2)
	lww.AddEdge(1, 2)
	if !lww.VertexInEdge(1) {
		t.Fatal()
	}
}

func TestEmptyPath(t *testing.T) {
	lww := New()
	lww.AddVertex(1)
	lww.AddVertex(2)
	if len(lww.Path(1, 2)) > 0 {
		t.Fatal()
	}
}

func TestSinglePath(t *testing.T) {
	lww := New()
	lww.AddVertex(1)
	path := lww.Path(1, 1)
	if len(path) != 1 || path[0] != 1 {
		t.Fatal()
	}
}

func TestPathUnknownVertices(t *testing.T) {
	lww := New()
	path := lww.Path(7, 8)
	if len(path) != 0 {
		t.Fatal()
	}
}

func TestPath(t *testing.T) {
	lww := New()
	lww.AddVertex(1)
	lww.AddVertex(2)
	lww.AddVertex(3)
	lww.AddVertex(4)
	lww.AddEdge(1, 2)
	lww.AddEdge(1, 4)
	lww.AddEdge(4, 3)
	path := lww.Path(1, 3)
	if len(path) != 3 ||
		path[0] != 1 ||
		path[1] != 4 ||
		path[2] != 3 {
		t.Fatal()
	}
}

func TestMerge(t *testing.T) {
	lww := New()
	other := New()
	MockTime(time.Now())
	ts1 := GetMockTime()
	lww.AddVertex(1)
	lww.AddVertex(2)
	lww.AddVertex(3)
	lww.AddEdge(1, 2)
	lww.AddEdge(2, 3)
	other.AddVertex(4)
	other.AddVertex(5)
	other.AddEdge(4, 5)
	MockTime(time.Now())
	ts2 := GetMockTime()
	lww.RemoveEdge(Edge{V1: 2, V2: 3})
	MockTime(time.Now())
	ts3 := GetMockTime()
	lww.RemoveVertex(3)

	lww.Merge(other)

	if lww.vertexAdded[1] != ts1 ||
		lww.vertexAdded[2] != ts1 ||
		lww.vertexAdded[3] != ts1 ||
		lww.vertexAdded[4] != ts1 ||
		lww.vertexAdded[5] != ts1 ||
		len(lww.vertexRemoved) != 1 ||
		lww.vertexRemoved[3] != ts3 ||
		lww.edgeAdded[Edge{V1: 1, V2: 2}] != ts1 ||
		lww.edgeAdded[Edge{V1: 2, V2: 3}] != ts1 ||
		lww.edgeAdded[Edge{V1: 4, V2: 5}] != ts1 ||
		lww.edgeRemoved[Edge{V1: 2, V2: 3}] != ts2 {
		t.Fatal()
	}
}
