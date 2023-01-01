package java

import (
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
)

var javaCode = `
package com.williamfzc.sibyl.core.listener.java8;

import com.williamfzc.sibyl.core.listener.Java8Parser;
import com.williamfzc.sibyl.core.listener.java8.base.Java8MethodLayerListener;
import com.williamfzc.sibyl.core.model.method.Method;

@ClassAnnotationA(argA="yes")
@ClassAnnotationB
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

func TestJavaExtractor_ExtractSymbols(t *testing.T) {
	parser := core.NewParser(core.LangJava)
	units, err := parser.Parse([]byte(javaCode))
	if err != nil {
		panic(err)
	}

	extractor := &Extractor{}
	_, err = extractor.ExtractSymbols(units)
	if err != nil {
		panic(err)
	}
}

func TestJavaExtractor_ExtractFunctions(t *testing.T) {
	parser := core.NewParser(core.LangJava)
	units, err := parser.Parse([]byte(javaCode))
	if err != nil {
		panic(err)
	}

	extractor := &Extractor{}
	data, err := extractor.ExtractFunctions(units)
	if err != nil {
		panic(err)
	}
	for _, each := range data {
		core.Log.Debugf("each: %s %s %s", each.Name, each.Extras, each.BodySpan.String())
		// check base info
		if each.Name == "enterMethodDeclarationWithoutMethodBody" {
			if each.BodySpan.String() != "14:71,17:5" {
				panic(nil)
			}
		}

		if each.Extras.(*FunctionExtras).ClassInfo.Annotations == nil {
			panic(err)
		}
	}
}

func TestExtractor_ExtractClasses(t *testing.T) {
	parser := core.NewParser(core.LangJava)
	units, err := parser.Parse([]byte(javaCode))
	if err != nil {
		panic(err)
	}

	extractor := &Extractor{}
	data, err := extractor.ExtractClasses(units)
	for _, each := range data {
		core.Log.Infof("find class: %v", each.GetSignature())
	}
}
