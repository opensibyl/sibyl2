package main

import (
	"bytes"
	"testing"
)

func TestRoot(t *testing.T) {
	cmd := rootCmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--lang", "GOLANG", "--ids", "12345", "--src", "../../../.."})
	cmd.Execute()
}

func TestRootShort(t *testing.T) {
	cmd := rootCmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--ids", "12345", "--src", "../../../.."})
	cmd.Execute()
}
