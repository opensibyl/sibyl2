package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"path/filepath"
)

var StoryTrack = &storyTrack{}

type storyTrack struct {
}

type TrackResult struct {
	Functions []*extractor.FileResult
}

type Rule = func(commit *object.Commit) bool

func (st *storyTrack) Track(gitDir string, targetRev string, ruleJudge Rule) (*TrackResult, error) {
	repo, err := st.loadRepo(gitDir)
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

	var results []*extractor.FileResult
	for _, each := range targetCommits {
		fIter, err := each.Files()
		if err != nil {
			return nil, err
		}
		fIter.ForEach(func(file *object.File) error {
			absFile := filepath.Join(gitDir, file.Name)
			core.Log.Infof("checking file: %s", absFile)

			// todo: incorrect path
			fileResults, err := sibyl2.Extract(absFile, &sibyl2.ExtractConfig{
				LangType:    core.LangGo,
				ExtractType: extractor.TypeExtractFunction,
			})
			if err != nil {
				return err
			}
			results = append(results, fileResults...)
			return nil
		})
	}

	final := &TrackResult{}
	final.Functions = results
	return final, nil
}

func (st *storyTrack) loadRepo(gitDir string) (*git.Repository, error) {
	repo, err := git.PlainOpen(gitDir)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
