package binding

import (
	"context"
	"regexp"
	"testing"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/stretchr/testify/assert"
)

var curMongoDriver Driver

func init() {
	config := object.DefaultExecuteConfig()
	curMongoDriver = initMongoDriver(config)
	curMongoDriver.InitDriver(context.Background())
}

func TestMongoClazz(t *testing.T) {
	ctx := context.Background()
	err := curMongoDriver.InitDriver(ctx)
	if err != nil {
		panic(err)
	}
	defer curMongoDriver.DeferDriver()
	defer curMongoDriver.DeleteWorkspace(wc, ctx)

	err = curMongoDriver.CreateWorkspace(wc, ctx)
	if err != nil {
		panic(err)
	}

	clazz := extractor.BaseFileResult[*extractor.Clazz]{
		Path:     "abc/de/f.go",
		Language: core.LangGo,
		Type:     extractor.TypeExtractFunction,
		Units: []*extractor.Clazz{
			{
				Name:   "clazz0",
				Span:   core.Span{},
				Extras: nil,
			},
			{
				Name:   "clazz1",
				Span:   core.Span{},
				Extras: nil,
			},
		},
	}

	err = curMongoDriver.CreateClazzFile(wc, &clazz, ctx)
	assert.Nil(t, err)

	// check
	classes, err := curMongoDriver.ReadClasses(wc, clazz.Path, ctx)
	assert.Nil(t, err)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(classes))
}

func TestMongoFunc(t *testing.T) {
	ctx := context.Background()
	err := curMongoDriver.InitDriver(ctx)
	if err != nil {
		panic(err)
	}
	defer curMongoDriver.DeferDriver()
	defer curMongoDriver.DeleteWorkspace(wc, ctx)

	function := extractor.BaseFileResult[*extractor.Function]{
		Path:     "abc/de/f.go",
		Language: core.LangGo,
		Type:     extractor.TypeExtractFunction,
		Units: []*extractor.Function{
			{
				Name:       "fn",
				Receiver:   "fr",
				Parameters: nil,
				Returns:    nil,
				Span:       core.Span{},
				BodySpan:   core.Span{},
				Extras:     nil,
			},
			{
				Name:       "fn1",
				Receiver:   "fr",
				Parameters: nil,
				Returns:    nil,
				Span:       core.Span{},
				BodySpan:   core.Span{},
				Extras:     nil,
			},
		},
	}

	err = curMongoDriver.CreateFuncFile(wc, &function, ctx)
	assert.Nil(t, err)
	functions, err := curMongoDriver.ReadFunctions(wc, function.Path, ctx)
	assert.Nil(t, err)
	assert.Len(t, functions, 2)

	// signatures
	signatures, err := curMongoDriver.ReadFunctionSignaturesWithRegex(wc, "fn1.*", ctx)
	assert.Nil(t, err)
	assert.Equal(t, len(signatures), 1)

	// add tag and query the tag
	tag := "thisIsATag"
	err = curMongoDriver.CreateFuncTag(wc, signatures[0], tag, ctx)
	assert.Nil(t, err)
	funcs, err := curMongoDriver.ReadFunctionsWithTag(wc, tag, ctx)
	assert.Nil(t, err)
	assert.Equal(t, len(funcs), 1)

	// regex fn
	newRegex, err := regexp.Compile(".*sIs.*")
	verify := func(s string) bool {
		return newRegex.Match([]byte(s))
	}
	ruleMap := make(Rule)
	ruleMap["tags"] = verify

	funcsWithTag, err := curMongoDriver.ReadFunctionsWithRule(wc, ruleMap, ctx)
	assert.Nil(t, err)
	assert.Equal(t, len(funcsWithTag), 1)
}

func TestMongoFuncCtx(t *testing.T) {
	ctx := context.Background()
	err := curMongoDriver.InitDriver(ctx)
	if err != nil {
		panic(err)
	}
	defer curMongoDriver.DeferDriver()
	defer curMongoDriver.DeleteWorkspace(wc, ctx)

	father := &extractor.Function{
		Name: "abcde",
		Lang: core.LangGo,
	}
	calledFunc := &extractor.Function{
		Name: "calledfunc",
		Lang: core.LangGo,
	}
	p := "abc/def.go"
	called := &extractor.FunctionFileResult{
		Path:     p,
		Language: core.LangGo,
		Units:    []*extractor.Function{calledFunc, father},
	}

	funcCtx := sibyl2.FunctionContext{
		FunctionWithPath: &extractor.FunctionWithPath{
			Function: father,
			Path:     p,
		},
		Calls: []*extractor.FunctionWithPath{
			{
				Function: calledFunc,
				Path:     p,
			},
		},
		ReverseCalls: []*extractor.FunctionWithPath{},
	}
	slimCtx := object.CompressFunctionContext(&funcCtx)

	err = curMongoDriver.CreateFuncFile(wc, called, ctx)
	assert.Nil(t, err)
	err = curMongoDriver.CreateFuncContext(wc, slimCtx, ctx)
	assert.Nil(t, err)

	newCtx, err := curMongoDriver.ReadFunctionContextWithSignature(wc, father.GetSignature(), ctx)
	assert.Nil(t, err)
	assert.Equal(t, newCtx.Function.Name, father.Name)
}
