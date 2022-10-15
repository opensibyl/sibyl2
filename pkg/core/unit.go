package core

type Span struct {
	Start Point `json:"start"`
	End   Point `json:"end"`
}

type Point struct {
	Row    uint32 `json:"row"`
	Column uint32 `json:"column"`
}

func (s *Span) ContainLine(lineNum int) bool {
	uint32Line := uint32(lineNum)
	return s.Start.Row <= uint32Line && uint32Line <= s.End.Row
}

func (s *Span) ContainAnyLine(lineNums []int) bool {
	for _, each := range lineNums {
		if s.ContainLine(each) {
			return true
		}
	}
	return false
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

type FileUnit struct {
	Path     string   `json:"path"`
	Language LangType `json:"language"`
	Units    []*Unit  `json:"units"`
}
