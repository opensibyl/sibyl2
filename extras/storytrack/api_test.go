package storytrack

import (
	"encoding/json"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/williamfzc/sibyl2/pkg/core"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApi(t *testing.T) {
	abs, _ := filepath.Abs("../..")
	trackResult, err := Track(abs, "c11370ebaefee9347c528effc2ab6f0a5f1648c2", func(commit *object.Commit) bool {
		if strings.Contains(commit.Message, "refactor") {
			return true
		}
		return false
	}, core.LangGo)
	if err != nil {
		return
	}
	output, err := json.MarshalIndent(&trackResult, "", "  ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("a.json", output, 0644)
	if err != nil {
		panic(err)
	}
}
