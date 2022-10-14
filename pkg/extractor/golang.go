package extractor

import (
	"errors"
	"github.com/williamfzc/sibyl2/pkg/core"
	"strings"

	"golang.org/x/exp/slices"
)

// https://github.com/tree-sitter/tree-sitter-go/blob/master/src/node-types.json
const (
	KindGolangMethodDecl      core.KindRepr = "method_declaration"
	KindGolangFuncDecl        core.KindRepr = "function_declaration"
	KindGolangIdentifier      core.KindRepr = "identifier"
	KindGolangFieldIdentifier core.KindRepr = "field_identifier"
	KindGolangTypeIdentifier  core.KindRepr = "type_identifier"
	KindGolangParameterList   core.KindRepr = "parameter_list"
	KindGolangParameterDecl   core.KindRepr = "parameter_declaration"
	KindGolangCallExpression  core.KindRepr = "call_expression"
	FieldGolangType           core.KindRepr = "type"
	FieldGolangName           core.KindRepr = "name"
	FieldGolangParameters     core.KindRepr = "parameters"
	FieldGolangFunction       core.KindRepr = "function"
	FieldGolangArguments      core.KindRepr = "arguments"
)

type GolangExtractor struct {
}

func (extractor *GolangExtractor) GetLang() core.LangType {
	return core.LangGo
}

func (extractor *GolangExtractor) IsSymbol(unit *core.Unit) bool {
	// todo: use grammar.js instead
	if strings.HasSuffix(unit.Kind, "identifier") {
		return true
	}
	return false
}

func (extractor *GolangExtractor) ExtractSymbols(unit []*core.Unit) ([]*Symbol, error) {
	var ret []*Symbol
	for _, eachUnit := range unit {
		if !extractor.IsSymbol(eachUnit) {
			continue
		}
		symbol := &Symbol{
			Symbol:    eachUnit.Content,
			Kind:      eachUnit.Kind,
			Span:      eachUnit.Span,
			FieldName: eachUnit.FieldName,
		}
		ret = append(ret, symbol)
	}
	return ret, nil
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
	if len(data) == 0 {
		return nil, err
	}
	return data[0], nil
}

func (extractor *GolangExtractor) IsCall(unit *core.Unit) bool {
	if unit.Kind == KindGolangCallExpression {
		return true
	}
	return false
}

func (extractor *GolangExtractor) ExtractCalls(units []*core.Unit) ([]*Call, error) {
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

func (extractor *GolangExtractor) unit2Call(unit *core.Unit) (*Call, error) {
	// todo: what about nested call
	funcUnit := core.FindFirstByOneOfKindInParent(unit, KindGolangFuncDecl, KindGolangMethodDecl)
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

	funcPart := core.FindFirstByFieldInSubsWithBfs(unit, FieldGolangFunction)
	argumentPart := core.FindFirstByFieldInSubsWithBfs(unit, FieldGolangArguments)

	// not perfect, eg: anonymous function call?
	var arguments []string
	for _, each := range argumentPart.SubUnits {
		if each.Kind == KindGolangIdentifier {
			arguments = append(arguments, each.Content)
		}
	}

	ret := &Call{
		Src:       srcFunc.GetSignature(),
		Caller:    funcPart.Content,
		Arguments: arguments,
		Span:      unit.Span,
	}
	return ret, nil
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

	return funcUnit, nil
}

func (extractor *GolangExtractor) funcUnit2Function(unit *core.Unit) (*Function, error) {
	funcUnit := &Function{}
	funcUnit.Span = unit.Span

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

	return funcUnit, nil
}
