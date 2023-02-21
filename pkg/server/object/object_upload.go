package object

import (
	"encoding/json"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
)

type FunctionUploadUnit struct {
	WorkspaceConfig *WorkspaceConfig              `json:"workspace"`
	FunctionResult  *extractor.FunctionFileResult `json:"funcResult"`
}

type FunctionContextUploadUnit struct {
	WorkspaceConfig  *WorkspaceConfig          `json:"workspace"`
	FunctionContexts []*sibyl2.FunctionContext `json:"functionContext"`
}

type ClazzUploadUnit struct {
	WorkspaceConfig *WorkspaceConfig           `json:"workspace"`
	ClazzFileResult *extractor.ClazzFileResult `json:"clazzFileResult"`
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

func DeserializeClazzUploadUnit(data []byte) (*ClazzUploadUnit, error) {
	u := &ClazzUploadUnit{}
	err := json.Unmarshal(data, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
