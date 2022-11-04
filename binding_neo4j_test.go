package sibyl2

import (
	"context"
	"sync"
	"testing"

	"github.com/williamfzc/sibyl2/pkg/core"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func TestNeo4jDriver_UploadFile(t *testing.T) {
	t.Skip("always skip in CI")
	wc := &WorkspaceConfig{
		RepoId:  "sibyl",
		RevHash: "12345f",
	}

	dbUri := "bolt://localhost:7687"
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", "williamfzc", ""))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	defer driver.Close(ctx)
	newDriver := &Neo4jDriver{driver}
	functions, _ := ExtractFunction(".", DefaultConfig())

	core.Log.Infof("start uploading")
	var wg sync.WaitGroup
	for _, each := range functions {
		wg.Add(1)
		a := each
		go func() {
			defer wg.Done()
			newDriver.UploadFileResultWithContext(wc, a, ctx)
		}()
	}
	wg.Wait()
	core.Log.Infof("upload finished")
}

func TestNeo4jDriver_UploadFuncContextWithContext(t *testing.T) {
	t.Skip("always skip in CI")
	wc := &WorkspaceConfig{
		RepoId:  "sibyl",
		RevHash: "12345f",
	}

	dbUri := "bolt://localhost:7687"
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", "williamfzc", ""))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	defer driver.Close(ctx)
	newDriver := &Neo4jDriver{driver}
	functions, _ := ExtractFunction(".", DefaultConfig())
	symbols, _ := ExtractSymbol(".", DefaultConfig())
	fg, _ := AnalyzeFuncGraph(functions, symbols)
	core.Log.Infof("target query done")
	for _, eachFunc := range functions {
		for _, eachFFF := range eachFunc.Units {
			fc := fg.FindRelated(eachFFF)
			err = newDriver.UploadFuncContextWithContext(wc, fc, ctx)
			if err != nil {
				panic(err)
			}
		}
	}
	core.Log.Infof("upload finished")
}
