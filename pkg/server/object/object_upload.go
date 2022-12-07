package object

import (
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/extractor"
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
