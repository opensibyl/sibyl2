package javascript

import (
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

func (extractor *Extractor) IsClass(unit *core.Unit) bool {
	if unit.Kind == KindJavaScriptClassDeclaration {
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

	nameNode := core.FindFirstByKindInSubsWithBfs(unit, KindJavaScriptIdentifier)
	if nameNode == nil {
		core.Log.Warnf("anonymous class: %v", unit)
	} else {
		clazz.Name = nameNode.Content
	}

	return clazz, nil
}
