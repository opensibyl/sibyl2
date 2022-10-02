package extractor

import (
	"errors"
	"sibyl2/pkg/core"
	"strings"

	"golang.org/x/exp/slices"
)

// https://github.com/tree-sitter/tree-sitter-go/blob/master/src/node-types.json
const (
	KindGolangMethodDecl      core.KindRepr = "method_declaration"
	KindGolangFuncDecl        core.KindRepr = "function_declaration"
	KindGolangIdentifier      core.KindRepr = "identifier"
	KindGolangFieldIdentifier core.KindRepr = "field_identifier"
	KindGolangParameterList   core.KindRepr = "parameter_list"
	KindGolangParameterDecl   core.KindRepr = "parameter_declaration"
	FieldGolangType           core.KindRepr = "type"
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

func (extractor *GolangExtractor) ExtractSymbols(unit []*core.Unit) ([]*core.Symbol, error) {
	var ret []*core.Symbol
	for _, eachUnit := range unit {
		if !extractor.IsSymbol(eachUnit) {
			continue
		}
		symbol := &core.Symbol{
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

func (extractor *GolangExtractor) ExtractFunctions(units []*core.Unit) ([]*core.Function, error) {
	var ret []*core.Function
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

func (extractor *GolangExtractor) unit2Function(unit *core.Unit) (*core.Function, error) {
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

func (extractor *GolangExtractor) methodUnit2Function(unit *core.Unit) (*core.Function, error) {
	funcUnit := &core.Function{}
	funcUnit.Span = unit.Span

	funcIdentifier := core.FindFirstByKindInSubsWithDfs(unit, KindGolangFieldIdentifier)
	if funcIdentifier == nil {
		return nil, errors.New("no func name found in " + unit.Content)
	}
	funcUnit.Name = funcIdentifier.Content

	parameterList := core.FindFirstByKindInSubsWithDfs(unit, KindGolangParameterList)
	parameterList = core.FindFirstByKindInSubsWithDfs(parameterList, KindGolangParameterList)
	receiverDecl := core.FindFirstByKindInSubsWithDfs(parameterList, KindGolangParameterDecl)
	typeDecl := core.FindFirstByFieldInSubsWithDfs(receiverDecl, FieldGolangType)
	if typeDecl == nil {
		return nil, errors.New("no receiver found in: " + typeDecl.Content)
	}
	funcUnit.Receiver = typeDecl.Content
	return funcUnit, nil
}

func (extractor *GolangExtractor) funcUnit2Function(unit *core.Unit) (*core.Function, error) {
	funcUnit := &core.Function{}
	funcUnit.Span = unit.Span
	funcIdentifier := core.FindFirstByKindInSubsWithDfs(unit, KindGolangIdentifier)
	if funcIdentifier == nil {
		return nil, errors.New("no func name found in " + unit.Content)
	}
	funcUnit.Name = funcIdentifier.Content
	return funcUnit, nil
}
