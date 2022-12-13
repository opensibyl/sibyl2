package extractor

import (
	"github.com/opensibyl/sibyl2/pkg/core"
)

type DataType interface {
	GetIndexName() string // for query and locate from outside
	GetDesc() string      // easy to understand what it actually contains
	GetSpan() *core.Span
}
type DataTypes = []DataType

func DataTypeOf[T DataType](dataList []T) []DataType {
	var retUnits []DataType
	for _, each := range dataList {
		retUnits = append(retUnits, each)
	}
	return retUnits
}
