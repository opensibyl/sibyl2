package python

import (
	"errors"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

type ClassExtras struct {
	Decorators []string `json:"decorators"`
}

func (extractor *Extractor) IsClass(unit *core.Unit) bool {
	if unit.Kind == KindPythonClassDefinition {
		return true
	}
	return false
}

func (extractor *Extractor) ExtractClasses(units []*core.Unit) ([]*object.Clazz, error) {
	var ret []*object.Clazz
	for _, eachUnit := range units {
		if !extractor.IsClass(eachUnit) {
			continue
		}
		eachClazz, err := extractor.ExtractClass(eachUnit)
		if err != nil {
			return nil, err
		}
		ret = append(ret, eachClazz)
	}
	return ret, nil
}

func (extractor *Extractor) ExtractClass(unit *core.Unit) (*object.Clazz, error) {
	clazz := object.NewClazz()
	clazz.Span = unit.Span
	clazz.Lang = extractor.GetLang()
	clazz.Unit = unit

	clazzName := core.FindFirstByKindInSubsWithBfs(unit, KindPythonIdentifier)
	if clazzName == nil {
		return nil, errors.New("class name not found")
	}
	clazz.Name = clazzName.Content

	extras := &ClassExtras{}
	// deco
	parentUnit := unit.ParentUnit
	if parentUnit.Kind == KindPythonDecoratedDefinition {
		decorators := core.FindAllByKindInSubs(parentUnit, KindPythonDecorator)
		if len(decorators) == 0 {
			return nil, errors.New("no deco found in " + parentUnit.Content)
		}

		var decoContents []string
		for _, each := range decorators {
			decoContents = append(decoContents, each.Content)
		}
		extras.Decorators = decoContents
	}
	clazz.Extras = extras

	return clazz, nil
}
