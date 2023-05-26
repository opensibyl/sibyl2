package object

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/opensibyl/sibyl2/pkg/core"
)

type FuncSignature = string
type ValueUnit struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

const DescSplit = "|,|"

type Function struct {
	Name       string       `json:"name" bson:"name"`
	Receiver   string       `json:"receiver" bson:"receiver"`
	Namespace  string       `json:"namespace" bson:"namespace"`
	Parameters []*ValueUnit `json:"parameters" bson:"parameters"`
	Returns    []*ValueUnit `json:"returns" bson:"returns"`

	// where the definition actually happened
	// avoid the impacts of annotations and decorations
	DefLine int `json:"defLine" bson:"defLine"`

	// this span will include header and content
	Span core.Span `json:"span" bson:"span"`

	// which includes only body
	BodySpan core.Span `json:"bodySpan" bson:"bodySpan"`

	// which contains language-specific contents
	Extras interface{} `json:"extras" bson:"extras,omitempty"`

	// ptr to origin Unit
	Unit *core.Unit `json:"-" bson:"-"`

	// language
	Lang core.LangType `json:"lang" bson:"lang"`
}

func NewFunction() *Function {
	return &Function{}
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

func (f *Function) GetSignature() FuncSignature {
	prefix := fmt.Sprintf("%s|%s|%s", f.Namespace, f.Receiver, f.Name)

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

func (f *Function) GetIndexName() string {
	return f.Name
}

func (f *Function) GetDesc() string {
	return fmt.Sprintf("%v%s%v", f.GetSpan(), DescSplit, f.GetSignature())
}

func (f *Function) GetSpan() *core.Span {
	return &f.Span
}

func (f *Function) GetUnit() *core.Unit {
	return f.Unit
}
