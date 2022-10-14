package extractor

import (
	"errors"
	"github.com/williamfzc/sibyl2/pkg/core"
)

func (extractor *JavaExtractor) IsCall(unit *core.Unit) bool {
	if unit.Kind == KindJavaMethodInvocation {
		return true
	}
	return false
}

func (extractor *JavaExtractor) ExtractCalls(units []*core.Unit) ([]*Call, error) {
	var ret []*Call
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

func (extractor *JavaExtractor) unit2Call(unit *core.Unit) (*Call, error) {
	funcUnit := core.FindFirstByOneOfKindInParent(unit, KindJavaMethodDeclaration)
	var srcFunc *Function
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

	var argumentPart *core.Unit
	var arguments []string
	var caller string

	callerPart := core.FindFirstByFieldInSubs(unit, FieldJavaObject)
	if callerPart == nil {
		// b()
		callerPart = core.FindFirstByFieldInSubs(unit, FieldJavaName)
		argumentPart = core.FindFirstByFieldInSubs(unit, FieldJavaArguments)
		caller = callerPart.Content
	} else {
		// a.b()
		identifiers := core.FindAllByKindInSubs(unit, KindJavaIdentifier)
		argumentPart = core.FindFirstByFieldInSubs(unit, FieldJavaName)

		if len(identifiers) == 0 {
			core.DebugDfs(unit, 0)
			return nil, errors.New("no id: " + unit.Content)
		}

		caller = callerPart.Content + "." + identifiers[len(identifiers)-1].Content
	}

	// not perfect, eg: anonymous function call?
	if argumentPart != nil {
		for _, each := range core.FindAllByKindInSubs(argumentPart, KindJavaIdentifier) {
			arguments = append(arguments, each.Content)
		}
	}

	ret := &Call{
		Src:       srcFunc.GetSignature(),
		Caller:    caller,
		Arguments: arguments,
		Span:      unit.Span,
	}
	return ret, nil
}
