package diff

import (
	"bytes"
	"testing"
)

func TestDiff(t *testing.T) {
	cmd := NewDiffCommand()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--src", "../../..", "--rev", "8143a59aef5ec5352e416265d44cd58abd89a461"})
	cmd.Execute()
}
