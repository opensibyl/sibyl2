package binding

import (
	"context"
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

var tikvDriver Driver

func init() {
	config := object.DefaultExecuteConfig()
	config.TikvAddrs = "127.0.0.1:2379"
	tikvDriver = initTikvDriver(config)
}

func TestTikv(t *testing.T) {
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
	if err != nil {
		panic(err)
	}

	// check
	files, err := tikvDriver.ReadFiles(wc, ctx)
	if len(files) != 1 {
		panic(nil)
	}

	functions, err := tikvDriver.ReadFunctions(wc, function.Path, ctx)
	if err != nil {
		return
	}

	if len(functions) != 2 {
		panic(nil)
	}

	if functions[0].Name != function.Units[0].Name {
		panic(nil)
	}
}
