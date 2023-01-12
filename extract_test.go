package sibyl2

import (
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
)

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

func BenchmarkExtractFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// no   cache: 499267 ns/op
		// with cache: 14275 ns/op
		_, err := ExtractFromString(javaCode, &ExtractConfig{
			LangType:    core.LangJava,
			ExtractType: extractor.TypeExtractFunction,
		})
		if err != nil {
			panic(err)
		}
	}
}
