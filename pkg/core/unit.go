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
	ParentUnit *Unit
	SubUnits   []*Unit
}

func (unit *Unit) ReverseLink() []*Unit {
	// include itself
	cur := unit
	var out []*Unit
	for cur != nil {
		out = append(out, cur)
		cur = cur.ParentUnit
	}
	return out
}

func (unit *Unit) Link() []*Unit {
	var out []*Unit
	out = append(out, unit)
	out = append(out, unit.SubUnits...)
	return out
}

type FileUnit struct {
	Path     string   `json:"path"`
	Language LangType `json:"language"`
	Units    []*Unit  `json:"units"`
}
