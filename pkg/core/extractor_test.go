package core

import (
	"fmt"
	"github.com/smacker/go-tree-sitter/java"
	"testing"
)

func TestJavaExtractor_IsSymbol(t *testing.T) {
	parser := NewParser(java.GetLanguage())
	units, err := parser.Parse([]byte(javaCode))
	if err != nil {
		panic(err)
	}

	extractor := GetExtractor(JAVA)
	symbols := extractor.ExtractSymbols(units)
	fmt.Printf("%+v\n", symbols)
}
