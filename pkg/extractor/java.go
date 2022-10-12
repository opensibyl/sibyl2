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
	KindJavaMethodInvocation     model.KindRepr = "method_invocation"
	FieldJavaType                model.KindRepr = "type"
	FieldJavaDimensions          model.KindRepr = "dimensions"
	FieldJavaObject              model.KindRepr = "object"
	FieldJavaName                model.KindRepr = "name"
	FieldJavaArguments           model.KindRepr = "arguments"
)

type JavaExtractor struct {
}

func (extractor *JavaExtractor) GetLang() model.LangType {
	return model.LangJava
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

func (extractor *JavaExtractor) ExtractFunction(unit *model.Unit) (*model.Function, error) {
	data, err := extractor.ExtractFunctions([]*model.Unit{unit})
	if len(data) == 0 {
		return nil, err
	}
	return data[0], nil
}

func (extractor *JavaExtractor) IsCall(unit *model.Unit) bool {
	if unit.Kind == KindJavaMethodInvocation {
		return true
	}
	return false
}

func (extractor *JavaExtractor) ExtractCalls(units []*model.Unit) ([]*model.Call, error) {
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

func (extractor *JavaExtractor) unit2Call(unit *model.Unit) (*model.Call, error) {
	funcUnit := model.FindFirstByOneOfKindInParent(unit, KindJavaMethodDeclaration)
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

	var argumentPart *model.Unit
	var arguments []string
	var caller string

	callerPart := model.FindFirstByFieldInSubs(unit, FieldJavaObject)
	if callerPart == nil {
		// b()
		callerPart = model.FindFirstByFieldInSubs(unit, FieldJavaName)
		argumentPart = model.FindFirstByFieldInSubs(unit, FieldJavaArguments)
		caller = callerPart.Content
	} else {
		// a.b()
		identifiers := model.FindAllByKindInSubs(unit, KindJavaIdentifier)
		argumentPart = model.FindFirstByFieldInSubs(unit, FieldJavaName)

		if len(identifiers) == 0 {
			core.DebugDfs(unit, 0)
			return nil, errors.New("no id: " + unit.Content)
		}

		caller = callerPart.Content + "." + identifiers[len(identifiers)-1].Content
	}

	// not perfect, eg: anonymous function call?
	if argumentPart != nil {
		for _, each := range model.FindAllByKindInSubs(argumentPart, KindJavaIdentifier) {
			arguments = append(arguments, each.Content)
		}
	}

	ret := &model.Call{
		Src:       srcFunc.GetSignature(),
		Caller:    caller,
		Arguments: arguments,
		Span:      unit.Span,
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
		core.Log.Warnf("no package found in %s", unit.Content)
	} else {
		pkgName = packageIdentifier.Content
	}

	// trace its class (the closest one
	clazzDecl := model.FindFirstByOneOfKindInParent(unit, KindJavaClassDeclaration, KindJavaEnumDeclaration, KindJavaInterfaceDeclaration)
	clazzIdentifier := model.FindFirstByKindInSubsWithBfs(clazzDecl, KindJavaIdentifier)
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
		// java has no named return value
		Name: "",
	}
	funcUnit.Returns = append(funcUnit.Returns, valueUnit)

	// params
	parameters := model.FindFirstByKindInSubsWithDfs(unit, KindJavaFormalParameters)
	if parameters != nil {
		for _, each := range model.FindAllByKindInSubsWithDfs(parameters, KindJavaFormalParameter) {
			typeName := model.FindFirstByFieldInSubsWithBfs(each, FieldJavaType)
			paramName := model.FindFirstByFieldInSubsWithBfs(each, FieldJavaDimensions)
			valueUnit = &model.ValueUnit{
				Type: typeName.Content,
				Name: paramName.Content,
			}
			funcUnit.Parameters = append(funcUnit.Parameters, valueUnit)
		}
	}

	return funcUnit, nil
}
