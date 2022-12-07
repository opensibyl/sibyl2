package object

import (
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"github.com/williamfzc/sibyl2/pkg/server/binding"
)

type FunctionWithSignature struct {
	*sibyl2.FunctionWithPath
	Signature string `json:"signature"`
}

type FunctionUploadUnit struct {
	WorkspaceConfig *binding.WorkspaceConfig      `json:"workspace"`
	FunctionResult  *extractor.FunctionFileResult `json:"funcResult"`
}

type FuncContextUploadUnit struct {
	WorkspaceConfig  *binding.WorkspaceConfig  `json:"workspace"`
	FunctionContexts []*sibyl2.FunctionContext `json:"functionContext"`
}
