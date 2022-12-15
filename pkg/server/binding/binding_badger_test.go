package binding

import (
	"context"
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
)

func TestWc(t *testing.T) {
	d := newBadgerDriver()
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
	revs, err := d.ReadRevs(wc.RepoId, ctx)
	if err != nil {
		panic(err)
	}
	if len(revs) != 1 {
		panic(nil)
	}
	for _, each := range revs {
		if each != wc.RevHash {
			panic(nil)
		}
	}
}

func TestReadWrite(t *testing.T) {
	d := newBadgerDriver()
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
	if err != nil {
		panic(err)
	}
	functions, err := d.ReadFunctions(wc, function.Path, ctx)
	if err != nil {
		return
	}
	if functions[0].Name != function.Units[0].Name {
		panic(nil)
	}
}
