package model

type DataType interface {
	Dt()
}

func (*Symbol) Dt() {
}

func (*Function) Dt() {
}

func (*Call) Dt() {
}
