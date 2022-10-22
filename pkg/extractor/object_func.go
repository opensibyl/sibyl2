package extractor

import (
	"fmt"
	"github.com/williamfzc/sibyl2/pkg/core"
	"strings"
)

type ValueUnit struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Function struct {
	Name       string       `json:"name"`
	Receiver   string       `json:"receiver"`
	Parameters []*ValueUnit `json:"parameters"`
	Returns    []*ValueUnit `json:"returns"`
	Span       core.Span    `json:"span"`

	// which contains language-specific contents
	Extras interface{} `json:"extras"`

	// ptr to origin unit
	unit *core.Unit
}

func NewFunction() *Function {
	return &Function{}
}

type FuncSignature = string

func (f *Function) GetSignature() FuncSignature {
	prefix := fmt.Sprintf("%s::%s", f.Receiver, f.Name)

	params := make([]string, len(f.Parameters))
	for i, each := range f.Parameters {
		params[i] = each.Type
	}
	paramPart := strings.Join(params, ",")

	rets := make([]string, len(f.Returns))
	for i, each := range f.Returns {
		rets[i] = each.Type
	}
	retPart := strings.Join(rets, ",")

	return fmt.Sprintf("%s|%s|%s", prefix, paramPart, retPart)
}

func (f *Function) GetSpan() *core.Span {
	return &f.Span
}

func (f *Function) GetUnit() *core.Unit {
	return f.unit
}
