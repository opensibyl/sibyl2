package model

import (
	"golang.org/x/exp/slices"
	"sibyl2/pkg/core"
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

func FindFirstByKindInParent(unit *Unit, kind KindRepr) *Unit {
	if unit == nil {
		return nil
	}
	if unit.Kind == kind {
		return unit
	}
	return FindFirstByKindInParent(unit.ParentUnit, kind)
}

func FindFirstByOneOfKindInParent(unit *Unit, kinds ...KindRepr) *Unit {
	if unit == nil {
		return nil
	}
	if slices.Contains(kinds, unit.Kind) {
		return unit
	}
	return FindFirstByOneOfKindInParent(unit.ParentUnit, kinds...)
}

func FindFirstByKindInSubsWithDfs(unit *Unit, kind KindRepr) *Unit {
	if unit == nil {
		return nil
	}
	if unit.Kind == kind {
		return unit
	}

	// dfs
	for _, each := range unit.SubUnits {
		eachResult := FindFirstByKindInSubsWithDfs(each, kind)
		if eachResult != nil {
			return eachResult
		}
	}
	return nil
}

func FindFirstByFieldInSubsWithDfs(unit *Unit, fieldName string) *Unit {
	if unit == nil {
		return nil
	}
	if unit.FieldName == fieldName {
		return unit
	}

	// dfs
	for _, each := range unit.SubUnits {
		eachResult := FindFirstByFieldInSubsWithDfs(each, fieldName)
		if eachResult != nil {
			return eachResult
		}
	}
	return nil
}

func FindAllByKindInSubsWithDfs(unit *Unit, kind KindRepr) []*Unit {
	var ret []*Unit
	if unit == nil {
		return ret
	}
	if unit.Kind == kind {
		ret = append(ret, unit)
	}

	// dfs
	for _, each := range unit.SubUnits {
		eachResult := FindAllByKindInSubsWithDfs(each, kind)
		ret = append(ret, eachResult...)
	}
	return ret
}

func FindFirstByKindInSubsWithBfs(unit *Unit, kind KindRepr) *Unit {
	if unit == nil {
		return nil
	}
	if unit.Kind == kind {
		return unit
	}

	// bfs
	queue := unit.SubUnits
	for len(queue) > 0 {
		var newQueue []*Unit
		for _, each := range queue {
			if each.Kind == kind {
				return each
			}
			newQueue = append(newQueue, each.SubUnits...)
		}
		queue = newQueue
	}
	return nil
}

type FileUnit struct {
	Path     string        `json:"path"`
	Language core.LangType `json:"language"`
	Units    []*Unit       `json:"units"`
}

type FileResult struct {
	Path     string        `json:"path"`
	Language core.LangType `json:"language"`
	Type     string        `json:"type"`
	Units    []DataType    `json:"units"`
}
