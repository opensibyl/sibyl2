package java

import (
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/stretchr/testify/assert"
)

var javaCode = `
package com.williamfzc.sibyl.core.listener.java8;

import com.williamfzc.sibyl.core.listener.Java8Parser;
import com.williamfzc.sibyl.core.listener.java8.base.Java8MethodLayerListener;
import com.williamfzc.sibyl.core.model.method.Method;

@ClassAnnotationA(argA="yes")
@ClassAnnotationB
public class Java8SnapshotListener extends Java8MethodLayerListener<Method> {
	private static final int ABC = 1;

	@InjectMocks
	private static final int ABCD = 1;

	private final String DBCA;

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
			assert.Equal(t, each.BodySpan.String(), "21:71,24:5")
		}

		assert.NotNil(t, each.Extras.(*FunctionExtras).ClassInfo.Annotations)
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
	assert.Nil(t, err)
	for _, each := range data {
		core.Log.Infof("find class: %v", each.GetSignature())
		for _, field := range each.Extras.(*ClassExtras).Fields {
			core.Log.Infof("field: %v", field)
		}
	}
}
