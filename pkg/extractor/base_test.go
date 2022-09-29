package extractor

import (
	"sibyl2/pkg/core"
	"testing"
)

var javaCode = `
package com.williamfzc.sibyl.core.listener.java8;

import com.williamfzc.sibyl.core.listener.Java8Parser;
import com.williamfzc.sibyl.core.listener.java8.base.Java8MethodLayerListener;
import com.williamfzc.sibyl.core.model.method.Method;

public class Java8SnapshotListener extends Java8MethodLayerListener<Method> {
    @Override
    public void enterMethodDeclarationWithoutMethodBody(
            Java8Parser.MethodDeclarationWithoutMethodBodyContext ctx) {
        super.enterMethodDeclarationWithoutMethodBody(ctx);
        this.storage.save(curMethodStack.peekLast());
    }

    @Override
    public void enterInterfaceMethodDeclaration(Java8Parser.InterfaceMethodDeclarationContext ctx) {
        super.enterInterfaceMethodDeclaration(ctx);
        this.storage.save(curMethodStack.peekLast());
    }
}
`

var goCode = `
type Parser struct {
	engine *sitter.Parser
}

func NewParser(lang *sitter.Language) *Parser {
	engine := sitter.NewParser()
	engine.SetLanguage(lang)
	return &Parser{
		engine,
	}
}
`

func TestJavaExtractor_ExtractSymbols(t *testing.T) {
	parser := core.NewParser(core.JAVA)
	units, err := parser.Parse([]byte(javaCode))
	if err != nil {
		panic(err)
	}

	extractor := GetExtractor(core.JAVA)
	_, err = extractor.ExtractSymbols(units)
	if err != nil {
		panic(err)
	}
}

func TestJavaExtractor_ExtractFunctions(t *testing.T) {
	parser := core.NewParser(core.JAVA)
	units, err := parser.Parse([]byte(javaCode))
	if err != nil {
		panic(err)
	}

	extractor := GetExtractor(core.JAVA)
	_, err = extractor.ExtractFunctions(units)
	if err != nil {
		panic(err)
	}
}

func TestGolangExtractor_ExtractFunctions(t *testing.T) {
	parser := core.NewParser(core.GOLANG)
	units, err := parser.Parse([]byte(goCode))
	if err != nil {
		panic(err)
	}

	extractor := GetExtractor(core.GOLANG)
	_, err = extractor.ExtractFunctions(units)
	if err != nil {
		panic(err)
	}
}
