package extractor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/opensibyl/sibyl2/pkg/core"
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
	// this span will include header and content
	Span core.Span `json:"span"`
	// which includes only body
	BodySpan core.Span `json:"bodySpan"`

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

// ToMap export a very simple map without any custom structs. It will lose ptr to origin unit.
func (f *Function) ToMap() (map[string]any, error) {
	b, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (f *Function) ToJson() ([]byte, error) {
	m, err := f.ToMap()
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// Map2Func reverse ToMap
func Map2Func(exported map[string]any) (*Function, error) {
	ret := &Function{}
	err := mapstructure.Decode(exported, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func Json2Func(exported []byte) (*Function, error) {
	var m map[string]any
	err := json.Unmarshal(exported, &m)
	if err != nil {
		return nil, err
	}
	return Map2Func(m)
}

func (f *Function) GetIndexName() string {
	return f.Name
}

func (f *Function) GetDesc() string {
	return fmt.Sprintf("<function %s %v>", f.GetSignature(), f.GetSpan())
}

func (f *Function) GetSpan() *core.Span {
	return &f.Span
}

func (f *Function) GetUnit() *core.Unit {
	return f.unit
}
