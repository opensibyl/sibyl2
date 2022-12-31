package golang

import (
	"errors"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

func (extractor *Extractor) IsCall(unit *core.Unit) bool {
	if unit.Kind == KindGolangCallExpression {
		return true
	}
	return false
}

func (extractor *Extractor) ExtractCalls(units []*core.Unit) ([]*object.Call, error) {
	var ret []*object.Call
	for _, eachUnit := range units {
		if !extractor.IsCall(eachUnit) {
			continue
		}

		eachCall, err := extractor.unit2Call(eachUnit)
		if err != nil {
			core.Log.Warnf("err: %v", err)
			continue
		}
		ret = append(ret, eachCall)
	}
	return ret, nil
}

func (extractor *Extractor) unit2Call(unit *core.Unit) (*object.Call, error) {
	// todo: what about nested call
	funcUnit := core.FindFirstByOneOfKindInParent(unit, KindGolangFuncDecl, KindGolangMethodDecl)
	var srcFunc *object.Function
	var err error
	if funcUnit != nil {
		srcFunc, err = extractor.ExtractFunction(funcUnit)
		if err != nil {
			return nil, errors.New("convert func failed: " + funcUnit.Content)
		}
	}

	// headless, give up (temp
	if srcFunc == nil {
		return nil, errors.New("headless call")
	}

	funcPart := core.FindFirstByFieldInSubsWithBfs(unit, FieldGolangFunction)
	argumentPart := core.FindFirstByFieldInSubsWithBfs(unit, FieldGolangArguments)

	// not perfect, eg: anonymous function call?
	var arguments []string
	for _, each := range argumentPart.SubUnits {
		if each.Kind == KindGolangIdentifier {
			arguments = append(arguments, each.Content)
		}
	}

	ret := &object.Call{
		Src:       srcFunc.GetSignature(),
		Caller:    funcPart.Content,
		Arguments: arguments,
		Span:      unit.Span,
	}
	return ret, nil
}
