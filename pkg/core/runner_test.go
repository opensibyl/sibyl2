package core

import (
	"fmt"
	"testing"
)

func TestRunner_HandleFile(t *testing.T) {
	runner := &Runner{}
	ret, err := runner.HandleFile(".", "GOLANG")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v+\n", ret)
}
