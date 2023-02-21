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

class D extends B {}

class B implements A, C {
	void abcd() {
	}
}

interface A {
	void abcd();
}

interface C {}
`

func TestJavaExtractor_ExtractSymbols(t *testing.T) {
	t.Parallel()
	parser := core.NewParser(core.LangJava)
	units, err := parser.Parse([]byte(javaCode))
	if err != nil {
		panic(err)
	}

	extractor := &Extractor{}
	symbols, err := extractor.ExtractSymbols(units)
	assert.Nil(t, err)
	assert.NotEmpty(t, symbols)
}

func TestJavaExtractor_ExtractFunctions(t *testing.T) {
	t.Parallel()
	parser := core.NewParser(core.LangJava)
	units, err := parser.Parse([]byte(javaCode))
	if err != nil {
		panic(err)
	}

	extractor := &Extractor{}
	data, err := extractor.ExtractFunctions(units)
	assert.Nil(t, err)
	assert.NotEmpty(t, data)
	for _, each := range data {
		core.Log.Debugf("each: %s %s %s", each.Name, each.Extras, each.BodySpan.String())
		// check base info
		if each.Name == "enterMethodDeclarationWithoutMethodBody" {
			assert.Equal(t, each.BodySpan.String(), "21:71,24:5")
			assert.NotNil(t, each.Extras.(*FunctionExtras).ClassInfo.Annotations)
			assert.Equal(t, each.Namespace, "com.williamfzc.sibyl.core.listener.java8")
		}
	}
}

func TestExtractor_ExtractClasses(t *testing.T) {
	t.Parallel()
	parser := core.NewParser(core.LangJava)
	units, err := parser.Parse([]byte(javaCode))
	if err != nil {
		panic(err)
	}

	extractor := &Extractor{}
	data, err := extractor.ExtractClasses(units)
	assert.Nil(t, err)
	for _, each := range data {
		core.Log.Debugf("find class: %v", each.GetSignature())
		for _, field := range each.Extras.(*ClassExtras).Fields {
			core.Log.Infof("field: %v", field)
		}
		core.Log.Debugf("class extends: %v", each.Extras.(*ClassExtras).Extends)
		core.Log.Debugf("class implements: %v", each.Extras.(*ClassExtras).Implements)
	}
}
