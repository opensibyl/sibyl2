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
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var httpClient = retryablehttp.NewClient()

func init() {
	httpClient.RetryMax = 3
	httpClient.RetryWaitMin = 500 * time.Millisecond
	httpClient.RetryWaitMax = 10 * time.Second
	httpClient.Logger = nil
}

func NewUploadCmd() *cobra.Command {
	var uploadSrc string
	var uploadLangType string
	var uploadUrl string
	var uploadWithCtx bool
	var uploadBatchLimit int
	var uploadDryRun bool

	uploadCmd := &cobra.Command{
		Use:    "upload",
		Short:  "test",
		Long:   `test`,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {
			config := defaultConfig()

			// read from config
			viper.AddConfigPath(configPath)
			viper.SetConfigFile(configFile)

			core.Log.Infof("trying to read config from: %s/%s", configPath, configFile)
			err := viper.ReadInConfig()
			if err != nil {
				core.Log.Warnf("no config file found, use default")
			} else {
				core.Log.Infof("found config file")
				err = viper.Unmarshal(&config)

				if err != nil {
					core.Log.Errorf("failed to parse config")
					panic(err)
				}
			}

			// read from cmd
			config.Src = uploadSrc
			config.Lang = uploadLangType
			config.Url = uploadUrl
			config.WithCtx = uploadWithCtx
			config.Batch = uploadBatchLimit
			config.Dry = uploadDryRun

			// execute
			execWithConfig(config)

			// save it back
			usedConfigMap, err := config.ToMap()
			if err != nil {
				panic(err)
			}
			err = viper.MergeConfigMap(usedConfigMap)
			if err != nil {
				panic(err)
			}
			err = viper.WriteConfigAs(viper.ConfigFileUsed())
			if err != nil {
				core.Log.Warnf("failed to write config back")
			}
		},
	}

	config := defaultConfig()
	uploadCmd.PersistentFlags().StringVar(&uploadSrc, "src", config.Src, "src dir path")
	uploadCmd.PersistentFlags().StringVar(&uploadLangType, "lang", config.Lang, "lang type of your source code")
	uploadCmd.PersistentFlags().StringVar(&uploadUrl, "url", config.Url, "backend url")
	uploadCmd.PersistentFlags().BoolVar(&uploadWithCtx, "withCtx", config.WithCtx, "with func context")
	uploadCmd.PersistentFlags().IntVar(&uploadBatchLimit, "batch", config.Batch, "with func context")
	uploadCmd.PersistentFlags().BoolVar(&uploadDryRun, "dry", config.Dry, "dry run without upload")

	return uploadCmd
}

func execWithConfig(c *uploadConfig) {
	uploadSrc, err := filepath.Abs(c.Src)
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

	fullUrl := fmt.Sprintf("%s/api/v1/func", c.Url)
	ctxUrl := fmt.Sprintf("%s/api/v1/funcctx", c.Url)
	core.Log.Infof("upload backend: %s", fullUrl)
	if !c.Dry {
		uploadFunctions(fullUrl, wc, f, c.Batch)
	}
	core.Log.Infof("upload functions finished, file count: %d", len(f))

	// building edges can be expensive
	// by default disabled
	if c.WithCtx {
		s, err := sibyl2.ExtractSymbol(uploadSrc, sibyl2.DefaultConfig())
		if err != nil {
			panic(err)
		}

		core.Log.Infof("start calculating func graph")
		g, err := sibyl2.AnalyzeFuncGraph(f, s)
		if err != nil {
			panic(err)
		}
		core.Log.Infof("graph ready")
		if !c.Dry {
			uploadGraph(ctxUrl, wc, f, g, c.Batch)
		}
		core.Log.Infof("upload graph finished")
	}

	core.Log.Infof("upload finished")
}

func loadRepo(gitDir string) (*git.Repository, error) {
	repo, err := git.PlainOpen(gitDir)
	if err != nil {
		return nil, err
	}
	return repo, nil
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
			uploadFuncUnits(url, fullUnits[ptr:ptr+batch])
		} else {
			uploadFuncUnits(url, fullUnits[ptr:])
		}

		ptr = newPtr
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

func uploadGraph(url string, wc *object.WorkspaceConfig, functions []*extractor.FunctionFileResult, g *sibyl2.FuncGraph, batch int) {
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

				var ctxs []*sibyl2.FunctionContext
				for _, eachFunc := range funcFile.Units {
					related := graph.FindRelated(eachFunc)
					ctxs = append(ctxs, related)
				}
				uploadFunctionContexts(url, wc, ctxs)
			}(eachFuncFile, &wg, g)
		}
		wg.Wait()
		ptr = newPtr
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
