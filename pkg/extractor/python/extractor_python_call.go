package python

import (
	"errors"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

func (extractor *Extractor) IsCall(_ *core.Unit) bool {
	//TODO implement me
	return false
}

func (extractor *Extractor) ExtractCalls(_ []*core.Unit) ([]*object.Call, error) {
	//TODO implement me
	return nil, errors.New("NOT IMPLEMENTED")
}
