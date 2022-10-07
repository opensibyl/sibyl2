package core

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
	Symbol    string `json:"symbol"`
	Kind      string `json:"kind"`
	Span      Span   `json:"span"`
	FieldName string `json:"fieldName"`
}

type ValueUnit struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Function struct {
	Name       string       `json:"name"`
	Receiver   string       `json:"receiver"`
	Parameters []*ValueUnit `json:"parameters"`
	Returns    []*ValueUnit `json:"returns"`
	Span       Span         `json:"span"`
}

type DataType interface {
	Dt()
}

func (*Symbol) Dt() {
}

func (*Function) Dt() {
}
