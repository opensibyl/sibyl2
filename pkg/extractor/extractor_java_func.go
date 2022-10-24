package extractor

import (
	"errors"

	"github.com/williamfzc/sibyl2/pkg/core"
)

type JavaFunctionExtras struct {
	Annotations []string `json:"annotations"`
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
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("unit not a func: " + unit.Content)
	}
	return data[0], nil
}

func (extractor *JavaExtractor) unit2Function(unit *core.Unit) (*Function, error) {
	funcUnit := NewFunction()
	funcUnit.Span = unit.Span
	funcUnit.unit = unit

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

	// extras
	extras := &JavaFunctionExtras{}
	modifiers := core.FindFirstByKindInSubsWithBfs(unit, KindJavaModifiers)
	if modifiers != nil {
		annotations := core.FindAllByKindsInSubs(modifiers, KindJavaMarkerAnnotation, KindJavaAnnotation)
		if len(annotations) != 0 {
			for _, each := range annotations {
				extras.Annotations = append(extras.Annotations, each.Content)
			}
		}
	}
	funcUnit.Extras = extras

	return funcUnit, nil
}
