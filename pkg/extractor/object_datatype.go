package extractor

import (
	"fmt"

	"github.com/williamfzc/sibyl2/pkg/core"
)

type DataType interface {
	GetDesc() string
	GetSpan() *core.Span
}

func (symbol *Symbol) GetDesc() string {
	return fmt.Sprintf("<symbol %s %s>", symbol.Kind, symbol.Symbol)
}

func (function *Function) GetDesc() string {
	return fmt.Sprintf("<function %s %v>", function.GetSignature(), function.GetSpan())
}

func (call *Call) GetDesc() string {
	return fmt.Sprintf("<call %s(%v) in %s>", call.Caller, call.Arguments, call.Src)
}

func DataTypeOf[T DataType](dataList []T) []DataType {
	var retUnits []DataType
	for _, each := range dataList {
		retUnits = append(retUnits, each)
	}
	return retUnits
}
