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
Symbol
https://github.com/github/semantic/blob/main/docs/examples.md#symbols
https://github.com/github/semantic/blob/main/proto/semantic.proto

	enum NodeType {
		ROOT_SCOPE = 0;
		JUMP_TO_SCOPE = 1;
		EXPORTED_SCOPE = 2;
		DEFINITION = 3;
		REFERENCE = 4;
	}
*/
type Symbol struct {
	// value
	Symbol string `json:"symbol"`

	// range
	Span Span `json:"span"`

	// type (lang specific
	// higher analyser will use this field
	NodeType string `json:"nodeType"`
}

type FileSymbol struct {
	Path     string   `json:"path"`
	Language LangType `json:"language"`
	Symbols  []Symbol `json:"symbols"`
}
