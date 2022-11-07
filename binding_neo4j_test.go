package sibyl2

import (
	"context"
	"sync"
	"testing"

	"github.com/williamfzc/sibyl2/pkg/core"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// set it true to run these tests
const hasNeo4jBackend = false
const dbUri = "bolt://localhost:7687"

var wc = &WorkspaceConfig{
	RepoId:  "sibyl",
	RevHash: "12345f",
}

// don't worry, fake password here :)
var authToken = neo4j.BasicAuth("neo4j", "williamfzc", "")

func TestNeo4jDriver_InitWorkspace(t *testing.T) {
	if !hasNeo4jBackend {
		t.Skip("always skip in CI")
	}

	driver, err := neo4j.NewDriverWithContext(dbUri, authToken)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	defer driver.Close(ctx)
	newDriver := &Neo4jDriver{driver}
	err = newDriver.InitWorkspace(wc, ctx)
	if err != nil {
		panic(err)
	}
}

func TestNeo4jDriver_UploadFile(t *testing.T) {
	if !hasNeo4jBackend {
		t.Skip("always skip in CI")
	}

	driver, err := neo4j.NewDriverWithContext(dbUri, authToken)
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
			err := newDriver.UploadFileResult(wc, a, ctx)
			if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
	core.Log.Infof("upload finished")
}

func TestNeo4jDriver_UploadFuncContextWithContext(t *testing.T) {
	if !hasNeo4jBackend {
		t.Skip("always skip in CI")
	}

	driver, err := neo4j.NewDriverWithContext(dbUri, authToken)
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
			err = newDriver.UploadFuncContext(wc, fc, ctx)
			if err != nil {
				panic(err)
			}
		}
	}
	core.Log.Infof("upload finished")
}

func TestNeo4jDriver_QueryFiles(t *testing.T) {
	if !hasNeo4jBackend {
		t.Skip("always skip in CI")
	}
	driver, err := neo4j.NewDriverWithContext(dbUri, authToken)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	defer driver.Close(ctx)
	newDriver := &Neo4jDriver{driver}
	files, err := newDriver.QueryFiles(wc, ctx)
	if err != nil {
		panic(err)
	}
	core.Log.Infof("files: %s", files)
}

func TestNeo4jDriver_QueryFunctions(t *testing.T) {
	if !hasNeo4jBackend {
		t.Skip("always skip in CI")
	}
	driver, err := neo4j.NewDriverWithContext(dbUri, authToken)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	defer driver.Close(ctx)
	newDriver := &Neo4jDriver{driver}
	files, err := newDriver.QueryFunctions(wc, "extract.go", ctx)
	if err != nil {
		panic(err)
	}
	for _, each := range files {
		core.Log.Infof("func: %v", each)
	}
}

func TestNeo4jDriver_QueryFunctionsWithLines(t *testing.T) {
	if !hasNeo4jBackend {
		t.Skip("always skip in CI")
	}
	driver, err := neo4j.NewDriverWithContext(dbUri, authToken)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	defer driver.Close(ctx)
	newDriver := &Neo4jDriver{driver}
	functions, err := newDriver.QueryFunctionsWithLines(wc, "extract.go", []int{32, 33}, ctx)
	if err != nil {
		panic(err)
	}
	for _, each := range functions {
		core.Log.Infof("func: %v", each)
	}
}

func TestNeo4jDriver_QueryFunctionWithSignature(t *testing.T) {
	if !hasNeo4jBackend {
		t.Skip("always skip in CI")
	}
	driver, err := neo4j.NewDriverWithContext(dbUri, authToken)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	defer driver.Close(ctx)
	newDriver := &Neo4jDriver{driver}
	ctxs, err := newDriver.QueryFunctionWithSignature(wc, "::ExtractFromString|string,*ExtractConfig|*extractor.FileResult,error", ctx)
	if err != nil {
		panic(err)
	}
	core.Log.Infof("ctx4: %v", ctxs.Name)
}

func TestNeo4jDriver_QueryFunctionContextWithSignature(t *testing.T) {
	if !hasNeo4jBackend {
		t.Skip("always skip in CI")
	}
	driver, err := neo4j.NewDriverWithContext(dbUri, authToken)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	defer driver.Close(ctx)
	newDriver := &Neo4jDriver{driver}
	ctxs, err := newDriver.QueryFunctionContextWithSignature(wc, "::ExtractFromString|string,*ExtractConfig|*extractor.FileResult,error", ctx)
	if err != nil {
		panic(err)
	}
	for _, each := range ctxs.ReverseCalls {
		core.Log.Infof("call: %v", each.GetIndexName())
	}
}

func TestNeo4jDriver_RemoveFileResult(t *testing.T) {
	if !hasNeo4jBackend {
		t.Skip("always skip in CI")
	}
	driver, err := neo4j.NewDriverWithContext(dbUri, authToken)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	defer driver.Close(ctx)
	newDriver := &Neo4jDriver{driver}
	err = newDriver.DeleteWorkspace(wc, ctx)
	if err != nil {
		panic(err)
	}
}

func TestNeo4jDriver_UpdateFuncProperties(t *testing.T) {
	if !hasNeo4jBackend {
		t.Skip("always skip in CI")
	}
	driver, err := neo4j.NewDriverWithContext(dbUri, authToken)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	defer driver.Close(ctx)
	newDriver := &Neo4jDriver{driver}
	err = newDriver.UpdateFuncProperties(wc, "::ExtractFromString|string,*ExtractConfig|*extractor.FileResult,error", "covered", 1, ctx)
	if err != nil {
		panic(err)
	}
}
