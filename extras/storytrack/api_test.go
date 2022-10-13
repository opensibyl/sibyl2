package main

import (
	"encoding/json"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os"
	"strings"
	"testing"
)

func TestApi(t *testing.T) {
	trackResult, err := StoryTrack.Track("../..", "c11370ebaefee9347c528effc2ab6f0a5f1648c2", func(commit *object.Commit) bool {
		if strings.Contains(commit.Message, "refactor") {
			return true
		}
		return false
	})
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
