package casedoctor

import (
	"fmt"
	"testing"
)

func TestCheckCases(t *testing.T) {
	cases, err := CheckCases("../..")
	if err != nil {
		return
	}
	for _, each := range cases.Functions {
		fmt.Println(each.Content)
	}
}
