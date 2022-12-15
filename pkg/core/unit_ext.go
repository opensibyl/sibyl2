package core

import (
	"golang.org/x/exp/slices"
)

func FindFirstByKindInParent(unit *Unit, kind KindRepr) *Unit {
	return NewQuery(unit).Bottom2Top().MatchKind(kind).First()
}

func FindAllByKindInParent(unit *Unit, kind KindRepr) []*Unit {
	return NewQuery(unit).Bottom2Top().MatchKind(kind).All()
}

func FindAllByOneOfKindInParent(unit *Unit, kinds ...KindRepr) []*Unit {
	query := NewQuery(unit).Bottom2Top()
	for _, each := range kinds {
		query.MatchKind(each)
	}
	return query.All()
}

func FindAll(unit *Unit) []*Unit {
	return NewQuery(unit).Top2Bottom().All()
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
	return NewQuery(unit).Top2Bottom().MatchKind(kind).Dfs().All()
}

func FindAllByKindInSubsWithBfs(unit *Unit, kind KindRepr) []*Unit {
	return NewQuery(unit).Top2Bottom().MatchKind(kind).Bfs().All()
}

func FindFirstByKindInSubsWithBfs(unit *Unit, kind KindRepr) *Unit {
	return NewQuery(unit).Bfs().Top2Bottom().MatchKind(kind).First()
}

func FindFirstByFieldInSubs(unit *Unit, fieldName string) *Unit {
	if unit == nil {
		return nil
	}

	for _, each := range unit.SubUnits {
		if each.FieldName == fieldName {
			return each
		}
	}
	return nil
}

func FindAllByKindInSubs(unit *Unit, kind string) []*Unit {
	var ret []*Unit
	if unit == nil {
		return ret
	}

	for _, each := range unit.SubUnits {
		if each.Kind == kind {
			ret = append(ret, each)
		}
	}
	return ret
}

func FindAllByKindsInSubs(unit *Unit, kinds ...string) []*Unit {
	var ret []*Unit
	if unit == nil {
		return ret
	}

	for _, each := range unit.SubUnits {
		if slices.Contains(kinds, each.Kind) {
			ret = append(ret, each)
		}
	}
	return ret
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
	}
	return q.bfsFirst(q.target)
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
	}
	return q.bfsAll(q.target)
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
