package extractor

import (
	"sibyl2/pkg/core"
	"sibyl2/pkg/model"
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

func TestJavaExtractor_ExtractSymbols(t *testing.T) {
	parser := core.NewParser(model.LangJava)
	units, err := parser.Parse([]byte(javaCode))
	if err != nil {
		panic(err)
	}

	extractor := GetExtractor(model.LangJava)
	_, err = extractor.ExtractSymbols(units)
	if err != nil {
		panic(err)
	}
}

func TestJavaExtractor_ExtractFunctions(t *testing.T) {
	parser := core.NewParser(model.LangJava)
	units, err := parser.Parse([]byte(javaCode))
	if err != nil {
		panic(err)
	}

	extractor := GetExtractor(model.LangJava)
	_, err = extractor.ExtractFunctions(units)
	if err != nil {
		panic(err)
	}
}
