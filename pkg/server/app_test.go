package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"github.com/williamfzc/sibyl2/pkg/server/binding"
)

var wc = &binding.WorkspaceConfig{
	RepoId:  "sibyl",
	RevHash: "12345f",
}

func uploadFunctions(wc *binding.WorkspaceConfig, f []*extractor.FunctionFileResult) {
	core.Log.Infof("uploading %v with files %d ...", wc, len(f))
	var wg sync.WaitGroup
	wg.Add(len(f))
	for _, each := range f {
		unit := &FunctionUploadUnit{
			WorkspaceConfig: wc,
			FunctionResult:  each,
		}

		go func() {
			defer wg.Done()
			jsonStr, err := json.Marshal(unit)
			if err != nil {
				panic(err)
			}
			resp, err := http.Post(
				"http://127.0.0.1:9876/api/v1/func",
				"application/json",
				bytes.NewBuffer(jsonStr))
			if err != nil {
				panic(err)
			}
			data, err := io.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusOK {
				core.Log.Errorf("upload %s resp: %v", unit.FunctionResult.Path, string(data))
			}
		}()
	}
	wg.Wait()
}

func TestFuncUpload(t *testing.T) {
	t.Skip("always skip in CI")
	functions, _ := sibyl2.ExtractFunction("../..", sibyl2.DefaultConfig())
	uploadFunctions(wc, functions)
}
