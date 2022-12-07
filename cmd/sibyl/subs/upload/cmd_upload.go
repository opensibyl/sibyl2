package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/spf13/cobra"
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"github.com/williamfzc/sibyl2/pkg/server/object"
)

var uploadSrc string
var uploadLangType string
var uploadUrl string
var uploadWithCtx bool
var uploadBatchLimit int
var uploadDryRun bool

var httpClient = retryablehttp.NewClient()

func init() {
	httpClient.RetryMax = 3
	httpClient.RetryWaitMin = 500 * time.Millisecond
	httpClient.RetryWaitMax = 10 * time.Second
	httpClient.Logger = nil
}

func NewUploadCmd() *cobra.Command {
	uploadCmd := &cobra.Command{
		Use:    "upload",
		Short:  "test",
		Long:   `test`,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {
			uploadSrc, err := filepath.Abs(uploadSrc)
			if err != nil {
				panic(err)
			}
			repo, err := loadRepo(uploadSrc)
			if err != nil {
				panic(err)
			}
			head, err := repo.Head()
			curRepo := filepath.Base(uploadSrc)
			curRev := head.Hash().String()

			wc := &object.WorkspaceConfig{
				RepoId:  curRepo,
				RevHash: curRev,
			}
			f, err := sibyl2.ExtractFunction(uploadSrc, sibyl2.DefaultConfig())
			if err != nil {
				panic(err)
			}

			s, err := sibyl2.ExtractSymbol(uploadSrc, sibyl2.DefaultConfig())
			if err != nil {
				panic(err)
			}

			fullUrl := fmt.Sprintf("%s/api/v1/func", uploadUrl)
			ctxUrl := fmt.Sprintf("%s/api/v1/funcctx", uploadUrl)
			core.Log.Infof("upload backend: %s", fullUrl)
			if !uploadDryRun {
				uploadFunctions(fullUrl, wc, f)
			}
			core.Log.Infof("upload functions finished")

			// building edges in neo4j can be very slow
			// by default disabled
			if uploadWithCtx {
				core.Log.Infof("start calculating func graph")
				g, err := sibyl2.AnalyzeFuncGraph(f, s)
				if err != nil {
					panic(err)
				}
				core.Log.Infof("graph ready")
				if !uploadDryRun {
					uploadGraph(ctxUrl, wc, f, g)
				}
				core.Log.Infof("upload graph finished")
			}

			core.Log.Infof("upload finished")
		},
	}
	uploadCmd.PersistentFlags().StringVar(&uploadSrc, "src", ".", "src dir path")
	uploadCmd.PersistentFlags().StringVar(&uploadLangType, "lang", "", "lang type of your source code")
	uploadCmd.PersistentFlags().StringVar(&uploadUrl, "url", "http://127.0.0.1:9876", "backend url")
	uploadCmd.PersistentFlags().BoolVar(&uploadWithCtx, "withCtx", false, "with func context")
	uploadCmd.PersistentFlags().IntVar(&uploadBatchLimit, "batch", 50, "with func context")
	uploadCmd.PersistentFlags().BoolVar(&uploadDryRun, "dry", false, "dry run without upload")

	return uploadCmd
}

func loadRepo(gitDir string) (*git.Repository, error) {
	repo, err := git.PlainOpen(gitDir)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func uploadFunctions(url string, wc *object.WorkspaceConfig, f []*extractor.FunctionFileResult) {
	core.Log.Infof("uploading %v with files %d ...", wc, len(f))

	// pack
	var fullUnits []*object.FunctionUploadUnit
	for _, each := range f {
		unit := &object.FunctionUploadUnit{
			WorkspaceConfig: wc,
			FunctionResult:  each,
		}
		fullUnits = append(fullUnits, unit)
	}
	// submit
	ptr := 0
	batch := uploadBatchLimit
	for ptr < len(fullUnits) {
		core.Log.Infof("upload batch: %d - %d", ptr, ptr+batch)
		uploadFuncUnits(url, fullUnits[ptr:ptr+batch])
		ptr += batch
	}
}

func uploadFuncUnits(url string, units []*object.FunctionUploadUnit) {
	var wg sync.WaitGroup
	for _, unit := range units {
		if unit == nil {
			continue
		}
		wg.Add(1)
		go func(u *object.FunctionUploadUnit, waitGroup *sync.WaitGroup) {
			defer waitGroup.Done()

			jsonStr, err := json.Marshal(u)
			if err != nil {
				panic(err)
			}
			resp, err := httpClient.Post(
				url,
				"application/json",
				bytes.NewBuffer(jsonStr))
			if err != nil {
				panic(err)
			}
			data, err := io.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusOK {
				core.Log.Errorf("upload failed: %v", string(data))
			}
		}(unit, &wg)
	}
	wg.Wait()
}

func uploadGraph(url string, wc *object.WorkspaceConfig, functions []*extractor.FunctionFileResult, g *sibyl2.FuncGraph) {
	var wg sync.WaitGroup
	ptr := 0
	batch := uploadBatchLimit
	for ptr < len(functions) {
		core.Log.Infof("upload batch: %d - %d", ptr, ptr+batch)
		for _, eachFuncFile := range functions[ptr : ptr+batch] {
			if eachFuncFile == nil {
				continue
			}
			wg.Add(1)
			go func(funcFile *extractor.FunctionFileResult, waitGroup *sync.WaitGroup, graph *sibyl2.FuncGraph) {
				defer waitGroup.Done()

				var ctxs []*sibyl2.FunctionContext
				for _, eachFunc := range funcFile.Units {
					related := graph.FindRelated(eachFunc)
					ctxs = append(ctxs, related)
				}
				uploadFunctionContexts(url, wc, ctxs)
			}(eachFuncFile, &wg, g)
		}
		wg.Wait()
		ptr += batch
	}
}

func uploadFunctionContexts(url string, wc *object.WorkspaceConfig, ctxs []*sibyl2.FunctionContext) {
	uploadUnit := &object.FunctionContextUploadUnit{WorkspaceConfig: wc, FunctionContexts: ctxs}
	jsonStr, err := json.Marshal(uploadUnit)
	if err != nil {
		panic(err)
	}
	resp, err := httpClient.Post(
		url,
		"application/json",
		bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		core.Log.Errorf("upload resp: %v", string(data))
	}
}
