package main

import (
	"bytes"
	"testing"
)

func TestRootShort(t *testing.T) {
	cmd := rootCmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--src", "../../../.."})
	cmd.Execute()
}
