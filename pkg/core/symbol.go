package core

type Span struct {
	Start Point `json:"start"`
	End   Point `json:"end"`
}

type Point struct {
	Row    uint32 `json:"row"`
	Column uint32 `json:"column"`
}

// Unit almost node
type Unit struct {
	Content   string `json:"content"`
	Kind      string `json:"kind"`
	Span      Span   `json:"span"`
	FieldName string `json:"fieldName"`
}

/*
Symbol
Units are named identifiers driven by the ASTs

https://github.com/github/semantic/blob/main/docs/examples.md#symbols
https://github.com/github/semantic/blob/main/proto/semantic.proto

	message Unit {
	  string symbol = 1;
	  string kind = 2;
	  Span span = 4;
	  NodeType node_type = 6;
	  SyntaxType syntax_type = 7;
	}

	enum NodeType {
	  ROOT_SCOPE = 0;
	  JUMP_TO_SCOPE = 1;
	  EXPORTED_SCOPE = 2;
	  DEFINITION = 3;
	  REFERENCE = 4;
	}

	enum SyntaxType {
	  FUNCTION = 0;
	  METHOD = 1;
	  CLASS = 2;
	  MODULE = 3;
	  CALL = 4;
	  TYPE = 5;
	  INTERFACE = 6;
	  IMPLEMENTATION = 7;
	}
*/
type Symbol struct {
	Symbol     string `json:"symbol"`
	Kind       string `json:"kind"`
	Span       Span   `json:"span"`
	FieldName  string `json:"fieldName"`
	NodeType   string `json:"nodeType"`
	SyntaxType string `json:"syntaxType"`
}

type FileUnit struct {
	Path     string   `json:"path"`
	Language LangType `json:"language"`
	Units    []Unit   `json:"units"`
}
