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

func getBadgerTestConfig() object.ExecuteConfig {
	conf := object.DefaultExecuteConfig()
	conf.DbType = object.DriverTypeInMemory
	return conf
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
		},
	}

	err = d.CreateFuncFile(wc, &function, ctx)
	assert.Nil(t, err)
	functions, err := d.ReadFunctions(wc, function.Path, ctx)
	assert.Nil(t, err)
	if functions[0].Name != function.Units[0].Name {
		panic(nil)
	}

	// signatures
	signatures, err := d.ReadFunctionSignaturesWithRegex(wc, ".*", ctx)
	assert.Nil(t, err)
	assert.Equal(t, len(signatures), 1)
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
	assert.Nil(t, err)
	assert.Equal(t, 2, len(classes))
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

	err = d.CreateFuncContext(wc, &funcCtx, ctx)
	assert.Nil(t, err)
	newCtx, err := d.ReadFunctionContextWithSignature(wc, father.GetSignature(), ctx)
	assert.Nil(t, err)
	assert.Equal(t, newCtx.Function, father)
}
