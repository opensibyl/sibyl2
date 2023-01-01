package java

import (
	"errors"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

// FunctionExtras JavaFunctionExtras
type FunctionExtras struct {
	Annotations []string   `json:"annotations"`
	ClassInfo   *ClassInfo `json:"classInfo"`
}

type ClassInfo struct {
	PackageName string   `json:"packageName"`
	ClassName   string   `json:"className"`
	Annotations []string `json:"annotations"`
}

func (extractor *Extractor) IsFunction(unit *core.Unit) bool {
	// no function in java
	if unit.Kind == KindJavaMethodDeclaration {
		return true
	}
	return false
}

func (extractor *Extractor) ExtractFunctions(units []*core.Unit) ([]*object.Function, error) {
	var ret []*object.Function
	for _, eachUnit := range units {
		if !extractor.IsFunction(eachUnit) {
			continue
		}
		eachFunc, err := extractor.ExtractFunction(eachUnit)
		if err != nil {
			return nil, err
		}
		ret = append(ret, eachFunc)
	}
	return ret, nil
}

func (extractor *Extractor) ExtractFunction(unit *core.Unit) (*object.Function, error) {
	funcUnit := object.NewFunction()
	funcUnit.Span = unit.Span
	funcUnit.Unit = unit
	// body scope
	funcBody := core.FindFirstByKindInSubsWithBfs(unit, KindJavaBlock)
	if funcBody != nil {
		funcUnit.BodySpan = funcBody.Span
	}

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
	valueUnit := &object.ValueUnit{
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
			valueUnit = &object.ValueUnit{
				Type: typeName.Content,
				Name: paramName.Content,
			}
			funcUnit.Parameters = append(funcUnit.Parameters, valueUnit)
		}
	}

	// extras
	extras := &FunctionExtras{}
	classInfo := &ClassInfo{
		PackageName: pkgName,
		ClassName:   clazzName,
		Annotations: nil,
	}
	extras.ClassInfo = classInfo

	// class annotations
	classModifiers := core.FindFirstByKindInSubsWithBfs(clazzDecl, KindJavaModifiers)
	if classModifiers != nil {
		classAnnotations := core.FindAllByKindsInSubs(classModifiers, KindJavaMarkerAnnotation, KindJavaAnnotation)
		if len(classAnnotations) != 0 {
			for _, each := range classAnnotations {
				classInfo.Annotations = append(classInfo.Annotations, each.Content)
			}
		}
	}
	// todo: inherit

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
