package binding

import (
	"context"
	"sync"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/core"
)

// set it true to run these tests
const hasNeo4jBackend = true
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
	newDriver, _ := NewNeo4jDriver(driver)
	err = newDriver.CreateWorkspace(wc, ctx)
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

	newDriver, _ := NewNeo4jDriver(driver)
	functions, _ := sibyl2.ExtractFunction(".", sibyl2.DefaultConfig())

	core.Log.Infof("start uploading")
	err = newDriver.CreateWorkspace(wc, ctx)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, each := range functions {
		wg.Add(1)
		a := each
		go func() {
			defer wg.Done()
			err := newDriver.CreateFuncFile(wc, a, ctx)
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
	newDriver, _ := NewNeo4jDriver(driver)
	functions, _ := sibyl2.ExtractFunction(".", sibyl2.DefaultConfig())
	symbols, _ := sibyl2.ExtractSymbol(".", sibyl2.DefaultConfig())
	fg, _ := sibyl2.AnalyzeFuncGraph(functions, symbols)
	core.Log.Infof("target query done")
	for _, eachFunc := range functions {
		for _, eachFFF := range eachFunc.Units {
			fc := fg.FindRelated(eachFFF)
			err = newDriver.CreateFuncContext(wc, fc, ctx)
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
	newDriver, _ := NewNeo4jDriver(driver)
	files, err := newDriver.ReadFiles(wc, ctx)
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
	newDriver, _ := NewNeo4jDriver(driver)
	files, err := newDriver.ReadFunctions(wc, "extract.go", ctx)
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
	newDriver, _ := NewNeo4jDriver(driver)
	functions, err := newDriver.ReadFunctionsWithLines(wc, "extract.go", []int{32, 33}, ctx)
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
	newDriver, _ := NewNeo4jDriver(driver)
	ctxs, err := newDriver.ReadFunctionWithSignature(wc, "::ExtractFromString|string,*ExtractConfig|*extractor.FileResult,error", ctx)
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
	newDriver, _ := NewNeo4jDriver(driver)
	ctxs, err := newDriver.ReadFunctionContextWithSignature(wc, "::ExtractFromString|string,*ExtractConfig|*extractor.FileResult,error", ctx)
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
	newDriver, _ := NewNeo4jDriver(driver)
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
	newDriver, _ := NewNeo4jDriver(driver)
	err = newDriver.UpdateFuncProperties(wc, "::ExtractFromString|string,*ExtractConfig|*extractor.FileResult,error", "covered", 1, ctx)
	if err != nil {
		panic(err)
	}

	err = newDriver.UpdateRevProperties(wc, "revK", "revV", ctx)
	if err != nil {
		panic(err)
	}

	err = newDriver.UpdateFileProperties(wc, "extract.go", "fileK", "fileV", ctx)
	if err != nil {
		panic(err)
	}
}
