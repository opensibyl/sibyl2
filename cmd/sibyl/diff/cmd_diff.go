package diff

import (
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/williamfzc/sibyl2"
)

type TrackResult struct {
	Before    string                    `json:"before"`
	After     string                    `json:"after"`
	Functions []*sibyl2.FunctionContext `json:"functions"`
}

func Track(gitDir string, before string, after string) (*TrackResult, error) {
	result := &TrackResult{before, after, nil}

	gitDir, err := filepath.Abs(gitDir)
	if err != nil {
		return nil, err
	}
	repo, err := loadRepo(gitDir)
	if err != nil {
		return nil, err
	}

	beforeRef, err := repo.Reference(plumbing.ReferenceName(before), true)
	if err != nil {
		return nil, err
	}
	afterRef, err := repo.Reference(plumbing.ReferenceName(after), true)
	if err != nil {
		return nil, err
	}
	beforeCommit, err := repo.CommitObject(beforeRef.Hash())
	if err != nil {
		return nil, err
	}
	afterCommit, err := repo.CommitObject(afterRef.Hash())
	if err != nil {
		return nil, err
	}
	_, err = afterCommit.Patch(beforeCommit)
	if err != nil {
		return nil, err
	}

	f, err := sibyl2.ExtractFunction(gitDir, sibyl2.DefaultConfig())
	if err != nil {
		panic(err)
	}

	s, err := sibyl2.ExtractSymbol(gitDir, sibyl2.DefaultConfig())
	if err != nil {
		panic(err)
	}
	_, err = sibyl2.AnalyzeFuncGraph(f, s)
	if err != nil {
		panic(err)
	}

	// todo
	return result, nil
}

func loadRepo(gitDir string) (*git.Repository, error) {
	repo, err := git.PlainOpen(gitDir)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
