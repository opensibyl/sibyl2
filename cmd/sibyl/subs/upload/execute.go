package upload

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	object2 "github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

var httpClient = retryablehttp.NewClient()

func init() {
	httpClient.RetryMax = 3
	httpClient.RetryWaitMin = 500 * time.Millisecond
	httpClient.RetryWaitMax = 10 * time.Second
	httpClient.Logger = nil
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func execWithConfig(c *uploadConfig) {
	startTime := time.Now()
	defer func() {
		core.Log.Infof("upload total cost: %d ms", time.Since(startTime).Milliseconds())
	}()

	configStr, err := c.ToJson()
	panicIfErr(err)
	core.Log.Infof("upload with config: %s", configStr)
	uploadSrc, err := filepath.Abs(c.Src)
	panicIfErr(err)
	repo, err := loadRepo(uploadSrc)
	panicIfErr(err)
	head, err := repo.Head()
	panicIfErr(err)

	var curRepo string
	if c.RepoId == "" {
		curRepo = filepath.Base(uploadSrc)
	} else {
		curRepo = c.RepoId
	}

	cIter, err := repo.Log(&git.LogOptions{From: head.Hash()})
	panicIfErr(err)
	var commits = make([]*object2.Commit, 0, c.Depth)
	count := 0
	_ = cIter.ForEach(func(commit *object2.Commit) error {
		commits = append(commits, commit)
		count += 1
		if count >= c.Depth {
			return errors.New("break")
		}
		return nil
	})

	tree, err := repo.Worktree()
	panicIfErr(err)

	for _, eachRev := range commits {
		if eachRev.Hash != head.Hash() {
			core.Log.Infof("checkout: %s", eachRev.Hash)
			err = tree.Checkout(&git.CheckoutOptions{
				Hash: eachRev.Hash,
				Keep: true,
			})
			panicIfErr(err)
		}

		wc := &object.WorkspaceConfig{
			RepoId:  curRepo,
			RevHash: eachRev.Hash.String(),
		}
		execCurRevWithConfig(uploadSrc, wc, c)
	}

	// recover
	if c.Depth != 1 {
		core.Log.Infof("recover checkout: %s", head)
		err = tree.Checkout(&git.CheckoutOptions{
			Hash: head.Hash(),
			Keep: true,
		})
		panicIfErr(err)
	}

	core.Log.Infof("upload finished")
}

func execCurRevWithConfig(uploadSrc string, wc *object.WorkspaceConfig, c *uploadConfig) {
	filterFunc, err := createFileFilter(c)
	panicIfErr(err)

	runner := &core.Runner{}
	var lang []string
	if len(c.Lang) == 0 {
		langFromDir, err := runner.GuessLangFromDir(c.Src, filterFunc)
		panicIfErr(err)
		lang = []string{string(langFromDir)}
	} else {
		lang = c.Lang
	}
	if len(lang) == 0 {
		panic(errors.New("no valid lang found"))
	}

	for _, eachLang := range lang {
		eachLangType := core.LangTypeValueOf(eachLang)
		if !eachLangType.IsSupported() {
			core.Log.Warnf("lang %v not supported, supported list: %v", eachLangType, core.SupportedLangs)
			continue
		}
		core.Log.Infof("scan lang: %v", eachLang)
		execCurRevCurLangWithConfig(uploadSrc, core.LangType(eachLang), filterFunc, wc, c)
	}
}

func execCurRevCurLangWithConfig(uploadSrc string, lang core.LangType, filterFunc func(path string) bool, wc *object.WorkspaceConfig, c *uploadConfig) {
	f, err := sibyl2.ExtractFunction(uploadSrc, &sibyl2.ExtractConfig{
		FileFilter: filterFunc,
		LangType:   lang,
	})
	panicIfErr(err)

	funcUrl := fmt.Sprintf("%s/api/v1/func", c.Url)
	funcCtxUrl := fmt.Sprintf("%s/api/v1/funcctx", c.Url)
	clazzUrl := fmt.Sprintf("%s/api/v1/clazz", c.Url)

	core.Log.Infof("upload backend: %s", funcUrl)
	if !c.Dry {
		uploadFunctions(funcUrl, wc, f, c.Batch)
	}
	core.Log.Infof("upload functions finished, file count: %d", len(f))

	// building edges can be expensive
	// by default disabled
	if c.WithCtx {
		s, err := sibyl2.ExtractSymbol(uploadSrc, &sibyl2.ExtractConfig{
			FileFilter: filterFunc,
			LangType:   lang,
		})
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
			uploadFunctionContexts(funcCtxUrl, wc, f, g, c.Batch)
		}
		core.Log.Infof("upload graph finished")
	}

	if c.WithClass {
		s, err := sibyl2.ExtractClazz(uploadSrc, &sibyl2.ExtractConfig{
			FileFilter: filterFunc,
			LangType:   lang,
		})
		if err != nil {
			panic(err)
		}
		core.Log.Infof("classes ready")
		if !c.Dry {
			uploadClazz(clazzUrl, wc, s, c.Batch)
		}
	}
}

func createFileFilter(c *uploadConfig) (func(path string) bool, error) {
	if c.IncludeRegex == "" && c.ExcludeRegex == "" {
		// need no filter
		return nil, nil
	}

	var include *regexp.Regexp
	var exclude *regexp.Regexp
	var err error
	if c.IncludeRegex != "" {
		include, err = regexp.Compile(c.IncludeRegex)
		if err != nil {
			core.Log.Errorf("failed to compile: %v", c.IncludeRegex)
			return nil, err
		}
	}
	if c.ExcludeRegex != "" {
		exclude, err = regexp.Compile(c.ExcludeRegex)
		if err != nil {
			core.Log.Errorf("failed to compile: %v", c.IncludeRegex)
			return nil, err
		}
	}

	core.Log.Infof("create file filter, include: %s, exclude: %s", c.IncludeRegex, c.ExcludeRegex)
	return func(path string) bool {
		var shouldInclude bool
		var shouldExclude bool
		if include == nil {
			shouldInclude = true
		} else {
			shouldInclude = include.MatchString(path)
		}

		if exclude == nil {
			shouldExclude = false
		} else {
			shouldExclude = exclude.MatchString(path)
		}

		return shouldInclude && !shouldExclude
	}, nil
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
			panicIfErr(err)
			data, err := io.ReadAll(resp.Body)
			panicIfErr(err)
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
	jsonStr, err := json.Marshal(uploadUnit)
	panicIfErr(err)
	resp, err := httpClient.Post(
		url,
		"application/json",
		bytes.NewBuffer(jsonStr))
	panicIfErr(err)
	data, err := io.ReadAll(resp.Body)
	panicIfErr(err)
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
		go func(u *object.ClazzUploadUnit, waitGroup *sync.WaitGroup) {
			defer waitGroup.Done()

			jsonStr, err := json.Marshal(u)
			if err != nil {
				panic(err)
			}
			resp, err := httpClient.Post(
				url,
				"application/json",
				bytes.NewBuffer(jsonStr))
			panicIfErr(err)
			data, err := io.ReadAll(resp.Body)
			panicIfErr(err)
			if resp.StatusCode != http.StatusOK {
				core.Log.Errorf("upload failed: %v", string(data))
			}
		}(unit, &wg)
	}
	wg.Wait()
}
