package extractor

import (
	"errors"

	"github.com/williamfzc/sibyl2/pkg/core"
	"golang.org/x/exp/slices"
)

type GolangFuncExtras struct {
}

func (extractor *GolangExtractor) IsFunction(unit *core.Unit) bool {
	allowed := []core.KindRepr{
		KindGolangMethodDecl,
		KindGolangFuncDecl,
	}
	return slices.Contains(allowed, unit.Kind)
}

func (extractor *GolangExtractor) ExtractFunctions(units []*core.Unit) ([]*Function, error) {
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

func (extractor *GolangExtractor) ExtractFunction(unit *core.Unit) (*Function, error) {
	data, err := extractor.ExtractFunctions([]*core.Unit{unit})
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("unit not a func: " + unit.Content)
	}
	return data[0], nil
}

func (extractor *GolangExtractor) unit2Function(unit *core.Unit) (*Function, error) {
	switch unit.Kind {
	case KindGolangFuncDecl:
		return extractor.funcUnit2Function(unit)
	case KindGolangMethodDecl:
		return extractor.methodUnit2Function(unit)
	default:
		// should not reach here
		return nil, errors.New("IMPOSSIBLE")
	}
}

func (extractor *GolangExtractor) methodUnit2Function(unit *core.Unit) (*Function, error) {
	funcUnit := &Function{}
	funcUnit.Span = unit.Span
	funcUnit.unit = unit

	// body scope
	funcBody := core.FindFirstByFieldInSubsWithBfs(unit, FieldGolangResult)
	if funcBody != nil {
		funcUnit.BodySpan = funcBody.Span
	}

	// name
	funcIdentifier := core.FindFirstByKindInSubsWithDfs(unit, KindGolangFieldIdentifier)
	if funcIdentifier == nil {
		return nil, errors.New("no func name found in " + unit.Content)
	}
	funcUnit.Name = funcIdentifier.Content

	// receiver
	parameterList := core.FindFirstByKindInSubsWithDfs(unit, KindGolangParameterList)
	parameterList = core.FindFirstByKindInSubsWithDfs(parameterList, KindGolangParameterList)
	receiverDecl := core.FindFirstByKindInSubsWithDfs(parameterList, KindGolangParameterDecl)
	typeDecl := core.FindFirstByFieldInSubsWithDfs(receiverDecl, FieldGolangType)
	if typeDecl == nil {
		return nil, errors.New("no receiver found in: " + typeDecl.Content)
	}
	funcUnit.Receiver = typeDecl.Content

	// params
	paramListList := core.FindAllByKindInSubsWithDfs(unit, KindGolangParameterList)
	// no param == empty slice, never nil
	paramList := paramListList[1]
	for _, each := range core.FindAllByKindInSubsWithDfs(paramList, KindGolangParameterDecl) {
		typeName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangType)
		paramName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangName)
		var paramNameContent string
		if paramName == nil {
			paramNameContent = ""
		} else {
			paramNameContent = paramName.Content
		}

		valueUnit := &ValueUnit{
			Type: typeName.Content,
			Name: paramNameContent,
		}
		funcUnit.Parameters = append(funcUnit.Parameters, valueUnit)
	}

	// returns
	// never nil
	retParams := core.FindFirstByFieldInSubsWithDfs(unit, FieldGolangParameters)
	switch retParams.Kind {
	case KindGolangParameterList:
		// multi params
		for _, each := range core.FindAllByKindInSubsWithDfs(retParams, KindGolangParameterDecl) {
			typeName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangType)
			paramName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangName)
			var paramNameContent string
			if paramName == nil {
				paramNameContent = ""
			} else {
				paramNameContent = paramName.Content
			}
			valueUnit := &ValueUnit{
				Type: typeName.Content,
				Name: paramNameContent,
			}
			funcUnit.Returns = append(funcUnit.Returns, valueUnit)
		}
	case KindGolangTypeIdentifier:
		// only one param, and anonymous
		valueUnit := &ValueUnit{
			Type: retParams.Content,
			Name: "",
		}
		funcUnit.Returns = append(funcUnit.Returns, valueUnit)
	default:
		// no returns
	}
	// extras
	funcUnit.Extras = &GolangFuncExtras{}

	return funcUnit, nil
}

func (extractor *GolangExtractor) funcUnit2Function(unit *core.Unit) (*Function, error) {
	funcUnit := &Function{}
	funcUnit.Span = unit.Span
	funcUnit.unit = unit
	// body scope
	funcBody := core.FindFirstByFieldInSubsWithBfs(unit, FieldGolangResult)
	if funcBody != nil {
		funcUnit.BodySpan = funcBody.Span
	}

	// name
	funcIdentifier := core.FindFirstByKindInSubsWithDfs(unit, KindGolangIdentifier)
	if funcIdentifier == nil {
		return nil, errors.New("no func name found in " + unit.Content)
	}
	funcUnit.Name = funcIdentifier.Content

	// params
	// no param == empty slice, never nil
	paramList := core.FindFirstByKindInSubsWithDfs(unit, KindGolangParameterList)
	for _, each := range core.FindAllByKindInSubsWithDfs(paramList, KindGolangParameterDecl) {
		typeName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangType)
		paramName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangName)
		var paramNameContent string
		if paramName == nil {
			paramNameContent = ""
		} else {
			paramNameContent = paramName.Content
		}
		valueUnit := &ValueUnit{
			Type: typeName.Content,
			Name: paramNameContent,
		}
		funcUnit.Parameters = append(funcUnit.Parameters, valueUnit)
	}

	// returns
	// never nil
	retParams := core.FindFirstByFieldInSubsWithDfs(unit, FieldGolangParameters)
	switch retParams.Kind {
	case KindGolangParameterList:
		// multi params
		for _, each := range core.FindAllByKindInSubsWithDfs(retParams, KindGolangParameterDecl) {
			typeName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangType)
			paramName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangName)
			var paramNameContent string
			if paramName == nil {
				paramNameContent = ""
			} else {
				paramNameContent = paramName.Content
			}
			valueUnit := &ValueUnit{
				Type: typeName.Content,
				Name: paramNameContent,
			}
			funcUnit.Returns = append(funcUnit.Returns, valueUnit)
		}
	case KindGolangTypeIdentifier:
		// only one param, and anonymous
		valueUnit := &ValueUnit{
			Type: retParams.Content,
			Name: "",
		}
		funcUnit.Returns = append(funcUnit.Returns, valueUnit)
	default:
		// no returns
	}

	// extras
	funcUnit.Extras = &GolangFuncExtras{}

	return funcUnit, nil
}
