package core

import (
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

func TestParser_Parse_Java(t *testing.T) {
	t.Parallel()
	parser := NewParser(LangJava)
	_, err := parser.Parse([]byte(javaCode))
	if err != nil {
		panic(err)
	}
}

func TestParser_Parse_Golang(t *testing.T) {
	t.Parallel()
	parser := NewParser(LangGo)
	_, err := parser.Parse([]byte(goCode))
	if err != nil {
		panic(err)
	}
}
