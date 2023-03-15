package binding

import (
	"context"
	"testing"

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
