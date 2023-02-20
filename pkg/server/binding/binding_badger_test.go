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

func getBadgerTestConfig() object.ExecuteConfig {
	conf := object.DefaultExecuteConfig()
	conf.DbType = object.DriverTypeInMemory
	return conf
}

var wc = &object.WorkspaceConfig{
	RepoId:  "sibyl",
	RevHash: "12345f",
}

func TestWc(t *testing.T) {
	d := initBadgerDriver(getBadgerTestConfig())
	ctx := context.TODO()
	err := d.InitDriver(ctx)
	if err != nil {
		panic(err)
	}

	defer d.DeferDriver()
	defer d.DeleteWorkspace(wc, ctx)
	err = d.CreateWorkspace(wc, ctx)
	assert.Nil(t, err)
	revs, err := d.ReadRevs(wc.RepoId, ctx)
	assert.Nil(t, err)
	if len(revs) != 1 {
		panic(nil)
	}
	for _, each := range revs {
		if each != wc.RevHash {
			panic(nil)
		}
	}

	info, err := d.ReadRevInfo(wc, ctx)
	assert.Nil(t, err)
	assert.Equal(t, info.Hash, wc.RevHash)
	assert.NotNil(t, info.CreateTime)
}

func TestBadgerFunc(t *testing.T) {
	d := initBadgerDriver(getBadgerTestConfig())
	ctx := context.TODO()
	err := d.InitDriver(ctx)
	if err != nil {
		panic(err)
	}

	defer d.DeferDriver()
	defer d.DeleteWorkspace(wc, ctx)
	err = d.CreateWorkspace(wc, ctx)
	if err != nil {
		panic(err)
	}

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

	err = d.CreateFuncFile(wc, &function, ctx)
	assert.Nil(t, err)
	functions, err := d.ReadFunctions(wc, function.Path, ctx)
	assert.Nil(t, err)
	assert.Len(t, functions, 2)

	// signatures
	signatures, err := d.ReadFunctionSignaturesWithRegex(wc, "fn1.*", ctx)
	assert.Nil(t, err)
	assert.Equal(t, len(signatures), 1)

	// rule
	rule := make(Rule)
	regex, err := regexp.Compile("fn1.*")
	assert.Nil(t, err)
	rule["name"] = func(s string) bool {
		return regex.Match([]byte(s))
	}
	funcs, err := d.ReadFunctionsWithRule(wc, rule, ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(funcs))

	// tag
	f := funcs[0]
	assert.Empty(t, f.Tags)
	newTag := "THIS_IS_A_TAG"
	err = d.CreateFuncTag(wc, f.GetSignature(), newTag, ctx)
	assert.Nil(t, err)
	newF, err := d.ReadFunctionWithSignature(wc, f.GetSignature(), ctx)
	assert.Nil(t, err)
	assert.Len(t, newF.Tags, 1)
	assert.Equal(t, newF.Tags[0], newTag)
}

func TestBadgerClazz(t *testing.T) {
	d := initBadgerDriver(getBadgerTestConfig())
	ctx := context.TODO()
	err := d.InitDriver(ctx)
	if err != nil {
		panic(err)
	}

	defer d.DeferDriver()
	defer d.DeleteWorkspace(wc, ctx)
	err = d.CreateWorkspace(wc, ctx)
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

	err = d.CreateClazzFile(wc, &clazz, ctx)
	assert.Nil(t, err)

	// check
	classes, err := d.ReadClasses(wc, clazz.Path, ctx)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(classes))

	rule := make(Rule)
	regex, err := regexp.Compile("clazz0")
	assert.Nil(t, err)
	rule["name"] = func(s string) bool {
		return regex.Match([]byte(s))
	}
	classes, err = d.ReadClassesWithRule(wc, rule, ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(classes))
}

func TestBadgerFuncCtx(t *testing.T) {
	d := initBadgerDriver(getBadgerTestConfig())
	ctx := context.TODO()
	err := d.InitDriver(ctx)
	if err != nil {
		panic(err)
	}

	defer d.DeferDriver()
	defer d.DeleteWorkspace(wc, ctx)
	err = d.CreateWorkspace(wc, ctx)
	if err != nil {
		panic(err)
	}

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
		FunctionWithPath: &sibyl2.FunctionWithPath{
			Function: father,
			Path:     p,
		},
		Calls: []*sibyl2.FunctionWithPath{
			{
				Function: calledFunc,
				Path:     p,
			},
		},
		ReverseCalls: []*sibyl2.FunctionWithPath{},
	}
	slimCtx := funcCtx.ToSlim()

	err = d.CreateFuncFile(wc, called, ctx)
	assert.Nil(t, err)
	err = d.CreateFuncContext(wc, slimCtx, ctx)
	assert.Nil(t, err)

	newCtx, err := d.ReadFunctionContextWithSignature(wc, father.GetSignature(), ctx)
	assert.Nil(t, err)
	assert.Equal(t, newCtx.Function.Name, father.Name)

	// check
	rule := make(Rule)
	regex, err := regexp.Compile("abc.*")
	assert.Nil(t, err)
	rule["name"] = func(s string) bool {
		return regex.Match([]byte(s))
	}

	f, err := d.ReadFunctionContextsWithRule(wc, rule, ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(f))

	// check its signature is valid
	fws, err := d.ReadFunctionWithSignature(wc, f[0].Calls[0], ctx)
	assert.Nil(t, err)
	assert.Equal(t, fws.Name, called.Units[0].Name)
}
