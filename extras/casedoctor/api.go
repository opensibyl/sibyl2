package casedoctor

import (
	"path/filepath"
	"strings"

	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/extractor"
)

type FunctionWithConclusion struct {
	*extractor.Function
	Path        string `json:"path"`
	Content     string `json:"content"`
	AssertCount int    `json:"assertCount"`
}

func (fwc *FunctionWithConclusion) HasAssertion() bool {
	return fwc.AssertCount > 0
}

type CheckStat struct {
	Files       int `json:"files"`
	Cases       int `json:"cases"`
	AssertCases int `json:"assertCases"`
}

type CheckResult struct {
	Stat      *CheckStat                `json:"stat"`
	Functions []*FunctionWithConclusion `json:"functions"`
}

func (cr *CheckResult) SyncStat() {
	fileSet := make(map[string]interface{})
	cases := 0
	assertCases := 0
	for _, each := range cr.Functions {
		fileSet[each.Path] = nil
		cases++
		if each.HasAssertion() {
			assertCases++
		}
	}
	cr.Stat = &CheckStat{
		Files:       len(fileSet),
		Cases:       cases,
		AssertCases: assertCases,
	}
}

func CheckCases(targetDir string) (*CheckResult, error) {
	filter := func(p string) bool {
		name := strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
		endsWithTest := strings.HasSuffix(name, "Test")
		endsWithTests := strings.HasSuffix(name, "Tests")
		endsWithTestSnack := strings.HasSuffix(name, "_test")
		startsWithTest := strings.HasPrefix(name, "Test")
		return endsWithTest || endsWithTests || endsWithTestSnack || startsWithTest
	}
	fileResults, err := sibyl2.ExtractFunction(targetDir, &sibyl2.ExtractConfig{
		FileFilter: filter,
	})
	if err != nil {
		return nil, err
	}
	var functions []*FunctionWithConclusion
	for _, each := range fileResults {
		for _, eachFunction := range each.Units {
			content := eachFunction.GetUnit().Content
			fwc := &FunctionWithConclusion{
				Path:     each.Path,
				Function: eachFunction,
				Content:  content,
				// todo: not good?
				AssertCount: strings.Count(content, "assert"),
			}
			functions = append(functions, fwc)
		}
	}
	result := &CheckResult{}
	result.Functions = functions
	result.SyncStat()
	return result, nil
}
