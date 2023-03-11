package upload

import (
	"errors"
	"path/filepath"
	"regexp"
	"time"

	"github.com/go-git/go-git/v5"
	object2 "github.com/go-git/go-git/v5/plumbing/object"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

type ExecuteCache struct {
	AnalyzeGraph *sibyl2.FuncGraph
}
type ExecuteCacheMap = map[core.LangType]*ExecuteCache

func ExecWithConfig(c *Config) error {
	startTime := time.Now()
	defer func() {
		core.Log.Infof("upload total cost: %d ms", time.Since(startTime).Milliseconds())
	}()

	configStr, err := c.ToJson()
	if err != nil {
		return err
	}
	core.Log.Infof("upload with config: %s", configStr)
	uploadSrc, err := filepath.Abs(c.Src)
	if err != nil {
		return err
	}

	// if repo id and rev hash has been set, do not access git.
	// https://github.com/opensibyl/sibyl2/issues/44
	if c.RepoId != "" && c.RevHash != "" {
		wc := &object.WorkspaceConfig{
			RepoId:  c.RepoId,
			RevHash: c.RevHash,
		}
		_, err := ExecCurRevWithConfig(uploadSrc, wc, c)
		if err != nil {
			return err
		}
	} else {
		err := execWithGit(uploadSrc, c)
		if err != nil {
			return err
		}
	}
	core.Log.Infof("upload finished")
	return nil
}

func execWithGit(uploadSrc string, c *Config) error {
	// extract from git
	repo, err := loadRepo(uploadSrc)
	if err != nil {
		return err
	}
	head, err := repo.Head()
	if err != nil {
		return err
	}

	var curRepo string
	if c.RepoId == "" {
		curRepo = filepath.Base(uploadSrc)
	} else {
		curRepo = c.RepoId
	}

	cIter, err := repo.Log(&git.LogOptions{From: head.Hash()})
	if err != nil {
		return err
	}
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
	if err != nil {
		return err
	}

	for _, eachRev := range commits {
		if eachRev.Hash != head.Hash() {
			core.Log.Infof("checkout: %s", eachRev.Hash)
			err = tree.Checkout(&git.CheckoutOptions{
				Hash: eachRev.Hash,
				Keep: true,
			})
			if err != nil {
				return err
			}
		}

		wc := &object.WorkspaceConfig{
			RepoId:  curRepo,
			RevHash: eachRev.Hash.String(),
		}
		_, err := ExecCurRevWithConfig(uploadSrc, wc, c)
		if err != nil {
			return err
		}
	}

	// recover
	if c.Depth != 1 {
		core.Log.Infof("recover checkout: %s", head)
		err = tree.Checkout(&git.CheckoutOptions{
			Hash: head.Hash(),
			Keep: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func ExecCurRevWithConfig(uploadSrc string, wc *object.WorkspaceConfig, c *Config) (ExecuteCacheMap, error) {
	cacheMap := make(ExecuteCacheMap)
	filterFunc, err := createFileFilter(c)
	if err != nil {
		return nil, err
	}

	runner := &core.Runner{}
	var lang []string
	if len(c.Lang) == 0 {
		langFromDir, err := runner.GuessLangFromDir(c.Src, filterFunc)
		if err != nil {
			return nil, err
		}
		lang = []string{string(langFromDir)}
	} else {
		lang = c.Lang
	}
	if len(lang) == 0 {
		return nil, errors.New("no valid lang found")
	}

	for _, eachLang := range lang {
		eachLangType := core.LangTypeValueOf(eachLang)
		if !eachLangType.IsSupported() {
			core.Log.Warnf("lang %v not supported, supported list: %v", eachLangType, core.SupportedLangs)
			continue
		}
		core.Log.Infof("scan lang: %v", eachLang)
		cache, err := execCurRevCurLangWithConfig(uploadSrc, eachLangType, filterFunc, wc, c)
		cacheMap[eachLangType] = cache
		if err != nil {
			return nil, err
		}
	}
	return cacheMap, nil
}

func execCurRevCurLangWithConfig(uploadSrc string, lang core.LangType, filterFunc func(path string) bool, wc *object.WorkspaceConfig, c *Config) (*ExecuteCache, error) {
	cache := &ExecuteCache{}
	f, err := sibyl2.ExtractFunction(uploadSrc, &sibyl2.ExtractConfig{
		FileFilter: filterFunc,
		LangType:   lang,
	})
	if err != nil {
		return nil, err
	}

	funcUrl := c.GetFuncUploadUrl()
	funcCtxUrl := c.GetFuncCtxUploadUrl()
	clazzUrl := c.GetClazzUploadUrl()

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
			return nil, err
		}

		core.Log.Infof("start calculating func graph")
		g, err := sibyl2.AnalyzeFuncGraph(f, s)
		if err != nil {
			return nil, err
		}
		cache.AnalyzeGraph = g

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
			return nil, err
		}
		core.Log.Infof("classes ready")
		if !c.Dry {
			uploadClazz(clazzUrl, wc, s, c.Batch)
		}
	}
	return cache, nil
}

func createFileFilter(c *Config) (func(path string) bool, error) {
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
