package extractor

import (
	"errors"
	"sibyl2/pkg/core"
	"sibyl2/pkg/model"
	"strings"

	"golang.org/x/exp/slices"
)

// https://github.com/tree-sitter/tree-sitter-go/blob/master/src/node-types.json
const (
	KindGolangMethodDecl      model.KindRepr = "method_declaration"
	KindGolangFuncDecl        model.KindRepr = "function_declaration"
	KindGolangIdentifier      model.KindRepr = "identifier"
	KindGolangFieldIdentifier model.KindRepr = "field_identifier"
	KindGolangTypeIdentifier  model.KindRepr = "type_identifier"
	KindGolangParameterList   model.KindRepr = "parameter_list"
	KindGolangParameterDecl   model.KindRepr = "parameter_declaration"
	KindGolangCallExpression  model.KindRepr = "call_expression"
	FieldGolangType           model.KindRepr = "type"
	FieldGolangName           model.KindRepr = "name"
	FieldGolangParameters     model.KindRepr = "parameters"
	FieldGolangFunction       model.KindRepr = "function"
	FieldGolangArguments      model.KindRepr = "arguments"
)

type GolangExtractor struct {
}

func (extractor *GolangExtractor) GetLang() model.LangType {
	return model.LangGo
}

func (extractor *GolangExtractor) IsSymbol(unit *model.Unit) bool {
	// todo: use grammar.js instead
	if strings.HasSuffix(unit.Kind, "identifier") {
		return true
	}
	return false
}

func (extractor *GolangExtractor) ExtractSymbols(unit []*model.Unit) ([]*model.Symbol, error) {
	var ret []*model.Symbol
	for _, eachUnit := range unit {
		if !extractor.IsSymbol(eachUnit) {
			continue
		}
		symbol := &model.Symbol{
			Symbol:    eachUnit.Content,
			Kind:      eachUnit.Kind,
			Span:      eachUnit.Span,
			FieldName: eachUnit.FieldName,
		}
		ret = append(ret, symbol)
	}
	return ret, nil
}

func (extractor *GolangExtractor) IsFunction(unit *model.Unit) bool {
	allowed := []model.KindRepr{
		KindGolangMethodDecl,
		KindGolangFuncDecl,
	}
	return slices.Contains(allowed, unit.Kind)
}

