package main

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"strings"
	"testing"
)

func TestApi(t *testing.T) {
	StoryTrack.Track("../..", "c11370ebaefee9347c528effc2ab6f0a5f1648c2", func(commit *object.Commit) bool {
		if strings.Contains(commit.Message, "refactor") {
			return true
		}
		return false
	})
}
