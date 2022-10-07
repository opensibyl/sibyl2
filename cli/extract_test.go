package cli

import (
	"bytes"
	"testing"
)

func Test_ExecuteCommand(t *testing.T) {
	cmd := NewExtractCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--lang", "GOLANG", "--type", "func"})
	cmd.Execute()
}
