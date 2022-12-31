package extractor

import (
	"errors"

	"github.com/opensibyl/sibyl2/pkg/core"
)

func (extractor *PythonExtractor) IsFunction(unit *core.Unit) bool {
	// python has only func type
	if unit.Kind == KindPythonFunctionDefinition {
		return true
	}
	return false
}

func (extractor *PythonExtractor) ExtractFunctions(units []*core.Unit) ([]*Function, error) {
	var ret []*Function
	for _, eachUnit := range units {
		if !extractor.IsFunction(eachUnit) {
			continue
		}
		eachFunc, err := extractor.unit2Function(eachUnit)
		if err != nil {
			return nil, err
		}
		ret = append(ret, eachFunc)
	}
	return ret, nil
}

func (extractor *PythonExtractor) ExtractFunction(unit *core.Unit) (*Function, error) {
	data, err := extractor.ExtractFunctions([]*core.Unit{unit})
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("unit not a func: " + unit.Content)
	}
	return data[0], nil
}

func (extractor *PythonExtractor) unit2Function(unit *core.Unit) (*Function, error) {
	funcUnit := NewFunction()
	funcUnit.Span = unit.Span
	funcUnit.unit = unit
	// body scope
	funcBody := core.FindFirstByKindInSubsWithBfs(unit, KindPythonBlock)
	if funcBody != nil {
		funcUnit.BodySpan = funcBody.Span
	}

	funcName := core.FindFirstByKindInSubsWithBfs(unit, KindPythonIdentifier)
	if funcName == nil {
		return nil, errors.New("no func name found in " + unit.Content)
	}
	funcUnit.Name = funcName.Content

	extras := &PythonFunctionExtras{}
	if unit.ParentUnit.Kind == KindPythonDecoratedDefinition {
		decoratedUnit := unit.ParentUnit
		decorators := core.FindAllByKindInSubs(decoratedUnit, KindPythonDecorator)
		if len(decorators) == 0 {
			return nil, errors.New("no deco found in " + decoratedUnit.Content)
		}

		var decoContents []string
		for _, each := range decorators {
			decoContents = append(decoContents, each.Content)
		}
		extras.Decorators = decoContents
	}
	funcUnit.Extras = extras

	// todo: returns and params?
	return funcUnit, nil
}
