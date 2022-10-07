package extractor

import (
	"errors"
	"sibyl2/pkg/core"
	"strings"
)

// https://github.com/tree-sitter/tree-sitter-java/tree/master/src
const (
	KindJavaProgram            core.KindRepr = "program"
	KindJavaProgramDeclaration core.KindRepr = "package_declaration"
	KindJavaScopeIdentifier    core.KindRepr = "scoped_identifier"
	KindJavaIdentifier         core.KindRepr = "identifier"
	KindJavaClassDeclaration   core.KindRepr = "class_declaration"
	KindJavaMethodDeclaration  core.KindRepr = "method_declaration"
	KindJavaFormalParameters   core.KindRepr = "formal_parameters"
	KindJavaFormalParameter    core.KindRepr = "formal_parameter"
	FieldJavaType              core.KindRepr = "type"
	FieldJavaDimensions        core.KindRepr = "dimensions"
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
	if unit.Kind == KindJavaMethodDeclaration {
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
		ret = append(ret, eachFunc)
	}
	return ret, nil
}

func (extractor *JavaExtractor) unit2Function(unit *core.Unit) (*core.Function, error) {
	funcUnit := &core.Function{}
	funcUnit.Span = unit.Span

	pkgName := ""
	clazzName := ""

	// trace its package
	program := core.FindFirstByKindInParent(unit, KindJavaProgram)
	packageDecl := core.FindFirstByKindInSubsWithDfs(program, KindJavaProgramDeclaration)
	packageIdentifier := core.FindFirstByKindInSubsWithDfs(packageDecl, KindJavaScopeIdentifier)
	if packageIdentifier == nil {
		return nil, errors.New("no package found in " + unit.Content)
	}
	pkgName = packageIdentifier.Content

	// trace its class (the closest one
	clazzDecl := core.FindFirstByKindInParent(unit, KindJavaClassDeclaration)
	clazzIdentifier := core.FindFirstByKindInSubsWithDfs(clazzDecl, KindJavaIdentifier)
	if clazzIdentifier == nil {
		return nil, errors.New("no class found in " + unit.Content)
	}
	clazzName = clazzIdentifier.Content
	funcUnit.Receiver = pkgName + "." + clazzName

	funcIdentifier := core.FindFirstByKindInSubsWithBfs(unit, KindJavaIdentifier)
	if funcIdentifier == nil {
		return nil, errors.New("no func id found in identifier" + unit.Content)
	}
	funcUnit.Name = funcIdentifier.Content

	// returns
	retUnit := core.FindFirstByFieldInSubsWithDfs(unit, FieldJavaDimensions)
	valueUnit := &core.ValueUnit{
		Type: retUnit.Content,
		Name: "",
	}
	funcUnit.Returns = append(funcUnit.Returns, valueUnit)

	// params
	parameters := core.FindFirstByKindInSubsWithDfs(unit, KindJavaFormalParameters)
	if parameters != nil {
		for _, each := range core.FindAllByKindInSubsWithDfs(parameters, KindJavaFormalParameter) {
			typeName := core.FindFirstByFieldInSubsWithDfs(each, FieldJavaType)
			paramName := core.FindFirstByFieldInSubsWithDfs(each, FieldJavaDimensions)
			valueUnit = &core.ValueUnit{
				Type: typeName.Content,
				Name: paramName.Content,
			}
			funcUnit.Parameters = append(funcUnit.Parameters, valueUnit)
		}
	}

	return funcUnit, nil
}
