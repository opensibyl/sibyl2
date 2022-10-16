package storytrack

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"golang.org/x/exp/slices"
	"path/filepath"
	"strings"
)

type TrackResult struct {
	Commits   []string                `json:"commits"`
	Files     map[string][]int        `json:"files"`
	Functions []*extractor.FileResult `json:"functions"`
}

type Rule = func(commit *object.Commit) bool

func TrackWithSharpId(gitDir string, targetRev string, ids []int, langType core.LangType) (*TrackResult, error) {
	rule := func(commit *object.Commit) bool {
		for _, each := range ids {
			matchKey := fmt.Sprintf("#%d", each)
			if strings.Contains(commit.Message, matchKey) {
				return true
			}
		}
		return false
	}
	return Track(gitDir, targetRev, rule, langType)
}

func Track(gitDir string, targetRev string, ruleJudge Rule, langType core.LangType) (*TrackResult, error) {
	gitDir, err := filepath.Abs(gitDir)
	if err != nil {
		return nil, err
	}
	repo, err := loadRepo(gitDir)
	if err != nil {
		return nil, err
	}

	var from plumbing.Hash
	if targetRev == "" {
		head, err := repo.Head()
		if err != nil {
			return nil, err
		}
		from = head.Hash()
	} else {
		from = plumbing.NewHash(targetRev)
	}
	cIter, err := repo.Log(&git.LogOptions{From: from})
	var targetCommits []*object.Commit
	cIter.ForEach(func(c *object.Commit) error {
		if ruleJudge(c) {
			targetCommits = append(targetCommits, c)
		}
		return nil
	})
	var targetHashes []string
	for _, each := range targetCommits {
		targetHashes = append(targetHashes, each.Hash.String())
	}

	var relatedFiles []string
	for _, each := range targetCommits {
		fIter, err := each.Files()
		if err != nil {
			return nil, err
		}
		fIter.ForEach(func(file *object.File) error {
			relatedFiles = append(relatedFiles, file.Name)
			return nil
		})
	}

	filter := func(path string) bool {
		relpath, err := filepath.Rel(gitDir, path)
		if err != nil {
			return false
		}
		return slices.Contains(relatedFiles, relpath)
	}

	fileResults, err := sibyl2.Extract(gitDir, &sibyl2.ExtractConfig{
		LangType:    langType,
		ExtractType: extractor.TypeExtractFunction,
		FileFilter:  filter,
	})

	// git blame
	var lineRange = make(map[string][]int)
	targetCommit, err := repo.CommitObject(from)
	if err != nil {
		return nil, err
	}
	for _, eachFile := range relatedFiles {
		// this file may be not existed in current version
		blame, err := git.Blame(targetCommit, eachFile)
		if err != nil {
			return nil, err
		}
		var eachFileLines []int
		for index, eachLine := range blame.Lines {
			// line num
			lineNum := index + 1
			if slices.Contains(targetHashes, eachLine.Hash.String()) {
				eachFileLines = append(eachFileLines, lineNum)
			}
		}
		lineRange[eachFile] = eachFileLines
	}
	for _, each := range fileResults {
		var newDt []extractor.DataType
		if lines, ok := lineRange[each.Path]; ok {
			for _, eachDt := range each.Units {
				if eachDt.GetSpan().ContainAnyLine(lines) {
					newDt = append(newDt, eachDt)
				}
			}
		}
		each.Units = newDt
	}
	var newFileResults []*extractor.FileResult
	for _, each := range fileResults {
		if len(each.Units) > 0 {
			newFileResults = append(newFileResults, each)
		}
	}

	final := &TrackResult{}
	final.Commits = targetHashes
	final.Files = lineRange
	final.Functions = newFileResults
	return final, nil
}

func loadRepo(gitDir string) (*git.Repository, error) {
	repo, err := git.PlainOpen(gitDir)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
