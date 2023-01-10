package ext

import (
	"github.com/go-git/go-git/v5"
	"github.com/opensibyl/sibyl2/pkg/core"
)

func LoadGitRepo(gitDir string) (*git.Repository, error) {
	repo, err := git.PlainOpen(gitDir)
	if err != nil {
		core.Log.Errorf("load repo from %s failed", gitDir)
		return nil, err
	}
	return repo, nil
}
