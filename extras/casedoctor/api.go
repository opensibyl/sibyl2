package casedoctor

import (
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"path/filepath"
	"strings"
)

type FunctionWithContent struct {
	*extractor.Function
	Content string `json:"content"`
}

type CheckResult struct {
	// todo: need some stats
	Functions []*FunctionWithContent `json:"functions"`
}

func CheckCases(targetDir string) (*CheckResult, error) {
	filter := func(p string) bool {
		name := strings.TrimSuffix(p, filepath.Ext(p))
		endsWithTest := strings.HasSuffix(name, "Test")
		endsWithTests := strings.HasSuffix(name, "Tests")
		endsWithTestSnack := strings.HasSuffix(name, "_test")
		startsWithTest := strings.HasPrefix(name, "Test")
		return endsWithTest || endsWithTests || endsWithTestSnack || startsWithTest
	}
	fileResults, err := sibyl2.Extract(targetDir, &sibyl2.ExtractConfig{
		ExtractType: extractor.TypeExtractFunction,
		FileFilter:  filter,
	})
	if err != nil {
		return nil, err
	}
	var funcs []*FunctionWithContent
	for _, each := range fileResults {
		for _, eachUnit := range each.Units {
			function, _ := eachUnit.(*extractor.Function)
			content := function.GetUnit().Content
			fwc := &FunctionWithContent{
				Function: function,
				Content:  content,
			}
			funcs = append(funcs, fwc)
		}
	}
	result := &CheckResult{}
	result.Functions = funcs
	return result, nil
}
