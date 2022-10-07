package extractor

import (
	"errors"
	"sibyl2/pkg/core"
	"sibyl2/pkg/model"
	"strings"
)

// https://github.com/tree-sitter/tree-sitter-java/tree/master/src
const (
	KindJavaProgram              model.KindRepr = "program"
	KindJavaProgramDeclaration   model.KindRepr = "package_declaration"
	KindJavaScopeIdentifier      model.KindRepr = "scoped_identifier"
	KindJavaIdentifier           model.KindRepr = "identifier"
	KindJavaClassDeclaration     model.KindRepr = "class_declaration"
	KindJavaEnumDeclaration      model.KindRepr = "enum_declaration"
	KindJavaInterfaceDeclaration model.KindRepr = "interface_declaration"
	KindJavaMethodDeclaration    model.KindRepr = "method_declaration"
	KindJavaFormalParameters     model.KindRepr = "formal_parameters"
	KindJavaFormalParameter      model.KindRepr = "formal_parameter"
	FieldJavaType                model.KindRepr = "type"
	FieldJavaDimensions          model.KindRepr = "dimensions"
)

type JavaExtractor struct {
}

func (extractor *JavaExtractor) GetLang() core.LangType {
	return core.LangJava
}

func (extractor *JavaExtractor) IsSymbol(unit *model.Unit) bool {
	// todo: use grammar.js instead
	if strings.HasSuffix(unit.Kind, "identifier") {
		return true
	}
	return false
}

func (extractor *JavaExtractor) ExtractSymbols(units []*model.Unit) ([]*model.Symbol, error) {
	var ret []*model.Symbol
	for _, eachUnit := range units {
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

func (extractor *JavaExtractor) IsFunction(unit *model.Unit) bool {
	// no function in java
	if unit.Kind == KindJavaMethodDeclaration {
		return true
	}
	return false
}

func (extractor *JavaExtractor) ExtractFunctions(units []*model.Unit) ([]*model.Function, error) {
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

func (extractor *JavaExtractor) unit2Function(unit *model.Unit) (*model.Function, error) {
	funcUnit := &model.Function{}
	funcUnit.Span = unit.Span

	pkgName := ""
	clazzName := ""

	// trace its package
	program := model.FindFirstByKindInParent(unit, KindJavaProgram)
	packageDecl := model.FindFirstByKindInSubsWithDfs(program, KindJavaProgramDeclaration)
	packageIdentifier := model.FindFirstByKindInSubsWithDfs(packageDecl, KindJavaScopeIdentifier)
	if packageIdentifier == nil {
		return nil, errors.New("no package found in " + unit.Content)
	}
	pkgName = packageIdentifier.Content

	// trace its class (the closest one
	clazzDecl := model.FindFirstByOneOfKindInParent(unit, KindJavaClassDeclaration, KindJavaEnumDeclaration, KindJavaInterfaceDeclaration)
	clazzIdentifier := model.FindFirstByKindInSubsWithDfs(clazzDecl, KindJavaIdentifier)
	if clazzIdentifier == nil {
		return nil, errors.New("no class found in " + unit.Content)
	}
	clazzName = clazzIdentifier.Content
	funcUnit.Receiver = pkgName + "." + clazzName

	funcIdentifier := model.FindFirstByKindInSubsWithBfs(unit, KindJavaIdentifier)
	if funcIdentifier == nil {
		return nil, errors.New("no func id found in identifier" + unit.Content)
	}
	funcUnit.Name = funcIdentifier.Content

	// returns
	retUnit := model.FindFirstByFieldInSubsWithDfs(unit, FieldJavaDimensions)
	valueUnit := &model.ValueUnit{
		Type: retUnit.Content,
		Name: "",
	}
	funcUnit.Returns = append(funcUnit.Returns, valueUnit)

	// params
	parameters := model.FindFirstByKindInSubsWithDfs(unit, KindJavaFormalParameters)
	if parameters != nil {
		for _, each := range model.FindAllByKindInSubsWithDfs(parameters, KindJavaFormalParameter) {
			typeName := model.FindFirstByFieldInSubsWithDfs(each, FieldJavaType)
			paramName := model.FindFirstByFieldInSubsWithDfs(each, FieldJavaDimensions)
			valueUnit = &model.ValueUnit{
				Type: typeName.Content,
				Name: paramName.Content,
			}
			funcUnit.Parameters = append(funcUnit.Parameters, valueUnit)
		}
	}

	return funcUnit, nil
}
