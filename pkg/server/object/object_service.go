package object

import (
	"github.com/opensibyl/sibyl2"
)

type FunctionWithSignature struct {
	*sibyl2.FunctionWithTag
	Signature string `json:"signature"`
}

type FunctionContextSlimWithSignature struct {
	*sibyl2.FunctionContextSlim
	Signature string `json:"signature"`
}
