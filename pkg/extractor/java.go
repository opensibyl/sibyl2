package extractor

import (
	"fmt"
	"sibyl2/pkg/core"
	"strings"
)

const (
	KindJavaProgram            core.KindRepr = "program"
	KindJavaProgramDeclaration core.KindRepr = "package_declaration"
	KindJavaScopeIdentifier    core.KindRepr = "scoped_identifier"
	KindJavaIdentifier         core.KindRepr = "identifier"
)

type JavaExtractor struct {
}

func (extractor *JavaExtractor) GetLang() core.LangType {
	return core.LangJava
}

func (extractor *JavaExtractor) IsSymbol(unit *core.Unit) bool {
	// todo: use grammar.js instead
	if strings.HasSuffix(unit.Kind, "identifier") {
		return true
	}
	return false
}

func (extractor *JavaExtractor) ExtractSymbols(units []*core.Unit) ([]*core.Symbol, error) {
	var ret []*core.Symbol
	for _, eachUnit := range units {
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

func (extractor *JavaExtractor) IsFunction(unit *core.Unit) bool {
	// no function in java
	if unit.Kind == "method_declaration" {
		return true
	}
	return false
}

func (extractor *JavaExtractor) ExtractFunctions(units []*core.Unit) ([]*core.Function, error) {
	var ret []*core.Function
	for _, eachUnit := range units {
		if !extractor.IsFunction(eachUnit) {
			continue
		}
		eachFunc, err := extractor.unit2Function(eachUnit)
		if err != nil {
			return nil, err
		}
		fmt.Printf("func: %v\n", eachFunc)
		ret = append(ret, eachFunc)
	}
	return ret, nil
}

func (extractor *JavaExtractor) unit2Function(unit *core.Unit) (*core.Function, error) {
	// todo: its receiver should contain package name and class name
	funcUnit := &core.Function{}
	funcUnit.Span = unit.Span

	for _, each := range unit.ReverseLink() {
		if each.Kind == KindJavaProgram {
			unitsInProgram := each.Link()
			for _, eachUnitInProgram := range unitsInProgram {
				if eachUnitInProgram.Kind == KindJavaProgramDeclaration {
					unitsInPackageDecl := eachUnitInProgram.Link()
					for _, eachUnitInPackageDecl := range unitsInPackageDecl {
						if eachUnitInPackageDecl.Kind == KindJavaScopeIdentifier {
							funcUnit.Receiver = eachUnitInPackageDecl.Content
							break
						}
					}
					break
				}
			}
		}
	}

	unitsInFunctions := unit.Link()
	for _, each := range unitsInFunctions {
		if each.Kind == KindJavaIdentifier {
			funcUnit.Name = each.Content
			break
		}
	}
	return funcUnit, nil
}
