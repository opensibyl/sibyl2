package object

import (
	"encoding/json"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
)

type FunctionWithSignature struct {
	*sibyl2.FunctionWithPath
	Signature string `json:"signature"`
}

type FunctionUploadUnit struct {
	WorkspaceConfig *WorkspaceConfig              `json:"workspace"`
	FunctionResult  *extractor.FunctionFileResult `json:"funcResult"`
}

type FunctionContextUploadUnit struct {
	WorkspaceConfig  *WorkspaceConfig          `json:"workspace"`
	FunctionContexts []*sibyl2.FunctionContext `json:"functionContext"`
}

func SerializeUploadUnit(u interface{}) ([]byte, error) {
	return json.Marshal(u)
}

func DeserializeFuncUploadUnit(data []byte) (*FunctionUploadUnit, error) {
	u := &FunctionUploadUnit{}
	err := json.Unmarshal(data, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func DeserializeFuncCtxUploadUnit(data []byte) (*FunctionContextUploadUnit, error) {
	u := &FunctionContextUploadUnit{}
	err := json.Unmarshal(data, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
