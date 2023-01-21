package core

import (
	"testing"
)

func TestRunner_HandleFile_Golang(t *testing.T) {
	t.Parallel()
	runner := &Runner{}
	_, err := runner.File2Units(".", "GOLANG", nil)
	if err != nil {
		panic(err)
	}
}
