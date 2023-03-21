package object

import (
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
)

type FunctionServiceDTO struct {
	*FunctionWithTag `bson:",inline"`
	Signature        string `json:"signature" bson:"signature"`
}

type ClazzServiceDTO struct {
	*extractor.ClazzWithPath `bson:",inline"`
	Signature                string `json:"signature" bson:"signature"`
}

type FuncCtxServiceDTO struct {
	*FunctionContextSlim `bson:",inline"`
	Signature            string `json:"signature" bson:"signature"`
}

type FuncTag = string

type FunctionWithTag struct {
	*extractor.FunctionWithPath `bson:",inline"`
	Tags                        []FuncTag `json:"tags" bson:"tags"`
}

func (fwt *FunctionWithTag) AddTag(tag FuncTag) {
	fwt.Tags = append(fwt.Tags, tag)
}

// FunctionContextSlim instead of whole object, slim will only keep the signature
type FunctionContextSlim struct {
	*extractor.FunctionWithPath `bson:",inline"`
	Calls                       []string `json:"calls" bson:"calls"`
	ReverseCalls                []string `json:"reverseCalls" bson:"reverseCalls"`
}

func CompressFunctionContext(f *sibyl2.FunctionContext) *FunctionContextSlim {
	slim := &FunctionContextSlim{
		FunctionWithPath: f.FunctionWithPath,
		Calls:            make([]string, 0, len(f.Calls)),
		ReverseCalls:     make([]string, 0, len(f.ReverseCalls)),
	}
	for _, each := range f.Calls {
		slim.Calls = append(slim.Calls, each.GetSignature())
	}
	for _, each := range f.ReverseCalls {
		slim.ReverseCalls = append(slim.ReverseCalls, each.GetSignature())
	}
	return slim
}
