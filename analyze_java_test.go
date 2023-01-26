package sibyl2

import (
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	extractor2 "github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/extractor/java"
	"github.com/stretchr/testify/assert"
)

var javaCodeForAnalyze = `
package com.williamfzc.sibyl.core.listener.java8;

public class Java8SnapshotListener {
	public void aaaaaaaa(){
		bbbbbbbb();
	}

	private String bbbbbbbb(){
		return "";
	}
}
`

func TestAnalyzeJava(t *testing.T) {
	t.Parallel()
	parser := core.NewParser(core.LangJava)
	units, err := parser.Parse([]byte(javaCodeForAnalyze))
	if err != nil {
		panic(err)
	}

	extractor := &java.Extractor{}
	symbols, err := extractor.ExtractSymbols(units)
	functions, err := extractor.ExtractFunctions(units)
	symbolWrap := &extractor2.SymbolFileResult{}
	symbolWrap.Units = symbols
	functionWrap := &extractor2.FunctionFileResult{}
	functionWrap.Units = functions

	if err != nil {
		panic(err)
	}

	g, err := AnalyzeFuncGraph([]*extractor2.FunctionFileResult{functionWrap}, []*extractor2.SymbolFileResult{symbolWrap})
	if err != nil {
		panic(err)
	}

	ctx := g.FindRelated(functions[1])
	// we ignore too short method in java
	assert.Equal(t, ctx.Name, "bbbbbbbb")
	assert.Empty(t, ctx.Calls)
	assert.Len(t, ctx.ReverseCalls, 1)
}
