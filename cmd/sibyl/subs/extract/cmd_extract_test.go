package extract

import (
	"bytes"
	"testing"
)

func Test_ExecuteCommand_Func(t *testing.T) {
	cmd := NewExtractCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--lang", "GOLANG", "--type", "func"})
	cmd.Execute()
}

func Test_ExecuteCommand_Call(t *testing.T) {
	cmd := NewExtractCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--lang", "GOLANG", "--type", "call"})
	cmd.Execute()
}

func Test_ExecuteCommand_Call_Java(t *testing.T) {
	cmd := NewExtractCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--lang", "JAVA", "--type", "call"})
	cmd.Execute()
}

func Test_ExecuteCommand_Call_Python(t *testing.T) {
	cmd := NewExtractCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--lang", "PYTHON", "--type", "func"})
	cmd.Execute()
}
