package extractor

type DataType interface {
	Dt()
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
