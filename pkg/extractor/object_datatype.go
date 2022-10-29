package extractor

import (
	"fmt"

	"github.com/williamfzc/sibyl2/pkg/core"
)

type DataType interface {
	GetIndexName() string // for query and locate from outside
	GetDesc() string      // easy to understand what it actually contains
	GetSpan() *core.Span
}

func (s *Symbol) GetIndexName() string {
	return s.Symbol
}

func (s *Symbol) GetDesc() string {
	return fmt.Sprintf("<symbol %s %s>", s.Kind, s.Symbol)
}

func (f *Function) GetIndexName() string {
	return f.Name
}

func (f *Function) GetDesc() string {
	return fmt.Sprintf("<function %s %v>", f.GetSignature(), f.GetSpan())
}

func (c *Call) GetIndexName() string {
	// hard to represent ...
	return fmt.Sprintf("%s->%s", c.Src, c.Caller)
}

func (c *Call) GetDesc() string {
	return fmt.Sprintf("<call %s(%v) in %s>", c.Caller, c.Arguments, c.Src)
}

func DataTypeOf[T DataType](dataList []T) []DataType {
	var retUnits []DataType
	for _, each := range dataList {
		retUnits = append(retUnits, each)
	}
	return retUnits
}
