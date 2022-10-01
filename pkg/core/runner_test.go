package core

import (
	"testing"
)

func TestRunner_HandleFile_Golang(t *testing.T) {
	runner := &Runner{}
	_, err := runner.File2Units(".", "GOLANG")
	if err != nil {
		panic(err)
	}
}
