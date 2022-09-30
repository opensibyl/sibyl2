package core

type Span struct {
	Start Point `json:"start"`
	End   Point `json:"end"`
}

type Point struct {
	Row    uint32 `json:"row"`
	Column uint32 `json:"column"`
}

/*
Unit

almost a node, but with enough data for analyzer.
no need to access raw byte data again
*/
type Unit struct {
	Kind      string `json:"kind"`
	Span      Span   `json:"span"`
	FieldName string `json:"fieldName"`
	Content   string `json:"content"`

	Parent *Unit
}

type FileUnit struct {
	Path     string   `json:"path"`
	Language LangType `json:"language"`
	Units    []Unit   `json:"units"`
}
