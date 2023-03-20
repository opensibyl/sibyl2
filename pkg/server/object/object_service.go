package object

import (
	"github.com/opensibyl/sibyl2"
)

type FunctionWithSignature struct {
	*sibyl2.FunctionWithTag `bson:",inline"`
	Signature               string `json:"signature" bson:"signature"`
}

type FunctionContextSlimWithSignature struct {
	*sibyl2.FunctionContextSlim `bson:",inline"`
	Signature                   string `json:"signature"`
}

type ClazzWithSignature struct {
	*sibyl2.ClazzWithPath `bson:",inline"`
	Signature             string `json:"signature"`
}
