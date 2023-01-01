package python

import (
	"errors"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

func (extractor *Extractor) IsClass(unit *core.Unit) bool {
	//TODO implement me
	return false
}

func (extractor *Extractor) ExtractClasses(units []*core.Unit) ([]*object.Clazz, error) {
	//TODO implement me
	return nil, errors.New("implement me")
}

func (extractor *Extractor) ExtractClass(unit *core.Unit) (*object.Clazz, error) {
	//TODO implement me
	return nil, errors.New("implement me")
}
