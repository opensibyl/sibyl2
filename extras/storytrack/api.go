package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"golang.org/x/exp/slices"
	"path/filepath"
)

type TrackResult struct {
	Functions []*extractor.FileResult
}

type Rule = func(commit *object.Commit) bool

func Track(gitDir string, targetRev string, ruleJudge Rule) (*TrackResult, error) {
	repo, err := loadRepo(gitDir)
	if err != nil {
		return nil, err
	}

	cIter, err := repo.Log(&git.LogOptions{From: plumbing.NewHash(targetRev)})
	var targetCommits []*object.Commit
	cIter.ForEach(func(c *object.Commit) error {
		if ruleJudge(c) {
			targetCommits = append(targetCommits, c)
		}
		return nil
	})

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
		LangType:    core.LangGo,
		ExtractType: extractor.TypeExtractFunction,
		FileFilter:  filter,
	})

	final := &TrackResult{}
	final.Functions = fileResults
	return final, nil
}

func loadRepo(gitDir string) (*git.Repository, error) {
	repo, err := git.PlainOpen(gitDir)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
