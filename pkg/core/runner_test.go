package core

import (
	"testing"
)

func TestRunner_HandleFile_Golang(t *testing.T) {
	runner := &Runner{}
	_, err := runner.HandleFile(".", "GOLANG")
	if err != nil {
		panic(err)
	}
}
