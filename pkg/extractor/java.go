package extractor

import (
	"errors"
	"github.com/williamfzc/sibyl2/pkg/core"
	"strings"
)

// https://github.com/tree-sitter/tree-sitter-java/tree/master/src
const (
	KindJavaProgram              core.KindRepr = "program"
	KindJavaProgramDeclaration   core.KindRepr = "package_declaration"
	KindJavaScopeIdentifier      core.KindRepr = "scoped_identifier"
	KindJavaIdentifier           core.KindRepr = "identifier"
	KindJavaClassDeclaration     core.KindRepr = "class_declaration"
	KindJavaEnumDeclaration      core.KindRepr = "enum_declaration"
	KindJavaInterfaceDeclaration core.KindRepr = "interface_declaration"
	KindJavaMethodDeclaration    core.KindRepr = "method_declaration"
	KindJavaFormalParameters     core.KindRepr = "formal_parameters"
	KindJavaFormalParameter      core.KindRepr = "formal_parameter"
	KindJavaMethodInvocation     core.KindRepr = "method_invocation"
	FieldJavaType                core.KindRepr = "type"
	FieldJavaDimensions          core.KindRepr = "dimensions"
	FieldJavaObject              core.KindRepr = "object"
	FieldJavaName                core.KindRepr = "name"
	FieldJavaArguments           core.KindRepr = "arguments"
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

func (extractor *JavaExtractor) ExtractSymbols(units []*core.Unit) ([]*Symbol, error) {
	var ret []*Symbol
	for _, eachUnit := range units {
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

func (extractor *JavaExtractor) IsFunction(unit *core.Unit) bool {
	// no function in java
	if unit.Kind == KindJavaMethodDeclaration {
		return true
	}
	return false
}

func (extractor *JavaExtractor) ExtractFunctions(units []*core.Unit) ([]*Function, error) {
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

func (extractor *JavaExtractor) ExtractFunction(unit *core.Unit) (*Function, error) {
	data, err := extractor.ExtractFunctions([]*core.Unit{unit})
	if len(data) == 0 {
		return nil, err
	}
	return data[0], nil
}

func (extractor *JavaExtractor) IsCall(unit *core.Unit) bool {
	if unit.Kind == KindJavaMethodInvocation {
		return true
	}
	return false
}

func (extractor *JavaExtractor) ExtractCalls(units []*core.Unit) ([]*Call, error) {
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

func (extractor *JavaExtractor) unit2Call(unit *core.Unit) (*Call, error) {
	funcUnit := core.FindFirstByOneOfKindInParent(unit, KindJavaMethodDeclaration)
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

	var argumentPart *core.Unit
	var arguments []string
	var caller string

	callerPart := core.FindFirstByFieldInSubs(unit, FieldJavaObject)
	if callerPart == nil {
		// b()
		callerPart = core.FindFirstByFieldInSubs(unit, FieldJavaName)
		argumentPart = core.FindFirstByFieldInSubs(unit, FieldJavaArguments)
		caller = callerPart.Content
	} else {
		// a.b()
		identifiers := core.FindAllByKindInSubs(unit, KindJavaIdentifier)
		argumentPart = core.FindFirstByFieldInSubs(unit, FieldJavaName)

		if len(identifiers) == 0 {
			core.DebugDfs(unit, 0)
			return nil, errors.New("no id: " + unit.Content)
		}

		caller = callerPart.Content + "." + identifiers[len(identifiers)-1].Content
	}

	// not perfect, eg: anonymous function call?
	if argumentPart != nil {
		for _, each := range core.FindAllByKindInSubs(argumentPart, KindJavaIdentifier) {
			arguments = append(arguments, each.Content)
		}
	}

	ret := &Call{
		Src:       srcFunc.GetSignature(),
		Caller:    caller,
		Arguments: arguments,
		Span:      unit.Span,
	}
	return ret, nil
}

func (extractor *JavaExtractor) unit2Function(unit *core.Unit) (*Function, error) {
	funcUnit := &Function{}
	funcUnit.Span = unit.Span

	pkgName := ""
	clazzName := ""

	// trace its package
	program := core.FindFirstByKindInParent(unit, KindJavaProgram)
	packageDecl := core.FindFirstByKindInSubsWithDfs(program, KindJavaProgramDeclaration)
	packageIdentifier := core.FindFirstByKindInSubsWithDfs(packageDecl, KindJavaScopeIdentifier)
	if packageIdentifier == nil {
		core.Log.Warnf("no package found in %s", unit.Content)
	} else {
		pkgName = packageIdentifier.Content
	}

	// trace its class (the closest one
	clazzDecl := core.FindFirstByOneOfKindInParent(unit, KindJavaClassDeclaration, KindJavaEnumDeclaration, KindJavaInterfaceDeclaration)
	clazzIdentifier := core.FindFirstByKindInSubsWithBfs(clazzDecl, KindJavaIdentifier)
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
	valueUnit := &ValueUnit{
		Type: retUnit.Content,
		// java has no named return value
		Name: "",
	}
	funcUnit.Returns = append(funcUnit.Returns, valueUnit)

	// params
	parameters := core.FindFirstByKindInSubsWithDfs(unit, KindJavaFormalParameters)
	if parameters != nil {
		for _, each := range core.FindAllByKindInSubsWithDfs(parameters, KindJavaFormalParameter) {
			typeName := core.FindFirstByFieldInSubsWithBfs(each, FieldJavaType)
			paramName := core.FindFirstByFieldInSubsWithBfs(each, FieldJavaDimensions)
			valueUnit = &ValueUnit{
				Type: typeName.Content,
				Name: paramName.Content,
			}
			funcUnit.Parameters = append(funcUnit.Parameters, valueUnit)
		}
	}

	return funcUnit, nil
}
