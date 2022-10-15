package extractor

import "github.com/williamfzc/sibyl2/pkg/core"

type DataType interface {
	Dt()
	GetSpan() *core.Span
}

func (*Symbol) Dt() {
}

func (*Function) Dt() {
}

func (*Call) Dt() {
}

func DataTypeOf[T DataType](dataList []T) []DataType {
	var retUnits []DataType
	for _, each := range dataList {
		retUnits = append(retUnits, each)
	}
	return retUnits
}
