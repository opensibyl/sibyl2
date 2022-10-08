package model

import (
	"golang.org/x/exp/slices"
)

type Span struct {
	Start Point `json:"start"`
	End   Point `json:"end"`
}

type Point struct {
	Row    uint32 `json:"row"`
	Column uint32 `json:"column"`
}

type KindRepr = string

/*
Unit

almost a node, but with enough data for analyzer.
no need to access raw byte data again
*/
type Unit struct {
	Kind      KindRepr `json:"kind"`
	Span      Span     `json:"span"`
	FieldName string   `json:"fieldName"`
	Content   string   `json:"content"`

	// double linked
	ParentUnit *Unit
	SubUnits   []*Unit
}

type Query struct {
	target       *Unit
	IsDfs        bool
	IsTop2Bottom bool
	Kinds        []string
	FieldNames   []string
}

func NewQuery(target *Unit) *Query {
	return &Query{
		target:       target,
		IsDfs:        true,
		IsTop2Bottom: true,
	}
}

func (q *Query) Dfs() *Query {
	q.IsDfs = true
	return q
}

func (q *Query) Bfs() *Query {
	q.IsDfs = false
	return q
}

func (q *Query) MatchKind(kind string) *Query {
	q.Kinds = append(q.Kinds, kind)
	return q
}

func (q *Query) MatchField(fieldName string) *Query {
	q.FieldNames = append(q.FieldNames, fieldName)
	return q
}

func (q *Query) Top2Bottom() *Query {
	q.IsTop2Bottom = true
	return q
}

func (q *Query) Bottom2Top() *Query {
	q.IsTop2Bottom = false
	return q
}

func (q *Query) First() *Unit {
	if !q.IsTop2Bottom {
		return q.parent(q.target)
	}

	if q.IsDfs {
		return q.dfsFirst(q.target)
	} else {
		return q.bfsFirst(q.target)
	}
}

func (q *Query) match(unit *Unit) bool {
	// compare
	matchField := false
	matchKind := false
	if len(q.FieldNames) == 0 {
		// no need to check
		matchField = true
	} else if slices.Contains(q.FieldNames, unit.FieldName) {
		matchField = true
	}

	if len(q.Kinds) == 0 {
		// no need to check
		matchKind = true
	} else if slices.Contains(q.Kinds, unit.Kind) {
		matchKind = true
	}

	return matchField && matchKind
}

func (q *Query) parent(unit *Unit) *Unit {
	if unit == nil {
		return nil
	}
	// compare
	if q.match(unit) {
		return unit
	}
	return q.parent(unit.ParentUnit)
}

func (q *Query) dfsFirst(unit *Unit) *Unit {
	if unit == nil {
		return nil
	}

	// compare
	if q.match(unit) {
		return unit
	}

	// dfs
	for _, each := range unit.SubUnits {
		eachResult := q.dfsFirst(each)
		if eachResult != nil {
			return eachResult
		}
	}
	return nil
}

func (q *Query) bfsFirst(unit *Unit) *Unit {
	if unit == nil {
		return nil
	}
	// compare
	if q.match(unit) {
		return unit
	}

	// bfs
	queue := unit.SubUnits
	for len(queue) > 0 {
		var newQueue []*Unit
		for _, each := range queue {
			if q.match(each) {
				return each
			}
			newQueue = append(newQueue, each.SubUnits...)
		}
		queue = newQueue
	}
	return nil
}

func (q *Query) All() []*Unit {
	if !q.IsTop2Bottom {
		return q.parentAll(q.target)
	}

	if q.IsDfs {
		return q.dfsAll(q.target)
	} else {
		return q.bfsAll(q.target)
	}
}

func (q *Query) parentAll(unit *Unit) []*Unit {
	panic("TODO")
}

func (q *Query) dfsAll(unit *Unit) []*Unit {
	var ret []*Unit
	if unit == nil {
		return ret
	}
	if q.match(unit) {
		ret = append(ret, unit)
	}

	// dfs
	for _, each := range unit.SubUnits {
		eachResult := q.dfsAll(each)
		ret = append(ret, eachResult...)
	}
	return ret
}

func (q *Query) bfsAll(unit *Unit) []*Unit {
	panic("TODO")
}

func FindFirstByKindInParent(unit *Unit, kind KindRepr) *Unit {
	return NewQuery(unit).Bottom2Top().MatchKind(kind).First()
}

func FindFirstByOneOfKindInParent(unit *Unit, kinds ...KindRepr) *Unit {
	query := NewQuery(unit).Bottom2Top()
	for _, each := range kinds {
		query.MatchKind(each)
	}
	return query.First()
}

func FindFirstByKindInSubsWithDfs(unit *Unit, kind KindRepr) *Unit {
	return NewQuery(unit).Top2Bottom().MatchKind(kind).First()
}

func FindFirstByFieldInSubsWithDfs(unit *Unit, fieldName string) *Unit {
	return NewQuery(unit).Top2Bottom().MatchField(fieldName).First()
}

func FindFirstByFieldInSubsWithBfs(unit *Unit, fieldName string) *Unit {
	return NewQuery(unit).Top2Bottom().Bfs().MatchField(fieldName).First()
}

func FindAllByKindInSubsWithDfs(unit *Unit, kind KindRepr) []*Unit {
	return NewQuery(unit).Top2Bottom().MatchKind(kind).All()
}

func FindFirstByKindInSubsWithBfs(unit *Unit, kind KindRepr) *Unit {
	return NewQuery(unit).Bfs().Top2Bottom().MatchKind(kind).First()
}
