package upload

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/vmihailenco/msgpack/v5"
)

var httpClient = retryablehttp.NewClient()

func init() {
	httpClient.RetryMax = 3
	httpClient.RetryWaitMin = 500 * time.Millisecond
	httpClient.RetryWaitMax = 10 * time.Second
	httpClient.Logger = nil
}

func msgpack2bytes(o interface{}) ([]byte, error) {
	var output bytes.Buffer
	enc := msgpack.NewEncoder(&output)
	enc.SetCustomStructTag("json")
	err := enc.Encode(o)
	if err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func uploadFunctions(url string, wc *object.WorkspaceConfig, f []*extractor.FunctionFileResult, batch int) {
	core.Log.Infof("uploading %v with files %d ...", wc, len(f))

	// pack
	fullUnits := make([]*object.FunctionUploadUnit, 0, len(f))
	for _, each := range f {
		unit := &object.FunctionUploadUnit{
			WorkspaceConfig: wc,
			FunctionResult:  each,
		}
		fullUnits = append(fullUnits, unit)
	}
	// submit
	ptr := 0
	for ptr < len(fullUnits) {
		core.Log.Infof("upload batch: %d - %d", ptr, ptr+batch)

		newPtr := ptr + batch
		if newPtr < len(fullUnits) {
			uploadFunctionUnits(url, fullUnits[ptr:ptr+batch])
		} else {
			uploadFunctionUnits(url, fullUnits[ptr:])
		}

		ptr = newPtr
	}
}

func uploadFunctionUnits(url string, units []*object.FunctionUploadUnit) {
	var wg sync.WaitGroup
	for _, unit := range units {
		if unit == nil {
			continue
		}
		wg.Add(1)
		go func(u *object.FunctionUploadUnit, waitGroup *sync.WaitGroup) {
			defer waitGroup.Done()

			uploadData, err := msgpack2bytes(u)
			if err != nil {
				core.Log.Errorf("error when upload: %v", err)
				return
			}
			resp, err := httpClient.Post(
				url,
				object.BodyTypeMsgpack,
				uploadData)
			if err != nil {
				core.Log.Errorf("error when upload: %v", err)
				return
			}
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				core.Log.Errorf("error when upload: %v", err)
				return
			}
			if resp.StatusCode != http.StatusOK {
				core.Log.Errorf("upload failed: %v", string(data))
			}
		}(unit, &wg)
	}
	wg.Wait()
}

func uploadFunctionContexts(url string, wc *object.WorkspaceConfig, functions []*extractor.FunctionFileResult, g *sibyl2.FuncGraph, batch int) {
	var wg sync.WaitGroup
	ptr := 0
	for ptr < len(functions) {
		core.Log.Infof("upload batch: %d - %d", ptr, ptr+batch)

		newPtr := ptr + batch
		var todoFuncs []*extractor.FunctionFileResult
		if newPtr < len(functions) {
			todoFuncs = functions[ptr:newPtr]
		} else {
			todoFuncs = functions[ptr:]
		}

		for _, eachFuncFile := range todoFuncs {
			if eachFuncFile == nil {
				continue
			}
			wg.Add(1)
			go func(funcFile *extractor.FunctionFileResult, waitGroup *sync.WaitGroup, graph *sibyl2.FuncGraph) {
				defer waitGroup.Done()

				contexts := make([]*sibyl2.FunctionContext, 0)
				for _, eachFunc := range funcFile.Units {
					eachFileWithPath := sibyl2.WrapFuncWithPath(eachFunc, funcFile.Path)
					related := graph.FindRelated(eachFileWithPath)
					contexts = append(contexts, related)
				}
				uploadFunctionContextUnits(url, wc, contexts)
			}(eachFuncFile, &wg, g)
		}
		wg.Wait()
		ptr = newPtr
	}
}

func uploadFunctionContextUnits(url string, wc *object.WorkspaceConfig, ctxs []*sibyl2.FunctionContext) {
	uploadUnit := &object.FunctionContextUploadUnit{WorkspaceConfig: wc, FunctionContexts: ctxs}
	uploadData, err := msgpack2bytes(uploadUnit)
	if err != nil {
		core.Log.Errorf("error when upload: %v", err)
		return
	}
	resp, err := httpClient.Post(
		url,
		object.BodyTypeMsgpack,
		uploadData)
	if err != nil {
		core.Log.Errorf("error when upload: %v", err)
		return
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		core.Log.Errorf("error when upload: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		core.Log.Errorf("upload resp: %v", string(data))
	}
}

func uploadClazz(url string, wc *object.WorkspaceConfig, classes []*extractor.ClazzFileResult, batch int) {
	core.Log.Infof("uploading %v with files %d ...", wc, len(classes))

	// pack
	fullUnits := make([]*object.ClazzUploadUnit, 0, len(classes))
	for _, each := range classes {
		unit := &object.ClazzUploadUnit{
			WorkspaceConfig: wc,
			ClazzFileResult: each,
		}
		fullUnits = append(fullUnits, unit)
	}
	// submit
	ptr := 0
	for ptr < len(fullUnits) {
		core.Log.Infof("upload batch: %d - %d", ptr, ptr+batch)

		newPtr := ptr + batch
		if newPtr < len(fullUnits) {
			uploadClazzUnits(url, fullUnits[ptr:ptr+batch])
		} else {
			uploadClazzUnits(url, fullUnits[ptr:])
		}

		ptr = newPtr
	}
}

func uploadClazzUnits(url string, units []*object.ClazzUploadUnit) {
	var wg sync.WaitGroup
	for _, unit := range units {
		if unit == nil {
			continue
		}
		wg.Add(1)
		go func(uploadUnit *object.ClazzUploadUnit, waitGroup *sync.WaitGroup) {
			defer waitGroup.Done()
			uploadData, err := msgpack2bytes(uploadUnit)
			if err != nil {
				core.Log.Errorf("error when upload: %v", err)
				return
			}
			resp, err := httpClient.Post(
				url,
				object.BodyTypeMsgpack,
				uploadData)
			if err != nil {
				core.Log.Errorf("error when upload: %v", err)
				return
			}
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				core.Log.Errorf("error when upload: %v", err)
				return
			}
			if resp.StatusCode != http.StatusOK {
				core.Log.Errorf("upload failed: %v", string(data))
			}
		}(unit, &wg)
	}
	wg.Wait()
}
