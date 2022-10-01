package core

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
	parentUnit *Unit
	subUnits   []*Unit
}

func FindFirstByKindInParent(unit *Unit, kind KindRepr) *Unit {
	if unit == nil {
		return nil
	}
	if unit.Kind == kind {
		return unit
	}
	return FindFirstByKindInParent(unit.parentUnit, kind)
}

func FindFirstByKindInSubsWithDfs(unit *Unit, kind KindRepr) *Unit {
	if unit == nil {
		return nil
	}
	if unit.Kind == kind {
		return unit
	}

	// dfs
	for _, each := range unit.subUnits {
		eachResult := FindFirstByKindInSubsWithDfs(each, kind)
		if eachResult != nil {
			return eachResult
		}
	}
	return nil
}

func FindFirstByKindInSubsWithBfs(unit *Unit, kind KindRepr) *Unit {
	if unit == nil {
		return nil
	}
	if unit.Kind == kind {
		return unit
	}

	// bfs
	queue := unit.subUnits
	for len(queue) > 0 {
		var newQueue []*Unit
		for _, each := range queue {
			if each.Kind == kind {
				return each
			}
			newQueue = append(newQueue, each.subUnits...)
		}
		queue = newQueue
	}
	return nil
}

type FileUnit struct {
	Path     string   `json:"path"`
	Language LangType `json:"language"`
	Units    []*Unit  `json:"units"`
}
