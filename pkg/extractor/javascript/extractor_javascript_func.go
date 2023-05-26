package javascript

import (
	"fmt"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
	"golang.org/x/exp/slices"
)

func (extractor *Extractor) IsFunction(unit *core.Unit) bool {
	allowed := []core.KindRepr{
		KindJavaScriptFunctionDeclaration,
		KindJavaScriptMethodDefinition,
	}
	return slices.Contains(allowed, unit.Kind)
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
	funcUnit.Lang = extractor.GetLang()

	bodyUnit := core.FindFirstByKindInSubsWithBfs(unit, KindJavaScriptStatementBlock)
	if bodyUnit != nil {
		funcUnit.BodySpan = bodyUnit.Span
	}

	var err error
	switch unit.Kind {
	case KindJavaScriptFunctionDeclaration:
		err = extractor.extractFromFunc(unit, funcUnit)
	case KindJavaScriptMethodDefinition:
		err = extractor.extractFromMethod(unit, funcUnit)
	}
	if err != nil {
		return nil, err
	}

	return funcUnit, nil
}

func (extractor *Extractor) extractFromFunc(unit *core.Unit, function *object.Function) error {
	nameNode := core.FindFirstByKindInSubsWithBfs(unit, KindJavaScriptIdentifier)
	if nameNode == nil {
		return fmt.Errorf("function without name: %s", unit.Content)
	}
	function.Name = nameNode.Content
	function.DefLine = int(nameNode.Span.Start.Row + 1)

	// params
	parametersNode := core.FindFirstByFieldInSubs(unit, KindJavaScriptFormalParameters)
	if parametersNode != nil {
		for _, each := range core.FindAllByKindInSubs(parametersNode, KindJavaScriptIdentifier) {
			valueUnit := &object.ValueUnit{
				Type: "",
				Name: each.Content,
			}
			function.Parameters = append(function.Parameters, valueUnit)
		}
	}

	return nil
}

func (extractor *Extractor) extractFromMethod(unit *core.Unit, function *object.Function) error {
	nameNode := core.FindFirstByFieldInSubsWithBfs(unit, FieldJavaScriptName)
	if nameNode == nil {
		return fmt.Errorf("anonymous function: %s", unit.Content)
	}
	function.Name = nameNode.Content

	// trace its class (the closest one
	clazzDecl := core.FindFirstByOneOfKindInParent(unit, KindJavaScriptClassDeclaration)
	clazzNameNode := core.FindFirstByKindInSubsWithBfs(clazzDecl, KindJavaScriptIdentifier)
	if clazzNameNode == nil {
		core.Log.Warnf("anonymous class: %v", unit)
	} else {
		function.Receiver = clazzNameNode.Content
	}

	// params
	parametersNode := core.FindFirstByFieldInSubs(unit, FieldJavaScriptParameters)
	if parametersNode != nil {
		for _, each := range core.FindAllByKindInSubs(parametersNode, KindJavaScriptIdentifier) {
			valueUnit := &object.ValueUnit{
				Type: "",
				Name: each.Content,
			}
			function.Parameters = append(function.Parameters, valueUnit)
		}
	}
	return nil
}
