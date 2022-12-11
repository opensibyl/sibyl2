package history

import (
	"testing"
)

func TestHistory(t *testing.T) {
	err := handle("../../../..", "output.html")
	if err != nil {
		panic(err)
	}
}
