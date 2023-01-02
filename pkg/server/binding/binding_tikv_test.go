package binding

import (
	"context"
	"testing"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/stretchr/testify/assert"
)

var tikvDriver Driver

func init() {
	config := object.DefaultExecuteConfig()
	config.TikvAddrs = "127.0.0.1:2379"
	tikvDriver = initTikvDriver(config)
}

func TestTikvFunc(t *testing.T) {
	ctx := context.Background()
	err := tikvDriver.InitDriver(ctx)
	if err != nil {
		panic(err)
	}
	defer tikvDriver.DeferDriver()
	defer tikvDriver.DeleteWorkspace(wc, ctx)

	err = tikvDriver.CreateWorkspace(wc, ctx)
	if err != nil {
		panic(err)
	}

	repos, err := tikvDriver.ReadRepos(ctx)
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

	err = tikvDriver.CreateFuncFile(wc, &function, ctx)
	assert.Nil(t, err)

	// check
	files, err := tikvDriver.ReadFiles(wc, ctx)
	assert.Equal(t, 1, len(files))

	// functions
	functions, err := tikvDriver.ReadFunctions(wc, function.Path, ctx)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(functions))
	assert.Equal(t, functions[0].Name, function.Units[0].Name)
}

func TestTikvClazz(t *testing.T) {
	ctx := context.Background()
	err := tikvDriver.InitDriver(ctx)
	if err != nil {
		panic(err)
	}
	defer tikvDriver.DeferDriver()
	defer tikvDriver.DeleteWorkspace(wc, ctx)

	err = tikvDriver.CreateWorkspace(wc, ctx)
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

	err = tikvDriver.CreateClazzFile(wc, &clazz, ctx)
	assert.Nil(t, err)

	// check
	classes, err := tikvDriver.ReadClasses(wc, clazz.Path, ctx)
	assert.Nil(t, err)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(classes))
}

func TestTikvFuncCtx(t *testing.T) {
	ctx := context.Background()
	err := tikvDriver.InitDriver(ctx)
	if err != nil {
		panic(err)
	}
	defer tikvDriver.DeferDriver()
	defer tikvDriver.DeleteWorkspace(wc, ctx)

	err = tikvDriver.CreateWorkspace(wc, ctx)
	if err != nil {
		panic(err)
	}

	father := &extractor.Function{
		Name: "abcde",
	}
	funcCtx := sibyl2.FunctionContext{
		FunctionWithPath: &sibyl2.FunctionWithPath{
			Function: father,
			Path:     "a/b/c.go",
			Language: core.LangGo,
		},
		Calls: []*sibyl2.FunctionWithPath{
			{
				Function: &extractor.Function{
					Name: "abcde",
				},
				Path:     "b/c/d.go",
				Language: core.LangGo,
			},
		},
		ReverseCalls: []*sibyl2.FunctionWithPath{},
	}

	err = tikvDriver.CreateFuncContext(wc, &funcCtx, ctx)
	assert.Nil(t, err)
	newCtx, err := tikvDriver.ReadFunctionContextWithSignature(wc, father.GetSignature(), ctx)
	assert.Nil(t, err)
	assert.Equal(t, newCtx.Function, father)
}
