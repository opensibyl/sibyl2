package diff

import (
	"bytes"
	"testing"
)

func TestDiff(t *testing.T) {
	t.Skip()
	cmd := NewDiffCommand()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--src", "../../..", "--patch", ""})
	cmd.Execute()
}
