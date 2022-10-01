package extractor

import (
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
)

type GolangExtractor struct {
}

func (extractor *GolangExtractor) GetLang() core.LangType {
	return core.GOLANG
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
	// todo: should not parse again
	unitsInFunctions, err := extractor.GetLang().GetParser().ParseString(unit.Content)
	if err != nil {
		return nil, err
	}
	funcUnit := &core.Function{}
	funcUnit.Span = unit.Span
	if unit.Kind == KindGolangFuncDecl {
		for _, each := range unitsInFunctions {
			if each.Kind == KindGolangIdentifier {
				funcUnit.Name = each.Content
				break
			}
		}
	} else {
		for _, each := range unitsInFunctions {
			if each.Kind == KindGolangFieldIdentifier {
				funcUnit.Name = each.Content
				break
			}
		}
		for _, each := range unitsInFunctions {
			if each.Kind == KindGolangParameterList {
				unitsInReceiver, err := extractor.GetLang().GetParser().ParseString(each.Content)
				if err != nil {
					return nil, err
				}
				for _, eachUnitInReceiver := range unitsInReceiver {
					if eachUnitInReceiver.FieldName == "operator" {
						funcUnit.Receiver = eachUnitInReceiver.Content
						break
					}
				}
				break
			}
		}
	}
	// todo: parameters and returns
	return funcUnit, nil
}
