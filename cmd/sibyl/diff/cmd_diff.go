package diff

import (
	"path/filepath"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/core"
)

type TrackResult struct {
	Before    string                    `json:"before"`
	After     string                    `json:"after"`
	Functions []*sibyl2.FunctionContext `json:"functions"`
}

func track(gitDir string, target string) (*TrackResult, error) {
	gitDir, err := filepath.Abs(gitDir)
	if err != nil {
		return nil, err
	}
	repo, err := loadRepo(gitDir)
	if err != nil {
		return nil, err
	}

	beforeCommitTxt, err := repo.Head()
	if err != nil {
		return nil, err
	}
	beforeCommit, err := repo.CommitObject(beforeCommitTxt.Hash())
	afterCommit, err := repo.CommitObject(plumbing.NewHash(target))
	if err != nil {
		return nil, err
	}
	result := &TrackResult{beforeCommit.String(), target, nil}

	patch, err := afterCommit.Patch(beforeCommit)
	if err != nil {
		return nil, err
	}

	f, err := sibyl2.ExtractFunction(gitDir, sibyl2.DefaultConfig())
	if err != nil {
		return nil, err
	}

	s, err := sibyl2.ExtractSymbol(gitDir, sibyl2.DefaultConfig())
	if err != nil {
		return nil, err
	}
	_, err = sibyl2.AnalyzeFuncGraph(f, s)
	if err != nil {
		return nil, err
	}

	core.Log.Infof("patch: %s", patch.String())
	parsed, _, err := gitdiff.Parse(strings.NewReader(patch.String()))
	if err != nil {
		return nil, err
	}
	affectedMap := make(map[string][]int)
	for _, each := range parsed {
		if each.IsBinary || each.IsDelete {
			continue
		}
		affectedMap[each.NewName] = make([]int, 0)
		fragments := each.TextFragments
		for _, eachF := range fragments {
			right := int(eachF.NewPosition+eachF.NewLines) - 1
			left := int(eachF.NewPosition)

			for i := left; i < right; i++ {
				if eachF.Lines[i].Op == gitdiff.OpAdd {
					affectedMap[each.NewName] = append(affectedMap[each.NewName], i)
				}
			}
		}
	}
	for fileName, lines := range affectedMap {
		for _, eachFunc := range f {
			if eachFunc.Path == fileName {
				for _, eachFuncUnit := range eachFunc.Units {
					if eachFuncUnit.GetSpan().ContainAnyLine(lines...) {
						core.Log.Infof("func %v changed", eachFuncUnit)
					}
				}
			}
		}
	}

	// todo
	return result, nil
}

var diffGitDir string
var diffRev string

func NewDiffCommand() *cobra.Command {
	diffCmd := &cobra.Command{
		Use:    "diff",
		Short:  "test",
		Long:   `test`,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := track(diffGitDir, diffRev)
			if err != nil {
				panic(err)
			}
		},
	}
	diffCmd.PersistentFlags().StringVar(&diffGitDir, "src", ".", "src dir path")
	diffCmd.PersistentFlags().StringVar(&diffRev, "rev", "", "rev")
	return diffCmd
}

func loadRepo(gitDir string) (*git.Repository, error) {
	repo, err := git.PlainOpen(gitDir)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
