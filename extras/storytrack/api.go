package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var StoryTrack = &storyTrack{}

type storyTrack struct {
}

type TrackResult struct {
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

	for _, each := range targetCommits {
		fmt.Println(each.Message)
	}

	return nil, nil
}

func (st *storyTrack) loadRepo(gitDir string) (*git.Repository, error) {
	repo, err := git.PlainOpen(gitDir)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
