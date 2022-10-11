package model

import (
	"fmt"
	"strings"
)

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