func (extractor *GolangExtractor) ExtractFunctions(units []*model.Unit) ([]*model.Function, error) {
	var ret []*model.Function
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

func (extractor *GolangExtractor) ExtractFunction(unit *model.Unit) (*model.Function, error) {
	data, err := extractor.ExtractFunctions([]*model.Unit{unit})
	if len(data) == 0 {
		return nil, err
	}
	return data[0], nil
}

func (extractor *GolangExtractor) IsCall(unit *model.Unit) bool {
	if unit.Kind == KindGolangCallExpression {
		return true
	}
	return false
}

func (extractor *GolangExtractor) ExtractCalls(units []*model.Unit) ([]*model.Call, error) {
	var ret []*model.Call
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

func (extractor *GolangExtractor) unit2Call(unit *model.Unit) (*model.Call, error) {
	// todo: what about nested call
	funcUnit := model.FindFirstByOneOfKindInParent(unit, KindGolangFuncDecl, KindGolangMethodDecl)
	var srcFunc *model.Function
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

	funcPart := model.FindFirstByFieldInSubsWithBfs(unit, FieldGolangFunction)
	argumentPart := model.FindFirstByFieldInSubsWithBfs(unit, FieldGolangArguments)

	// not perfect, eg: anonymous function call?
	var arguments []string
	for _, each := range argumentPart.SubUnits {
		if each.Kind == KindGolangIdentifier {
			arguments = append(arguments, each.Content)
		}
	}

	ret := &model.Call{
		Src:       srcFunc.GetSignature(),
		Caller:    funcPart.Content,
		Arguments: arguments,
	}
	return ret, nil
}

func (extractor *GolangExtractor) unit2Function(unit *model.Unit) (*model.Function, error) {
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

func (extractor *GolangExtractor) methodUnit2Function(unit *model.Unit) (*model.Function, error) {
	funcUnit := &model.Function{}
	funcUnit.Span = unit.Span

	// name
	funcIdentifier := model.FindFirstByKindInSubsWithDfs(unit, KindGolangFieldIdentifier)
	if funcIdentifier == nil {
		return nil, errors.New("no func name found in " + unit.Content)
	}
	funcUnit.Name = funcIdentifier.Content

	// receiver
	parameterList := model.FindFirstByKindInSubsWithDfs(unit, KindGolangParameterList)
	parameterList = model.FindFirstByKindInSubsWithDfs(parameterList, KindGolangParameterList)
	receiverDecl := model.FindFirstByKindInSubsWithDfs(parameterList, KindGolangParameterDecl)
	typeDecl := model.FindFirstByFieldInSubsWithDfs(receiverDecl, FieldGolangType)
	if typeDecl == nil {
		return nil, errors.New("no receiver found in: " + typeDecl.Content)
	}
	funcUnit.Receiver = typeDecl.Content

	// params
	paramListList := model.FindAllByKindInSubsWithDfs(unit, KindGolangParameterList)
	// no param == empty slice, never nil
	paramList := paramListList[1]
	for _, each := range model.FindAllByKindInSubsWithDfs(paramList, KindGolangParameterDecl) {
		typeName := model.FindFirstByFieldInSubsWithDfs(each, FieldGolangType)
		paramName := model.FindFirstByFieldInSubsWithDfs(each, FieldGolangName)
		var paramNameContent string
		if paramName == nil {
			paramNameContent = ""
		} else {
			paramNameContent = paramName.Content
		}

		valueUnit := &model.ValueUnit{
			Type: typeName.Content,
			Name: paramNameContent,
		}
		funcUnit.Parameters = append(funcUnit.Parameters, valueUnit)
	}

	// returns
	// never nil
	retParams := model.FindFirstByFieldInSubsWithDfs(unit, FieldGolangParameters)
	switch retParams.Kind {
	case KindGolangParameterList:
		// multi params
		for _, each := range model.FindAllByKindInSubsWithDfs(retParams, KindGolangParameterDecl) {
			typeName := model.FindFirstByFieldInSubsWithDfs(each, FieldGolangType)
			paramName := model.FindFirstByFieldInSubsWithDfs(each, FieldGolangName)
			var paramNameContent string
			if paramName == nil {
				paramNameContent = ""
			} else {
				paramNameContent = paramName.Content
			}
			valueUnit := &model.ValueUnit{
				Type: typeName.Content,
				Name: paramNameContent,
			}
			funcUnit.Returns = append(funcUnit.Returns, valueUnit)
		}
	case KindGolangTypeIdentifier:
		// only one param, and anonymous
		valueUnit := &model.ValueUnit{
			Type: retParams.Content,
			Name: "",
		}
		funcUnit.Returns = append(funcUnit.Returns, valueUnit)
	default:
		// no returns
	}

	return funcUnit, nil
}

func (extractor *GolangExtractor) funcUnit2Function(unit *model.Unit) (*model.Function, error) {
	funcUnit := &model.Function{}
	funcUnit.Span = unit.Span

	// name
	funcIdentifier := model.FindFirstByKindInSubsWithDfs(unit, KindGolangIdentifier)
	if funcIdentifier == nil {
		return nil, errors.New("no func name found in " + unit.Content)
	}
	funcUnit.Name = funcIdentifier.Content

	// params
	// no param == empty slice, never nil
	paramList := model.FindFirstByKindInSubsWithDfs(unit, KindGolangParameterList)
	for _, each := range model.FindAllByKindInSubsWithDfs(paramList, KindGolangParameterDecl) {
		typeName := model.FindFirstByFieldInSubsWithDfs(each, FieldGolangType)
		paramName := model.FindFirstByFieldInSubsWithDfs(each, FieldGolangName)
		var paramNameContent string
		if paramName == nil {
			paramNameContent = ""
		} else {
			paramNameContent = paramName.Content
		}
		valueUnit := &model.ValueUnit{
			Type: typeName.Content,
			Name: paramNameContent,
		}
		funcUnit.Parameters = append(funcUnit.Parameters, valueUnit)
	}

	// returns
	// never nil
	retParams := model.FindFirstByFieldInSubsWithDfs(unit, FieldGolangParameters)
	switch retParams.Kind {
	case KindGolangParameterList:
		// multi params
		for _, each := range model.FindAllByKindInSubsWithDfs(retParams, KindGolangParameterDecl) {
			typeName := model.FindFirstByFieldInSubsWithDfs(each, FieldGolangType)
			paramName := model.FindFirstByFieldInSubsWithDfs(each, FieldGolangName)
			var paramNameContent string
			if paramName == nil {
				paramNameContent = ""
			} else {
				paramNameContent = paramName.Content
			}
			valueUnit := &model.ValueUnit{
				Type: typeName.Content,
				Name: paramNameContent,
			}
			funcUnit.Returns = append(funcUnit.Returns, valueUnit)
		}
	case KindGolangTypeIdentifier:
		// only one param, and anonymous
		valueUnit := &model.ValueUnit{
			Type: retParams.Content,
			Name: "",
		}
		funcUnit.Returns = append(funcUnit.Returns, valueUnit)
	default:
		// no returns
	}

	return funcUnit, nil
}
