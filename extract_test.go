package sibyl2

import (
	"testing"

	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
)

func TestExtract(t *testing.T) {
	fileResult, err := Extract(".", &ExtractConfig{
		LangType:    core.LangGo,
		ExtractType: extractor.TypeExtractFunction,
	})
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		core.Log.Debugf("path: %v, %v", each.Path, each.Units)
	}
}

func TestExtractWithGuess(t *testing.T) {
	fileResult, err := Extract(".", &ExtractConfig{
		ExtractType: extractor.TypeExtractFunction,
	})
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		core.Log.Debugf("path: %v, %v", each.Path, each.Units)
	}
}

func TestExtractFunctionWithGuess(t *testing.T) {
	fileResult, err := ExtractFunction(".", &ExtractConfig{})
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		core.Log.Debugf("path: %v, %v", each.Path, each.Units)
	}
}

func TestExtractSymbol(t *testing.T) {
	fileResult, err := ExtractSymbol(".", &ExtractConfig{})
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		core.Log.Debugf("path: %s, %v", each.Path, each.Units)
	}
}

func TestExtractFile(t *testing.T) {
	fileResult, err := Extract("./extract.go", &ExtractConfig{
		LangType:    core.LangGo,
		ExtractType: extractor.TypeExtractFunction,
	})
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		core.Log.Debugf("path: %v, %v", each.Path, each.Units)
	}
}

var javaCode = `
package com.williamfzc.sibyl.core.listener.java8;

import com.williamfzc.sibyl.core.listener.Java8Parser;
import com.williamfzc.sibyl.core.listener.java8.base.Java8MethodLayerListener;
import com.williamfzc.sibyl.core.model.method.Method;

public class Java8SnapshotListener extends Java8MethodLayerListener<Method> {
    @Override
	@abcde
	@adeflkjbg(abc = "dfff")
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

func TestExtractString(t *testing.T) {
	fileResult, err := ExtractFromString(javaCode, &ExtractConfig{
		LangType:    core.LangJava,
		ExtractType: extractor.TypeExtractFunction,
	})
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult.Units {
		core.Log.Debugf("result: %s", each.GetDesc())
	}
}
