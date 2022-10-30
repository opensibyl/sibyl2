package core

type Span struct {
	Start Point `json:"start"`
	End   Point `json:"end"`
}

type Point struct {
	// NOTICE: INDEX, NOT REAL LINE NUMBER
	Row    uint32 `json:"row"`
	Column uint32 `json:"column"`
}

func (s *Span) Lines() []int {
	var ret = make([]int, s.End.Row-s.Start.Row+1)
	for i := s.Start.Row; i <= s.End.Row; i++ {
		ret = append(ret, int(i))
	}
	return ret
}

func (s *Span) ContainLine(lineNum int) bool {
	// real line number
	uint32Line := uint32(lineNum) + 1
	return s.Start.Row <= uint32Line && uint32Line <= s.End.Row
}

func (s *Span) ContainAnyLine(lineNums ...int) bool {
	for _, each := range lineNums {
		if s.ContainLine(each) {
			return true
		}
	}
	return false
}

func (s *Span) MatchAny(target *Span) bool {
	return !(target.End.Row < s.Start.Row || target.Start.Row > s.End.Row)
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
