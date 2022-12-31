package extractor

import (
	"errors"

	"github.com/opensibyl/sibyl2/pkg/core"
)

func (extractor *PythonExtractor) IsCall(_ *core.Unit) bool {
	//TODO implement me
	return false
}

func (extractor *PythonExtractor) ExtractCalls(_ []*core.Unit) ([]*Call, error) {
	//TODO implement me
	return nil, errors.New("NOT IMPLEMENTED")
}
