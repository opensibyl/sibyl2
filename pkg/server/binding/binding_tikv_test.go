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

// ./.tiup/bin/tiup playground --mode tikv-slim --host 0.0.0.0
var tikvTestDriver Driver

func init() {
	config := object.DefaultExecuteConfig()
	config.TikvAddrs = "127.0.0.1:2379"
	tikvTestDriver = initTikvDriver(config)
}

func TestTikvWc(t *testing.T) {
	ctx := context.Background()
	err := tikvTestDriver.InitDriver(ctx)
	if err != nil {
		panic(err)
	}
	defer tikvTestDriver.DeferDriver()
	defer tikvTestDriver.DeleteWorkspace(wc, ctx)
	err = tikvTestDriver.CreateWorkspace(wc, ctx)
	assert.Nil(t, err)
	revs, err := tikvTestDriver.ReadRevs(wc.RepoId, ctx)
	assert.Nil(t, err)
	if len(revs) != 1 {
		panic(nil)
	}
	for _, each := range revs {
		if each != wc.RevHash {
			panic(nil)
		}
	}

	info, err := tikvTestDriver.ReadRevInfo(wc, ctx)
	assert.Nil(t, err)
	assert.Equal(t, info.Hash, wc.RevHash)
	assert.NotNil(t, info.CreateTime)
}

func TestTikvFunc(t *testing.T) {
	ctx := context.Background()
	err := tikvTestDriver.InitDriver(ctx)
	if err != nil {
		panic(err)
	}
	defer tikvTestDriver.DeferDriver()
	defer tikvTestDriver.DeleteWorkspace(wc, ctx)

	err = tikvTestDriver.CreateWorkspace(wc, ctx)
	if err != nil {
		panic(err)
	}

	repos, err := tikvTestDriver.ReadRepos(ctx)
	if err != nil {
		panic(err)
	}
	core.Log.Debugf("repos: %v", repos)

	function := extractor.BaseFileResult[*extractor.Function]{
		Path:     "abc/de/f.go",
		Language: core.LangGo,
		Type:     extractor.TypeExtractFunction,
		Units: []*extractor.Function{
			{
				Name:       "fn0",
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

	err = tikvTestDriver.CreateFuncFile(wc, &function, ctx)
	assert.Nil(t, err)

	// check
	files, err := tikvTestDriver.ReadFiles(wc, ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(files))

	// functions
	functions, err := tikvTestDriver.ReadFunctions(wc, function.Path, ctx)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(functions))
	assert.Equal(t, functions[0].Name, function.Units[0].Name)

	// signatures
	signatures, err := tikvTestDriver.ReadFunctionSignaturesWithRegex(wc, "fn1.*", ctx)
	assert.Nil(t, err)
	assert.Equal(t, len(signatures), 1)

	// rule
	rule := make(Rule)
	regex, err := regexp.Compile("fn1.*")
	assert.Nil(t, err)
	rule["name"] = func(s string) bool {
		return regex.Match([]byte(s))
	}
	funcs, err := tikvTestDriver.ReadFunctionsWithRule(wc, rule, ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(funcs))
}

func TestTikvClazz(t *testing.T) {
	ctx := context.Background()
	err := tikvTestDriver.InitDriver(ctx)
	if err != nil {
		panic(err)
	}
	defer tikvTestDriver.DeferDriver()
	defer tikvTestDriver.DeleteWorkspace(wc, ctx)

	err = tikvTestDriver.CreateWorkspace(wc, ctx)
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

	err = tikvTestDriver.CreateClazzFile(wc, &clazz, ctx)
	assert.Nil(t, err)

	// check
	classes, err := tikvTestDriver.ReadClasses(wc, clazz.Path, ctx)
	assert.Nil(t, err)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(classes))

	// rule
	rule := make(Rule)
	regex, err := regexp.Compile("clazz0")
	assert.Nil(t, err)
	rule["name"] = func(s string) bool {
		return regex.Match([]byte(s))
	}
	classes, err = tikvTestDriver.ReadClassesWithRule(wc, rule, ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(classes))
}

func TestTikvFuncCtx(t *testing.T) {
	ctx := context.Background()
	err := tikvTestDriver.InitDriver(ctx)
	if err != nil {
		panic(err)
	}
	defer tikvTestDriver.DeferDriver()
	defer tikvTestDriver.DeleteWorkspace(wc, ctx)

	err = tikvTestDriver.CreateWorkspace(wc, ctx)
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

	err = tikvTestDriver.CreateFuncContext(wc, slimCtx, ctx)
	assert.Nil(t, err)
	newCtx, err := tikvTestDriver.ReadFunctionContextWithSignature(wc, father.GetSignature(), ctx)
	assert.Nil(t, err)
	assert.Equal(t, newCtx.Function, father)

	// rule
	rule := make(Rule)
	regex, err := regexp.Compile("abc.*")
	assert.Nil(t, err)
	rule["name"] = func(s string) bool {
		return regex.Match([]byte(s))
	}
	funcs, err := tikvTestDriver.ReadFunctionContextsWithRule(wc, rule, ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(funcs))

	// check its signature is valid
	fws, err := tikvTestDriver.ReadFunctionWithSignature(wc, funcs[0].Calls[0], ctx)
	assert.Nil(t, err)
	assert.Equal(t, fws.Name, called.Units[0].Name)
}
