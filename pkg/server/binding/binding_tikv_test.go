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

var tikvTestDriver Driver

func init() {
	config := object.DefaultExecuteConfig()
	config.TikvAddrs = "127.0.0.1:2379"
	tikvTestDriver = initTikvDriver(config)
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
	funcCtx := sibyl2.FunctionContext{
		FunctionWithPath: &sibyl2.FunctionWithPath{
			Function: father,
			Path:     "a/b/c.go",
		},
		Calls: []*sibyl2.FunctionWithPath{
			{
				Function: &extractor.Function{
					Name: "abcde",
					Lang: core.LangGo,
				},
				Path: "b/c/d.go",
			},
		},
		ReverseCalls: []*sibyl2.FunctionWithPath{},
	}

	err = tikvTestDriver.CreateFuncContext(wc, &funcCtx, ctx)
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
}
